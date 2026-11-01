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

type Standard struct {
	StandardVersion    string   `json:"standard_version"`
	GovernedGlobs      []string `json:"governed_globs"`
	DependableStatuses []string `json:"dependable_statuses"`
	Identity           string   `json:"identity"`
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
	if *phase < 0 || *phase > 2 {
		fmt.Fprintln(stderr, "phase must be 0, 1 or 2")
		return 2
	}

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
	if len(std.DependableStatuses) == 0 {
		fmt.Fprintln(stderr, "policy.load_error: dependable_statuses is empty")
		return 2
	}
	if len(std.GovernedGlobs) == 0 {
		fmt.Fprintln(stderr, "policy.load_error: governed_globs is empty")
		return 2
	}

	corpus, err := Extract(rootAbs, std.GovernedGlobs)
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
		fmt.Fprintf(stderr, "coverage error: governed_globs matched no files: %v\n", std.GovernedGlobs)
		return 2
	}
	if *phase >= 1 && corpus.Stats.Parsed == 0 {
		fmt.Fprintln(stderr, "coverage error: no governed document parsed")
		return 2
	}
	if *phase >= 1 && corpus.Stats.Edges == 0 {
		fmt.Fprintln(stderr, "coverage error: the corpus has zero dependency edges")
		return 2
	}

	findings, err := evalPolicy(policyDir, corpus)
	if err != nil {
		fmt.Fprintf(stderr, "policy.load_error: %v\n", err)
		return 2
	}
	sortFindings(findings)

	rep := Report{
		SchemaVersion: findingsSchemaVersion,
		Findings:      findings,
		ExtractErrors: corpus.ExtractErrors,
		Stats:         corpus.Stats,
	}

	blocking := *phase >= 1 && !dirtyPolicy
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

func evalPolicy(policyDir string, corpus Corpus) ([]Finding, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
