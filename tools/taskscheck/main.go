// taskscheck validates a dependency-ordered tasks.md against its sibling
// spec.md. It is deliberately a separate binary: mdstandards keeps its
// document-reference graph contract, and taskscheck owns the task graph.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	schemaVersion = "1"
	tableHeader   = "| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |"
)

var (
	version = "dev"

	taskIDPattern  = regexp.MustCompile(`^T[A-Za-z0-9_-]*$`)
	goalPattern    = regexp.MustCompile(`^\s*(?:[-*]\s+)?\*\*(G\d+)\b`)
	reqPattern     = regexp.MustCompile(`^\s*(?:[-*]\s+)?\*\*(REQ-\d{3})\b`)
	shaPattern     = regexp.MustCompile(`^[0-9a-f]{64}$`)
	separatorCell  = regexp.MustCompile(`^:?-{3,}:?$`)
	dateLayout     = "2006-01-02"
	bannedDoneText = map[string]bool{
		"完成": true, "works": true, "tested": true, "tbd": true, "待补充": true,
		"功能可用": true, "已修复": true, "手动看看": true, "跑一下": true, "看看": true,
	}
	verifyPrefixes = []string{
		"go test", "go build", "go vet", "make ", "./", "sh ", "bash ",
		"curl ", "test ", "python", "node ", "npm ", "cargo ", "shasum", "sha256sum",
	}
)

type task struct {
	ID      string
	Title   string
	Covers  []string
	Depends []string
	Writes  []string
	Done    string
	Verify  string
	Line    int
}

type tasksDoc struct {
	Path            string
	SpecField       string
	SpecHash        string
	Status          string
	NotApplicable   string
	Tasks           []task
	HasTable        bool
	FrontmatterErrs []Finding
}

type specInfo struct {
	Path   string
	Status string
	Body   string
}

type Finding struct {
	Rule     string `json:"rule"`
	Spec     string `json:"spec"`
	Task     string `json:"task,omitempty"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type CheckError struct {
	Rule    string `json:"rule"`
	Spec    string `json:"spec"`
	Task    string `json:"task,omitempty"`
	Message string `json:"message"`
}

type Stats struct {
	SpecsChecked int `json:"specs_checked"`
	TasksChecked int `json:"tasks_checked"`
	Warnings     int `json:"warnings"`
	Errors       int `json:"errors"`
}

type HotPath struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

type OverCovered struct {
	Task   string `json:"task"`
	Covers int    `json:"covers"`
}

type PlanShape struct {
	CriticalPath          []string      `json:"critical_path"`
	ParallelGroups        [][]string    `json:"parallel_groups"`
	HotPaths              []HotPath     `json:"hot_paths"`
	OverCoveredTasks      []OverCovered `json:"over_covered_tasks"`
	UncoveredRequirements []string      `json:"uncovered_requirements"`
}

type Report struct {
	SchemaVersion string       `json:"schema_version"`
	Findings      []Finding    `json:"findings"`
	Errors        []CheckError `json:"errors"`
	Stats         Stats        `json:"stats"`
	PlanShape     PlanShape    `json:"plan_shape"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fsx := flag.NewFlagSet("taskscheck", flag.ContinueOnError)
	fsx.SetOutput(stderr)
	showVersion := fsx.Bool("version", false, "print version and exit")
	root := fsx.String("root", ".", "repository root")
	specFlag := fsx.String("spec", "", "comma-separated spec.md paths; empty discovers docs/specs/**/spec.md")
	phase := fsx.Int("phase", 1, "rollout phase: 0 report-only, 1 block")
	format := fsx.String("format", "table", "output format: table|json|github")
	outPath := fsx.String("out", "", "write the report to this file instead of stdout")
	printHash := fsx.String("print-spec-hash", "", "print the sha256 of a spec.md and exit")
	if err := fsx.Parse(args); err != nil {
		return 2
	}
	if *showVersion {
		fmt.Fprintf(stdout, "taskscheck %s\n", version)
		return 0
	}
	if *phase != 0 && *phase != 1 {
		fmt.Fprintln(stderr, "tasks.usage: phase must be 0 or 1")
		return 2
	}
	switch *format {
	case "table", "json", "github":
	default:
		fmt.Fprintf(stderr, "tasks.usage: unknown format %q\n", *format)
		return 2
	}
	if *printHash != "" {
		sum, err := fileHash(*printHash)
		if err != nil {
			fmt.Fprintf(stderr, "tasks.parse_error: %v\n", err)
			return 2
		}
		fmt.Fprintln(stdout, sum)
		return 0
	}

	rootAbs, err := filepath.Abs(*root)
	if err != nil {
		fmt.Fprintf(stderr, "tasks.usage: %v\n", err)
		return 2
	}

	var specs []string
	var errorsList []CheckError
	if strings.TrimSpace(*specFlag) != "" {
		for _, raw := range strings.Split(*specFlag, ",") {
			p := strings.TrimSpace(raw)
			if p == "" {
				errorsList = append(errorsList, CheckError{Rule: "tasks.invalid_spec_flag", Message: "empty -spec entry"})
				continue
			}
			abs, err := filepath.Abs(p)
			if err != nil {
				errorsList = append(errorsList, CheckError{Rule: "tasks.invalid_spec_flag", Message: p + ": " + err.Error()})
				continue
			}
			if !inside(rootAbs, abs) {
				errorsList = append(errorsList, CheckError{Rule: "tasks.invalid_spec_flag", Spec: p, Message: "-spec must be inside -root"})
				continue
			}
			info, err := os.Stat(abs)
			if err != nil || info.IsDir() || filepath.Base(abs) != "spec.md" {
				errorsList = append(errorsList, CheckError{Rule: "tasks.invalid_spec_flag", Spec: p, Message: "-spec must name an existing spec.md"})
				continue
			}
			specs = append(specs, abs)
		}
	} else {
		specs, err = discoverSpecs(rootAbs)
		if err != nil {
			fmt.Fprintf(stderr, "tasks.parse_error: %v\n", err)
			return 2
		}
		if len(specs) == 0 {
			errorsList = append(errorsList, CheckError{Rule: "tasks.no_specs", Message: "no docs/specs/**/spec.md found"})
		}
	}
	specs = uniqueSorted(specs)

	orphans, err := discoverOrphans(rootAbs)
	if err != nil {
		fmt.Fprintf(stderr, "tasks.parse_error: %v\n", err)
		return 2
	}

	var findings []Finding
	for _, orphan := range orphans {
		findings = append(findings, Finding{
			Rule: "tasks.orphan_tasks", Spec: orphan, Severity: "error",
			Message: "tasks.md has no sibling spec.md",
		})
	}

	stats := Stats{SpecsChecked: len(specs)}
	var planShape *PlanShape
	for _, specPath := range specs {
		spec, err := readSpec(specPath)
		if err != nil {
			errorsList = append(errorsList, CheckError{Rule: "tasks.parse_error", Spec: specPath, Message: err.Error()})
			continue
		}
		tasksPath := filepath.Join(filepath.Dir(specPath), "tasks.md")
		info, statErr := os.Stat(tasksPath)
		if statErr != nil || info.IsDir() {
			if spec.Status == "Approved" || spec.Status == "Implemented" {
				findings = append(findings, Finding{
					Rule: "tasks.missing_file", Spec: specPath, Severity: "error",
					Message: "tasks.md is required for status " + spec.Status,
				})
			}
			continue
		}

		doc, parseErr := parseTasks(tasksPath)
		if parseErr != nil {
			errorsList = append(errorsList, CheckError{Rule: "tasks.parse_error", Spec: specPath, Message: parseErr.Error()})
			continue
		}
		stats.TasksChecked += len(doc.Tasks)

		for _, f := range doc.FrontmatterErrs {
			f.Spec = specPath
			findings = append(findings, f)
		}
		refFindings, refErrors := validateSpecRef(doc, tasksPath, specPath, rootAbs)
		findings = append(findings, refFindings...)
		errorsList = append(errorsList, refErrors...)
		if doc.NotApplicable != "" {
			continue
		}
		if !doc.HasTable {
			continue
		}

		goals, reqs, err := extractSpecIDs(spec.Body)
		if err != nil {
			errorsList = append(errorsList, CheckError{Rule: "spec.no_goals", Spec: specPath, Message: err.Error()})
			continue
		}
		checkFindings, ps := validateTasks(doc, specPath, goals, reqs)
		findings = append(findings, checkFindings...)
		if ps != nil {
			planShape = ps
		}
	}

	if planShape == nil {
		planShape = &PlanShape{
			CriticalPath: []string{}, ParallelGroups: [][]string{}, HotPaths: []HotPath{},
			OverCoveredTasks: []OverCovered{}, UncoveredRequirements: []string{},
		}
	}
	_ = planShape
	for _, f := range findings {
		if f.Severity == "warning" {
			stats.Warnings++
		}
	}
	stats.Errors = len(errorsList)

	sortFindings(findings)
	sortErrors(errorsList)
	if findings == nil {
		findings = []Finding{}
	}
	if errorsList == nil {
		errorsList = []CheckError{}
	}
	rep := Report{
		SchemaVersion: schemaVersion,
		Findings:      findings,
		Errors:        errorsList,
		Stats:         stats,
		PlanShape:     *planShape,
	}

	var out io.Writer = stdout
	if *outPath != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		defer f.Close()
		out = f
	}
	switch *format {
	case "json":
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(rep); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
	case "github":
		writeGitHub(out, rep)
	default:
		writeTable(out, rep)
	}

	if len(errorsList) > 0 {
		return 2
	}
	if *phase >= 1 && hasErrorFinding(findings) {
		return 1
	}
	return 0
}

func hasErrorFinding(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == "error" {
			return true
		}
	}
	return false
}

func sortFindings(f []Finding) {
	sort.Slice(f, func(i, j int) bool {
		if f[i].Spec != f[j].Spec {
			return f[i].Spec < f[j].Spec
		}
		if f[i].Rule != f[j].Rule {
			return f[i].Rule < f[j].Rule
		}
		return f[i].Task < f[j].Task
	})
}

func sortErrors(e []CheckError) {
	sort.Slice(e, func(i, j int) bool {
		if e[i].Spec != e[j].Spec {
			return e[i].Spec < e[j].Spec
		}
		if e[i].Rule != e[j].Rule {
			return e[i].Rule < e[j].Rule
		}
		return e[i].Task < e[j].Task
	})
}

func uniqueSorted(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

func inside(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != "..")
}

func fileHash(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func discoverSpecs(root string) ([]string, error) {
	return walkNamed(filepath.Join(root, "docs", "specs"), "spec.md")
}

func discoverOrphans(root string) ([]string, error) {
	all, err := walkNamed(filepath.Join(root, "docs", "specs"), "tasks.md")
	if err != nil {
		return nil, err
	}
	var out []string
	for _, p := range all {
		spec := filepath.Join(filepath.Dir(p), "spec.md")
		if _, err := os.Stat(spec); err != nil {
			out = append(out, p)
		}
	}
	return out, nil
}

func walkNamed(root, name string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			if d.Name() == "archive" {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == name {
			out = append(out, path)
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

func readSpec(path string) (specInfo, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return specInfo{}, err
	}
	fm, body, err := splitFrontmatter(b)
	if err != nil {
		return specInfo{}, err
	}
	return specInfo{Path: path, Status: yamlString(fm["status"]), Body: body}, nil
}

func parseTasks(path string) (tasksDoc, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return tasksDoc{}, err
	}
	fm, body, err := splitFrontmatter(b)
	if err != nil {
		return tasksDoc{}, err
	}
	doc := tasksDoc{Path: path}
	doc.SpecField = yamlString(fm["spec"])
	doc.SpecHash = yamlString(fm["spec_sha256"])
	doc.Status = yamlString(fm["status"])
	doc.NotApplicable = yamlString(fm["not_applicable"])

	if yamlString(fm["type"]) != "tasks" {
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_frontmatter", "type must be tasks"))
	}
	switch doc.Status {
	case "Draft", "In Progress", "Done":
	default:
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_frontmatter", "status must be Draft, In Progress or Done"))
	}
	if yamlString(fm["owner"]) == "" {
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_frontmatter", "owner is required"))
	}
	if !emptyList(fm["depends_on"]) {
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_frontmatter", "depends_on must be an empty list"))
	}
	if doc.SpecField == "" {
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_frontmatter", "spec is required"))
	}
	if !shaPattern.MatchString(doc.SpecHash) {
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_frontmatter", "spec_sha256 must be 64 lowercase hex characters"))
	}
	if !validDate(yamlString(fm["updated_at"])) {
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_frontmatter", "updated_at must be a real YYYY-MM-DD date"))
	}
	if _, ok := fm["not_applicable"]; ok && doc.NotApplicable == "" {
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_not_applicable", "not_applicable must be a non-empty string"))
	}

	tasks, hasTable, err := parseTable(body)
	if err != nil {
		return tasksDoc{}, err
	}
	doc.Tasks = tasks
	doc.HasTable = hasTable
	if doc.NotApplicable != "" && hasTable {
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_not_applicable", "not_applicable cannot be combined with a task table"))
	}
	if doc.NotApplicable == "" && !hasTable {
		doc.FrontmatterErrs = append(doc.FrontmatterErrs, fmFinding("tasks.invalid_frontmatter", "missing task table"))
	}
	return doc, nil
}

func fmFinding(rule, msg string) Finding {
	return Finding{Rule: rule, Severity: "error", Message: msg}
}

func emptyList(v any) bool {
	items, ok := v.([]any)
	return ok && len(items) == 0
}

func validDate(s string) bool {
	_, err := time.Parse(dateLayout, s)
	return err == nil
}

func yamlString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case time.Time:
		return x.Format(dateLayout)
	default:
		return ""
	}
}

func splitFrontmatter(b []byte) (map[string]any, string, error) {
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return nil, "", errors.New("missing frontmatter")
	}
	rest := s[4:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return nil, "", errors.New("frontmatter has no closing delimiter")
	}
	fmText := rest[:idx]
	body := strings.TrimPrefix(rest[idx+4:], "\n")
	if strings.Count(fmText, "*") > 100 || strings.Count(fmText, "<<") > 10 {
		return nil, "", errors.New("frontmatter rejected: possible YAML alias expansion")
	}
	var m map[string]any
	if err := yaml.Unmarshal([]byte(fmText), &m); err != nil {
		return nil, "", fmt.Errorf("invalid YAML: %w", err)
	}
	if m == nil {
		return nil, "", errors.New("frontmatter is not a mapping")
	}
	return m, body, nil
}

func parseTable(body string) ([]task, bool, error) {
	lines := strings.Split(body, "\n")
	inFence := false
	fence := ""
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if m, ok := fenceMarker(line); ok {
			if !inFence {
				inFence = true
				fence = m
			} else if m == fence {
				inFence = false
				fence = ""
			}
			continue
		}
		if inFence {
			continue
		}
		if strings.TrimSpace(line) != tableHeader {
			continue
		}
		if i+1 >= len(lines) {
			return nil, true, errors.New("task table has no separator row")
		}
		sep := splitTableRow(lines[i+1])
		if len(sep) != 7 {
			return nil, true, errors.New("task table separator has the wrong column count")
		}
		for _, c := range sep {
			if !separatorCell.MatchString(c) {
				return nil, true, errors.New("task table has an invalid separator row")
			}
		}
		var tasks []task
		for j := i + 2; j < len(lines); j++ {
			row := lines[j]
			if !strings.HasPrefix(strings.TrimSpace(row), "|") {
				break
			}
			cells := splitTableRow(row)
			if len(cells) != 7 {
				return nil, true, fmt.Errorf("task row %d has %d columns, want 7", j+1, len(cells))
			}
			tasks = append(tasks, task{
				ID: cells[0], Title: cells[1],
				Covers: splitList(cells[2]), Depends: splitList(cells[3]),
				Writes: splitList(cells[4]), Done: cells[5], Verify: cells[6],
				Line: j + 1,
			})
		}
		return tasks, true, nil
	}
	return nil, false, nil
}

func fenceMarker(line string) (string, bool) {
	t := strings.TrimSpace(line)
	if strings.HasPrefix(t, "```") {
		return "```", true
	}
	if strings.HasPrefix(t, "~~~") {
		return "~~~", true
	}
	return "", false
}

func splitTableRow(line string) []string {
	const pipePlaceholder = "\x00PIPE\x00"
	s := strings.ReplaceAll(line, `\|`, pipePlaceholder)
	cells := strings.Split(s, "|")
	if len(cells) > 0 && strings.TrimSpace(cells[0]) == "" {
		cells = cells[1:]
	}
	if len(cells) > 0 && strings.TrimSpace(cells[len(cells)-1]) == "" {
		cells = cells[:len(cells)-1]
	}
	for i := range cells {
		cells[i] = strings.ReplaceAll(strings.TrimSpace(cells[i]), pipePlaceholder, "|")
	}
	return cells
}

func splitList(cell string) []string {
	cell = strings.TrimSpace(cell)
	if cell == "" || cell == "—" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(cell, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func extractSpecIDs(body string) (map[string]bool, map[string]bool, error) {
	goals := map[string]bool{}
	reqs := map[string]bool{}
	for _, line := range sectionLines(body, "### 2.1 目标", 3) {
		if m := goalPattern.FindStringSubmatch(line); m != nil {
			goals[m[1]] = true
		}
	}
	if len(goals) == 0 {
		return nil, nil, errors.New("spec has no bold Goal declarations under ### 2.1 目标")
	}
	for _, line := range sectionLines(body, "## 4. 详细设计", 2) {
		if m := reqPattern.FindStringSubmatch(line); m != nil {
			reqs[m[1]] = true
		}
	}
	return goals, reqs, nil
}

func sectionLines(body, heading string, level int) []string {
	lines := strings.Split(body, "\n")
	var out []string
	in := false
	for _, line := range lines {
		if lvl := headingLevel(line); lvl > 0 {
			if in && lvl <= level {
				break
			}
			if strings.TrimSpace(line) == heading {
				in = true
				continue
			}
		}
		if in {
			out = append(out, line)
		}
	}
	return out
}

func headingLevel(line string) int {
	t := strings.TrimLeft(line, " ")
	if !strings.HasPrefix(t, "#") {
		return 0
	}
	n := 0
	for n < len(t) && t[n] == '#' {
		n++
	}
	if n > 0 && n < len(t) && t[n] != ' ' {
		return 0
	}
	return n
}

func validateSpecRef(doc tasksDoc, tasksPath, expectedSpec, root string) ([]Finding, []CheckError) {
	var findings []Finding
	var errs []CheckError
	specField := strings.TrimSpace(doc.SpecField)
	if specField == "" || strings.Contains(specField, "://") || filepath.IsAbs(specField) || escapes(specField) {
		findings = append(findings, Finding{
			Rule: "tasks.invalid_frontmatter", Spec: expectedSpec, Severity: "error",
			Message: "spec must be a relative path without ..",
		})
		return findings, errs
	}
	resolved := filepath.Clean(filepath.Join(filepath.Dir(tasksPath), specField))
	if _, err := os.Stat(resolved); err != nil {
		findings = append(findings, Finding{
			Rule: "tasks.missing_spec", Spec: expectedSpec, Severity: "error",
			Message: "spec target does not exist: " + specField,
		})
		return findings, errs
	}
	if filepath.Clean(resolved) != filepath.Clean(expectedSpec) {
		findings = append(findings, Finding{
			Rule: "tasks.wrong_spec_path", Spec: expectedSpec, Severity: "error",
			Message: "spec must point to the sibling spec.md",
		})
		return findings, errs
	}
	if doc.NotApplicable != "" {
		return findings, errs
	}
	sum, err := fileHash(expectedSpec)
	if err != nil {
		errs = append(errs, CheckError{Rule: "tasks.parse_error", Spec: expectedSpec, Message: err.Error()})
		return findings, errs
	}
	if doc.SpecHash != sum {
		findings = append(findings, Finding{
			Rule: "tasks.stale_plan", Spec: expectedSpec, Severity: "error",
			Message: fmt.Sprintf("spec_sha256 mismatch: expected %s got %s; run taskscheck -print-spec-hash %s", sum, doc.SpecHash, expectedSpec),
		})
	}
	return findings, errs
}

func escapes(p string) bool {
	for _, part := range strings.Split(filepath.ToSlash(p), "/") {
		if part == ".." {
			return true
		}
	}
	return false
}

func normalizeWritePath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", errors.New("empty write path")
	}
	if strings.Contains(p, "://") {
		return "", errors.New("write path must not be a URL")
	}
	p = strings.ReplaceAll(p, "\\", "/")
	if strings.HasPrefix(p, "/") || (len(p) >= 2 && p[1] == ':') {
		return "", errors.New("write path must be repo-relative")
	}
	trailing := strings.HasSuffix(p, "/")
	var parts []string
	for _, s := range strings.Split(p, "/") {
		if s == "" || s == "." {
			continue
		}
		if s == ".." {
			return "", errors.New("write path must not escape the repository")
		}
		parts = append(parts, s)
	}
	if len(parts) == 0 {
		return "", errors.New("empty write path")
	}
	res := strings.Join(parts, "/")
	if trailing {
		res += "/"
	}
	return res, nil
}

func isEmptyCell(s string) bool {
	s = strings.TrimSpace(s)
	return s == "" || s == "—"
}

func normalizeDone(s string) string {
	s = strings.TrimSpace(s)
	for {
		before := s
		s = strings.Trim(s, "*_`")
		s = strings.TrimRight(s, "。．.!！")
		s = strings.TrimSpace(s)
		if s == before {
			break
		}
	}
	return strings.ToLower(s)
}

func verifyExecutable(s string) bool {
	s = strings.TrimSpace(strings.Trim(s, "`"))
	if s == "" {
		return false
	}
	if strings.HasPrefix(s, "Test") {
		return true
	}
	for _, p := range verifyPrefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	if !strings.ContainsAny(s, " \t") && strings.Contains(s, "/") {
		return true
	}
	return false
}

func validateTasks(doc tasksDoc, specPath string, goals, reqs map[string]bool) ([]Finding, *PlanShape) {
	var findings []Finding
	byID := map[string]task{}
	for _, t := range doc.Tasks {
		if !taskIDPattern.MatchString(t.ID) {
			findings = append(findings, Finding{Rule: "tasks.invalid_task_id", Spec: specPath, Task: t.ID, Severity: "error", Message: "task ID must match T[A-Za-z0-9_-]*"})
		}
		if _, dup := byID[t.ID]; dup {
			findings = append(findings, Finding{Rule: "tasks.duplicate_id", Spec: specPath, Task: t.ID, Severity: "error", Message: "duplicate task ID"})
		}
		byID[t.ID] = t
		if len(t.Covers) == 0 {
			findings = append(findings, Finding{Rule: "tasks.empty_coverage", Spec: specPath, Task: t.ID, Severity: "error", Message: "task covers no Goal or REQ"})
		}
		for _, c := range t.Covers {
			if !goals[c] && !reqs[c] {
				findings = append(findings, Finding{Rule: "tasks.unknown_coverage", Spec: specPath, Task: t.ID, Severity: "error", Message: "covered ID does not exist in spec.md: " + c})
			}
		}
		for _, d := range t.Depends {
			if d == t.ID {
				findings = append(findings, Finding{Rule: "tasks.self_dependency", Spec: specPath, Task: t.ID, Severity: "error", Message: "task depends on itself"})
			} else if _, ok := byID[d]; !ok {
				findings = append(findings, Finding{Rule: "tasks.unknown_dependency", Spec: specPath, Task: t.ID, Severity: "error", Message: "unknown dependency: " + d})
			}
		}
		done := normalizeDone(t.Done)
		if isEmptyCell(t.Done) {
			done = ""
		}
		if done == "" {
			findings = append(findings, Finding{Rule: "tasks.empty_done_criterion", Spec: specPath, Task: t.ID, Severity: "error", Message: "done criterion is empty"})
		} else if bannedDoneText[done] {
			findings = append(findings, Finding{Rule: "tasks.banned_done_criterion", Spec: specPath, Task: t.ID, Severity: "error", Message: "done criterion is generic: " + t.Done})
		}
		verify := strings.TrimSpace(t.Verify)
		if verify == "—" {
			verify = ""
		}
		if verify == "" {
			findings = append(findings, Finding{Rule: "tasks.empty_verify", Spec: specPath, Task: t.ID, Severity: "error", Message: "verify is empty"})
		} else if !verifyExecutable(verify) {
			findings = append(findings, Finding{Rule: "tasks.verify_not_executable", Spec: specPath, Task: t.ID, Severity: "error", Message: "verify must be a command, test ID or artifact path"})
		}
		var writes []string
		for _, w := range t.Writes {
			n, err := normalizeWritePath(w)
			if err != nil {
				findings = append(findings, Finding{Rule: "tasks.invalid_write_path", Spec: specPath, Task: t.ID, Severity: "error", Message: "write path: " + err.Error()})
				continue
			}
			writes = append(writes, n)
		}
		t.Writes = writes
		byID[t.ID] = t
	}

	if cycle := findCycle(byID); cycle != "" {
		findings = append(findings, Finding{Rule: "tasks.dependency_cycle", Spec: specPath, Task: cycle, Severity: "error", Message: "task dependency cycle"})
		return findings, nil
	}

	dedicated := map[string]bool{}
	coveredReqs := map[string]bool{}
	for _, t := range byID {
		var g []string
		for _, c := range t.Covers {
			if goals[c] {
				g = append(g, c)
			}
			if reqs[c] {
				coveredReqs[c] = true
			}
		}
		if len(g) == 1 {
			dedicated[g[0]] = true
		}
	}
	for g := range goals {
		if !dedicated[g] {
			findings = append(findings, Finding{Rule: "tasks.goal_not_dedicated", Spec: specPath, Task: "", Severity: "error", Message: "Goal has no dedicated task: " + g})
		}
	}
	for r := range reqs {
		if !coveredReqs[r] {
			findings = append(findings, Finding{Rule: "tasks.uncovered_requirement", Spec: specPath, Task: "", Severity: "warning", Message: "REQ has no task: " + r})
		}
	}

	reach := buildReach(byID)
	conflicts := map[string]bool{}
	caseWarned := map[string]bool{}
	for aID, a := range byID {
		for bID, b := range byID {
			if aID >= bID || len(a.Writes) == 0 || len(b.Writes) == 0 {
				continue
			}
			for _, wa := range a.Writes {
				for _, wb := range b.Writes {
					if strings.EqualFold(wa, wb) && wa != wb {
						key := aID + "\x00" + bID
						if !caseWarned[key] {
							caseWarned[key] = true
							findings = append(findings, Finding{Rule: "tasks.case_collision", Spec: specPath, Task: aID, Severity: "warning", Message: "write paths differ only by case: " + wa + " vs " + wb})
						}
					}
					if !overlap(wa, wb) {
						continue
					}
					if reach[aID][bID] || reach[bID][aID] {
						continue
					}
					key := aID + "\x00" + bID
					if conflicts[key] {
						continue
					}
					conflicts[key] = true
					findings = append(findings, Finding{Rule: "tasks.write_conflict", Spec: specPath, Task: aID, Severity: "error", Message: "independent tasks write the same path: " + wa + " vs " + wb + " (" + bID + ")"})
				}
			}
		}
	}

	shape := computePlanShape(byID, reqs)
	return findings, &shape
}

func findCycle(byID map[string]task) string {
	color := map[string]int{}
	var visit func(string) string
	visit = func(id string) string {
		color[id] = 1
		for _, d := range byID[id].Depends {
			if _, ok := byID[d]; !ok {
				continue
			}
			switch color[d] {
			case 1:
				return id
			case 0:
				if c := visit(d); c != "" {
					return c
				}
			}
		}
		color[id] = 2
		return ""
	}
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if color[id] == 0 {
			if c := visit(id); c != "" {
				return c
			}
		}
	}
	return ""
}

func buildReach(byID map[string]task) map[string]map[string]bool {
	out := map[string]map[string]bool{}
	for id := range byID {
		seen := map[string]bool{}
		var walk func(string)
		walk = func(n string) {
			for _, d := range byID[n].Depends {
				if _, ok := byID[d]; !ok || seen[d] {
					continue
				}
				seen[d] = true
				walk(d)
			}
		}
		walk(id)
		out[id] = seen
	}
	return out
}

func overlap(a, b string) bool {
	if a == b {
		return true
	}
	if strings.HasSuffix(a, "/") && strings.HasPrefix(b, a) {
		return true
	}
	if strings.HasSuffix(b, "/") && strings.HasPrefix(a, b) {
		return true
	}
	return false
}

func computePlanShape(byID map[string]task, reqs map[string]bool) PlanShape {
	crit := criticalPath(byID)
	groups := parallelGroups(byID)
	hot := map[string]map[string]bool{}
	for id, t := range byID {
		for _, w := range t.Writes {
			if hot[w] == nil {
				hot[w] = map[string]bool{}
			}
			hot[w][id] = true
		}
	}
	var hotPaths []HotPath
	for p, ids := range hot {
		hotPaths = append(hotPaths, HotPath{Path: p, Count: len(ids)})
	}
	sort.Slice(hotPaths, func(i, j int) bool {
		if hotPaths[i].Count != hotPaths[j].Count {
			return hotPaths[i].Count > hotPaths[j].Count
		}
		return hotPaths[i].Path < hotPaths[j].Path
	})
	var over []OverCovered
	for id, t := range byID {
		if len(t.Covers) > 5 {
			over = append(over, OverCovered{Task: id, Covers: len(t.Covers)})
		}
	}
	sort.Slice(over, func(i, j int) bool {
		if over[i].Covers != over[j].Covers {
			return over[i].Covers > over[j].Covers
		}
		return over[i].Task < over[j].Task
	})
	var uncovered []string
	for r := range reqs {
		found := false
		for _, t := range byID {
			for _, c := range t.Covers {
				if c == r {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			uncovered = append(uncovered, r)
		}
	}
	sort.Strings(uncovered)
	if crit == nil {
		crit = []string{}
	}
	if groups == nil {
		groups = [][]string{}
	}
	if hotPaths == nil {
		hotPaths = []HotPath{}
	}
	if over == nil {
		over = []OverCovered{}
	}
	if uncovered == nil {
		uncovered = []string{}
	}
	return PlanShape{
		CriticalPath: crit, ParallelGroups: groups, HotPaths: hotPaths,
		OverCoveredTasks: over, UncoveredRequirements: uncovered,
	}
}

func criticalPath(byID map[string]task) []string {
	type info struct {
		n    int
		path []string
	}
	memo := map[string]info{}
	var dfs func(string) info
	dfs = func(id string) info {
		if v, ok := memo[id]; ok {
			return v
		}
		best := info{n: 1, path: []string{id}}
		for _, d := range byID[id].Depends {
			if _, ok := byID[d]; !ok {
				continue
			}
			c := dfs(d)
			if c.n+1 > best.n || (c.n+1 == best.n && lessPath(append(c.path, id), best.path)) {
				best = info{n: c.n + 1, path: append(append([]string{}, c.path...), id)}
			}
		}
		memo[id] = best
		return best
	}
	var best info
	first := true
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		c := dfs(id)
		if first || c.n > best.n || (c.n == best.n && lessPath(c.path, best.path)) {
			best = c
			first = false
		}
	}
	return best.path
}

func parallelGroups(byID map[string]task) [][]string {
	depth := map[string]int{}
	var dfs func(string) int
	dfs = func(id string) int {
		if v, ok := depth[id]; ok {
			return v
		}
		d := 0
		for _, dep := range byID[id].Depends {
			if _, ok := byID[dep]; !ok {
				continue
			}
			if x := dfs(dep) + 1; x > d {
				d = x
			}
		}
		depth[id] = d
		return d
	}
	groups := map[int][]string{}
	for id := range byID {
		groups[dfs(id)] = append(groups[dfs(id)], id)
	}
	var keys []int
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	var out [][]string
	for _, k := range keys {
		sort.Strings(groups[k])
		out = append(out, groups[k])
	}
	return out
}

func lessPath(a, b []string) bool {
	return strings.Join(a, "\x00") < strings.Join(b, "\x00")
}

func writeTable(w io.Writer, rep Report) {
	fmt.Fprintf(w, "taskscheck: specs=%d tasks=%d findings=%d errors=%d warnings=%d\n",
		rep.Stats.SpecsChecked, rep.Stats.TasksChecked, len(rep.Findings), len(rep.Errors), rep.Stats.Warnings)
	for _, f := range rep.Findings {
		fmt.Fprintf(w, "[%s] %s %s %s: %s\n", f.Severity, f.Spec, f.Task, f.Rule, f.Message)
	}
	for _, e := range rep.Errors {
		fmt.Fprintf(w, "[error] %s %s %s: %s\n", e.Spec, e.Task, e.Rule, e.Message)
	}
	fmt.Fprintf(w, "plan: critical_path=%v parallel_groups=%v\n", rep.PlanShape.CriticalPath, rep.PlanShape.ParallelGroups)
	for _, h := range rep.PlanShape.HotPaths {
		fmt.Fprintf(w, "hot: %s (%d)\n", h.Path, h.Count)
	}
	for _, o := range rep.PlanShape.OverCoveredTasks {
		fmt.Fprintf(w, "over-covered: %s (%d)\n", o.Task, o.Covers)
	}
	for _, r := range rep.PlanShape.UncoveredRequirements {
		fmt.Fprintf(w, "uncovered-req: %s\n", r)
	}
}

func writeGitHub(w io.Writer, rep Report) {
	for _, f := range rep.Findings {
		level := "error"
		if f.Severity == "warning" {
			level = "warning"
		}
		fmt.Fprintf(w, "::%s file=%s::[%s] %s\n", level, f.Spec, f.Rule, f.Message)
	}
	for _, e := range rep.Errors {
		fmt.Fprintf(w, "::error file=%s::[%s] %s\n", e.Spec, e.Rule, e.Message)
	}
	fmt.Fprintf(w, "::notice::taskscheck specs=%d tasks=%d findings=%d errors=%d\n",
		rep.Stats.SpecsChecked, rep.Stats.TasksChecked, len(rep.Findings), len(rep.Errors))
}
