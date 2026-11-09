---
type: changelog
status: Implemented
owner: example
depends_on: []
updated_at: 2026-11-06
---

# Changelog: SPEC-002 idea 类型

## [2026-11-06] - 可选「外部的声音」章节

- **修改类型**: Added
- **影响范围**: 1.1、2.1 G6、2.2 NG6、4.2、4.8、`policy/data/standard.json`、`skills/idea-writer/`、`testdata/idea-repo/`
- **变更原因**: 想法库已经在多篇想法里手工记录外部反馈、同类尝试和反例，但模板与结构闸门没有对应位置；需要把它固化为可选章节，同时保持「记想法不求完整」的原则。
- **具体内容**:
  - 模板在「为什么有意思」与「下一步」之间增加可选二级标题「外部的声音」，用于外部反馈、同类尝试、反例或行业标准。
  - `standard.json` 的 `placeholder_markers` 增加短哨兵 `optional:outside-voice`；模板注释里包含该哨兵。整节缺失不产生 finding；哨兵残留仍命中 `idea.placeholder_remaining`，改文案不会绕过检查。
  - 必需标题仍为四个；rubric 与 headline 不变，「外部的声音」只作为评审上下文。
  - 同步更新 `~/Workspace/ideas` 的模板、CLAUDE.md 与 README。
  - `docs/quality/SPEC-002-idea-type/` 下的报告与 findings 由 T10 重新独立评审生成，覆盖 2026-11-06 版。
- **状态变更**: 按 Spec 标准先回退到 `In Review`；所有者确认发布后重新批准，回到 `Implemented`，`tasks.md` 回到 `Done`。

## [2026-11-03] - 结构补全与实现收口

- **修改类型**: Changed
- **影响范围**: 3.4、第 4 章、第 6 章、G3、NFR-2、`tasks.md`、独立质量评审
- **变更原因**: 章节结构与 rubric 期望不一致；实现已随 v0.2.0 落地，状态与任务需要收口。
- **具体内容**:
  - 新增 `3.4 关键选型与约束`；第 4 章重排为术语、数据模型、结构规则、接口、核心流程、错误码、日期路径、内容评分；第 6 章拆为发布、迁移、验证、回滚。
  - G3 从 5 篇更新为 2026-11-03 实测的 9 篇真实语料；NFR-2 记录 0.01s 实测。
  - `tasks.md` 的 T1–T7 完成，新增 T8 收口与质量复审；新增独立质量评审产物。
- **状态变更**: `spec.md` 与 `changelog.md` 的 `status` 由 `In Review` 改为 `Implemented`；`tasks.md` 为 `Done`。

## [2026-11-02] - 初稿

- **修改类型**: Added
- **影响范围**: 新增 `docs/specs/SPEC-002-idea-type/spec.md`；`policy/idea.rego`、`extract.go`、`main.go`、`tools/qualitycheck`、`skills/idea-writer/`、`.claude-plugin/`。
- **变更原因**: 想法库有模板但没有结构验证和内容评估；需要给 mdstandards 增加 `idea` 类型。
- **具体内容**:
  - 仓库级 `.mdstandards.json` 声明 `{"type":"idea"}`；globs 只在 policy 的 `document_types.idea` 里。
  - 结构规则覆盖 ID、日期、路径、状态、tags、标题、占位符、必需章节。
  - `date` 同时兼容 YAML 字符串和 `time.Time`；非零时钟的 datetime 报错。
  - 内容评分复用 `qualitycheck`，增加 `guess` 字段和 `-title` / `-target-label`。
  - 结构错误阻断；内容分只诊断，不阻止入库。
- **状态变更**: 新建，`In Review`，等待独立评审。
