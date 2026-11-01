package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const miniRubric = `version: "test"
dimensions:
  - id: test
    name: Test
    weight: 1
    items:
      - id: test.item
        question: "test"
        tier: core
        weight: 1
        section: "1. Test"
        aliases: ["Test"]
        na_allowed: false
        levels:
          "0": "none"
          "1": "partial"
          "2": "full"
`

const miniDoc = `## 1. Test

Evidence here.
`

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func baseFindings(r Rubric, rubricPath, docPath string) Findings {
	rubricBytes, _ := os.ReadFile(rubricPath)
	docBytes, _ := os.ReadFile(docPath)
	f := Findings{
		SchemaVersion: 1,
		RubricVersion: r.Version,
		RubricSHA256:  sha256Hex(rubricBytes),
		Target:        "spec.md",
		TargetSHA256:  sha256Hex(docBytes),
		Reviewer: Reviewer{
			Role:      "independent",
			Model:     "test",
			SessionID: "reviewer-1",
			StartedAt: "2026-11-01T00:00:00Z",
		},
	}
	for _, d := range r.Dimensions {
		for _, it := range d.Items {
			f.Items = append(f.Items, ItemFinding{
				ID:        it.ID,
				Status:    "scored",
				Level:     0,
				Rationale: "fixture",
				Missing:   "not provided in this fixture",
			})
		}
	}
	return f
}

func writeFindings(t *testing.T, f Findings) string {
	t.Helper()
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	return writeTemp(t, "findings.json", string(b))
}

func setItem(f *Findings, id string, mut func(*ItemFinding)) {
	for i := range f.Items {
		if f.Items[i].ID == id {
			mut(&f.Items[i])
			return
		}
	}
}

func TestValidateRealRubric(t *testing.T) {
	path := filepath.Join("..", "..", ".claude", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
	r, err := LoadRubric(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateRubric(r); err != nil {
		t.Fatalf("real rubric invalid: %v", err)
	}
}

func TestValidateRejectsCoreNA(t *testing.T) {
	r := Rubric{Version: "x", Dimensions: []Dimension{{
		ID: "d", Name: "D", Weight: 1,
		Items: []Item{{
			ID: "i", Tier: "core", Weight: 1, Section: "1. X", Aliases: []string{"X"}, NaAllowed: true,
			Levels: map[string]string{"0": "a", "1": "b", "2": "c"},
		}},
	}}}
	if err := ValidateRubric(r); err == nil {
		t.Fatal("expected error for core item with na_allowed true")
	}
}

func TestAuditAllZeroSucceeds(t *testing.T) {
	rubricPath := filepath.Join("..", "..", ".claude", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
	docPath := writeTemp(t, "spec.md", "# Spec\n\nNothing here.\n")
	r, err := LoadRubric(rubricPath)
	if err != nil {
		t.Fatal(err)
	}
	f := baseFindings(r, rubricPath, docPath)
	findingsPath := writeFindings(t, f)
	res, err := Audit(rubricPath, findingsPath, docPath, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", res.Failures)
	}
	if res.Headline != 0 {
		t.Fatalf("headline = %v, want 0", res.Headline)
	}
}

func TestAuditQuoteNotFound(t *testing.T) {
	rubricPath := filepath.Join("..", "..", ".claude", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
	docPath := writeTemp(t, "spec.md", "## 1.2 核心技术痛点\n\nSomething real.\n")
	r, _ := LoadRubric(rubricPath)
	f := baseFindings(r, rubricPath, docPath)
	setItem(&f, "bg.pain_evidence", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1.2 核心技术痛点", Quote: "this quote does not exist"}}
	})
	res, err := Audit(rubricPath, writeFindings(t, f), docPath, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "not found") {
		t.Fatalf("expected quote-not-found failure, got %v", res.Failures)
	}
}

func TestAuditQuoteWrongSection(t *testing.T) {
	rubricPath := filepath.Join("..", "..", ".claude", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
	docPath := writeTemp(t, "spec.md", "## 1.1 系统现状\n\nAlpha.\n\n## 1.2 核心技术痛点\n\nBeta.\n")
	r, _ := LoadRubric(rubricPath)
	f := baseFindings(r, rubricPath, docPath)
	setItem(&f, "bg.pain_evidence", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1.2 核心技术痛点", Quote: "Alpha."}}
	})
	res, err := Audit(rubricPath, writeFindings(t, f), docPath, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "not found") {
		t.Fatalf("expected wrong-section failure, got %v", res.Failures)
	}
}

func TestAuditRequiresWriterSession(t *testing.T) {
	rubricPath := filepath.Join("..", "..", ".claude", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
	docPath := writeTemp(t, "spec.md", "# Spec\n")
	r, _ := LoadRubric(rubricPath)
	f := baseFindings(r, rubricPath, docPath)
	res, err := Audit(rubricPath, writeFindings(t, f), docPath, "")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "-writer-session is required") {
		t.Fatalf("expected writer-session failure, got %v", res.Failures)
	}
}

func TestAuditRejectsSameSession(t *testing.T) {
	rubricPath := filepath.Join("..", "..", ".claude", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
	docPath := writeTemp(t, "spec.md", "# Spec\n")
	r, _ := LoadRubric(rubricPath)
	f := baseFindings(r, rubricPath, docPath)
	res, err := Audit(rubricPath, writeFindings(t, f), docPath, "reviewer-1")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "reviewer session equals the writer session") {
		t.Fatalf("expected same-session failure, got %v", res.Failures)
	}
}

func TestAuditRejectsNAOnCore(t *testing.T) {
	rubricPath := filepath.Join("..", "..", ".claude", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
	docPath := writeTemp(t, "spec.md", "# Spec\n")
	r, _ := LoadRubric(rubricPath)
	f := baseFindings(r, rubricPath, docPath)
	setItem(&f, "bg.pain_evidence", func(it *ItemFinding) {
		it.Status = "na"
		it.NaReason = "not applicable"
	})
	res, err := Audit(rubricPath, writeFindings(t, f), docPath, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "na is not allowed") {
		t.Fatalf("expected na-on-core failure, got %v", res.Failures)
	}
}

func TestAuditRejectsAllNASection(t *testing.T) {
	rubricPath := filepath.Join("..", "..", ".claude", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
	docPath := writeTemp(t, "spec.md", "## 3.1 总体架构\n\nTopology.\n")
	r, _ := LoadRubric(rubricPath)
	f := baseFindings(r, rubricPath, docPath)
	setItem(&f, "arch.topology", func(it *ItemFinding) {
		it.Status = "na"
		it.NaReason = "no diagram"
	})
	res, err := Audit(rubricPath, writeFindings(t, f), docPath, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "all items marked na") {
		t.Fatalf("expected all-na section failure, got %v", res.Failures)
	}
}

func TestScoreArithmetic(t *testing.T) {
	rubricPath := writeTemp(t, "rubric.yaml", miniRubric)
	docPath := writeTemp(t, "spec.md", miniDoc)
	r, _ := LoadRubric(rubricPath)
	f := baseFindings(r, rubricPath, docPath)
	setItem(&f, "test.item", func(it *ItemFinding) {
		it.Level = 1
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1. Test", Quote: "Evidence here."}}
	})
	res, err := Audit(rubricPath, writeFindings(t, f), docPath, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", res.Failures)
	}
	if res.Headline != 50 {
		t.Fatalf("headline = %v, want 50", res.Headline)
	}
}

func TestNFKCMatch(t *testing.T) {
	rubricPath := writeTemp(t, "rubric.yaml", miniRubric)
	docPath := writeTemp(t, "spec.md", "## 1. Test\n\nＡＢＣ，ＤＥＦ\n")
	r, _ := LoadRubric(rubricPath)
	f := baseFindings(r, rubricPath, docPath)
	setItem(&f, "test.item", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1. Test", Quote: "ABC,DEF"}}
	})
	res, err := Audit(rubricPath, writeFindings(t, f), docPath, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Failures) != 0 {
		t.Fatalf("expected NFKC match, got %v", res.Failures)
	}
}

func TestEscapeCell(t *testing.T) {
	got := escapeCell("a|b\nc`d")
	if got != "a\\|b c'd" {
		t.Fatalf("escapeCell = %q", got)
	}
}

func containsFailure(failures []string, want string) bool {
	for _, f := range failures {
		if strings.Contains(f, want) {
			return true
		}
	}
	return false
}
