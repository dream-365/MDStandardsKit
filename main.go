package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-policy-agent/opa/rego"
)

const agentHookMaxDocs = 200

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
	PlaceholderMarkers []string `json:"placeholder_markers"`
}

type repoConfig struct {
	Type string `json:"type"`
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
	if std.StandardVersion != "1" {
		fmt.Fprintf(stderr, "policy.load_error: unsupported standard_version %q\n", std.StandardVersion)
		return 2
	}

	docType, err := loadRepoConfig(rootAbs)
	if err != nil {
		fmt.Fprintf(stderr, "policy.load_error: %v\n", err)
		return 2
	}

	globs := std.GovernedGlobs
	graph := true
	var profile DocumentType
	if docType == "spec" {
		if len(std.DependableStatuses) == 0 {
			fmt.Fprintln(stderr, "policy.load_error: dependable_statuses is empty")
			return 2
		}
		if len(std.GovernedGlobs) == 0 {
			fmt.Fprintln(stderr, "policy.load_error: governed_globs is empty")
			return 2
		}
	} else {
		p, ok := std.DocumentTypes[docType]
		if !ok {
			fmt.Fprintf(stderr, "policy.load_error: no document type %q in standard.json\n", docType)
			return 2
		}
		profile = p
		globs = p.Globs
		graph = p.Graph
		if len(globs) == 0 {
			fmt.Fprintf(stderr, "policy.load_error: document type %q has no globs\n", docType)
			return 2
		}
	}

	corpus, err := Extract(rootAbs, globs, docType, profile)
	if err != nil {
		fmt.Fprintf(stderr, "extract error: %v\n", err)
		return 2
	}
	if docType != "spec" {
		corpus.DocumentType = docType
	}
	if *dumpCorpus {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(corpus)
		return 0
	}

	// Coverage assertions. A gate that scanned nothing must not report success.
	if corpus.Stats.Matched == 0 {
		fmt.Fprintf(stderr, "coverage error: governed_globs matched no files: %v\n", std.GovernedGlobs)
		return 2
	}
	if *phase >= 1 && corpus.Stats.Parsed == 0 {
		fmt.Fprintln(stderr, "coverage error: no governed document parsed")
		return 2
	}
	if *phase >= 1 && graph && corpus.Stats.Edges == 0 {
		fmt.Fprintln(stderr, "coverage error: the corpus has zero dependency edges")
		return 2
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
	var sf standardFile
	if err := json.Unmarshal(b, &sf); err != nil {
		return Standard{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return sf.Standard, nil
}

// loadRepoConfig reads the repository-level type declaration. A missing file
// means spec, preserving the pre-idea behavior. Malformed JSON, unknown keys,
// an empty type, or an unknown type fail closed before extraction.
func loadRepoConfig(root string) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, ".mdstandards.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return "spec", nil
		}
		return "", err
	}
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	var c repoConfig
	if err := dec.Decode(&c); err != nil {
		return "", fmt.Errorf("invalid .mdstandards.json: %w", err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		return "", fmt.Errorf("invalid .mdstandards.json: trailing data")
	}
	if c.Type == "" {
		return "", fmt.Errorf(".mdstandards.json: type is required")
	}
	return c.Type, nil
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
// judge it (D11). The dirty flag tells the caller to run report-only.
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
