---
type: changelog
status: Implemented
owner: example
depends_on: []
updated_at: 2026-11-12
---

# Changelog: SPEC-003 多文档类型路由

## [2026-11-12] - 实现与收口

- **修改类型**: Changed
- **影响范围**: 状态、`tasks.md`、`docs/quality/SPEC-003-multi-type/`、根仓库 dogfood、CI 与发布 smoke。
- **变更原因**: Phase B 实现完成；需要把 Spec 生命周期、任务哈希与质量产物收口到同一版本。
- **具体内容**:
  - 实现 `.mdstandards.json` 双模式、`buildRunPlan`、`ExtractMulti`、type/path 路由、`policy/types.rego` 作用域与 `--policy-from-head` 配置保护。
  - `standard_version` 升到 2；根仓库新增 `.mdstandards.json` dogfood；CI/release 增加 mixed-repo smoke。
  - 新增 `testdata/mixed-repo` 与配置/路由/Rego/兼容性测试矩阵；legacy golden 与 idea 行为不变。
  - `tasks.md` 标记 `Done`，`spec_sha256` 与最终 `spec.md` 一致；质量 findings/report 按最终 SHA 重新生成。
- **状态变更**: `In Review` → `Approved`（Gate 1）→ `Implemented`；`tasks.md` 为 `Done`。

## [2026-11-12] - Gate 1 复核收口

- **修改类型**: Changed
- **影响范围**: 4.1 术语、4.2 配置契约、4.3 数据模型、4.5 核心流程、4.6 异常与退出码、2.2 非目标、G6、6.2 迁移、质量产物。
- **变更原因**: 第二轮需求/规范复核要求补齐 legacy Rego 作用域的推导依据，并修正降级、`rule_set` 兼容、`-dump-corpus` 覆盖断言、NG8 边界与质量产物哈希。
- **具体内容**:
  - `§4.5` 增加 legacy 作用域段：legacy spec 的 `is_spec_scope`/空 `profile_for`/顶层 `dependable_for`，legacy idea 的 `profile_for`，以及 multi 缺 `rule_set` 的 `policy.load_error` 守卫；明确这是 G5 的推导依据。
  - `§4.1` 区分性能降级（仅 phase 1）与 `--policy-from-head` 配置降级（phase ≥ 1，含 phase 2）；`§4.6` 规则 finding 行补“未降级”。
  - `§4.3` 明确 legacy v1 档案允许 `rule_set` 缺省，multi 启用类型必须存在且非空；补 `graph` 缺省 false、spec 需显式 true 的说明。
  - `§4.6` 覆盖断言限定为门禁运行；`-dump-corpus` 直接输出并 exit 0，不执行零文档/零解析/零边断言（配置/策略加载错误仍 exit 2）。
  - NG8 收窄为“状态流转校验”，审批来源归 NG3；G6 明确架构/README/standard 是 Phase B 交付物，状态改为 `Approved` 时创建 `tasks.md` 并写入最终 `spec_sha256`。
  - `§4.2` 明确 legacy `type` 精确匹配、不 trim；`§6.2` 明确用 `types` 替换而不是追加 legacy `type`。
  - `quality-findings.json` 与 `quality-report.md` 按当前 `spec.md` 的 SHA-256 重新生成；audit 通过，headline 100.0 / diagnostic 98.5。
- **状态变更**: 保持 `In Review`，等待 Gate 1 批准。

## [2026-11-12] - Gate 1 评审修订

- **修改类型**: Changed
- **影响范围**: 2.1 Goals、4.2 配置契约、4.3 数据模型、4.5 核心流程、4.6 退出码与降级、4.7 信任边界、4.8 REQ、5 NFR、6.2 迁移、6.3 冒烟、6.4 回滚。
- **变更原因**: 工程评审、需求/规范评审与质量评审提出阻塞项：配置 null/trim 语义不清、回滚步骤会导致 `{}`、降级语义未定义、零边断言聚合与 per-type 表述冲突、`standard.json` 示例与校验规则冲突、Goals 数量超模板上限、章节标题与 rubric 锚点不一致。
- **具体内容**:
  - Goals 从 7 个合并为 6 个；新增 NG8 状态流转校验（审批来源归 NG3），统一 1.4 与 2.2 的标签。
  - 配置契约明确“`type`/`types` 恰好一个非 null；值为 null 视为未提供”；补 `{"type":""}` 与 trim 匹配语义。
  - `standard.json` 示例补全 `adr` 与 idea 结构字段；`date_field` 改为可选、缺省 `date`；补顶层 globs/statuses 非空校验与 version 拒绝矩阵。
  - 定义两类降级：性能降级仅影响 phase 1；`--policy-from-head` 配置降级影响 phase ≥ 1（含 phase 2）。routing error 与解析错误同语义；配置/覆盖错误降级时仍 exit 2。
  - 零边断言改为聚合语义：仅当语料中至少有一篇 `graph: true` 文档时才要求全语料非零边；单类型零文档只 warning。
  - 回滚改为“移除配置文件或改写为 legacy `{"type":...}`”，不再写“删除 types 键”。
  - 信任边界限定顶层 globs 覆盖范围；补 `standard.json`/`.mdstandards.json` 的 load-bearing 说明与残留风险。
  - 迁移 recipe 增加 policy bundle `standard_version: 2` 步骤；冒烟增加 routing error 数与耗时指标。
  - 章节标题对齐 rubric 锚点（3.4、4.3、4.4、4.5、4.6、6.2）；补 NFR 验收方式与数字来源。
- **状态变更**: 保持 `In Review`，等待下一轮 Gate 1 复核。

## [2026-11-12] - 初稿

- **修改类型**: Added
- **影响范围**: 新增 `docs/specs/SPEC-003-multi-type/spec.md`；后续实施涉及 `.mdstandards.json`、`policy/data/standard.json`、`main.go`、`extract.go`、`policy/types.rego`、`policy/refs.rego`、`policy/status.rego`、`policy/idea.rego`、`testdata/mixed-repo/`、README 与架构文档。
- **变更原因**: 仓库同时存在 spec、architecture、changelog、tasks、index、idea 等多种文档；现有仓库级单类型配置一次只能启用一套规则，frontmatter `type:` 只是被记录，不参与路由。
- **具体内容**:
  - 新增 opt-in 多类型模式 `{"types":[...]}`，与 legacy `{"type":"..."}` 互斥；文件缺失仍是 legacy spec。
  - multi 模式按 frontmatter `type` 逐文件判型；缺失、未知、未启用、与路径不兼容都产生 routing error。
  - 每个类型用 `rule_set` 选择 `spec` / `idea` / `none`；顶层 `governed_globs` 始终参与发现，防止缩小 `types` 后文件静默消失。
  - `standard_version` 升到 2；新二进制接受 1（legacy）与 2（multi 要求 2）。
  - legacy 单类型模式的错误码、rule ID 与 `-dump-corpus` 字节输出保持不变。
  - 新增 `Corpus.Mode`（仅 multi 输出）、`testdata/mixed-repo/` 与多类型测试矩阵。
- **状态变更**: 新建，`In Review`，等待独立评审；批准前不修改实现代码。
