package main

import (
	"strings"
	"testing"
)

func TestCalibrateAgreement(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", twoItemRubric)
	doc := writeTemp(t, "spec.md", "## 1. S\n\nEvidence.\n")
	r, _ := LoadRubric(rubric)

	ref := baseFindings(r, rubric, doc)
	setItem(&ref, "i1", func(it *ItemFinding) { it.Level = 2; it.Missing = "" })
	setItem(&ref, "i2", func(it *ItemFinding) { it.Level = 1; it.Missing = "" })
	cand := baseFindings(r, rubric, doc)
	setItem(&cand, "i1", func(it *ItemFinding) { it.Level = 2; it.Missing = "" })
	setItem(&cand, "i2", func(it *ItemFinding) { it.Level = 0; it.Missing = "missing" })

	res, err := Calibrate(rubric, writeFindings(t, ref), writeFindings(t, cand))
	if err != nil {
		t.Fatal(err)
	}
	if res.Compared != 2 || res.Exact != 1 || res.Agreement != 50 {
		t.Fatalf("compared=%d exact=%d agreement=%.1f", res.Compared, res.Exact, res.Agreement)
	}
	if len(res.Disagreements) != 1 || res.Disagreements[0].ID != "i2" {
		t.Fatalf("disagreements=%v", res.Disagreements)
	}
	if res.CoreCompared != 1 || res.CoreExact != 1 || res.CoreAgreement != 100 || res.CoreWeightedKappa != 1 {
		t.Fatalf("core: compared=%d exact=%d agreement=%.1f kappa=%v", res.CoreCompared, res.CoreExact, res.CoreAgreement, res.CoreWeightedKappa)
	}
	if !strings.Contains(RenderCalibration(res), "Inter-rater agreement") {
		t.Fatal("calibration report missing boundary note")
	}
}

func TestCalibratePerfectAgreement(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", twoItemRubric)
	doc := writeTemp(t, "spec.md", "## 1. S\n\nEvidence.\n")
	r, _ := LoadRubric(rubric)
	f := baseFindings(r, rubric, doc)
	setItem(&f, "i1", func(it *ItemFinding) { it.Level = 2; it.Missing = "" })
	setItem(&f, "i2", func(it *ItemFinding) { it.Level = 1; it.Missing = "" })
	path := writeFindings(t, f)
	res, err := Calibrate(rubric, path, path)
	if err != nil {
		t.Fatal(err)
	}
	if res.Compared != 2 || res.Exact != 2 || res.Agreement != 100 {
		t.Fatalf("compared=%d exact=%d agreement=%.1f", res.Compared, res.Exact, res.Agreement)
	}
	if res.WeightedKappa != 1 {
		t.Fatalf("weighted kappa = %v, want 1", res.WeightedKappa)
	}
}

func TestCalibrateStatusMismatch(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", twoItemRubric)
	doc := writeTemp(t, "spec.md", "## 1. S\n\nEvidence.\n")
	r, _ := LoadRubric(rubric)
	ref := baseFindings(r, rubric, doc)
	setItem(&ref, "i2", func(it *ItemFinding) { it.Status = "na"; it.NaReason = "n/a" })
	cand := baseFindings(r, rubric, doc)
	setItem(&cand, "i2", func(it *ItemFinding) { it.Level = 1; it.Missing = "" })
	res, err := Calibrate(rubric, writeFindings(t, ref), writeFindings(t, cand))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.StatusMismatches) != 1 || res.StatusMismatches[0] != "i2" {
		t.Fatalf("status mismatches = %v", res.StatusMismatches)
	}
	if res.Compared != 1 {
		t.Fatalf("compared = %d, want 1", res.Compared)
	}
}

func TestCalibrateRejectsDifferentTargets(t *testing.T) {
	rubric := writeTemp(t, "rubric.yaml", twoItemRubric)
	doc := writeTemp(t, "spec.md", "## 1. S\n\nEvidence.\n")
	r, _ := LoadRubric(rubric)
	ref := baseFindings(r, rubric, doc)
	cand := baseFindings(r, rubric, doc)
	cand.Target = "other.md"
	if _, err := Calibrate(rubric, writeFindings(t, ref), writeFindings(t, cand)); err == nil || !strings.Contains(err.Error(), "targets differ") {
		t.Fatalf("err=%v", err)
	}
}
