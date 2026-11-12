package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/open-policy-agent/opa/rego"
)

const (
	agentHookMaxDocs = 200
	multiModeName    = "multi"
)

var (
	version    = "dev"
	policyHash = "unknown"

	// Test seams: production uses wall-clock time and the documented budgets.
	elapsedSince      = time.Since
	agentHookBudget   = 2 * time.Second
	policyEvalTimeout = 10 * time.Second
)

// performanceOverBudget reports whether an Agent-hook-sized corpus exceeded the
// per-run time budget. A single process cannot measure a p95; this guard only
// catches the current run. The p95 target is measured by the opt-in benchmark.
func performanceOverBudget(elapsed time.Duration, matched int) bool {
	return matched > 0 && matched <= agentHookMaxDocs && elapsed > agentHookBudget
}

type Standard struct {
	StandardVersion    string                  `json:"standard_version"`
	GovernedGlobs      []string                `json:"governed_globs"`
	DependableStatuses []string                `json:"dependable_statuses"`
	Identity           string                  `json:"identity"`
	DocumentTypes      map[string]DocumentType `json:"document_types"`
}

// DocumentType is the policy profile for a governed document type. The spec
// type keeps using the top-level governed_globs; other types carry their own
// globs and structural facts.
type DocumentType struct {
	Globs              []string `json:"globs"`
	RequiredFields     []string `json:"required_fields"`
	StatusValues       []string `json:"status_values"`
	RequiredHeadings   []string `json:"required_headings"`
	IDPattern          string   `json:"id_pattern"`
	DateField          string   `json:"date_field"`
	Graph              bool     `json:"graph"`
	RuleSet            string   `json:"rule_set"`
	PlaceholderMarkers []string `json:"placeholder_markers"`
}

type repoMode int

const (
	legacyMode repoMode = iota
	multiMode
)

type repoConfigFile struct {
	Type  *string   `json:"type"`
	Types *[]string `json:"types"`
}

type RepoConfig struct {
	Mode  repoMode
	Type  string
	Types []string
}

type runPlan struct {
	Mode          repoMode
	LegacyType    string
	LegacyProfile DocumentType
	Globs         []string
	Graph         bool
	Profiles      map[string]DocumentType
	SpecGlobs     []string
}

type standardFile struct {
	Standard Standard `json:"standard"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("mdstandards", flag.ContinueOnError)
	fs.SetOutput(stderr)
	showVersion := fs.Bool("version", false, "print version and exit")
	root := fs.String("root", ".", "repository root")
	policyDirFlag := fs.String("policy", "policy", "policy bundle directory")
	format := fs.String("format", "table", "output format: table|json|github")
	baselinePath := fs.String("baseline", "", "baseline file of accepted violations")
	phase := fs.Int("phase", 1, "rollout phase: 0 report-only, 1 baseline+block, 2 full block")
	outPath := fs.String("out", "", "write the report to this file instead of stdout")
	dumpCorpus := fs.Bool("dump-corpus", false, "print the extracted corpus and exit")
	policyFromHead := fs.Bool("policy-from-head", false, "load the policy bundle from git HEAD")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showVersion {
		fmt.Fprintf(stdout, "mdstandards %s policy %s\n", version, policyHash)
		return 0
	}
	if *phase < 0 || *phase > 2 {
		fmt.Fprintln(stderr, "phase must be 0, 1 or 2")
		return 2
	}

	// The budget covers policy loading, extraction, coverage assertions, policy
	// evaluation and sorting. Report rendering and baseline filtering happen
	// after the measurement.
	checkStart := time.Now()

	rootAbs, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	// --policy resolves against the current working directory, like any other
	// CLI flag; --root only scopes which files are governed.
	policyDir, err := filepath.Abs(*policyDirFlag)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	dirtyPolicy := false
	if *policyFromHead {
		tmpDir, dirty, cleanup, err := materializePolicyFromHEAD(rootAbs, policyDir, *baselinePath)
		if err != nil {
			fmt.Fprintf(stderr, "policy.load_error: %v\n", err)
			return 2
		}
		defer cleanup()
		policyDir = tmpDir
		dirtyPolicy = dirty
		if dirty {
			fmt.Fprintln(stderr, "warning: policy or baseline files have uncommitted changes; running report-only")
		}
	}

	std, err := loadStandard(policyDir)
	if err != nil {
		fmt.Fprintf(stderr, "policy.load_error: %v\n", err)
		return 2
	}
	if std.StandardVersion != "1" && std.StandardVersion != "2" {
		fmt.Fprintf(stderr, "policy.load_error: unsupported standard_version %q\n", std.StandardVersion)
		return 2
	}

	var cfg RepoConfig
	if *policyFromHead {
		var configDirty bool
		cfg, configDirty, err = loadRepoConfigFromHEAD(rootAbs)
		if err != nil {
			fmt.Fprintf(stderr, "policy.load_error: %v\n", err)
			return 2
		}
		if configDirty {
			dirtyPolicy = true
			fmt.Fprintln(stderr, "warning: .mdstandards.json differs from HEAD; running report-only")
		}
	} else {
		cfg, err = loadRepoConfig(rootAbs)
		if err != nil {
			fmt.Fprintf(stderr, "policy.load_error: %v\n", err)
			return 2
		}
	}

	plan, err := buildRunPlan(std, cfg)
	if err != nil {
		fmt.Fprintf(stderr, "policy.load_error: %v\n", err)
		return 2
	}
	if plan.Mode == multiMode && std.StandardVersion != "2" {
		fmt.Fprintln(stderr, "policy.load_error: multi-type mode requires standard_version 2")
		return 2
	}

	var corpus Corpus
	if plan.Mode == multiMode {
		corpus, err = ExtractMulti(rootAbs, plan.Globs, plan.Profiles, plan.SpecGlobs)
	} else {
		corpus, err = Extract(rootAbs, plan.Globs, plan.LegacyType, plan.LegacyProfile)
		if err == nil && plan.LegacyType != "spec" {
			corpus.DocumentType = plan.LegacyType
		}
	}
	if err != nil {
		fmt.Fprintf(stderr, "extract error: %v\n", err)
		return 2
	}
	if *dumpCorpus {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(corpus)
		return 0
	}

	// Coverage assertions. A gate that scanned nothing must not report success.
	if corpus.Stats.Matched == 0 {
		fmt.Fprintf(stderr, "coverage error: governed_globs matched no files: %v\n", plan.Globs)
		return 2
	}
	if *phase >= 1 && corpus.Stats.Parsed == 0 {
		fmt.Fprintln(stderr, "coverage error: no governed document parsed")
		return 2
	}
	graphDocs := graphDocumentCount(corpus, plan)
	if *phase >= 1 && plan.Graph && graphDocs > 0 && corpus.Stats.Edges == 0 {
		fmt.Fprintln(stderr, "coverage error: the corpus has zero dependency edges")
		return 2
	}
	if plan.Mode == multiMode {
		warnEmptyTypes(stderr, corpus, plan)
	}

	findings, err := evalPolicy(context.Background(), policyDir, corpus)
	if err != nil {
		fmt.Fprintf(stderr, "policy.load_error: %v\n", err)
		return 2
	}
	sortFindings(findings)

	elapsed := elapsedSince(checkStart)
	overBudget := performanceOverBudget(elapsed, corpus.Stats.Matched)
	performanceDegrade := overBudget && *phase == 1
	if overBudget {
		action := "this phase does not degrade"
		if *phase == 1 {
			action = "phase 1 degrades to warning-only"
		}
		fmt.Fprintf(stderr,
			"warning: performance budget exceeded: %s for %d matched documents (limit %s); %s\n",
			elapsed.Round(time.Millisecond), corpus.Stats.Matched, agentHookBudget, action)
	}

	rep := Report{
		SchemaVersion: findingsSchemaVersion,
		Findings:      findings,
		ExtractErrors: corpus.ExtractErrors,
		Stats:         corpus.Stats,
	}

	blocking := *phase >= 1 && !dirtyPolicy && !performanceDegrade
	if blocking && *baselinePath != "" {
		base, err := loadBaseline(*baselinePath)
		if err != nil {
			fmt.Fprintf(stderr, "baseline error: %v\n", err)
			return 2
		}
		var suppressed int
		rep.Findings, suppressed = filterBaseline(rep.Findings, base)
		_ = suppressed
	}

	var out io.Writer = stdout
	if *outPath != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		defer f.Close()
		out = f
	}

	switch *format {
	case "json":
		if err := writeJSON(out, rep); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
	case "github":
		writeGitHub(out, rep)
	case "table":
		writeTable(out, rep)
	default:
		fmt.Fprintf(stderr, "unknown format %q\n", *format)
		return 2
	}

	if !blocking {
		return 0
	}
	if len(rep.ExtractErrors) > 0 {
		return 2
	}
	if len(rep.Findings) > 0 {
		return 1
	}
	return 0
}

func loadStandard(policyDir string) (Standard, error) {
	b, err := os.ReadFile(filepath.Join(policyDir, "data", "standard.json"))
	if err != nil {
		return Standard{}, err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var sf standardFile
	if err := dec.Decode(&sf); err != nil {
		return Standard{}, fmt.Errorf("invalid JSON: %w", err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		return Standard{}, fmt.Errorf("invalid JSON: trailing data")
	}
	return sf.Standard, nil
}

// loadRepoConfig reads the repository-level type declaration. A missing file
// means legacy spec. A present file must declare exactly one non-null form.
func loadRepoConfig(root string) (RepoConfig, error) {
	b, err := os.ReadFile(filepath.Join(root, ".mdstandards.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return RepoConfig{Mode: legacyMode, Type: "spec"}, nil
		}
		return RepoConfig{}, err
	}
	return loadRepoConfigBytes(b)
}

func loadRepoConfigBytes(b []byte) (RepoConfig, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var c repoConfigFile
	if err := dec.Decode(&c); err != nil {
		return RepoConfig{}, fmt.Errorf("invalid .mdstandards.json: %w", err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		return RepoConfig{}, fmt.Errorf("invalid .mdstandards.json: trailing data")
	}
	if c.Type == nil && c.Types == nil {
		return RepoConfig{}, fmt.Errorf(".mdstandards.json: one of type or types is required")
	}
	if c.Type != nil && c.Types != nil {
		return RepoConfig{}, fmt.Errorf(".mdstandards.json: type and types are mutually exclusive")
	}
	if c.Type != nil {
		t := *c.Type
		if t == "" {
			return RepoConfig{}, fmt.Errorf(".mdstandards.json: type is required")
		}
		if t != "spec" && t != "idea" {
			return RepoConfig{}, fmt.Errorf(".mdstandards.json: type %q requires multi-type mode; use types", t)
		}
		return RepoConfig{Mode: legacyMode, Type: t}, nil
	}

	types := *c.Types
	if len(types) == 0 {
		return RepoConfig{}, fmt.Errorf(".mdstandards.json: types must not be empty")
	}
	seen := make(map[string]bool, len(types))
	out := make([]string, 0, len(types))
	for _, raw := range types {
		t := strings.TrimSpace(raw)
		if t == "" {
			return RepoConfig{}, fmt.Errorf(".mdstandards.json: types entries must be non-empty")
		}
		if seen[t] {
			return RepoConfig{}, fmt.Errorf(".mdstandards.json: duplicate type %q", t)
		}
		seen[t] = true
		out = append(out, t)
	}
	return RepoConfig{Mode: multiMode, Types: out}, nil
}

func buildRunPlan(std Standard, cfg RepoConfig) (runPlan, error) {
	if cfg.Mode == legacyMode {
		return buildLegacyPlan(std, cfg)
	}
	return buildMultiPlan(std, cfg)
}

func buildLegacyPlan(std Standard, cfg RepoConfig) (runPlan, error) {
	p := runPlan{Mode: legacyMode, LegacyType: cfg.Type}
	switch cfg.Type {
	case "spec":
		if len(std.DependableStatuses) == 0 {
			return runPlan{}, fmt.Errorf("dependable_statuses is empty")
		}
		if len(std.GovernedGlobs) == 0 {
			return runPlan{}, fmt.Errorf("governed_globs is empty")
		}
		p.Globs = std.GovernedGlobs
		p.Graph = true
		p.SpecGlobs = std.GovernedGlobs
		p.LegacyProfile = DocumentType{Graph: true, RuleSet: "spec"}
	case "idea":
		profile, ok := std.DocumentTypes["idea"]
		if !ok {
			return runPlan{}, fmt.Errorf("no document type %q in standard.json", cfg.Type)
		}
		if len(profile.Globs) == 0 {
			return runPlan{}, fmt.Errorf("document type %q has no globs", cfg.Type)
		}
		p.Globs = profile.Globs
		p.Graph = profile.Graph
		p.LegacyProfile = profile
	default:
		return runPlan{}, fmt.Errorf("legacy type %q requires multi-type mode", cfg.Type)
	}
	return p, nil
}

func buildMultiPlan(std Standard, cfg RepoConfig) (runPlan, error) {
	if len(cfg.Types) == 0 {
		return runPlan{}, fmt.Errorf("types must not be empty")
	}
	if len(std.DependableStatuses) == 0 {
		return runPlan{}, fmt.Errorf("dependable_statuses is empty")
	}

	profiles := make(map[string]DocumentType, len(cfg.Types))
	globs := make([]string, 0, len(std.GovernedGlobs))
	seenGlobs := make(map[string]bool)
	addGlob := func(g string) {
		if g == "" || seenGlobs[g] {
			return
		}
		seenGlobs[g] = true
		globs = append(globs, g)
	}
	// The top-level baseline always participates in discovery, even when spec
	// is not enabled. Files of unenabled types then surface as routing errors
	// instead of disappearing.
	for _, g := range std.GovernedGlobs {
		addGlob(g)
	}

	seenTypes := make(map[string]bool, len(cfg.Types))
	graph := false
	for _, t := range cfg.Types {
		if seenTypes[t] {
			return runPlan{}, fmt.Errorf("duplicate type %q", t)
		}
		seenTypes[t] = true
		profile, ok := std.DocumentTypes[t]
		if !ok {
			return runPlan{}, fmt.Errorf("no document type %q in standard.json", t)
		}
		switch profile.RuleSet {
		case "spec", "idea", "none":
		case "":
			return runPlan{}, fmt.Errorf("document type %q has no rule_set", t)
		default:
			return runPlan{}, fmt.Errorf("document type %q has invalid rule_set %q", t, profile.RuleSet)
		}
		if profile.Graph && profile.RuleSet != "spec" {
			return runPlan{}, fmt.Errorf("document type %q: graph requires rule_set spec", t)
		}
		if t == "spec" {
			if len(std.GovernedGlobs) == 0 {
				return runPlan{}, fmt.Errorf("spec is enabled but governed_globs is empty")
			}
		} else {
			if len(profile.Globs) == 0 {
				return runPlan{}, fmt.Errorf("document type %q has no globs", t)
			}
			for _, g := range profile.Globs {
				addGlob(g)
			}
		}
		if profile.RuleSet == "idea" {
			if err := validateIdeaProfile(t, &profile); err != nil {
				return runPlan{}, err
			}
		}
		if profile.Graph {
			graph = true
		}
		profiles[t] = profile
	}
	if len(globs) == 0 {
		return runPlan{}, fmt.Errorf("discovery globs are empty")
	}
	return runPlan{
		Mode:      multiMode,
		Globs:     globs,
		Graph:     graph,
		Profiles:  profiles,
		SpecGlobs: std.GovernedGlobs,
	}, nil
}

func validateIdeaProfile(name string, profile *DocumentType) error {
	if len(profile.RequiredFields) == 0 {
		return fmt.Errorf("document type %q: idea profile needs required_fields", name)
	}
	if len(profile.StatusValues) == 0 {
		return fmt.Errorf("document type %q: idea profile needs status_values", name)
	}
	if len(profile.RequiredHeadings) == 0 {
		return fmt.Errorf("document type %q: idea profile needs required_headings", name)
	}
	if profile.IDPattern == "" {
		return fmt.Errorf("document type %q: idea profile needs id_pattern", name)
	}
	if len(profile.PlaceholderMarkers) == 0 {
		return fmt.Errorf("document type %q: idea profile needs placeholder_markers", name)
	}
	if profile.DateField == "" {
		profile.DateField = "date"
	}
	return nil
}

func graphDocumentCount(corpus Corpus, plan runPlan) int {
	if plan.Mode != multiMode {
		if plan.Graph {
			return len(corpus.Documents)
		}
		return 0
	}
	n := 0
	for _, doc := range corpus.Documents {
		if p, ok := plan.Profiles[doc.Type]; ok && p.Graph {
			n++
		}
	}
	return n
}

func warnEmptyTypes(stderr io.Writer, corpus Corpus, plan runPlan) {
	counts := make(map[string]int, len(plan.Profiles))
	for _, doc := range corpus.Documents {
		counts[doc.Type]++
	}
	var empty []string
	for t := range plan.Profiles {
		if counts[t] == 0 {
			empty = append(empty, t)
		}
	}
	sort.Strings(empty)
	for _, t := range empty {
		fmt.Fprintf(stderr, "warning: enabled document type %q matched no parsed documents\n", t)
	}
}

func loadRepoConfigFromHEAD(root string) (RepoConfig, bool, error) {
	headBytes, headExists, err := gitShowFile(root, ".mdstandards.json")
	if err != nil {
		return RepoConfig{}, false, fmt.Errorf("git show HEAD:.mdstandards.json: %w", err)
	}
	workBytes, workErr := os.ReadFile(filepath.Join(root, ".mdstandards.json"))
	workExists := workErr == nil
	if workErr != nil && !os.IsNotExist(workErr) {
		return RepoConfig{}, false, workErr
	}

	dirty := false
	switch {
	case headExists && !workExists:
		dirty = true
	case headExists && workExists && !bytes.Equal(headBytes, workBytes):
		dirty = true
	case !headExists && workExists:
		dirty = true
	}

	if !headExists {
		return RepoConfig{Mode: legacyMode, Type: "spec"}, dirty, nil
	}
	cfg, err := loadRepoConfigBytes(headBytes)
	if err != nil {
		return RepoConfig{}, dirty, err
	}
	return cfg, dirty, nil
}

func gitShowFile(root, rel string) ([]byte, bool, error) {
	out, err := exec.Command("git", "-C", root, "show", "HEAD:"+rel).Output()
	if err == nil {
		return out, true, nil
	}
	if exec.Command("git", "-C", root, "cat-file", "-e", "HEAD:"+rel).Run() != nil {
		return nil, false, nil
	}
	return nil, false, err
}

func evalPolicy(parent context.Context, policyDir string, corpus Corpus) ([]Finding, error) {
	ctx, cancel := context.WithTimeout(parent, policyEvalTimeout)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	paths, err := filepath.Glob(filepath.Join(policyDir, "*.rego"))
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no .rego files in %s", policyDir)
	}
	paths = append(paths, filepath.Join(policyDir, "data", "standard.json"))

	q, err := rego.New(
		rego.Query("data.mdstandards.deny"),
		rego.Load(paths, nil),
		rego.Input(corpus),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	rs, err := q.Eval(ctx)
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return nil, nil
	}
	raw, ok := rs[0].Expressions[0].Value.([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected policy result type %T", rs[0].Expressions[0].Value)
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var out []Finding
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func loadBaseline(path string) ([]BaselineEntry, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []BaselineEntry
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// gitPathDirty reports whether path has uncommitted changes in the git worktree
// rooted at root. Paths outside root are treated as dirty: the HEAD check cannot
// verify them, and the safe default is report-only.
func gitPathDirty(root, path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return true
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return true
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return true
	}
	status, err := exec.Command("git", "-C", root, "status", "--porcelain", "--", filepath.ToSlash(rel)).Output()
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(status)) != ""
}

// materializePolicyFromHEAD writes the committed policy bundle into a temp
// directory so an Agent editing the working tree cannot change the rules that
// judge it (D15). The dirty flag tells the caller to run report-only.
func materializePolicyFromHEAD(root, policyDir, baselinePath string) (string, bool, func(), error) {
	noop := func() {}
	rel, err := filepath.Rel(root, policyDir)
	if err != nil {
		return "", false, noop, err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false, noop, fmt.Errorf("--policy must be inside --root when --policy-from-head is set")
	}
	rel = filepath.ToSlash(rel)

	dirty := gitPathDirty(root, policyDir)
	if baselinePath != "" && gitPathDirty(root, baselinePath) {
		dirty = true
	}

	list, err := exec.Command("git", "-C", root, "ls-tree", "-r", "--name-only", "HEAD", "--", rel).Output()
	if err != nil {
		return "", dirty, noop, fmt.Errorf("git ls-tree: %w", err)
	}

	tmp, err := os.MkdirTemp("", "mdstandards-policy-")
	if err != nil {
		return "", dirty, noop, err
	}
	cleanup := func() { _ = os.RemoveAll(tmp) }

	for _, name := range strings.Split(strings.TrimSpace(string(list)), "\n") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		content, err := exec.Command("git", "-C", root, "show", "HEAD:"+name).Output()
		if err != nil {
			cleanup()
			return "", dirty, noop, fmt.Errorf("git show %s: %w", name, err)
		}
		dst := filepath.Join(tmp, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			cleanup()
			return "", dirty, noop, err
		}
		if err := os.WriteFile(dst, content, 0o644); err != nil {
			cleanup()
			return "", dirty, noop, err
		}
	}
	return filepath.Join(tmp, filepath.FromSlash(rel)), dirty, cleanup, nil
}
