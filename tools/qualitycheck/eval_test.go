package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These are instrument evals: deterministic checks of the rubric, the audit
// pipeline and the reviewer contract. They do not call a model. The LLM
// reviewer calibration (gold set, kappa, thresholds) is a follow-up.

const evalRubric = `version: "eval-1"
dimensions:
  - id: core
    name: Core
    weight: 1
    items:
      - id: eval.background
        question: background
        tier: core
        weight: 1
        section: "1. Background"
        aliases: ["Background"]
        na_allowed: false
        levels: {"0": "none", "1": "partial", "2": "full"}
      - id: eval.pain
        question: pain
        tier: core
        weight: 1
        section: "2. Pain"
        aliases: ["Pain"]
        na_allowed: false
        levels: {"0": "none", "1": "partial", "2": "full"}
      - id: eval.goal
        question: goal
        tier: core
        weight: 1
        section: "3. Goals"
        aliases: ["Goals"]
        na_allowed: false
        levels: {"0": "none", "1": "partial", "2": "full"}
      - id: eval.rollback
        question: rollback
        tier: core
        weight: 1
        section: "4. Rollback"
        aliases: ["Rollback"]
        na_allowed: false
        levels: {"0": "none", "1": "partial", "2": "full"}
      - id: eval.hygiene
        question: hygiene
        tier: core
        weight: 1
        section: "5. Hygiene"
        aliases: ["Hygiene"]
        na_allowed: false
        levels: {"0": "none", "1": "partial", "2": "full"}
  - id: extended
    name: Extended
    weight: 1
    items:
      - id: eval.arch
        question: arch
        tier: extended
        weight: 1
        section: "6. Architecture"
        aliases: ["Architecture"]
        na_allowed: false
        levels: {"0": "none", "1": "partial", "2": "full"}
      - id: eval.detail
        question: detail
        tier: extended
        weight: 1
        section: "7. Detailed design"
        aliases: ["Detailed design"]
        na_allowed: false
        levels: {"0": "none", "1": "partial", "2": "full"}
`

const strongDoc = `## 1. Background

The service parses PDFs and stores nodes.

## 2. Pain

P99 latency is 4.2s; incident on 2026-10-12 lost 37 files.

## 3. Goals

G1 reduces P99 to 500ms; acceptance: benchmark.

## 4. Rollback

Rollback when error rate > 2% for 3 minutes; steps: flip flag, restore image.

## 5. Hygiene

No placeholders; numbers sourced.

## 6. Architecture

Topology: API -> pipeline -> storage.

## 7. Detailed design

Data model: ParseTask(id string, status enum); interface: POST /tasks.
`

const weakDoc = `## 1. Background

Existing system.

## 2. Pain

Performance is bad.

## 3. Goals

Improve performance.

## 4. Rollback

Can rollback.

## 5. Hygiene

待补充

## 6. Architecture

...

## 7. Detailed design

...
`

var strongEvidence = map[string]string{
	"eval.background": "The service parses PDFs and stores nodes.",
	"eval.pain":       "P99 latency is 4.2s; incident on 2026-10-12 lost 37 files.",
	"eval.goal":       "G1 reduces P99 to 500ms; acceptance: benchmark.",
	"eval.rollback":   "Rollback when error rate > 2% for 3 minutes; steps: flip flag, restore image.",
	"eval.hygiene":    "No placeholders; numbers sourced.",
	"eval.arch":       "Topology: API -> pipeline -> storage.",
	"eval.detail":     "Data model: ParseTask(id string, status enum); interface: POST /tasks.",
}

func evalAudit(t *testing.T, doc string, level func(string) int) AuditResult {
	t.Helper()
	rubricPath := writeTemp(t, "rubric.yaml", evalRubric)
	docPath := writeTemp(t, "spec.md", doc)
	r, err := LoadRubric(rubricPath)
	if err != nil {
		t.Fatal(err)
	}
	f := baseFindings(r, rubricPath, docPath)
	for i := range f.Items {
		l := level(f.Items[i].ID)
		f.Items[i].Level = l
		if l > 0 {
			f.Items[i].Missing = ""
			f.Items[i].Evidence = []Evidence{{
				Section: sectionFor(f.Items[i].ID),
				Quote:   strongEvidence[f.Items[i].ID],
			}}
		}
	}
	res, err := Audit(rubricPath, writeFindings(t, f), docPath, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func sectionFor(id string) string {
	switch id {
	case "eval.background":
		return "1. Background"
	case "eval.pain":
		return "2. Pain"
	case "eval.goal":
		return "3. Goals"
	case "eval.rollback":
		return "4. Rollback"
	case "eval.hygiene":
		return "5. Hygiene"
	case "eval.arch":
		return "6. Architecture"
	case "eval.detail":
		return "7. Detailed design"
	}
	return ""
}

func TestEvalStrongScoresHigh(t *testing.T) {
	res := evalAudit(t, strongDoc, func(string) int { return 2 })
	if len(res.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", res.Failures)
	}
	if res.Headline < 90 {
		t.Fatalf("strong fixture headline = %.1f, want >= 90", res.Headline)
	}
}

// Value: protects=level aggregation; fails_when=scoring arithmetic changes; why_new=instrument eval; seam=none; no_claim=levels are fixture-assigned, so green does not prove reviewer judgment or document-quality discrimination.
func TestEvalWeakScoresLow(t *testing.T) {
	res := evalAudit(t, weakDoc, func(string) int { return 0 })
	if len(res.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", res.Failures)
	}
	if res.Headline > 10 {
		t.Fatalf("weak fixture headline = %.1f, want <= 10", res.Headline)
	}
}

func TestEvalSeparatesStrongFromWeak(t *testing.T) {
	strong := evalAudit(t, strongDoc, func(string) int { return 2 })
	weak := evalAudit(t, weakDoc, func(string) int { return 0 })
	if strong.Headline-weak.Headline < 80 {
		t.Fatalf("separation = %.1f, want >= 80", strong.Headline-weak.Headline)
	}
}

func TestReviewerContract(t *testing.T) {
	path := filepath.Join("..", "..", "skills", "doc-writer", "spec", "quality", "reviewer.md")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{"不可信输入", "逐字", "不输出总分", "na_allowed", "level"} {
		if !strings.Contains(s, want) {
			t.Fatalf("reviewer.md missing %q", want)
		}
	}
}

func TestFindingsSchemaContract(t *testing.T) {
	path := filepath.Join("..", "..", "skills", "doc-writer", "spec", "quality", "findings.schema.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(b, &schema); err != nil {
		t.Fatal(err)
	}
	req, _ := schema["required"].([]any)
	want := map[string]bool{"schema_version": false, "rubric_version": false, "rubric_sha256": false, "target": false, "target_sha256": false, "reviewer": false, "items": false}
	for _, r := range req {
		if s, ok := r.(string); ok {
			want[s] = true
		}
	}
	for k, seen := range want {
		if !seen {
			t.Fatalf("schema missing required field %q", k)
		}
	}
}

func TestRealRubricHasChapterCoreItems(t *testing.T) {
	path := filepath.Join("..", "..", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
	r, err := LoadRubric(path)
	if err != nil {
		t.Fatal(err)
	}
	found := map[string]bool{"arch.solution_shape": false, "design.contract": false}
	for _, d := range r.Dimensions {
		for _, it := range d.Items {
			if _, ok := found[it.ID]; ok {
				found[it.ID] = it.Tier == "core"
			}
		}
	}
	for id, ok := range found {
		if !ok {
			t.Fatalf("%s must be a core item", id)
		}
	}
}
