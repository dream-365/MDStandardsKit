package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHeadingIsNotEvidence(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", miniRubric)
	doc := writeTemp(t, "spec.md", "## 1. Test\n\nEvidence here.\n")
	r, _ := LoadRubric(rubric)

	f := baseFindings(r, rubric, doc)
	setItem(&f, "test.item", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1. Test", Quote: "1. Test"}}
	})
	res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "not found") {
		t.Fatalf("heading-only evidence should fail: %v", res.Failures)
	}

	f2 := baseFindings(r, rubric, doc)
	setItem(&f2, "test.item", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1. Test", Quote: "Evidence here."}}
	})
	res, err = Audit(rubric, writeFindings(t, f2), doc, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Failures) != 0 {
		t.Fatalf("body evidence should pass: %v", res.Failures)
	}
}

func TestFencedHeadingLineIsEvidence(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", miniRubric)
	doc := writeTemp(t, "spec.md", "## 1. Test\n\n```\n# install deps\n```\n")
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	setItem(&f, "test.item", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1. Test", Quote: "# install deps"}}
	})
	res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Failures) != 0 {
		t.Fatalf("fenced heading-like line should stay in the body: %v", res.Failures)
	}
}

func TestIllegalNADoesNotRaiseHeadline(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", `version: "x"
dimensions:
  - id: d
    name: D
    weight: 1
    items:
      - id: a
        question: q
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
      - id: b
        question: q
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
`)
	doc := writeTemp(t, "spec.md", "## 1. S\n\nEvidence here.\n")
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	setItem(&f, "a", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1. S", Quote: "Evidence here."}}
	})
	setItem(&f, "b", func(it *ItemFinding) {
		it.Status = "na"
		it.NaReason = "illegal core na"
	})
	res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "na is not allowed") {
		t.Fatalf("expected illegal-na failure: %v", res.Failures)
	}
	if res.Headline >= 100 {
		t.Fatalf("illegal na raised headline to %.1f", res.Headline)
	}
}

func TestNoCoreRubricRejected(t *testing.T) {
	p := writeTemp(t, "rubric.yaml", `version: "x"
dimensions:
  - id: d
    name: D
    weight: 1
    items:
      - id: i
        question: q
        tier: extended
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
	if err := ValidateRubric(r); err == nil || !strings.Contains(err.Error(), "no core items") {
		t.Fatalf("err=%v", err)
	}
}

func TestEscapeCellBackslashPipe(t *testing.T) {
	if got := escapeCell(`a\|b`); got != "a\\\\\\|b" {
		t.Fatalf("escapeCell = %q", got)
	}
}

func TestOutSymlinkCollision(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", miniRubric)
	doc := writeTemp(t, "spec.md", miniDoc)
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	findings := writeFindings(t, f)
	link := filepath.Join(t.TempDir(), "link.md")
	if err := os.Symlink(doc, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	before, _ := os.ReadFile(doc)
	code, _, errs := runQC(t, "-audit", "-rubric", rubric, "-findings", findings, "-doc", doc, "-writer-session", "writer-1", "-out", link)
	if code != 2 || !strings.Contains(errs, "must not overwrite") {
		t.Fatalf("code=%d errs=%q", code, errs)
	}
	after, _ := os.ReadFile(doc)
	if string(before) != string(after) {
		t.Fatal("symlink let the output overwrite the input")
	}
}

func TestFailedOrMissingItemsDoNotRaiseHeadline(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", `version: "x"
dimensions:
  - id: d
    name: D
    weight: 1
    items:
      - id: a
        question: q
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
      - id: b
        question: q
        tier: core
        weight: 1
        section: "1. S"
        aliases: ["S"]
        na_allowed: false
        levels: {"0": "a", "1": "b", "2": "c"}
`)
	doc := writeTemp(t, "spec.md", "## 1. S\n\nEvidence here.\n")
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	setItem(&f, "a", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1. S", Quote: "fabricated quote"}}
	})
	// Remove item b entirely: a missing finding must count as 0.
	var kept []ItemFinding
	for _, it := range f.Items {
		if it.ID != "b" {
			kept = append(kept, it)
		}
	}
	f.Items = kept
	res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Failures) == 0 {
		t.Fatal("expected audit failures")
	}
	if res.Headline >= 100 {
		t.Fatalf("failed/missing items raised headline to %.1f", res.Headline)
	}
	if !strings.Contains(RenderReport(res), "not computed") {
		t.Fatal("failed audit should not render a headline score")
	}
}

func TestOutMustNotOverwriteInput(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", miniRubric)
	doc := writeTemp(t, "spec.md", miniDoc)
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	findings := writeFindings(t, f)
	before, _ := os.ReadFile(doc)
	code, _, errs := runQC(t, "-audit", "-rubric", rubric, "-findings", findings, "-doc", doc, "-writer-session", "writer-1", "-out", doc)
	if code != 2 || !strings.Contains(errs, "must not overwrite") {
		t.Fatalf("code=%d errs=%q", code, errs)
	}
	after, _ := os.ReadFile(doc)
	if string(before) != string(after) {
		t.Fatal("doc was overwritten")
	}
}

func TestTargetMustMatchDoc(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", miniRubric)
	doc := writeTemp(t, "spec.md", miniDoc)
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	f.Target = "other.md"
	res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "does not match -doc") {
		t.Fatalf("target mismatch: %v", res.Failures)
	}
}

func TestDuplicateHeadingRejected(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", miniRubric)
	doc := writeTemp(t, "spec.md", "## 1. Test\n\nA.\n\n## 1. Test\n\nB.\n")
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if !containsFailure(res.Failures, "duplicate heading") {
		t.Fatalf("duplicate heading: %v", res.Failures)
	}
}

func TestTrailingJSONRejected(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", miniRubric)
	doc := writeTemp(t, "spec.md", miniDoc)
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	b, _ := json.Marshal(f)
	bad := writeTemp(t, "findings.json", string(b)+"GARBAGE")
	if _, err := Audit(rubric, bad, doc, "writer-1"); err == nil || !strings.Contains(err.Error(), "trailing data") {
		t.Fatalf("err=%v", err)
	}
}

func TestUnknownRubricKeyRejected(t *testing.T) {
	p := writeTemp(t, "rubric.yaml", "version: \"x\"\nbogus: true\ndimensions: []\n")
	if _, err := LoadRubric(p); err == nil {
		t.Fatal("unknown rubric key should fail decoding")
	}
}

func TestSmokeReportsReproduce(t *testing.T) {
	root := filepath.Join("..", "..", "docs", "quality")
	matches, err := filepath.Glob(filepath.Join(root, "*", "quality-findings.smoke.json"))
	if err != nil || len(matches) == 0 {
		t.Fatalf("no smoke findings: %v", err)
	}
	for _, findingsPath := range matches {
		b, err := os.ReadFile(findingsPath)
		if err != nil {
			t.Fatal(err)
		}
		var f Findings
		if err := json.Unmarshal(b, &f); err != nil {
			t.Fatal(err)
		}
		docPath := filepath.Join("..", "..", f.Target)
		res, err := Audit(realRubricPath(), findingsPath, docPath, "writer-smoke")
		if err != nil {
			t.Fatalf("%s: %v", findingsPath, err)
		}
		if len(res.Failures) != 0 {
			t.Fatalf("%s: %v", findingsPath, res.Failures)
		}
		reportPath := strings.Replace(findingsPath, "quality-findings.smoke.json", "quality-report.smoke.md", 1)
		want, err := os.ReadFile(reportPath)
		if err != nil {
			t.Fatal(err)
		}
		if RenderReport(res) != string(want) {
			t.Fatalf("%s is stale; regenerate the smoke report", reportPath)
		}
	}
}
