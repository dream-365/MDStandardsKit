package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func doc(status string, deps ...string) string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("type: spec\n")
	b.WriteString("status: " + status + "\n")
	b.WriteString("owner: example\n")
	if len(deps) > 0 {
		b.WriteString("depends_on:\n")
		for _, d := range deps {
			b.WriteString("  - " + d + "\n")
		}
	}
	b.WriteString("---\n\n# Body\n")
	return b.String()
}

func policyDir(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs("policy")
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

func runCLI(t *testing.T, root string, extra ...string) (int, string, string) {
	t.Helper()
	args := append([]string{"--root", root, "--policy", policyDir(t)}, extra...)
	var out, errb bytes.Buffer
	code := run(args, &out, &errb)
	return code, out.String(), errb.String()
}

func decodeReport(t *testing.T, s string) Report {
	t.Helper()
	var r Report
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, s)
	}
	return r
}

func TestCleanCorpusPasses(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, out, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s", code, errs)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 0 {
		t.Fatalf("findings = %v, want none", r.Findings)
	}
	if r.Stats.Edges != 1 {
		t.Fatalf("edges = %d, want 1", r.Stats.Edges)
	}
}

func TestMissingTargetBlocks(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 1 || r.Findings[0].Rule != "refs.missing_target" {
		t.Fatalf("findings = %+v, want one refs.missing_target", r.Findings)
	}
}

func TestNonDependableStatusBlocks(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Draft"))

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 1 || r.Findings[0].Rule != "refs.status_not_dependable" {
		t.Fatalf("findings = %+v, want one refs.status_not_dependable", r.Findings)
	}
}

func TestLowercaseStatusIsNotDependable(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("approved"))

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 1 {
		t.Fatalf("exit = %d, want 1 (exact literals only)", code)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 1 || r.Findings[0].Rule != "refs.status_not_dependable" {
		t.Fatalf("findings = %+v, want one refs.status_not_dependable", r.Findings)
	}
}

func TestNormalizedDependsOnResolves(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", `./docs/architecture/../architecture/003.md/`))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, out, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s out=%s", code, errs, out)
	}
}

func TestZeroGovernedDocumentsFailsClosed(t *testing.T) {
	root := t.TempDir()
	code, _, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(errs, "matched no files") {
		t.Fatalf("stderr = %q, want coverage error", errs)
	}
}

func TestZeroEdgesFailsClosed(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, _, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "zero dependency edges") {
		t.Fatalf("stderr = %q, want zero-edge coverage error", errs)
	}
}

func TestExtractErrorFailsClosedInPhase1(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	writeFile(t, root, "docs/specs/bad.md", "# no frontmatter\n")

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	r := decodeReport(t, out)
	if len(r.ExtractErrors) != 1 || r.ExtractErrors[0].Path != "docs/specs/bad.md" {
		t.Fatalf("extract errors = %+v", r.ExtractErrors)
	}
}

func TestExtractErrorIsReportOnlyInPhase0(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	writeFile(t, root, "docs/specs/bad.md", "# no frontmatter\n")

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "0")
	if code != 0 {
		t.Fatalf("exit = %d, want 0 in phase 0", code)
	}
	r := decodeReport(t, out)
	if len(r.ExtractErrors) != 1 {
		t.Fatalf("extract errors = %+v, want 1", r.ExtractErrors)
	}
}

func TestPhase0NeverBlocks(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "0")
	if code != 0 {
		t.Fatalf("exit = %d, want 0 in phase 0", code)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 1 {
		t.Fatalf("findings = %+v, want 1 reported but not blocking", r.Findings)
	}
}

func TestBaselineSuppressesKnownViolation(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	baseline := `[{"rule":"refs.missing_target","document":"docs/specs/001.md","target":"docs/architecture/missing.md"}]`
	basePath := filepath.Join(root, "baseline.json")
	if err := os.WriteFile(basePath, []byte(baseline), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "1", "--baseline", basePath)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 with baseline", code)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 0 {
		t.Fatalf("findings = %+v, want none after baseline", r.Findings)
	}
}

func TestBaselineDoesNotSuppressRenamedTarget(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	baseline := `[{"rule":"refs.missing_target","document":"docs/specs/001.md","target":"docs/architecture/old-name.md"}]`
	basePath := filepath.Join(root, "baseline.json")
	if err := os.WriteFile(basePath, []byte(baseline), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "1", "--baseline", basePath)
	if code != 1 {
		t.Fatalf("exit = %d, want 1 (rename must surface)", code)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 1 {
		t.Fatalf("findings = %+v, want one", r.Findings)
	}
}

func TestJSONSchemaAndStats(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	_, out, _ := runCLI(t, root, "--format", "json", "--phase", "0")
	var raw map[string]any
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"schema_version", "findings", "extract_errors", "stats"} {
		if _, ok := raw[k]; !ok {
			t.Fatalf("output has no %q: %s", k, out)
		}
	}
	stats, ok := raw["stats"].(map[string]any)
	if !ok || stats["edges_evaluated"].(float64) != 1 {
		t.Fatalf("stats = %v", raw["stats"])
	}
}

func TestExampleRepoPasses(t *testing.T) {
	root, err := filepath.Abs("testdata/example-repo")
	if err != nil {
		t.Fatal(err)
	}
	code, out, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s out=%s", code, errs, out)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 0 || len(r.ExtractErrors) != 0 {
		t.Fatalf("example repo not clean: findings=%+v errors=%+v", r.Findings, r.ExtractErrors)
	}
	if r.Stats.Edges != 1 {
		t.Fatalf("edges = %d, want 1", r.Stats.Edges)
	}
}

func TestPolicyFromHeadDowngradesOnDirtyBaseline(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("cp", "-R", "policy", filepath.Join(root, "policy")).CombinedOutput(); err != nil {
		t.Fatalf("copy policy: %v\n%s", err, out)
	}
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	baseline := `[{"rule":"refs.missing_target","document":"docs/specs/001.md","target":"docs/architecture/missing.md"}]`
	writeFile(t, root, "baseline.json", baseline)

	gitRun(t, root, "init")
	gitRun(t, root, "add", ".")
	gitRun(t, root, "commit", "-m", "init")

	// Dirty baseline: the worktree baseline differs from HEAD. The HEAD policy
	// check must downgrade to report-only instead of blocking on the empty file.
	writeFile(t, root, "baseline.json", "[]")
	code, out, errs := runCLI(t, root,
		"--policy", filepath.Join(root, "policy"),
		"--policy-from-head",
		"--phase", "1",
		"--baseline", filepath.Join(root, "baseline.json"),
		"--format", "json",
	)
	if code != 0 {
		t.Fatalf("dirty baseline: exit = %d, want 0; stderr=%s out=%s", code, errs, out)
	}
	if !strings.Contains(errs, "report-only") {
		t.Fatalf("dirty baseline: stderr = %q, want report-only warning", errs)
	}

	// Clean baseline: the committed suppression applies and the run stays clean.
	gitRun(t, root, "checkout", "--", "baseline.json")
	code, out, errs = runCLI(t, root,
		"--policy", filepath.Join(root, "policy"),
		"--policy-from-head",
		"--phase", "1",
		"--baseline", filepath.Join(root, "baseline.json"),
		"--format", "json",
	)
	if code != 0 {
		t.Fatalf("clean baseline: exit = %d, want 0; stderr=%s out=%s", code, errs, out)
	}
	if strings.Contains(errs, "report-only") {
		t.Fatalf("clean baseline: unexpected report-only warning: %s", errs)
	}
}

func TestQualityReportUnderDocsQualityIsNotGoverned(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	// A Markdown report without frontmatter under docs/quality/ must not be
	// extracted as a governed document. Placing it under docs/specs/ would.
	writeFile(t, root, "docs/quality/001/quality-report.md", "# Spec Quality Report\n\nNo frontmatter.\n")

	code, out, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s out=%s", code, errs, out)
	}
	r := decodeReport(t, out)
	if len(r.ExtractErrors) != 0 {
		t.Fatalf("quality report was extracted: %+v", r.ExtractErrors)
	}
}

func gitRun(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
