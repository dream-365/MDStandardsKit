package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const findingsSchemaVersion = 1

// wholeDocSection marks a rubric item whose evidence may come from any section.
const wholeDocSection = "全文"

// Rubric is the machine-readable scoring standard. The rubric is the single
// source of item identity, anchors, weights and section aliases.
type Rubric struct {
	Version      string       `yaml:"version"`
	ScoringRules ScoringRules `yaml:"scoring_rules"`
	Dimensions   []Dimension  `yaml:"dimensions"`
}

type ScoringRules struct {
	CoreNAWithReason bool   `yaml:"core_na_with_reason"`
	Note             string `yaml:"note"`
}

type Dimension struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	Weight int    `yaml:"weight"`
	Items  []Item `yaml:"items"`
}

type Item struct {
	ID           string            `yaml:"id"`
	Question     string            `yaml:"question"`
	Tier         string            `yaml:"tier"`
	Weight       int               `yaml:"weight"`
	Section      string            `yaml:"section"`
	Aliases      []string          `yaml:"aliases"`
	NaAllowed    bool              `yaml:"na_allowed"`
	Levels       map[string]string `yaml:"levels"`
	EvidenceHint string            `yaml:"evidence_hint"`
}

type Findings struct {
	SchemaVersion int           `json:"schema_version"`
	RubricVersion string        `json:"rubric_version"`
	RubricSHA256  string        `json:"rubric_sha256"`
	Target        string        `json:"target"`
	TargetSHA256  string        `json:"target_sha256"`
	Reviewer      Reviewer      `json:"reviewer"`
	Items         []ItemFinding `json:"items"`
}

type Reviewer struct {
	Role      string `json:"role"`
	Model     string `json:"model"`
	SessionID string `json:"session_id"`
	StartedAt string `json:"started_at"`
}

type ItemFinding struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	Level     int        `json:"level"`
	Evidence  []Evidence `json:"evidence"`
	Rationale string     `json:"rationale"`
	Missing   string     `json:"missing"`
	NaReason  string     `json:"na_reason"`
}

type Evidence struct {
	Section string `json:"section"`
	Quote   string `json:"quote"`
}

type ItemResult struct {
	ID            string
	Question      string
	Tier          string
	Weight        int
	DimensionID   string
	DimensionName string
	Status        string
	Level         int
	Evidence      []Evidence
	Rationale     string
	Missing       string
	NaReason      string
	Failures      []string
}

type DimensionResult struct {
	ID         string
	Name       string
	Weight     int
	Headline   float64 // -1 when the dimension has no scorable core item
	Diagnostic float64 // -1 when the dimension has no scorable item
	Items      []ItemResult
}

type AuditResult struct {
	RubricVersion string
	RubricSHA256  string
	Target        string
	TargetSHA256  string
	ReviewerID    string
	Dimensions    []DimensionResult
	Headline      float64
	Diagnostic    float64
	Failures      []string
	Notes         []string
}

func loadFindings(path string) (Findings, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Findings{}, err
	}
	var f Findings
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&f); err != nil {
		return Findings{}, fmt.Errorf("invalid findings JSON: %w", err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		return Findings{}, fmt.Errorf("invalid findings JSON: trailing data after the object")
	}
	return f, nil
}

func LoadRubric(path string) (Rubric, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Rubric{}, err
	}
	var r Rubric
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	if err := dec.Decode(&r); err != nil {
		return Rubric{}, fmt.Errorf("invalid rubric YAML: %w", err)
	}
	return r, nil
}

func ValidateRubric(r Rubric) error {
	if strings.TrimSpace(r.Version) == "" {
		return fmt.Errorf("rubric.version is empty")
	}
	if len(r.Dimensions) == 0 {
		return fmt.Errorf("rubric has no dimensions")
	}
	seenDims := map[string]bool{}
	seenItems := map[string]bool{}
	for _, d := range r.Dimensions {
		if strings.TrimSpace(d.ID) == "" || strings.TrimSpace(d.Name) == "" {
			return fmt.Errorf("dimension id and name are required")
		}
		if seenDims[d.ID] {
			return fmt.Errorf("duplicate dimension id %q", d.ID)
		}
		seenDims[d.ID] = true
		if d.Weight <= 0 {
			return fmt.Errorf("dimension %q weight must be positive", d.ID)
		}
		hasScorable := false
		for _, it := range d.Items {
			if strings.TrimSpace(it.ID) == "" {
				return fmt.Errorf("dimension %q has an item with an empty id", d.ID)
			}
			if seenItems[it.ID] {
				return fmt.Errorf("duplicate item id %q", it.ID)
			}
			seenItems[it.ID] = true
			if it.Tier != "core" && it.Tier != "extended" {
				return fmt.Errorf("item %q tier must be core or extended", it.ID)
			}
			if it.Weight <= 0 {
				return fmt.Errorf("item %q weight must be positive", it.ID)
			}
			if strings.TrimSpace(it.Section) == "" || len(it.Aliases) == 0 {
				return fmt.Errorf("item %q section and aliases are required", it.ID)
			}
			for _, level := range []string{"0", "1", "2"} {
				if strings.TrimSpace(it.Levels[level]) == "" {
					return fmt.Errorf("item %q is missing level %s anchor", it.ID, level)
				}
			}
			if it.Tier == "core" && it.NaAllowed {
				return fmt.Errorf("core item %q must set na_allowed: false", it.ID)
			}
			if it.Section == wholeDocSection {
				found := false
				for _, a := range it.Aliases {
					if a == wholeDocSection {
						found = true
					}
				}
				if !found {
					return fmt.Errorf("item %q uses the whole-document section without the %q alias", it.ID, wholeDocSection)
				}
			}
			if !it.NaAllowed {
				hasScorable = true
			}
		}
		if !hasScorable {
			return fmt.Errorf("dimension %q has no na_allowed: false item", d.ID)
		}
	}
	hasCore := false
	for _, d := range r.Dimensions {
		for _, it := range d.Items {
			if it.Tier == "core" {
				hasCore = true
			}
		}
	}
	if !hasCore {
		return fmt.Errorf("rubric has no core items")
	}
	return nil
}

func Audit(rubricPath, findingsPath, docPath, writerSession string) (AuditResult, error) {
	r, err := LoadRubric(rubricPath)
	if err != nil {
		return AuditResult{}, err
	}
	if err := ValidateRubric(r); err != nil {
		return AuditResult{}, fmt.Errorf("rubric invalid: %w", err)
	}
	rubricBytes, err := os.ReadFile(rubricPath)
	if err != nil {
		return AuditResult{}, err
	}
	docBytes, err := os.ReadFile(docPath)
	if err != nil {
		return AuditResult{}, err
	}
	f, err := loadFindings(findingsPath)
	if err != nil {
		return AuditResult{}, err
	}

	res := AuditResult{
		RubricVersion: r.Version,
		RubricSHA256:  sha256Hex(rubricBytes),
		Target:        f.Target,
		TargetSHA256:  sha256Hex(docBytes),
		ReviewerID:    f.Reviewer.SessionID,
	}

	if f.SchemaVersion != findingsSchemaVersion {
		res.Failures = append(res.Failures, fmt.Sprintf("schema_version = %d, want %d", f.SchemaVersion, findingsSchemaVersion))
	}
	if f.RubricVersion != r.Version {
		res.Failures = append(res.Failures, fmt.Sprintf("rubric_version = %q, want %q", f.RubricVersion, r.Version))
	}
	if f.RubricSHA256 != res.RubricSHA256 {
		res.Failures = append(res.Failures, "rubric_sha256 does not match the rubric file")
	}
	if f.TargetSHA256 != res.TargetSHA256 {
		res.Failures = append(res.Failures, "target_sha256 does not match the target document")
	}
	if strings.TrimSpace(f.Target) == "" {
		res.Failures = append(res.Failures, "target is empty")
	} else {
		cleanTarget := filepath.Clean(f.Target)
		cleanDoc := filepath.Clean(docPath)
		if cleanTarget != cleanDoc && !strings.HasSuffix(filepath.ToSlash(cleanDoc), "/"+filepath.ToSlash(cleanTarget)) {
			res.Failures = append(res.Failures, fmt.Sprintf("target %q does not match -doc %q", f.Target, docPath))
		}
	}
	if f.Reviewer.Role != "independent" {
		res.Failures = append(res.Failures, `reviewer.role must be "independent"`)
	}
	if strings.TrimSpace(f.Reviewer.SessionID) == "" {
		res.Failures = append(res.Failures, "reviewer.session_id is empty")
	}
	if strings.TrimSpace(f.Reviewer.Model) == "" {
		res.Failures = append(res.Failures, "reviewer.model is empty")
	}
	if strings.TrimSpace(f.Reviewer.StartedAt) == "" {
		res.Failures = append(res.Failures, "reviewer.started_at is empty")
	}
	if strings.TrimSpace(writerSession) == "" {
		res.Failures = append(res.Failures, "-writer-session is required")
	} else if f.Reviewer.SessionID == writerSession {
		res.Failures = append(res.Failures, "reviewer session equals the writer session")
	}

	sections := ParseSections(string(docBytes))
	seenHeadings := map[string]bool{}
	for _, s := range sections {
		n := normalizeForMatch(s.Text)
		if n == "" {
			continue
		}
		if seenHeadings[n] {
			res.Failures = append(res.Failures, fmt.Sprintf("duplicate heading %q", s.Text))
		}
		seenHeadings[n] = true
	}
	seenNotes := map[string]bool{}
	sectionAliases := map[string][]string{}
	for _, d := range r.Dimensions {
		for _, item := range d.Items {
			if item.Section == wholeDocSection || seenNotes[item.Section] {
				continue
			}
			sectionAliases[item.Section] = append(sectionAliases[item.Section], item.Aliases...)
			if !sectionExistsAny(sections, item.Section, item.Aliases) {
				res.Notes = append(res.Notes, fmt.Sprintf("section %q not found", item.Section))
				seenNotes[item.Section] = true
			}
		}
	}
	byID := map[string]ItemFinding{}
	for _, it := range f.Items {
		if _, ok := byID[it.ID]; ok {
			res.Failures = append(res.Failures, fmt.Sprintf("duplicate finding for item %q", it.ID))
			continue
		}
		byID[it.ID] = it
	}

	itemToDim := map[string]string{}
	naBySection := map[string][]ItemFinding{}
	for _, d := range r.Dimensions {
		dr := DimensionResult{ID: d.ID, Name: d.Name, Weight: d.Weight}
		for _, item := range d.Items {
			itemToDim[item.ID] = d.ID
			finding, ok := byID[item.ID]
			ir := ItemResult{
				ID:            item.ID,
				Question:      item.Question,
				Tier:          item.Tier,
				Weight:        item.Weight,
				DimensionID:   d.ID,
				DimensionName: d.Name,
				Status:        finding.Status,
				Level:         finding.Level,
				Evidence:      finding.Evidence,
				Rationale:     finding.Rationale,
				Missing:       finding.Missing,
				NaReason:      finding.NaReason,
			}
			if !ok {
				ir.Failures = append(ir.Failures, "missing finding")
			} else {
				ir.Failures = auditItem(item, finding, sections, string(docBytes))
			}
			if finding.Status == "na" {
				naBySection[item.Section] = append(naBySection[item.Section], finding)
			}
			dr.Items = append(dr.Items, ir)
		}
		dr.Headline = scoreItems(dr.Items, true)
		dr.Diagnostic = scoreItems(dr.Items, false)
		res.Dimensions = append(res.Dimensions, dr)
	}

	for _, it := range f.Items {
		if _, ok := itemToDim[it.ID]; !ok {
			res.Failures = append(res.Failures, fmt.Sprintf("unknown item id %q", it.ID))
		}
	}

	for _, d := range res.Dimensions {
		for _, it := range d.Items {
			for _, failure := range it.Failures {
				res.Failures = append(res.Failures, fmt.Sprintf("%s: %s", it.ID, failure))
			}
		}
	}

	// A present section whose entire rubric group is na is a whole-section
	// evasion. Missing sections may be na only for extended items.
	var naSections []string
	for section := range naBySection {
		naSections = append(naSections, section)
	}
	sort.Strings(naSections)
	for _, section := range naSections {
		if section == wholeDocSection {
			continue
		}
		if !sectionExistsAny(sections, section, sectionAliases[section]) {
			continue
		}
		allNA := true
		for _, d := range r.Dimensions {
			for _, item := range d.Items {
				if item.Section != section {
					continue
				}
				f, ok := byID[item.ID]
				if !ok || f.Status != "na" {
					allNA = false
					break
				}
			}
			if !allNA {
				break
			}
		}
		if allNA {
			res.Failures = append(res.Failures, fmt.Sprintf("section %q has all items marked na", section))
		}
	}

	res.Headline = weightedAverage(res.Dimensions, true)
	res.Diagnostic = weightedAverage(res.Dimensions, false)
	return res, nil
}

func auditItem(item Item, f ItemFinding, sections []Section, doc string) []string {
	var failures []string
	switch f.Status {
	case "scored":
		if strings.TrimSpace(f.Rationale) == "" {
			failures = append(failures, "scored item requires a rationale")
		}
		if f.Level < 0 || f.Level > 2 {
			failures = append(failures, fmt.Sprintf("level %d out of range", f.Level))
			return failures
		}
		if f.Level == 0 {
			if strings.TrimSpace(f.Missing) == "" {
				failures = append(failures, "level 0 requires a missing reason")
			}
			return failures
		}
		if len(f.Evidence) == 0 {
			failures = append(failures, "level > 0 requires evidence")
			return failures
		}
		for _, ev := range f.Evidence {
			if strings.TrimSpace(ev.Quote) == "" {
				failures = append(failures, "evidence quote is empty")
				continue
			}
			if !evidenceAllowed(item, ev.Section) {
				failures = append(failures, fmt.Sprintf("evidence section %q is not declared by item %q", ev.Section, item.ID))
				continue
			}
			text := doc
			if item.Section != wholeDocSection {
				text = sectionTextForEvidence(item, ev.Section, sections)
				if text == "" {
					failures = append(failures, fmt.Sprintf("section %q not found", ev.Section))
					continue
				}
			}
			if !strings.Contains(normalizeForMatch(text), normalizeForMatch(ev.Quote)) {
				failures = append(failures, fmt.Sprintf("quote not found in section %q", ev.Section))
			}
		}
	case "na":
		if !item.NaAllowed {
			failures = append(failures, "na is not allowed for this item")
		}
		if strings.TrimSpace(f.NaReason) == "" {
			failures = append(failures, "na requires a reason")
		}
	default:
		failures = append(failures, fmt.Sprintf("unknown status %q", f.Status))
	}
	return failures
}

func evidenceAllowed(item Item, evidenceSection string) bool {
	if item.Section == wholeDocSection {
		return true
	}
	norm := normalizeForMatch(evidenceSection)
	candidates := append([]string{item.Section}, item.Aliases...)
	for _, c := range candidates {
		if normalizeForMatch(c) == norm {
			return true
		}
	}
	return false
}

func sectionText(sections []Section, name string) string {
	norm := normalizeForMatch(name)
	for _, s := range sections {
		if normalizeForMatch(s.Text) == norm {
			return s.Body
		}
	}
	return ""
}

func sectionTextForItem(item Item, sections []Section) string {
	candidates := append([]string{item.Section}, item.Aliases...)
	for _, s := range sections {
		for _, c := range candidates {
			if normalizeForMatch(s.Text) == normalizeForMatch(c) {
				return s.Body
			}
		}
	}
	return ""
}

// sectionTextForEvidence prefers the heading the reviewer cited. Only when the
// cited alias has no heading does it fall back to the item's declared heading.
func sectionTextForEvidence(item Item, evidenceSection string, sections []Section) string {
	if body := sectionText(sections, evidenceSection); body != "" {
		return body
	}
	return sectionTextForItem(item, sections)
}

func sectionExistsAny(sections []Section, section string, aliases []string) bool {
	candidates := append([]string{section}, aliases...)
	for _, s := range sections {
		for _, c := range candidates {
			if normalizeForMatch(s.Text) == normalizeForMatch(c) {
				return true
			}
		}
	}
	return false
}

func scoreItems(items []ItemResult, coreOnly bool) float64 {
	var num, den float64
	for _, it := range items {
		if it.Status == "na" && len(it.Failures) == 0 {
			continue
		}
		if coreOnly && it.Tier != "core" {
			continue
		}
		// Missing findings, failed evidence and illegal na score 0, never
		// leave the denominator: an unverifiable item must not raise the score.
		den += float64(it.Weight * 2)
		if it.Status != "scored" || len(it.Failures) > 0 {
			continue
		}
		num += float64(it.Weight * it.Level)
	}
	if den == 0 {
		return -1
	}
	return 100 * num / den
}

func weightedAverage(dims []DimensionResult, headline bool) float64 {
	var num, den float64
	for _, d := range dims {
		score := d.Diagnostic
		if headline {
			score = d.Headline
		}
		if score < 0 {
			continue
		}
		num += float64(d.Weight) * score
		den += float64(d.Weight)
	}
	if den == 0 {
		return -1
	}
	return num / den
}

type Section struct {
	Text string
	Rank int // heading level; 0 means the preamble
	Body string
}

var headingRE = regexp.MustCompile(`^(#{1,6})\s+(.*?)(?:\s+#+)?\s*$`)

// ParseSections returns ATX headings with their body text. Headings inside
// backtick or tilde fences are ignored. Setext headings are not parsed; a
// setext document reports its sections as missing rather than misbinding them.
func ParseSections(doc string) []Section {
	lines := strings.Split(strings.ReplaceAll(doc, "\r\n", "\n"), "\n")
	type heading struct {
		section Section
		line    int
	}
	var hs []heading
	inFence := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		m := headingRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		hs = append(hs, heading{Section{Text: strings.TrimSpace(m[2]), Rank: len(m[1])}, i})
	}
	out := make([]Section, len(hs))
	for i := range hs {
		end := len(lines)
		for j := i + 1; j < len(hs); j++ {
			if hs[j].section.Rank <= hs[i].section.Rank {
				end = hs[j].line
				break
			}
		}
		s := hs[i].section
		var parts []string
		inBodyFence := false
		for _, line := range lines[hs[i].line+1 : end] {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
				inBodyFence = !inBodyFence
				parts = append(parts, line)
				continue
			}
			if !inBodyFence && headingRE.MatchString(line) {
				continue
			}
			parts = append(parts, line)
		}
		s.Body = strings.Join(parts, "\n")
		out[i] = s
	}
	return out
}

// normalizeForMatch is the NFKC-compatible subset the tool uses without a new
// dependency: NBSP and ideographic space fold to a space, and full-width ASCII
// folds to ASCII. True Unicode NFKC (ligatures, circled characters, half-width
// katakana) is out of scope; the report states this boundary.
func normalizeForMatch(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == ' ' || r == '　':
			b.WriteRune(' ')
		case r >= '！' && r <= '～':
			b.WriteRune(r - 0xFEE0)
		default:
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

type CalibrationDisagreement struct {
	ID        string
	Reference int
	Candidate int
}

type CalibrationDimension struct {
	ID            string
	Name          string
	Compared      int
	Exact         int
	Agreement     float64
	WeightedKappa float64
}

type CalibrationResult struct {
	Reference         string
	Candidate         string
	Compared          int
	Exact             int
	Agreement         float64
	WeightedKappa     float64
	CoreCompared      int
	CoreExact         int
	CoreAgreement     float64
	CoreWeightedKappa float64
	Confusion         [3][3]int
	Dimensions        []CalibrationDimension
	Disagreements     []CalibrationDisagreement
	StatusMismatches  []string
}

// Calibrate compares two independent reviewer passes. It measures agreement,
// not truth: the reference is another reviewer, not a human gold standard.
func Calibrate(rubricPath, referencePath, candidatePath string) (CalibrationResult, error) {
	r, err := LoadRubric(rubricPath)
	if err != nil {
		return CalibrationResult{}, err
	}
	if err := ValidateRubric(r); err != nil {
		return CalibrationResult{}, fmt.Errorf("rubric invalid: %w", err)
	}
	ref, err := loadFindings(referencePath)
	if err != nil {
		return CalibrationResult{}, err
	}
	cand, err := loadFindings(candidatePath)
	if err != nil {
		return CalibrationResult{}, err
	}
	if ref.RubricVersion != r.Version || cand.RubricVersion != r.Version {
		return CalibrationResult{}, fmt.Errorf("findings rubric_version does not match %s", r.Version)
	}
	if filepath.Clean(ref.Target) != filepath.Clean(cand.Target) {
		return CalibrationResult{}, fmt.Errorf("reference and candidate targets differ")
	}
	refBy := map[string]ItemFinding{}
	for _, it := range ref.Items {
		refBy[it.ID] = it
	}
	candBy := map[string]ItemFinding{}
	for _, it := range cand.Items {
		candBy[it.ID] = it
	}

	res := CalibrationResult{Reference: referencePath, Candidate: candidatePath}
	dimConf := map[string]*[3][3]int{}
	dimCompared := map[string]int{}
	var coreConf [3][3]int
	for _, d := range r.Dimensions {
		dimConf[d.ID] = &[3][3]int{}
		for _, item := range d.Items {
			rf, ok1 := refBy[item.ID]
			cf, ok2 := candBy[item.ID]
			if !ok1 || !ok2 {
				continue
			}
			if rf.Status == "na" || cf.Status == "na" {
				if rf.Status != cf.Status {
					res.StatusMismatches = append(res.StatusMismatches, item.ID)
				}
				continue
			}
			if rf.Status != "scored" || cf.Status != "scored" {
				continue
			}
			if rf.Level < 0 || rf.Level > 2 || cf.Level < 0 || cf.Level > 2 {
				continue
			}
			res.Confusion[rf.Level][cf.Level]++
			dimConf[d.ID][rf.Level][cf.Level]++
			res.Compared++
			dimCompared[d.ID]++
			if item.Tier == "core" {
				coreConf[rf.Level][cf.Level]++
				res.CoreCompared++
			}
			if rf.Level == cf.Level {
				res.Exact++
				if item.Tier == "core" {
					res.CoreExact++
				}
			} else {
				res.Disagreements = append(res.Disagreements, CalibrationDisagreement{ID: item.ID, Reference: rf.Level, Candidate: cf.Level})
			}
		}
	}
	if res.Compared == 0 {
		return res, nil
	}
	res.Agreement = 100 * float64(res.Exact) / float64(res.Compared)
	res.WeightedKappa = weightedKappa(res.Confusion, res.Compared)
	if res.CoreCompared > 0 {
		res.CoreAgreement = 100 * float64(res.CoreExact) / float64(res.CoreCompared)
		res.CoreWeightedKappa = weightedKappa(coreConf, res.CoreCompared)
	}
	for _, d := range r.Dimensions {
		if dimCompared[d.ID] == 0 {
			continue
		}
		cd := CalibrationDimension{ID: d.ID, Name: d.Name, Compared: dimCompared[d.ID]}
		for i := 0; i < 3; i++ {
			cd.Exact += dimConf[d.ID][i][i]
		}
		cd.Agreement = 100 * float64(cd.Exact) / float64(cd.Compared)
		cd.WeightedKappa = weightedKappa(*dimConf[d.ID], cd.Compared)
		res.Dimensions = append(res.Dimensions, cd)
	}
	return res, nil
}

func weightedKappa(conf [3][3]int, n int) float64 {
	if n == 0 {
		return 0
	}
	var row, col [3]float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			row[i] += float64(conf[i][j])
			col[j] += float64(conf[i][j])
		}
	}
	var observed, expected float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			w := 1 - float64((i-j)*(i-j))/4
			observed += float64(conf[i][j]) * w
			expected += row[i] * col[j] * w
		}
	}
	observed /= float64(n)
	expected /= float64(n * n)
	if 1-expected == 0 {
		if observed == 1 {
			return 1
		}
		return 0
	}
	return (observed - expected) / (1 - expected)
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
