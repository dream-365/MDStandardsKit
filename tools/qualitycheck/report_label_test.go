package main

import (
	"strings"
	"testing"
)

func TestRenderReportCustomLabels(t *testing.T) {
	res := AuditResult{
		ReportTitle: "Idea Quality Report",
		TargetLabel: "Idea",
		Target:      "docs/ideas/a.md",
		Dimensions: []DimensionResult{{
			ID:   "claim",
			Name: "一句话",
			Items: []ItemResult{{
				ID:    "claim.standalone",
				Guess: "它想做一个提醒工具",
			}},
		}},
	}
	out := RenderReport(res)
	for _, want := range []string{
		"# Idea Quality Report",
		"- Idea: `docs/ideas/a.md`",
		"blind guess: 它想做一个提醒工具",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("report missing %q:\n%s", want, out)
		}
	}
}

func TestRenderReportDefaultsToSpec(t *testing.T) {
	out := RenderReport(AuditResult{Target: "spec.md"})
	if !strings.Contains(out, "# Spec Quality Report") {
		t.Fatalf("default title missing:\n%s", out)
	}
	if !strings.Contains(out, "- Target: `spec.md`") {
		t.Fatalf("default target label missing:\n%s", out)
	}
}
