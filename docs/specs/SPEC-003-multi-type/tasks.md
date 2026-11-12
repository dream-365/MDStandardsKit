---
type: tasks
status: Done
owner: example
depends_on: []
spec: ./spec.md
spec_sha256: "53f186ffa57eb09df3ebf0d47c8b6c95e941d850ecbcd39d6165ffc588dea4f0"
updated_at: 2026-11-12
---

# MDStandardsKit 多文档类型路由 执行拆解

| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |
|---|---|---|---|---|---|---|
| T1 | 更新架构、README 与 idea-writer SKILL，记录多类型模式与迁移语义 | G6, REQ-013 | — | docs/architecture/mdstandards-architecture.md, README.md, skills/idea-writer/SKILL.md | 三份文档与已批准 Spec 一致，包含配置示例、rule_set 矩阵、迁移与信任边界 | `./mdstandards -root . -policy ./policy -phase 1 -format table` |
| T2 | 将标准升级到 `standard_version: 2`，补 `rule_set`/`graph`/新类型档案 | G1, G3, REQ-003, REQ-011 | T1 | policy/data/standard.json | JSON 可解析；spec/architecture/adr/changelog/tasks/index/idea 档案齐全；idea 结构字段完整 | `TestBuildRunPlan*` |
| T3 | 实现配置双模式与运行计划：`type`/`types`、null/trim、启用类型校验 | G1, REQ-001, REQ-002, REQ-003, REQ-011 | T2 | main.go | 缺失/legacy/multi/`{}`/null/双键/空/重复/未知 profile 各走指定退出码 | `TestLoadRepoConfig*` |
| T4 | 实现 `ExtractMulti`、frontmatter 判型、type/path 兼容、canonical type 与 HEAD 配置保护 | G2, REQ-004, REQ-005, REQ-009 | T3 | extract.go, main.go | 缺 type/未知/未启用/非字符串/路径不兼容均产生 routing error；HEAD 配置 dirty 时降级 | `TestMultiRouting*` |
| T5 | 新增 `policy/types.rego` 并改造 refs/status/idea 为逐文档作用域 | G3, REQ-006, REQ-010 | T4 | policy/types.rego, policy/refs.rego, policy/status.rego, policy/idea.rego | legacy spec/idea 行为不变；multi 按 `rule_set` 分流；缺 `rule_set` 产出 `policy.load_error` | `TestMultiRuleScope*` |
| T6 | 收口覆盖断言、边统计与零文档 warning | G4, REQ-007, REQ-008, REQ-012 | T4 | main.go | 顶层 globs 始终参与发现；零边按 graph 文档统计；零文档 warning 按类型名排序 | `TestMultiCoverage*` |
| T7 | 新增 `testdata/mixed-repo` 混合语料 fixture | G1, G2, G3, G4, REQ-004, REQ-006, REQ-008 | T2 | testdata/mixed-repo/.mdstandards.json, testdata/mixed-repo/docs/, testdata/mixed-repo/2026/ | fixture phase 1 exit 0；覆盖 spec/architecture/changelog/tasks/index/idea 六种类型 | `./mdstandards -root ./testdata/mixed-repo -policy ./policy -phase 1 -format table` |
| T8 | 增加配置、路由、Rego 与兼容性测试矩阵 | G1, G2, G3, G4, G5, REQ-001, REQ-002, REQ-003, REQ-004, REQ-005, REQ-006, REQ-007, REQ-008, REQ-009, REQ-010, REQ-011, REQ-012 | T3, T4, T5, T6, T7 | plan_test.go, multi_test.go, multi_rego_test.go | 新增测试全绿；既有 golden 与 idea 测试不变 | `go test ./...` |
| T9 | 接入 CI、dogfood、CODEOWNERS 与发布 smoke | G6, REQ-013 | T8 | Makefile, .github/workflows/ci.yml, .github/workflows/release.yml, .github/CODEOWNERS, .mdstandards.json, .claude-plugin/plugin.json | `make check` 跑 mixed smoke；根仓库 dogfood 通过；release smoke 覆盖 mixed-repo | `make check` |
| T10 | 收口状态、任务哈希与质量产物 | G5, REQ-010, REQ-013 | T9 | docs/specs/SPEC-003-multi-type/spec.md, docs/specs/SPEC-003-multi-type/changelog.md, docs/specs/SPEC-003-multi-type/tasks.md, docs/quality/SPEC-003-multi-type/quality-findings.json, docs/quality/SPEC-003-multi-type/quality-report.md | Spec 转 Implemented；tasks 的 `spec_sha256` 与质量产物 target SHA 与最终 spec 一致 | `make check` |
