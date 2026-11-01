package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

const findingsSchemaVersion = "1"

type Finding struct {
	Rule     string `json:"rule"`
	Document string `json:"document"`
	Target   string `json:"target,omitempty"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type Report struct {
	SchemaVersion string         `json:"schema_version"`
	Findings      []Finding      `json:"findings"`
	ExtractErrors []ExtractError `json:"extract_errors"`
	Stats         Stats          `json:"stats"`
}

type BaselineEntry struct {
	Rule     string `json:"rule"`
	Document string `json:"document"`
	Target   string `json:"target,omitempty"`
}

func sortFindings(fs []Finding) {
	sort.SliceStable(fs, func(i, j int) bool {
		a, b := fs[i], fs[j]
		if a.Rule != b.Rule {
			return a.Rule < b.Rule
		}
		if a.Document != b.Document {
			return a.Document < b.Document
		}
		return a.Target < b.Target
	})
}

// filterBaseline removes findings whose (rule, document, target) tuple is in
// the baseline. Extraction errors are never baseline-eligible.
func filterBaseline(findings []Finding, base []BaselineEntry) (kept []Finding, suppressed int) {
	if len(base) == 0 {
		return findings, 0
	}
	accepted := make(map[BaselineEntry]bool, len(base))
	for _, b := range base {
		accepted[b] = true
	}
	for _, f := range findings {
		if accepted[BaselineEntry{Rule: f.Rule, Document: f.Document, Target: f.Target}] {
			suppressed++
			continue
		}
		kept = append(kept, f)
	}
	return kept, suppressed
}

func writeJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func writeGitHub(w io.Writer, r Report) {
	for _, e := range r.ExtractErrors {
		fmt.Fprintf(w, "::error file=%s::[extract.parse_error] %s\n", e.Path, e.Reason)
	}
	for _, f := range r.Findings {
		fmt.Fprintf(w, "::error file=%s::[%s] %s\n", f.Document, f.Rule, f.Message)
	}
	fmt.Fprintf(w, "::notice::scanned %d governed documents, evaluated %d dependency edges\n",
		r.Stats.Parsed, r.Stats.Edges)
}

func writeTable(w io.Writer, r Report) {
	if len(r.ExtractErrors) == 0 && len(r.Findings) == 0 {
		fmt.Fprintf(w, "OK: scanned %d documents, evaluated %d edges\n", r.Stats.Parsed, r.Stats.Edges)
		return
	}
	fmt.Fprintf(w, "%-10s %-8s %-36s %-36s %s\n", "SEVERITY", "RULE", "DOCUMENT", "TARGET", "MESSAGE")
	for _, e := range r.ExtractErrors {
		fmt.Fprintf(w, "%-10s %-8s %-36s %-36s %s\n", "error", "extract.parse_error", e.Path, "", e.Reason)
	}
	for _, f := range r.Findings {
		fmt.Fprintf(w, "%-10s %-8s %-36s %-36s %s\n", f.Severity, f.Rule, f.Document, f.Target, f.Message)
	}
	fmt.Fprintf(w, "\nscanned %d documents, evaluated %d dependency edges\n", r.Stats.Parsed, r.Stats.Edges)
}
