package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("qualitycheck", flag.ContinueOnError)
	fs.SetOutput(stderr)
	rubricPath := fs.String("rubric", "", "rubric YAML path")
	findingsPath := fs.String("findings", "", "findings JSON path")
	docPath := fs.String("doc", "", "target spec.md path")
	writerSession := fs.String("writer-session", "", "writer session id; required for -audit")
	outPath := fs.String("out", "", "report output path (audit defaults to quality-report.md)")
	referencePath := fs.String("reference", "", "reference findings JSON for -calibrate")
	validate := fs.Bool("validate-rubric", false, "validate the rubric and exit")
	audit := fs.Bool("audit", false, "audit findings against the document")
	calibrate := fs.Bool("calibrate", false, "compare two independent reviewer passes")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	switch {
	case *validate:
		if *rubricPath == "" {
			fmt.Fprintln(stderr, "-rubric is required")
			return 2
		}
		r, err := LoadRubric(*rubricPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		if err := ValidateRubric(r); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "OK: rubric %s\n", r.Version)
		return 0
	case *audit:
		if *rubricPath == "" || *findingsPath == "" || *docPath == "" || *writerSession == "" {
			fmt.Fprintln(stderr, "usage: qualitycheck -audit -rubric <rubric.yaml> -findings <findings.json> -doc <spec.md> -writer-session <id>")
			return 2
		}
		if *outPath != "" {
			for _, input := range []string{*rubricPath, *findingsPath, *docPath} {
				if samePath(*outPath, input) {
					fmt.Fprintf(stderr, "-out %q must not overwrite an input file\n", *outPath)
					return 2
				}
			}
		}
		res, err := Audit(*rubricPath, *findingsPath, *docPath, *writerSession)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		report := RenderReport(res)
		out := *outPath
		if out == "" {
			out = "quality-report.md"
		}
		if err := os.WriteFile(out, []byte(report), 0o644); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		fmt.Fprint(stdout, report)
		if len(res.Failures) > 0 {
			return 1
		}
		return 0
	case *calibrate:
		if *rubricPath == "" || *referencePath == "" || *findingsPath == "" {
			fmt.Fprintln(stderr, "usage: qualitycheck -calibrate -rubric <rubric.yaml> -reference <findings.json> -findings <findings.json> [-out <report.md>]")
			return 2
		}
		res, err := Calibrate(*rubricPath, *referencePath, *findingsPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		report := RenderCalibration(res)
		if *outPath != "" {
			if err := os.WriteFile(*outPath, []byte(report), 0o644); err != nil {
				fmt.Fprintln(stderr, err)
				return 2
			}
		}
		fmt.Fprint(stdout, report)
		return 0
	default:
		fmt.Fprintln(stderr, "one of -validate-rubric or -audit is required")
		return 2
	}
}

func RenderReport(res AuditResult) string {
	var b strings.Builder
	failed := len(res.Failures) > 0
	fmt.Fprintf(&b, "# Spec Quality Report\n\n")
	fmt.Fprintf(&b, "- Target: `%s`\n", escapeCell(res.Target))
	fmt.Fprintf(&b, "- Rubric: %s (`%s`)\n", escapeCell(res.RubricVersion), escapeCell(res.RubricSHA256))
	fmt.Fprintf(&b, "- Target SHA-256: `%s`\n", res.TargetSHA256)
	fmt.Fprintf(&b, "- Reviewer session: `%s`\n\n", escapeCell(res.ReviewerID))
	if failed {
		fmt.Fprintf(&b, "**Headline score: not computed (audit failed)**\n\n")
	} else {
		fmt.Fprintf(&b, "**Headline score (core): %s**  |  Diagnostic score (all): %s\n\n", scoreCell(res.Headline), scoreCell(res.Diagnostic))
	}
	fmt.Fprintf(&b, "> The auditor verifies evidence presence, section placement, and score arithmetic. It does not judge whether evidence is sufficient or true.\n\n")
	fmt.Fprintf(&b, "> Normalization covers full-width ASCII, NBSP, and ideographic space; it is not full Unicode NFKC.\n\n")
	fmt.Fprintf(&b, "| Dimension | Headline | Diagnostic | Weight |\n|---|---:|---:|---:|\n")
	for _, d := range res.Dimensions {
		headline, diagnostic := d.Headline, d.Diagnostic
		if failed {
			headline, diagnostic = -1, -1
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %d |\n", escapeCell(d.Name), scoreCell(headline), scoreCell(diagnostic), d.Weight)
	}
	fmt.Fprintf(&b, "\n")

	for _, d := range res.Dimensions {
		fmt.Fprintf(&b, "## %s (%s)\n\n", d.Name, d.ID)
		for _, it := range d.Items {
			status := fmt.Sprintf("level %d", it.Level)
			if it.Status == "na" {
				status = "na"
			}
			fmt.Fprintf(&b, "- **%s** [%s, %s] — %s\n", it.ID, it.Tier, status, escapeCell(it.Question))
			for _, ev := range it.Evidence {
				fmt.Fprintf(&b, "  - evidence `%s`: \"%s\"\n", escapeCell(ev.Section), escapeCell(ev.Quote))
				fmt.Fprintf(&b, "  - normalized: \"%s\"\n", escapeCell(normalizeForMatch(ev.Quote)))
			}
			if it.Rationale != "" {
				fmt.Fprintf(&b, "  - rationale: %s\n", escapeCell(it.Rationale))
			}
			if it.Missing != "" {
				fmt.Fprintf(&b, "  - missing: %s\n", escapeCell(it.Missing))
			}
			if it.NaReason != "" {
				fmt.Fprintf(&b, "  - na: %s\n", escapeCell(it.NaReason))
			}
			for _, failure := range it.Failures {
				fmt.Fprintf(&b, "  - **audit failure:** %s\n", escapeCell(failure))
			}
		}
		fmt.Fprintf(&b, "\n")
	}

	if len(res.Notes) > 0 {
		fmt.Fprintf(&b, "## Chapter notes\n\n")
		for _, note := range res.Notes {
			fmt.Fprintf(&b, "- %s\n", escapeCell(note))
		}
		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "## Audit failures\n\n")
	if len(res.Failures) == 0 {
		fmt.Fprintf(&b, "None.\n")
	} else {
		for _, failure := range res.Failures {
			fmt.Fprintf(&b, "- %s\n", escapeCell(failure))
		}
	}
	return b.String()
}

func RenderCalibration(res CalibrationResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Calibration Report\n\n")
	fmt.Fprintf(&b, "- Reference: `%s`\n", escapeCell(res.Reference))
	fmt.Fprintf(&b, "- Candidate: `%s`\n", escapeCell(res.Candidate))
	fmt.Fprintf(&b, "- Compared: %d\n", res.Compared)
	fmt.Fprintf(&b, "- Exact agreement: %.1f%%\n", res.Agreement)
	fmt.Fprintf(&b, "- Weighted kappa (quadratic): %.3f\n", res.WeightedKappa)
	if res.CoreCompared > 0 {
		fmt.Fprintf(&b, "- Core items: compared %d, exact %.1f%%, weighted kappa %.3f\n", res.CoreCompared, res.CoreAgreement, res.CoreWeightedKappa)
	}
	fmt.Fprintf(&b, "\n")
	if len(res.StatusMismatches) > 0 {
		fmt.Fprintf(&b, "Signal mismatches (na vs scored): %s\n\n", escapeCell(strings.Join(res.StatusMismatches, ", ")))
	}
	fmt.Fprintf(&b, "| ref\\cand | 0 | 1 | 2 |\n|---|---:|---:|---:|\n")
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&b, "| %d | %d | %d | %d |\n", i, res.Confusion[i][0], res.Confusion[i][1], res.Confusion[i][2])
	}
	fmt.Fprintf(&b, "\n## Dimensions\n\n| Dimension | Compared | Exact | Agreement | Weighted kappa |\n|---|---:|---:|---:|---:|\n")
	for _, d := range res.Dimensions {
		fmt.Fprintf(&b, "| %s | %d | %d | %.1f%% | %.3f |\n", escapeCell(d.Name), d.Compared, d.Exact, d.Agreement, d.WeightedKappa)
	}
	if len(res.Disagreements) > 0 {
		fmt.Fprintf(&b, "\n## Disagreements\n\n")
		for _, d := range res.Disagreements {
			fmt.Fprintf(&b, "- `%s`: reference %d, candidate %d\n", d.ID, d.Reference, d.Candidate)
		}
	}
	fmt.Fprintf(&b, "\n> Inter-rater agreement, not human truth. A high kappa means the instrument is reproducible, not that its levels are correct.\n")
	return b.String()
}

func scoreCell(score float64) string {
	if score < 0 {
		return "—"
	}
	return fmt.Sprintf("%.1f", score)
}

func samePath(a, b string) bool {
	if ai, errA := os.Stat(a); errA == nil {
		if bi, errB := os.Stat(b); errB == nil && os.SameFile(ai, bi) {
			return true
		}
	}
	aa, errA := filepath.Abs(a)
	bb, errB := filepath.Abs(b)
	if errA == nil && errB == nil && filepath.Clean(aa) == filepath.Clean(bb) {
		return true
	}
	ra, errA := filepath.EvalSymlinks(a)
	rb, errB := filepath.EvalSymlinks(b)
	return errA == nil && errB == nil && filepath.Clean(ra) == filepath.Clean(rb)
}

func escapeCell(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "`", "'")
	return s
}
