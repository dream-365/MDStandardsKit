package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func realRubricPath() string {
	return filepath.Join("..", "..", "skills", "doc-writer", "spec", "quality", "rubric.v1.yaml")
}

func runQC(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	code := run(args, &out, &errb)
	return code, out.String(), errb.String()
}

func TestRunModes(t *testing.T) {
	rubric := realRubricPath()
	doc := writeTemp(t, "spec.md", "# Spec\n")

	t.Run("no mode", func(t *testing.T) {
		code, _, errs := runQC(t)
		if code != 2 || !strings.Contains(errs, "one of") {
			t.Fatalf("code=%d errs=%q", code, errs)
		}
	})
	t.Run("validate without rubric", func(t *testing.T) {
		code, _, errs := runQC(t, "-validate-rubric")
		if code != 2 || !strings.Contains(errs, "-rubric is required") {
			t.Fatalf("code=%d errs=%q", code, errs)
		}
	})
	t.Run("validate missing file", func(t *testing.T) {
		code, _, errs := runQC(t, "-validate-rubric", "-rubric", "/no/such/rubric.yaml")
		if code != 2 || errs == "" {
			t.Fatalf("code=%d errs=%q", code, errs)
		}
	})
	t.Run("validate invalid rubric", func(t *testing.T) {
		bad := writeTemp(t, "rubric.yaml", "version: \"\"\ndimensions: []\n")
		code, _, _ := runQC(t, "-validate-rubric", "-rubric", bad)
		if code != 1 {
			t.Fatalf("code=%d, want 1", code)
		}
	})
	t.Run("validate ok", func(t *testing.T) {
		code, out, _ := runQC(t, "-validate-rubric", "-rubric", rubric)
		r, err := LoadRubric(rubric)
		if err != nil {
			t.Fatal(err)
		}
		if want := "OK: rubric " + r.Version + "\n"; code != 0 || out != want {
			t.Fatalf("code=%d out=%q, want %q", code, out, want)
		}
	})
	t.Run("audit missing flags", func(t *testing.T) {
		code, _, errs := runQC(t, "-audit", "-rubric", rubric)
		if code != 2 || !strings.Contains(errs, "usage:") {
			t.Fatalf("code=%d errs=%q", code, errs)
		}
	})
	t.Run("audit clean", func(t *testing.T) {
		r, _ := LoadRubric(rubric)
		f := baseFindings(r, rubric, doc)
		findings := writeFindings(t, f)
		out := filepath.Join(t.TempDir(), "report.md")
		code, stdout, _ := runQC(t, "-audit", "-rubric", rubric, "-findings", findings, "-doc", doc, "-writer-session", "writer-1", "-out", out)
		if code != 0 || !strings.Contains(stdout, "# Spec Quality Report") {
			t.Fatalf("code=%d stdout=%q", code, stdout)
		}
	})
	t.Run("audit failures", func(t *testing.T) {
		r, _ := LoadRubric(rubric)
		f := baseFindings(r, rubric, doc)
		setItem(&f, "bg.pain_evidence", func(it *ItemFinding) {
			it.Level = 2
			it.Missing = ""
			it.Evidence = []Evidence{{Section: "1.2 核心技术痛点", Quote: "not in the doc"}}
		})
		findings := writeFindings(t, f)
		out := filepath.Join(t.TempDir(), "report.md")
		code, _, _ := runQC(t, "-audit", "-rubric", rubric, "-findings", findings, "-doc", doc, "-writer-session", "writer-1", "-out", out)
		if code != 1 {
			t.Fatalf("code=%d, want 1", code)
		}
	})
	t.Run("audit write error", func(t *testing.T) {
		r, _ := LoadRubric(rubric)
		f := baseFindings(r, rubric, doc)
		findings := writeFindings(t, f)
		code, _, _ := runQC(t, "-audit", "-rubric", rubric, "-findings", findings, "-doc", doc, "-writer-session", "writer-1", "-out", "/no/such/dir/report.md")
		if code != 2 {
			t.Fatalf("code=%d, want 2", code)
		}
	})
	t.Run("audit invalid findings", func(t *testing.T) {
		bad := writeTemp(t, "findings.json", "{not json")
		out := filepath.Join(t.TempDir(), "report.md")
		code, _, _ := runQC(t, "-audit", "-rubric", rubric, "-findings", bad, "-doc", doc, "-writer-session", "writer-1", "-out", out)
		if code != 2 {
			t.Fatalf("code=%d, want 2", code)
		}
	})
}

func TestRenderReport(t *testing.T) {
	res := AuditResult{
		RubricVersion: "1.0.0",
		RubricSHA256:  "abc",
		Target:        "spec.md",
		TargetSHA256:  "def",
		ReviewerID:    "reviewer-1",
		Headline:      50,
		Diagnostic:    60,
		Dimensions: []DimensionResult{
			{
				ID: "d", Name: "Dim", Weight: 1, Headline: 50, Diagnostic: 60,
				Items: []ItemResult{
					{ID: "i1", Question: "q1", Tier: "core", Weight: 1, Status: "scored", Level: 1, Evidence: []Evidence{{Section: "S", Quote: "a|b"}}, Missing: "m"},
					{ID: "i2", Question: "q2", Tier: "extended", Weight: 1, Status: "na", NaReason: "n/a"},
					{ID: "i3", Question: "q3", Tier: "core", Weight: 1, Status: "scored", Level: 0, Failures: []string{"bad"}},
				},
			},
			{ID: "empty", Name: "Empty", Weight: 1, Headline: -1, Diagnostic: -1},
		},
		Failures: []string{"global failure"},
	}
	out := RenderReport(res)
	for _, want := range []string{"# Spec Quality Report", "Headline score", "evidence", "missing:", "na:", "audit failure", "global failure", "—", "a\\|b"} {
		if !strings.Contains(out, want) {
			t.Fatalf("report missing %q:\n%s", want, out)
		}
	}
	res.Failures = nil
	if !strings.Contains(RenderReport(res), "None.") {
		t.Fatal("empty failures should render None.")
	}
}

func TestRubricValidationErrors(t *testing.T) {
	base := `version: "x"
dimensions:
  - id: d
    name: D
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
`
	cases := map[string]string{
		"empty version":      strings.Replace(base, `version: "x"`, `version: ""`, 1),
		"no dimensions":      "version: \"x\"\ndimensions: []\n",
		"empty dimension id": strings.Replace(base, "id: d", "id: \"\"", 1),
		"duplicate dim":      strings.Replace(base, "  - id: d", "  - id: d\n    name: D2\n    weight: 1\n    items: []\n  - id: d", 1),
		"dim weight":         strings.Replace(base, "    weight: 1\n    items:", "    weight: 0\n    items:", 1),
		"item empty id":      strings.Replace(base, "id: i", "id: \"\"", 1),
		"bad tier":           strings.Replace(base, "tier: core", "tier: other", 1),
		"item weight":        strings.Replace(base, "        weight: 1", "        weight: 0", 1),
		"missing section":    strings.Replace(base, `section: "1. S"`, `section: ""`, 1),
		"missing anchor":     strings.Replace(base, `"2": "c"`, `"2": ""`, 1),
	}
	for name, yaml := range cases {
		t.Run(name, func(t *testing.T) {
			p := writeTemp(t, "rubric.yaml", yaml)
			r, err := LoadRubric(p)
			if err != nil {
				t.Fatalf("load: %v", err)
			}
			if err := ValidateRubric(r); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}

	t.Run("dim no scorable", func(t *testing.T) {
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
        na_allowed: true
        levels: {"0": "a", "1": "b", "2": "c"}
`)
		r, err := LoadRubric(p)
		if err != nil {
			t.Fatal(err)
		}
		err = ValidateRubric(r)
		if err == nil || !strings.Contains(err.Error(), "has no na_allowed: false item") {
			t.Fatalf("err=%v", err)
		}
	})
}

func TestAuditRejections(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", miniRubric)
	doc := writeTemp(t, "spec.md", miniDoc)
	r, _ := LoadRubric(rubric)
	mutate := func(id string, f func(*ItemFinding)) Findings {
		x := baseFindings(r, rubric, doc)
		if id != "" {
			setItem(&x, id, f)
		}
		return x
	}
	cases := map[string]struct {
		f       Findings
		wantSub string
	}{
		"schema version": {func() Findings { x := mutate("", nil); x.SchemaVersion = 2; return x }(), "schema_version"},
		"rubric version": {func() Findings { x := mutate("", nil); x.RubricVersion = "nope"; return x }(), "rubric_version"},
		"rubric sha":     {func() Findings { x := mutate("", nil); x.RubricSHA256 = strings.Repeat("0", 64); return x }(), "rubric_sha256"},
		"target sha":     {func() Findings { x := mutate("", nil); x.TargetSHA256 = strings.Repeat("0", 64); return x }(), "target_sha256"},
		"empty target":   {func() Findings { x := mutate("", nil); x.Target = ""; return x }(), "target is empty"},
		"role":           {func() Findings { x := mutate("", nil); x.Reviewer.Role = "author"; return x }(), "reviewer.role"},
		"empty session":  {func() Findings { x := mutate("", nil); x.Reviewer.SessionID = ""; return x }(), "reviewer.session_id"},
		"unknown id": {func() Findings {
			x := mutate("", nil)
			x.Items = append(x.Items, ItemFinding{ID: "nope", Status: "scored", Missing: "m"})
			return x
		}(), "unknown item id"},
		"level range":    {mutate("test.item", func(it *ItemFinding) { it.Level = 3 }), "out of range"},
		"level0 missing": {mutate("test.item", func(it *ItemFinding) { it.Missing = "" }), "missing reason"},
		"no evidence":    {mutate("test.item", func(it *ItemFinding) { it.Level = 2; it.Missing = "" }), "requires evidence"},
		"empty quote": {mutate("test.item", func(it *ItemFinding) {
			it.Level = 2
			it.Missing = ""
			it.Evidence = []Evidence{{Section: "1. Test", Quote: ""}}
		}), "quote is empty"},
		"bad section": {mutate("test.item", func(it *ItemFinding) {
			it.Level = 2
			it.Missing = ""
			it.Evidence = []Evidence{{Section: "nope", Quote: "Evidence here."}}
		}), "not declared"},
		"na no reason": {func() Findings {
			x := mutate("", nil)
			setItem(&x, "test.item", func(it *ItemFinding) { it.Status = "na" })
			return x
		}(), "na requires a reason"},
		"unknown status": {mutate("test.item", func(it *ItemFinding) { it.Status = "maybe" }), "unknown status"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			res, err := Audit(rubric, writeFindings(t, tc.f), doc, "writer-1")
			if err != nil {
				t.Fatalf("audit: %v", err)
			}
			if !containsFailure(res.Failures, tc.wantSub) {
				t.Fatalf("want %q in %v", tc.wantSub, res.Failures)
			}
		})
	}

	t.Run("duplicate and missing ids", func(t *testing.T) {
		x := mutate("", nil)
		dup := x
		dup.Items = append(dup.Items, x.Items[0])
		res, _ := Audit(rubric, writeFindings(t, dup), doc, "writer-1")
		if !containsFailure(res.Failures, "duplicate finding") {
			t.Fatalf("duplicate: %v", res.Failures)
		}
		miss := mutate("", nil)
		miss.Items = nil
		res, _ = Audit(rubric, writeFindings(t, miss), doc, "writer-1")
		if !containsFailure(res.Failures, "missing finding") {
			t.Fatalf("missing: %v", res.Failures)
		}
	})

	t.Run("missing declared section", func(t *testing.T) {
		x := mutate("test.item", func(it *ItemFinding) {
			it.Level = 2
			it.Missing = ""
			it.Evidence = []Evidence{{Section: "1. Test", Quote: "Evidence here."}}
		})
		empty := writeTemp(t, "spec.md", "# empty\n")
		// The target hash in x is for the original doc; recompute for the new doc.
		docBytes := []byte("# empty\n")
		x.TargetSHA256 = sha256Hex(docBytes)
		res, _ := Audit(rubric, writeFindings(t, x), empty, "writer-1")
		if !containsFailure(res.Failures, "section") {
			t.Fatalf("missing section: %v", res.Failures)
		}
	})

	t.Run("full text scope", func(t *testing.T) {
		realDoc := writeTemp(t, "spec.md", "# Spec\n\nNo placeholders here.\n")
		rr, _ := LoadRubric(realRubricPath())
		x := baseFindings(rr, realRubricPath(), realDoc)
		setItem(&x, "x.no_placeholders", func(it *ItemFinding) {
			it.Level = 2
			it.Missing = ""
			it.Evidence = []Evidence{{Section: "全文", Quote: "No placeholders here."}}
		})
		res, err := Audit(realRubricPath(), writeFindings(t, x), realDoc, "writer-1")
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Failures) != 0 {
			t.Fatalf("full-text evidence should pass: %v", res.Failures)
		}
	})
}

func TestParseAndNormalize(t *testing.T) {
	doc := "# Top\n\n```\n## not a heading\n```\n\n## 1. A\n\nbody\n\n### 1.1 B\n\nchild\n\n## 2. C\n\nlast\n"
	secs := ParseSections(doc)
	if len(secs) != 4 {
		t.Fatalf("sections = %d, want 4: %+v", len(secs), secs)
	}
	if secs[1].Text != "1. A" || !strings.Contains(secs[1].Body, "child") {
		t.Fatalf("unexpected section: %+v", secs[1])
	}
	if secs[3].Text != "2. C" {
		t.Fatalf("unexpected section: %+v", secs[3])
	}
	crlf := ParseSections("# A\r\n\r\nbody\r\n")
	if len(crlf) != 1 || crlf[0].Text != "A" {
		t.Fatalf("crlf parse: %+v", crlf)
	}
	trailing := ParseSections("## A ##\n\nbody\n")
	if len(trailing) != 1 || trailing[0].Text != "A" {
		t.Fatalf("trailing # parse: %+v", trailing)
	}
	if got := normalizeForMatch("A\u00a0B\u3000C\n\nD"); got != "A B C D" {
		t.Fatalf("normalize = %q", got)
	}
	if got := normalizeForMatch("ＡＢＣ，Ｄ"); got != "ABC,D" {
		t.Fatalf("full-width normalize = %q", got)
	}
}

func TestScoringPaths(t *testing.T) {
	if got := scoreItems(nil, false); got != -1 {
		t.Fatalf("empty score = %v, want -1", got)
	}
	items := []ItemResult{
		{Tier: "core", Weight: 1, Status: "scored", Level: 2},
		{Tier: "extended", Weight: 1, Status: "scored", Level: 0},
		{Tier: "extended", Weight: 1, Status: "na"},
	}
	if got := scoreItems(items, true); got != 100 {
		t.Fatalf("coreOnly = %v, want 100", got)
	}
	if got := scoreItems(items, false); got != 50 {
		t.Fatalf("all items = %v, want 50", got)
	}
	if got := weightedAverage([]DimensionResult{{Weight: 1, Headline: -1, Diagnostic: -1}, {Weight: 1, Headline: 100, Diagnostic: 100}}, true); got != 100 {
		t.Fatalf("weighted average = %v", got)
	}
	if got := weightedAverage([]DimensionResult{{Weight: 1, Headline: -1, Diagnostic: -1}}, false); got != -1 {
		t.Fatalf("expected no scorable dimensions, got %v", got)
	}
}

func TestAuditDiagnosticAndInputErrors(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", miniRubric)
	doc := writeTemp(t, "spec.md", miniDoc)
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	setItem(&f, "test.item", func(it *ItemFinding) {
		it.Level = 2
		it.Missing = ""
		it.Evidence = []Evidence{{Section: "1. Test", Quote: "Evidence here."}}
	})
	res, err := Audit(rubric, writeFindings(t, f), doc, "writer-1")
	if err != nil {
		t.Fatal(err)
	}
	if res.Headline != 100 || res.Diagnostic != 100 {
		t.Fatalf("headline=%v diagnostic=%v", res.Headline, res.Diagnostic)
	}

	if _, err := Audit("/no/such/rubric.yaml", "/no/such/findings.json", doc, "w"); err == nil {
		t.Fatal("expected rubric load error")
	}
	badFindings := writeTemp(t, "findings.json", "{bad")
	if _, err := Audit(rubric, badFindings, doc, "w"); err == nil {
		t.Fatal("expected findings error")
	}
	if _, err := Audit(rubric, writeFindings(t, f), "/no/such/spec.md", "w"); err == nil {
		t.Fatal("expected doc read error")
	}
}

func TestJSONShape(t *testing.T) {
	// The findings struct must round-trip the schema's required keys.
	f := Findings{
		SchemaVersion: 1, RubricVersion: "1", RubricSHA256: strings.Repeat("a", 64),
		Target: "t", TargetSHA256: strings.Repeat("b", 64),
		Reviewer: Reviewer{Role: "independent", Model: "m", SessionID: "s", StartedAt: "d"},
		Items:    []ItemFinding{{ID: "i", Status: "scored", Level: 0, Missing: "m"}},
	}
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"schema_version", "rubric_version", "rubric_sha256", "target", "target_sha256", "reviewer", "items"} {
		if !strings.Contains(string(b), `"`+key+`"`) {
			t.Fatalf("missing key %q in %s", key, b)
		}
	}
}
