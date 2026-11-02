---
type: tasks
status: Done
owner: alang
depends_on: []
spec: ./spec.md
spec_sha256: "bb8d4081ca8b4738feeeea3a8c22353dabf47fdeb466803ef4593b8b4d0f1d56"
updated_at: 2026-11-03
---

# SPEC-001 MDStandardsKit 元数据闸门 执行拆解

| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |
|---|---|---|---|---|---|---|
| T1 | depends_on 归一化去重与引用图回归 | G1, REQ-001, REQ-005, REQ-006, REQ-010 | — | extract_test.go | 重复目标只保留首次出现，归一化后边计数等于唯一目标数；缺失目标仍产出 refs.missing_target | go test ./... -run 'TestDependsOnDedupedPreservesFirstOrder' |
| T2 | status 缺失与标准守卫 | G2, REQ-011, REQ-012, REQ-013 | T1 | status_test.go | 缺失 status 与大小写不符均产出 refs.status_not_dependable；标准缺 dependable_statuses 时 Rego 产出 policy.load_error | go test ./... -run 'TestMissingStatusIsNotDependable\|TestStatusRegoGuardEmitsPolicyLoadError' |
| T3 | 覆盖断言 fail-closed 回归 | G3, REQ-020, REQ-021, REQ-022 | T2 | coverage_test.go | 零受管文档在 phase 0/1 均退出 2；phase ≥ 1 的零解析文档与零依赖边均退出 2 | go test ./... -run 'TestZeroParsedFailsClosedInPhase1\|TestZeroGovernedDocumentsFailsClosedPhase0' |
| T4 | 解析错误与 1 MiB/别名上限 | G3, REQ-002, REQ-004, REQ-007 | T1, T3 | extract_test.go | 超过 1 MiB、`*` 超过 100、`<<` 超过 10、非法 YAML 与非映射 frontmatter 均产出 extract_error，文件不被静默丢弃 | go test ./... -run 'TestFileSizeLimitProducesExtractError\|TestYAMLAliasGuardProducesExtractError\|TestInvalidYAMLProducesExtractError\|TestNonMappingFrontmatterProducesExtractError' |
| T5 | 策略降级、HEAD 内容与 phase/baseline 回归 | G4, REQ-030, REQ-031, REQ-032, REQ-033, REQ-034, REQ-050, REQ-051, REQ-052 | T4 | policyhead_test.go | HEAD 物化读取已提交策略内容；脏策略或 `-root` 外 baseline 降级只报告；脏策略含 extract_error 退出 0；策略目录在 root 外退出 2 | go test ./... -run 'Test.*Policy' |
| T6 | NFR-002 性能预算降级 | G4, REQ-031 | T5 | main.go, perf_test.go | ≤200 篇受管匹配文档单次门禁超过 2 秒时打印警告并在 phase 1 降级为只警告；phase 2 不降级；覆盖断言仍退出 2 | go test ./... -run 'TestPerformanceBudget' |
| T7 | 稳定输出、GitHub 注解与 -out | G5, REQ-041, REQ-042, REQ-043 | T6 | report_test.go | findings 按 (rule, document, target) 稳定排序；github 每条 finding/extract_error 一行注解并输出 notice；-out 写文件且 stdout 为空 | go test ./... -run 'TestSortFindingsOrder\|TestGitHubFormat\|TestOutWritesReportToFile' |
| T8 | JSON 契约与字节稳定 | G5, REQ-003, REQ-023, REQ-040 | T7 | json_test.go | stats 含三个计数且数值正确；同一输入连续运行的 JSON 字节一致 | go test ./... -run 'TestJSON' |
| T9 | 策略求值超时与错误传播 | G3 | T6 | main.go, eval_test.go | 取消或已过期的求值上下文返回错误，run 以 policy.load_error 退出 2 | go test ./... -run 'TestEval' |
| T10 | 1000 篇性能实测 | G1, REQ-001 | T9 | perf_test.go | 1000 篇受管文档单次全量检查低于 10 秒，测试日志记录实际耗时 | sh -c 'MDSTANDARDS_PERF=1 go test ./... -run TestNFR001ThousandDocsBudget -count=1' |
| T11 | 验收表、状态与索引收口 | G5, REQ-040 | T10 | spec.md, changelog.md, tasks.md, docs/specs/README.md, README.md, docs/architecture/mdstandards-architecture.md, docs/quality/SPEC-001-mdstandards/quality-findings.smoke.json, docs/quality/SPEC-001-mdstandards/quality-report.smoke.md | 验收表与新增测试一致且未覆盖项保留诚实标注；spec 状态为 Implemented、tasks 状态为 Done；索引同步；smoke 的 target_sha256 与当前 spec 一致；taskscheck 无 stale_plan | make check |
| T12 | 独立质量复审与审计 | G5 | T11 | spec.md, changelog.md, tasks.md, docs/quality/SPEC-001-mdstandards/quality-findings.json, docs/quality/SPEC-001-mdstandards/quality-report.md, docs/quality/README.md, docs/quality/SPEC-001-mdstandards/quality-findings.smoke.json, docs/quality/SPEC-001-mdstandards/quality-report.smoke.md | 若复审要求改 spec，重算 spec_sha256 并重跑 taskscheck；最终 findings target_sha256 等于当前 spec；qualitycheck 审计无 failure | docs/quality/SPEC-001-mdstandards/quality-report.md |
