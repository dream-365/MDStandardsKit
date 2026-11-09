---
type: spec
status: Draft
owner: example
depends_on:
  - docs/architecture/mdstandards-architecture.md
---

# MDStandardsKit v1 规格

## 1. 目的

规定 `mdstandards` v1 的可观察行为：输入、输出、退出码、规则语义、失败处理与验收标准。
实现细节属于架构文档；本文只写可以被测试断言的东西。

## 2. 术语

| 术语 | 含义 |
|---|---|
| 受管文档 | 路径匹配 `standard.governed_globs` 的 Markdown 文件 |
| 语料 | 由受管文档抽取出的 `{documents, extract_errors, stats}` |
| 依赖边 | 某文档 `depends_on` 中的一项，指向另一篇受管文档 |
| 可依赖状态 | `standard.dependable_statuses` 中的字面量，逐字匹配 |

## 3. 前置条件

- PREC-001：`policy/data/standard.json` 存在且可解析。
- PREC-002：`standard_version` 为 `"1"`；未知版本必须拒绝。
- PREC-003：`dependable_statuses` 与 `governed_globs` 均非空。

违反任一前置条件时，进程以退出码 2 结束，且不产出 findings。

## 4. 功能需求

### 4.1 抽取

- **REQ-001** 工具应遍历每个 `governed_globs` 模式，把匹配的文件解析为语料中的一条 document。
- **REQ-002** 语料中的每个受管文件应恰好产出一条 document 或一条 extract_error，不得静默丢弃。
- **REQ-003** frontmatter 解析应读取 `type`、`status`、`owner`、`depends_on`；缺失的可选字段留空。
- **REQ-004** 文件没有 frontmatter、没有结束分隔符、YAML 非法或不是映射时，应产出 `extract.parse_error`。
- **REQ-005** `depends_on` 的每一项应按有界规则归一化：反斜杠转正斜杠、折叠重复斜杠与尾斜杠、解析 `..`；
  解析后仍在仓库根内则接受，越界、绝对路径、URL、空路径均产生 extract_error。
- **REQ-006** 归一化后重复的 `depends_on` 应去重并保持首次出现顺序。
- **REQ-007** 单个文件超过 1 MiB，或 frontmatter 疑似 YAML 别名展开时，应产生 extract_error。

### 4.2 规则

- **REQ-010** 对每条依赖边，若目标路径不在语料中，应产出 `refs.missing_target`。
- **REQ-011** 对每条依赖边，若目标文档的 `status` 不等于任何 `dependable_statuses` 字面量，
  应产出 `refs.status_not_dependable`；`status` 缺失按同一规则处理。
- **REQ-012** 状态比较应逐字、大小写敏感。
- **REQ-013** `data.standard.dependable_statuses` 不存在时，策略应产出 `policy.load_error`，
  不得静默通过。

### 4.3 覆盖断言

- **REQ-020** 受管匹配数为 0 时，进程应退出 2。
- **REQ-021** `phase >= 1` 且解析出的文档数为 0 时，进程应退出 2。
- **REQ-022** `phase >= 1` 且依赖边总数为 0 时，进程应退出 2。
- **REQ-023** 每次输出应包含 `documents_matched`、`documents_parsed`、`edges_evaluated` 三个计数。

### 4.4 基线与阶段

- **REQ-030** `phase = 0` 时，进程退出码恒为 0，但 findings 与 extract_errors 仍应完整输出。
- **REQ-031** `phase >= 1` 且未指定 `-baseline` 时，任何 finding 都导致退出 1。
- **REQ-032** 指定 `-baseline` 时，`(rule, document, target)` 与基线条目精确匹配的 finding 应被抑制。
- **REQ-033** extract_error 不得被基线抑制；存在 extract_error 时 `phase >= 1` 退出 2。
- **REQ-034** 目标路径重命名后，旧基线条目失配，违规应重新暴露。

### 4.5 输出

- **REQ-040** `-format json` 输出应包含 `schema_version`、`findings`、`extract_errors`、`stats`。
- **REQ-041** findings 应按 `(rule, document, target)` 排序，保证字节稳定。
- **REQ-042** `-format github` 每条 finding 输出一行 `::error file=<path>::[<rule>] <message>`，
  并输出一行 `::notice::` 统计。
- **REQ-043** `-out <path>` 应把报告写入文件而不是 stdout。

### 4.6 策略完整性

- **REQ-050** `-policy-from-head` 时，策略包应从 `git HEAD` 读取。
- **REQ-051** `-policy-from-head` 时，若 `policy/`、`standard.json` 或 `baseline.json` 有未提交改动，
  应打印警告并降级为只报告（退出 0）。
- **REQ-052** `-policy-from-head` 时，若策略目录不在 `-root` 之内，应退出 2。

## 5. 非功能需求

- **NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒（`[estimated]`，待实测）。
- **NFR-002** Agent hook 场景下，≤200 篇文档的 p95 应 < 2 秒，超预算时降级为只警告（尚未实现）。
- **NFR-003** 策略求值应有 10 秒上限。
- **NFR-004** 构建目标为 darwin-arm64 与 linux-amd64。

## 6. 验收标准

| ID | 验收方式 | 状态 |
|---|---|---|
| REQ-001 | `TestCleanCorpusPasses` | 已实现并测试 |
| REQ-002 | `TestExtractErrorFailsClosedInPhase1` | 已实现并测试 |
| REQ-003 | `TestCleanCorpusPasses`（字段写入并回流） | 已实现并测试 |
| REQ-004 | `TestExtractErrorFailsClosedInPhase1` | 已实现并测试 |
| REQ-005 | `TestNormalizeRepoPath`、`TestNormalizedDependsOnResolves` | 已实现并测试 |
| REQ-006 | `TestNormalizeRepoPath` 间接覆盖 | 已实现，去重无独立用例 |
| REQ-007 | 无 | 已实现（大小上限与别名守卫），无测试 |
| REQ-010 | `TestMissingTargetBlocks` | 已实现并测试 |
| REQ-011 | `TestNonDependableStatusBlocks` | 已实现并测试 |
| REQ-012 | `TestLowercaseStatusIsNotDependable` | 已实现并测试 |
| REQ-013 | `TestStandardDataResolvesAsDataStandard`、`refs.rego` 守卫 | 已实现并测试 |
| REQ-020 | `TestZeroGovernedDocumentsFailsClosed` | 已实现并测试 |
| REQ-021 | 无独立用例 | 已实现，无测试 |
| REQ-022 | `TestZeroEdgesFailsClosed` | 已实现并测试 |
| REQ-023 | `TestJSONSchemaAndStats` | 已实现并测试 |
| REQ-030 | `TestPhase0NeverBlocks`、`TestExtractErrorIsReportOnlyInPhase0` | 已实现并测试 |
| REQ-031 | `TestMissingTargetBlocks` | 已实现并测试 |
| REQ-032 | `TestBaselineSuppressesKnownViolation` | 已实现并测试 |
| REQ-033 | `TestExtractErrorFailsClosedInPhase1` | 已实现并测试 |
| REQ-034 | `TestBaselineDoesNotSuppressRenamedTarget` | 已实现并测试 |
| REQ-040 | `TestJSONSchemaAndStats` | 已实现并测试 |
| REQ-041 | 无独立用例 | 已实现，无测试 |
| REQ-042 | 烟测 | 已实现，无自动化测试 |
| REQ-043 | 无 | 已实现，无测试 |
| REQ-050 | 无 | 已实现，**无测试** |
| REQ-051 | 无 | 已实现，**无测试** |
| REQ-052 | 无 | 已实现，**无测试** |
| NFR-001 | 无 | 未实测 |
| NFR-002 | 无 | **未实现** |
| NFR-003 | 无 | 已实现，无测试 |
| NFR-004 | `.github/workflows/release.yml` | 已定义，尚未在 CI 上实际跑过 |

## 7. 已知限制

1. **语料尚未回填。** 真实语料没有 `depends_on` 字段，直接对真实仓库运行会因"零依赖边"退出 2。
   这是设计中的 fail-closed 行为，不是缺陷。
2. **审批来源不可信。** `status: Approved` 由 Agent 与文档作者共用同一条写入路径，v1 不校验来源。
3. **循环引用未检测**，与 PDS-STD-004 §8.5.6 存在偏离。
4. **Rego 非本地内置函数未禁用**（例如 `http.send`），依靠 10 秒求值上限兜底。
5. **Agent hook 未安装时是静默无操作**，当前无检测手段。

## 8. 变更影响

修改本规格中任意 REQ 编号的语义，都必须同步更新
`docs/architecture/mdstandards-architecture.md` 的对应章节、`policy/` 中的规则，
以及 `main_test.go` 的验收用例。三处不一致即为规格漂移。
