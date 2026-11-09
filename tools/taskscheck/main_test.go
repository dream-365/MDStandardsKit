package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validTable = `| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |
|---|---|---|---|---|---|---|
| T1 | 任务一 | G1, REQ-001 | — | a.go | 输出包含 alpha | go test ./... -run TestAlpha |
| T2 | 任务二 | G2, REQ-002 | T1 | b.go | 输出包含 beta | go test ./... -run TestBeta |
`

func writeSpec(t *testing.T, path, status string, withGoals bool) {
	t.Helper()
	body := "---\ntype: spec\nstatus: " + status + "\nowner: example\ndepends_on: []\nupdated_at: 2026-11-02\n---\n\n# Test\n\n## 2. 目标与非目标\n\n### 2.1 目标\n\n"
	if withGoals {
		body += "- **G1 目标一**：验收一。\n- **G2 目标二**：验收二。\n"
	}
	body += "\n### 2.2 非目标\n\n- **NG1** 不做。\n\n## 4. 详细设计\n\n- **REQ-001** 需求一。\n- **REQ-002** 需求二。\n\n## 5. 非功能性需求\n"
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeTasksWith(t *testing.T, tasksPath, specPath, frontmatter, body string) {
	t.Helper()
	sum, err := fileHash(specPath)
	if err != nil {
		t.Fatal(err)
	}
	if frontmatter == "" {
		frontmatter = "type: tasks\nstatus: In Progress\nowner: example\ndepends_on: []\nspec: ./spec.md\nspec_sha256: {{HASH}}\nupdated_at: 2026-11-02"
	}
	frontmatter = strings.ReplaceAll(frontmatter, "{{HASH}}", sum)
	content := "---\n" + frontmatter + "\n---\n\n" + body
	if err := os.WriteFile(tasksPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) (root, specPath, tasksPath string) {
	t.Helper()
	root = t.TempDir()
	dir := filepath.Join(root, "docs", "specs", "SPEC-001")
	specPath = filepath.Join(dir, "spec.md")
	writeSpec(t, specPath, "Approved", true)
	tasksPath = filepath.Join(dir, "tasks.md")
	return
}

func runCLI(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	code := run(args, &out, &errb)
	return code, out.String(), errb.String()
}

func assertRule(t *testing.T, out, rule string) {
	t.Helper()
	if !strings.Contains(out, rule) {
		t.Fatalf("expected rule %q in output:\n%s", rule, out)
	}
}

func TestCleanRepoPasses(t *testing.T) {
	root, spec, tasks := setup(t)
	writeTasksWith(t, tasks, spec, "", validTable)
	code, out, errb := runCLI(t, "-root", root, "-phase", "1")
	if code != 0 {
		t.Fatalf("exit=%d out=%s err=%s", code, out, errb)
	}
}

func TestDraftWithoutTasksIsOptional(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "docs", "specs", "SPEC-001")
	writeSpec(t, filepath.Join(dir, "spec.md"), "Draft", true)
	code, out, errb := runCLI(t, "-root", root, "-phase", "1")
	if code != 0 {
		t.Fatalf("exit=%d out=%s err=%s", code, out, errb)
	}
}

func TestApprovedWithoutTasksBlocks(t *testing.T) {
	root, _, _ := setup(t)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.missing_file")
}

func TestDuplicateID(t *testing.T) {
	root, spec, tasks := setup(t)
	body := `| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |
|---|---|---|---|---|---|---|
| T1 | 一 | G1, REQ-001 | — | a.go | alpha 输出 | go test ./... -run TestA |
| T1 | 二 | G2, REQ-002 | — | b.go | beta 输出 | go test ./... -run TestB |
`
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.duplicate_id")
}

func TestUnknownDependency(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "| T1 | b.go", "| T1 | b.go", 1)
	body = strings.Replace(body, "| T2 | 任务二 | G2, REQ-002 | T1", "| T2 | 任务二 | G2, REQ-002 | T9", 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.unknown_dependency")
}

func TestSelfDependency(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "| T2 | 任务二 | G2, REQ-002 | T1", "| T2 | 任务二 | G2, REQ-002 | T2", 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.self_dependency")
}

func TestCycle(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "| T2 | 任务二 | G2, REQ-002 | T1", "| T2 | 任务二 | G2, REQ-002 | T1", 1)
	body = strings.Replace(body, "| T1 | 任务一 | G1, REQ-001 | —", "| T1 | 任务一 | G1, REQ-001 | T2", 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.dependency_cycle")
}

func TestUnknownCoverage(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "G1, REQ-001", "G9, REQ-001", 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.unknown_coverage")
}

func TestEmptyCoverage(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "G1, REQ-001", "—", 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.empty_coverage")
}

func TestGoalNotDedicated(t *testing.T) {
	root, spec, tasks := setup(t)
	body := `| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |
|---|---|---|---|---|---|---|
| T1 | 一 | G1, G2, REQ-001 | — | a.go | alpha 输出 | go test ./... -run TestA |
| T2 | 二 | G1, G2, REQ-002 | — | b.go | beta 输出 | go test ./... -run TestB |
`
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.goal_not_dedicated")
}

func TestBannedDone(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "输出包含 alpha", "完成", 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.banned_done_criterion")
}

func TestEmptyVerify(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "go test ./... -run TestAlpha", "—", 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.empty_verify")
}

func TestVerifyNotExecutable(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "go test ./... -run TestAlpha", "手动看看", 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.verify_not_executable")
}

func TestWriteConflict(t *testing.T) {
	root, spec, tasks := setup(t)
	body := `| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |
|---|---|---|---|---|---|---|
| T1 | 一 | G1, REQ-001 | — | a.go | alpha 输出 | go test ./... -run TestA |
| T2 | 二 | G2, REQ-002 | — | a.go | beta 输出 | go test ./... -run TestB |
`
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.write_conflict")
}

func TestCaseCollisionIsWarning(t *testing.T) {
	root, spec, tasks := setup(t)
	body := `| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |
|---|---|---|---|---|---|---|
| T1 | 一 | G1, REQ-001 | — | a.go | alpha 输出 | go test ./... -run TestA |
| T2 | 二 | G2, REQ-002 | — | A.go | beta 输出 | go test ./... -run TestB |
`
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 0 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.case_collision")
}

func TestStaleHash(t *testing.T) {
	root, spec, tasks := setup(t)
	fm := "type: tasks\nstatus: In Progress\nowner: example\ndepends_on: []\nspec: ./spec.md\nspec_sha256: " + strings.Repeat("a", 64) + "\nupdated_at: 2026-11-02"
	writeTasksWith(t, tasks, spec, fm, validTable)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.stale_plan")
}

func TestWrongSpecPath(t *testing.T) {
	root, spec, tasks := setup(t)
	other := filepath.Join(filepath.Dir(spec), "other.md")
	if err := os.WriteFile(other, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	fm := "type: tasks\nstatus: In Progress\nowner: example\ndepends_on: []\nspec: ./other.md\nspec_sha256: {{HASH}}\nupdated_at: 2026-11-02"
	writeTasksWith(t, tasks, spec, fm, validTable)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.wrong_spec_path")
}

func TestMissingSpec(t *testing.T) {
	root, spec, tasks := setup(t)
	fm := "type: tasks\nstatus: In Progress\nowner: example\ndepends_on: []\nspec: ./nope.md\nspec_sha256: {{HASH}}\nupdated_at: 2026-11-02"
	writeTasksWith(t, tasks, spec, fm, validTable)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.missing_spec")
	if strings.Contains(out, "tasks.wrong_spec_path") {
		t.Fatalf("missing_spec must not also emit wrong_spec_path:\n%s", out)
	}
}

func TestNotApplicableSkipsHash(t *testing.T) {
	root, spec, tasks := setup(t)
	fm := "type: tasks\nstatus: Draft\nowner: example\ndepends_on: []\nspec: ./spec.md\nspec_sha256: " + strings.Repeat("a", 64) + "\nupdated_at: 2026-11-02\nnot_applicable: 单批可完成"
	writeTasksWith(t, tasks, spec, fm, "")
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 0 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
}

func TestNotApplicableWithTable(t *testing.T) {
	root, spec, tasks := setup(t)
	fm := "type: tasks\nstatus: Draft\nowner: example\ndepends_on: []\nspec: ./spec.md\nspec_sha256: {{HASH}}\nupdated_at: 2026-11-02\nnot_applicable: 单批可完成"
	writeTasksWith(t, tasks, spec, fm, validTable)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.invalid_not_applicable")
}

func TestNoGoals(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "docs", "specs", "SPEC-001")
	spec := filepath.Join(dir, "spec.md")
	writeSpec(t, spec, "Approved", false)
	tasks := filepath.Join(dir, "tasks.md")
	writeTasksWith(t, tasks, spec, "", validTable)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 2 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "spec.no_goals")
}

func TestNoSpecs(t *testing.T) {
	root := t.TempDir()
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 2 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.no_specs")
}

func TestInvalidSpecFlag(t *testing.T) {
	root := t.TempDir()
	code, out, _ := runCLI(t, "-root", root, "-spec", filepath.Join(root, "nope", "spec.md"), "-phase", "1")
	if code != 2 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.invalid_spec_flag")
}

func TestPhase0ReportsWithoutBlocking(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "| T2 | 任务二 | G2, REQ-002 | T1", "| T2 | 任务二 | G2, REQ-002 | T9", 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "0")
	if code != 0 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.unknown_dependency")
}

func TestFencedExampleIgnored(t *testing.T) {
	root, spec, tasks := setup(t)
	body := "```\n" + validTable + "```\n\n" + validTable
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 0 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
}

func TestEscapedPipe(t *testing.T) {
	root, spec, tasks := setup(t)
	body := strings.Replace(validTable, "go test ./... -run TestAlpha", `go test ./... -run 'TestAlpha\|TestBeta'`, 1)
	writeTasksWith(t, tasks, spec, "", body)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 0 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
}

func TestPrintSpecHash(t *testing.T) {
	root, spec, _ := setup(t)
	want, err := fileHash(spec)
	if err != nil {
		t.Fatal(err)
	}
	code, out, _ := runCLI(t, "-root", root, "-print-spec-hash", spec)
	if code != 0 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	if strings.TrimSpace(out) != want {
		t.Fatalf("got %q want %q", strings.TrimSpace(out), want)
	}
}

func TestOrphanTasks(t *testing.T) {
	root, spec, tasks := setup(t)
	writeTasksWith(t, tasks, spec, "", validTable)
	orphan := filepath.Join(root, "docs", "specs", "ORPHAN", "tasks.md")
	if err := os.MkdirAll(filepath.Dir(orphan), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(orphan, []byte("---\ntype: tasks\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 1 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	assertRule(t, out, "tasks.orphan_tasks")
}

func TestJSONShape(t *testing.T) {
	root, spec, tasks := setup(t)
	writeTasksWith(t, tasks, spec, "", validTable)
	code, out, errb := runCLI(t, "-root", root, "-phase", "1", "-format", "json")
	if code != 0 {
		t.Fatalf("exit=%d out=%s err=%s", code, out, errb)
	}
	var rep Report
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatal(err)
	}
	if rep.SchemaVersion != "1" {
		t.Fatalf("schema_version=%q", rep.SchemaVersion)
	}
	if rep.Findings == nil || rep.Errors == nil {
		t.Fatalf("findings/errors must be arrays, got %#v %#v", rep.Findings, rep.Errors)
	}
	if len(rep.PlanShape.CriticalPath) == 0 {
		t.Fatalf("plan_shape.critical_path is empty")
	}
}

func TestPlanShape(t *testing.T) {
	root, spec, tasks := setup(t)
	writeTasksWith(t, tasks, spec, "", validTable)
	code, out, _ := runCLI(t, "-root", root, "-phase", "1")
	if code != 0 {
		t.Fatalf("exit=%d out=%s", code, out)
	}
	if !strings.Contains(out, "critical_path=[T1 T2]") {
		t.Fatalf("unexpected plan shape:\n%s", out)
	}
}
