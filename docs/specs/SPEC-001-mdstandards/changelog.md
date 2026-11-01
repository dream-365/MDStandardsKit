---
type: changelog
status: In Review
owner: alang
depends_on: []
updated_at: 2026-11-01
---

# Changelog: SPEC-001 MDStandardsKit 元数据闸门

## [2026-11-01] - 独立评审后修订

- **修改类型**: Fixed（评审修订）
- **影响范围**: G3/G5、REQ-023、REQ-030、REQ-032、REQ-033、REQ-042、REQ-051、NFR-005..009、4.1 术语、4.4 接口、4.6 退出码、6.3 冒烟、6.4 工具错误触发条件、验收表、第 1.2 节证据口径、新增 `docs/specs/README.md`
- **变更原因**: 独立评审发现 4 个阻断项：依赖边定义与规则矛盾；策略降级与覆盖断言的退出码表述与 `main.go` 不一致；验收表把未覆盖行为写成已测试；changelog 对改动的描述失实
- **具体内容**:
  - 4.1 “依赖边”改为“目标可以不在语料中，是否存在由规则判定”，与 `policy/refs.rego` 和 REQ-010 一致。
  - G3、REQ-030/032/033/051、3.3 数据流第 6-7 步、4.6 退出码与规则表补上“未触发策略降级”限定；明确降级只跳过 baseline 过滤和 finding/extract_error 阻断，覆盖断言（零受管文档、零解析文档、零依赖边）仍退出 2。
  - NFR-005 与 6.4 工具错误触发条件区分覆盖断言与降级边界：零解析文档、零依赖边在 `phase >= 1` 均不返回成功；extract_error 仅在未降级时阻断。
  - 验收表改为诚实标注：REQ-006/REQ-023/NFR-006 不再写“已测试”；REQ-033/050/051、NFR-005/007/009 注明测试只覆盖部分路径；G5 注明 `TestJSONSchemaAndStats` 不覆盖排序与字节稳定性。
  - 6.3 冒烟第 1 条改为本仓库 `make check`，示例语料单独用 `-root ./testdata/example-repo`；第 4 条补上覆盖断言例外。
  - 4.4 注明 `-dump-corpus` 是调试路径，不执行覆盖断言与规则求值；GitHub 输出与 REQ-042 补上 extract_error 行。
  - 1.2 证据口径写明额外排除 `crawlers/` 抓取语料，使 411 篇可复现。
  - 全文统一使用“受管文档”。
  - 新增 `docs/specs/README.md` 索引，满足 Spec 标准对总索引的要求。
- **状态变更**: 保持 `In Review`，等待复审。

## [2026-11-01] - 按最新 doc-writer 规范重写

- **修改类型**: Changed
- **影响范围**: 全文结构；REQ-001..052 编号保留，REQ-007 补充别名阈值；PREC-001..003 编号与语义保留；NFR-001/002/004 语义保留，NFR-003 补充超时处理，新增 NFR-005..009
- **变更原因**: 按最新 doc-writer 技能的分章规范重写，补齐第 3 章“新增/改造/不变”、第 4 章数据模型与接口时序、第 5 章五类 NFR、第 6 章冒烟清单与可观测回滚条件、第 7 章权衡矩阵
- **具体内容**:
  - 第 1 章把现状写成三步核心流程；1.4 明确不做的后果。
  - 第 2 章保留 G1-G5 与 NG1-NG7。
  - 第 3 章标注新增、改造、保持不变，模块职责改为“一句话职责 + 代码布局”两层；3.4 策略完整性补充 `-baseline` 脏检查。
  - 第 4 章保留已有术语与前置条件，补充“基线”术语；新增数据模型表、CLI 参数表、时序图、规则与退出码表；REQ-007 补上 `*` > 100 与 `<<` > 10 的别名阈值。
  - 第 5 章拆为性能、可靠性、可观测性、安全、构建五类；NFR-004 保持构建目标不变，新增 NFR-005..009，并说明指标/Trace/认证不适用的原因。
  - 第 6 章新增冒烟清单与观察指标；回滚触发条件改为可观测的门槛；验收表补齐新增 NFR。
  - 第 7 章改为四个方案（含 Status Quo）、权衡矩阵、已知限制与变更影响。
- **状态变更**: `spec.md` 与 `changelog.md` 的 `status` 由 `Approved` 回退为 `In Review`，等待重新评审。

## [2026-11-01] - Agent 复核通过（baseline 脏检查修复）

- **修改类型**: Changed（评审通过）
- **影响范围**: REQ-051、验收表 REQ-051、7.2 已知限制；复核架构文档 §5/§6 与 CODEOWNERS 表述（文档未改）
- **复核原因**: `-policy-from-head` 原先只检查策略目录的未提交改动，不检查 `-baseline` 指向的文件，Agent 可通过改基线改变判定而不触发降级。修复后 Spec 回退为 `In Review`，由用户 alang 委托独立 agent 代审。
- **复核结论**: 通过，无阻断项。
  - REQ-051 与 `main.go` 一致：`gitPathDirty` 同时检查策略目录与 `-baseline` 文件（含未跟踪改动）；`-baseline` 在 `-root` 之外或 git 检查失败时按不可验证处理。`dirtyPolicy` 为真时 `blocking` 为假，跳过 baseline 过滤但保留 findings，打印 report-only 警告并退出 0。
  - 回归测试覆盖脏 baseline（工作区改为 `[]` 后退出 0 且带 report-only 警告）与干净 baseline（`git checkout -- baseline.json` 后提交的条目仍抑制 finding，退出 0 且无警告）。测试依赖真实 git，`ubuntu-latest` 的 CI 环境具备 git 与 `cp`。
  - changelog 修复条目已写明修改类型、影响范围、变更原因与具体内容；架构文档 §5/§6 与代码一致；CODEOWNERS 过期表述已修正为已由 `.github/CODEOWNERS` 保护，实际覆盖 `/policy/` 与 `/baseline.json`。
- **测试结果**: `go test -run TestPolicyFromHeadDowngradesOnDirtyBaseline -v` PASS；`go test -count=1 ./...` 通过；`make check` 退出 0（`OK: scanned 3 documents, evaluated 1 edges`）；`-phase 0` 与 `-phase 1` JSON 均退出 0、无 findings。手工验证：脏 baseline 降级时 findings 仍输出；`-baseline` 在 `-root` 外同样降级。
- **非阻断观察**: §3.4 的摘要句仍只提"工作区策略"，REQ-051 已完整覆盖 baseline；降级警告文本对"确实脏"与"在 root 外不可验证"不作区分。均不影响语义。
- **评审人**: 独立 agent（用户 alang 委托代审）
- **状态变更**: `spec.md` 与 `changelog.md` 的 `status` 由 `In Review` 改为 `Approved`。

## [2026-11-01] - 修复 -policy-from-head 的 baseline 脏检查

- **修改类型**: Fixed
- **影响范围**: REQ-051；7.2 移除"基线脏改动未检查"；验收表 REQ-051
- **变更原因**: 架构文档承诺 baseline 脏改动也降级，但代码只检查策略目录；Agent 可通过修改 baseline 改变判定而不触发降级。
- **具体内容**: `main.go` 的 `materializePolicyFromHEAD` 同时检查 `-baseline` 指向文件的未提交改动；`-baseline` 在 `-root` 之外时视为不可验证并降级为只报告。新增 `TestPolicyFromHeadDowngradesOnDirtyBaseline`，覆盖脏 baseline 与干净 baseline 两条路径。REQ-051 恢复完整语义，7.2 删除旧的已知限制。状态由 `Approved` 回退为 `In Review`，等待代审复核。

## [2026-11-01] - Agent 代审与修订

- **修改类型**: Changed（评审修订）
- **影响范围**: 全文；REQ-001..052、PREC-001..003、NFR-001..004 的编号与验收语义保留
- **变更原因**: 独立 agent 按用户 alang 的委托代审，逐条核对实现、测试、架构文档与写作规范后修正不一致项，并推进状态
- **具体内容**:
  - §1.2 的证据度量更新为 2026-11-01 可复核的快照：排除本仓库与 `.claude` 技能目录后，411 篇 Markdown 中只有 1 篇声明了空的 `depends_on`；§1.3 补充触发点（office-hours 设计批准后进入 v1 实现与接入阶段）。
  - 澄清 REQ-005 的归一化行为（去掉 `.` 段；URL 以 `://` 判定），与 `normalize.go` 和 `TestNormalizeRepoPath` 一致。
  - 补充 REQ-030/REQ-031 的阶段与降级边界：phase 0 不因 finding 或 extract_error 阻断，但工具级错误仍退出 2；策略降级优先于阻断判定。与 `main.go` 的退出码一致。
  - 将 REQ-051 收窄为代码实际行为（仅策略目录脏检查），并在 7.2 增加"基线脏改动未检查"的已知限制。
  - 修正验收表：`refs.rego` 守卫改为 `status.rego`；`烟测` 改为 `无`；为部分覆盖的 REQ-003/004/011/013/020 注明测试边界；REQ-042 标注无自动化测试。
  - 同步修正 §3.4、§4.4、§6.1、§6.2、§6.4 与附录 A 中与 phase 0、零依赖边、证据度量相关的表述，避免与 REQ-030/REQ-051 和 `main.go` 的退出码冲突。
  - 非目标统一使用 `本次不做` / `Deferred to Phase 2` 词汇；7.2 补充 PDS 状态词汇与 `doc_id` 偏离的已知限制。
  - 评审人字段改为独立 agent（用户 alang 委托代审）。
- **评审结论**: Agent 代审通过，评审人：独立 agent（用户 alang 明确委托）；v1 不校验审批来源。
- **状态变更**: `spec.md` 与 `changelog.md` 的 `status` 由 `In Review` 改为 `Approved`。

## [2026-11-01] - 迁移到 SPEC-001 目录并按新模板重写

- **修改类型**: Changed
- **影响范围**: 全文结构；REQ-001..052、NFR-001..004、验收标准保留
- **变更原因**: 用 doc-writer 技能重写，补齐背景、痛点、目标/非目标、替代方案与回滚触发条件
- **具体内容**: 从 `docs/specs/mdstandards-v1.md` 迁移到 `docs/specs/SPEC-001-mdstandards/spec.md`；新增 `changelog.md`；状态从 Draft 改为 In Review。
