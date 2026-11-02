---
type: tasks
status: Done
owner: alang
depends_on: []
spec: ./spec.md
spec_sha256: "52181d85687c1b9cc525c35fd0ca1799b695ec284f158af541e4146e3865fcdf"
updated_at: 2026-11-03
---

# SPEC-002 idea 类型 执行拆解

| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |
|---|---|---|---|---|---|---|
| T1 | 类型声明与配置 fail-closed | G1 | — | main.go, idea_test.go | 合法 ideas 仓库 phase 0 退出 0；损坏配置、未知 type、未知键均退出 2 | `go test ./... -run 'TestIdeaConfigFailsClosed\|TestIdeaValidRepo'` |
| T2 | idea 结构规则 | G2 | T1 | extract.go, policy/idea.rego, idea_test.go | ID、日期、目录、状态、tags、标题、占位符、路径形状各有一条规则，并在 fixture 上产出对应 rule | `go test ./... -run 'TestIdeaStructuralFindings\|TestIdeaDuplicateID'` |
| T3 | 记录不误伤 | G3 | T2 | extract.go, policy/idea.rego, testdata/idea-repo/, idea_test.go | 空正文、`暂无`、`目前说不清，因为……` 不产生 finding；真实 ideas 语料（2026-11-03 实测 9 篇）phase 1 退出 0 | `go test ./... -run TestIdeaValidRepo` 与 `./mdstandards -root ~/Workspace/ideas -policy ./policy -phase 1` |
| T4 | 内容评分与独立评审 | G4 | T2 | skills/idea-writer/, tools/qualitycheck/ | rubric 校验通过；reviewer 产出 guess 与逐字证据；qualitycheck 复算分数并拒绝同一 writer 会话 | `go test ./tools/qualitycheck/... -run 'TestValidateRealRubric\|TestAuditRejectsSameSession'` |
| T5 | 五篇校准与真人标签 | G4, G3 | T3, T4 | skills/idea-writer/idea/quality/calibration/ | 5 篇人工分诊标签已记录并留出 1 篇；校准引擎测试通过 | `go test ./tools/qualitycheck/... -run TestCalibrate` 与 `test -s skills/idea-writer/idea/quality/calibration/ideas.v1.yaml` |
| T6 | spec 输出不回归 | G5 | T2, T3, T4 | testdata/spec-corpus.golden.json, idea_test.go | spec 的 `-dump-corpus` 与 golden 逐字节一致 | `go test ./... -run TestSpecDumpCorpusGolden` |
| T7 | v0.2.0 发布与冒烟 | G5 | T5, T6 | .claude-plugin/plugin.json, .github/workflows/release.yml, README.md, docs/architecture/mdstandards-architecture.md | tag 与 plugin.json 版本一致；四平台资产与 checksums 齐全；冒烟安装通过；make check 绿 | `make check` 与 `sh ./scripts/shim-test.sh` |
| T8 | 状态收口与独立质量复审 | G5 | T7 | spec.md, changelog.md, tasks.md, docs/specs/README.md, docs/quality/SPEC-002-idea-type/quality-findings.json, docs/quality/SPEC-002-idea-type/quality-report.md | spec 状态为 Implemented、tasks 状态为 Done、索引同步；qualitycheck 审计无 failure 且 target_sha256 与当前 spec 一致 | docs/quality/SPEC-002-idea-type/quality-report.md |
