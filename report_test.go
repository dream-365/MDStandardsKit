package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSortFindingsOrder(t *testing.T) {
	findings := []Finding{
		{Rule: "refs.missing_target", Document: "b", Target: "z", Message: "b/z"},
		{Rule: "refs.missing_target", Document: "a", Target: "y", Message: "a/y"},
		{Rule: "policy.load_error", Message: "policy"},
		{Rule: "refs.missing_target", Document: "a", Target: "x", Message: "a/x first"},
		{Rule: "refs.missing_target", Document: "a", Target: "x", Message: "a/x second"},
	}
	sortFindings(findings)
	want := []string{"policy", "a/x first", "a/x second", "a/y", "b/z"}
	for i, w := range want {
		if findings[i].Message != w {
			t.Fatalf("findings[%d] = %q, want %q; all=%+v", i, findings[i].Message, w, findings)
		}
	}
}

func TestGitHubFormat(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/specs/bad.md", "# no frontmatter\n")

	code, out, errs := runCLI(t, root, "--format", "github", "--phase", "0")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s", code, errs)
	}
	for _, want := range []string{
		"::error file=docs/specs/bad.md::[extract.parse_error]",
		"::error file=docs/specs/001.md::[refs.missing_target]",
		"::notice::scanned 1 governed documents, evaluated 1 dependency edges",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}

func TestOutWritesReportToFile(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	outPath := filepath.Join(t.TempDir(), "findings.json")

	code, stdout, errs := runCLI(t, root, "--format", "json", "--phase", "0", "--out", outPath)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s", code, errs)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty when -out is set", stdout)
	}
	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	r := decodeReport(t, string(b))
	if len(r.Findings) != 1 {
		t.Fatalf("findings = %+v, want 1", r.Findings)
	}
}

func TestJSONByteStableAcrossRuns(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	var first string
	for i := 0; i < 3; i++ {
		code, out, errs := runCLI(t, root, "--format", "json", "--phase", "0")
		if code != 0 {
			t.Fatalf("run %d: exit = %d; stderr=%s", i, code, errs)
		}
		if i == 0 {
			first = out
			continue
		}
		if out != first {
			t.Fatalf("run %d JSON differs:\nfirst=%s\nnow=%s", i, first, out)
		}
	}
}
