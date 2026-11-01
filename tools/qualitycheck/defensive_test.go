package main

import (
	"strings"
	"testing"
)

const twoItemRubric = `version: "x"
dimensions:
  - id: d
    name: D
    weight: 1
    items:
      - id: i1
        question: q1
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
      - id: i2
        question: q2
        tier: extended
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: true
        levels: {"0": "a", "1": "b", "2": "c"}
`

const fullTextRubric = `version: "x"
dimensions:
  - id: d
    name: D
    weight: 1
    items:
      - id: i1
        question: q1
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
      - id: i2
        question: q2
        tier: extended
        weight: 1
        section: "全文"
        aliases: ["全文"]
        na_allowed: true
        levels: {"0": "a", "1": "b", "2": "c"}
`

func TestDefensiveBranches(t *testing.T) {
	t.Run("invalid flag", func(t *testing.T) {
		code, _, errs := runQC(t, "-bogus")
		if code != 2 || !strings.Contains(errs, "flag provided but not defined") {
			t.Fatalf("code=%d errs=%q", code, errs)
		}
	})

	t.Run("malformed rubric yaml", func(t *testing.T) {
		p := writeTemp(t, "rubric.yaml", "version: [unclosed\n")
		if _, err := LoadRubric(p); err == nil || !strings.Contains(err.Error(), "invalid rubric YAML") {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("duplicate dimension id", func(t *testing.T) {
		p := writeTemp(t, "rubric.yaml", `version: "x"
dimensions:
  - id: d
    name: D
    weight: 1
    items:
      - id: i1
        question: q
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
  - id: d
    name: D2
    weight: 1
    items:
      - id: i2
        question: q
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
`)
		r, err := LoadRubric(p)
		if err != nil {
			t.Fatal(err)
		}
		if err := ValidateRubric(r); err == nil || !strings.Contains(err.Error(), "duplicate dimension") {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("duplicate item id", func(t *testing.T) {
		p := writeTemp(t, "rubric.yaml", `version: "x"
dimensions:
  - id: d1
    name: D1
    weight: 1
    items:
      - id: i
        question: q
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
  - id: d2
    name: D2
    weight: 1
    items:
      - id: i
        question: q
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
`)
		r, err := LoadRubric(p)
		if err != nil {
			t.Fatal(err)
		}
		if err := ValidateRubric(r); err == nil || !strings.Contains(err.Error(), "duplicate item") {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("audit invalid rubric", func(t *testing.T) {
		p := writeTemp(t, "rubric.yaml", "version: \"\"\ndimensions: []\n")
		doc := writeTemp(t, "spec.md", "# Spec\n")
		f := writeTemp(t, "findings.json", "{}")
		if _, err := Audit(p, f, doc, "w"); err == nil || !strings.Contains(err.Error(), "rubric invalid") {
			t.Fatalf("err=%v", err)
		}
	})

	t.Run("unreadable findings", func(t *testing.T) {
		doc := writeTemp(t, "spec.md", "# Spec\n")
		if _, err := Audit(realRubricPath(), "/no/such/findings.json", doc, "w"); err == nil {
			t.Fatal("expected findings read error")
		}
	})

	t.Run("mixed all-na section does not fail", func(t *testing.T) {
		rubric := writeTemp(t, "rubric.yaml", twoItemRubric)
		doc := writeTemp(t, "spec.md", "## 1. S\n\nEvidence.\n")
		r, _ := LoadRubric(rubric)
		f := baseFindings(r, rubric, doc)
		setItem(&f, "i2", func(it *ItemFinding) {
			it.Status = "na"
			it.NaReason = "not applicable"
		})
		res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
		if err != nil {
			t.Fatal(err)
		}
		if containsFailure(res.Failures, "all items marked na") {
			t.Fatalf("mixed section should not be all-na: %v", res.Failures)
		}
	})

	t.Run("full-text na skips section sweep", func(t *testing.T) {
		rubric := writeTemp(t, "rubric.yaml", fullTextRubric)
		doc := writeTemp(t, "spec.md", "## 1. S\n\nEvidence.\n")
		r, _ := LoadRubric(rubric)
		f := baseFindings(r, rubric, doc)
		setItem(&f, "i2", func(it *ItemFinding) {
			it.Status = "na"
			it.NaReason = "not applicable"
		})
		res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
		if err != nil {
			t.Fatal(err)
		}
		if containsFailure(res.Failures, "all items marked na") {
			t.Fatalf("full-text na should be skipped: %v", res.Failures)
		}
	})
}

func TestTwoItemScoring(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", twoItemRubric)
	doc := writeTemp(t, "spec.md", "## 1. S\n\nEvidence.\n")
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	setItem(&f, "i1", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1. S", Quote: "Evidence."}}
	})
	res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if res.Headline != 100 || res.Diagnostic != 50 {
		t.Fatalf("headline=%v diagnostic=%v, want 100/50", res.Headline, res.Diagnostic)
	}
}

