package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func findingFor(findings []Finding, rule, document string) bool {
	for _, f := range findings {
		if f.Rule == rule && f.Document == document {
			return true
		}
	}
	return false
}

func TestMultiRegoScope(t *testing.T) {
	corpus := Corpus{
		Mode: "multi",
		Documents: []Document{
			{Path: "docs/specs/001.md", Type: "spec", Status: "Approved",
				DependsOn: []string{"docs/architecture/missing.md"}},
			{Path: "docs/specs/002/changelog.md", Type: "changelog",
				DependsOn: []string{"docs/architecture/missing.md"}},
			{Path: "2026/2026-11-01/idea.md", Type: "idea", Status: "萌芽"},
		},
	}
	findings, err := evalPolicy(context.Background(), policyDir(t), corpus)
	if err != nil {
		t.Fatal(err)
	}
	if !findingFor(findings, "refs.missing_target", "docs/specs/001.md") {
		t.Fatalf("spec missing_target not produced: %+v", findings)
	}
	if findingFor(findings, "refs.missing_target", "docs/specs/002/changelog.md") {
		t.Fatalf("none-scope document produced refs finding: %+v", findings)
	}
	if !findingFor(findings, "idea.empty_title", "2026/2026-11-01/idea.md") {
		t.Fatalf("idea rule did not run for idea doc: %+v", findings)
	}
}

// TestLegacySpecStatusRegression guards the legacy profile_for fallback: an
// undefined profile_for would make every dependency look non-dependable.
func TestLegacySpecStatusRegression(t *testing.T) {
	corpus := Corpus{
		Documents: []Document{
			{Path: "docs/specs/001.md", Type: "spec",
				DependsOn: []string{"docs/architecture/003.md"}},
			{Path: "docs/architecture/003.md", Type: "architecture", Status: "Approved"},
		},
	}
	findings, err := evalPolicy(context.Background(), policyDir(t), corpus)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range findings {
		if f.Rule == "refs.status_not_dependable" {
			t.Fatalf("legacy spec status regression: %+v", findings)
		}
	}
}

func TestMultiMissingRuleSetGuard(t *testing.T) {
	root := t.TempDir()
	copyPolicy(t, root)
	writeFile(t, root, "policy/data/standard.json", `{
  "standard": {
    "standard_version": "2",
    "governed_globs": ["docs/**/*.md"],
    "dependable_statuses": ["Approved"],
    "identity": "test",
    "document_types": {
      "changelog": {"globs": ["docs/**/*.md"]}
    }
  }
}`)
	corpus := Corpus{
		Mode:      "multi",
		Documents: []Document{{Path: "docs/specs/001/changelog.md", Type: "changelog"}},
	}
	findings, err := evalPolicy(context.Background(), filepath.Join(root, "policy"), corpus)
	if err != nil {
		t.Fatal(err)
	}
	if !findingFor(findings, "policy.load_error", "docs/specs/001/changelog.md") {
		t.Fatalf("findings = %+v, want policy.load_error", findings)
	}
}

func TestMultiEmptyDocumentTypesGuard(t *testing.T) {
	root := t.TempDir()
	copyPolicy(t, root)
	writeFile(t, root, "policy/data/standard.json", `{
  "standard": {
    "standard_version": "2",
    "governed_globs": ["docs/**/*.md"],
    "dependable_statuses": ["Approved"],
    "identity": "test",
    "document_types": {}
  }
}`)
	findings, err := evalPolicy(context.Background(), filepath.Join(root, "policy"), Corpus{Mode: "multi"})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, f := range findings {
		if f.Rule == "policy.load_error" {
			found = true
		}
	}
	if !found {
		t.Fatalf("findings = %+v, want policy.load_error", findings)
	}
}

func TestMultiRegoCompile(t *testing.T) {
	if _, err := evalPolicy(context.Background(), policyDir(t), Corpus{Mode: "multi"}); err != nil {
		t.Fatalf("multi policy does not compile: %v", err)
	}
}

func TestMultiPolicyFilesExist(t *testing.T) {
	if _, err := os.Stat(filepath.Join(policyDir(t), "types.rego")); err != nil {
		t.Fatalf("types.rego missing: %v", err)
	}
}
