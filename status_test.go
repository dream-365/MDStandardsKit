package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestMissingStatusIsNotDependable(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", "---\ntype: architecture\nowner: example\n---\n\n# Body\n")

	code, out, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stderr=%s out=%s", code, errs, out)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 1 || r.Findings[0].Rule != "refs.status_not_dependable" {
		t.Fatalf("findings = %+v, want refs.status_not_dependable", r.Findings)
	}
}

func TestStatusRegoGuardEmitsPolicyLoadError(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0o755); err != nil {
		t.Fatal(err)
	}
	regos, err := filepath.Glob(filepath.Join("policy", "*.rego"))
	if err != nil {
		t.Fatal(err)
	}
	for _, src := range regos {
		b, err := os.ReadFile(src)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, filepath.Base(src)), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// standard.json without dependable_statuses keeps the Rego guard honest.
	if err := os.WriteFile(filepath.Join(dir, "data", "standard.json"),
		[]byte(`{"standard":{"standard_version":"1","governed_globs":["docs/**/*.md"],"identity":"test"}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	findings, err := evalPolicy(context.Background(), dir, Corpus{})
	if err != nil {
		t.Fatalf("evalPolicy: %v", err)
	}
	for _, f := range findings {
		if f.Rule == "policy.load_error" {
			return
		}
	}
	t.Fatalf("findings = %+v, want policy.load_error", findings)
}
