# Spec Quality Report

- Target: `docs/specs/SPEC-001-mdstandards/spec.md`
- Rubric: 1.0.0 (`acb806c85c7a357f1cdfeaea80a8dc951333aefe5b630fe8576e0e0ce7f25947`)
- Target SHA-256: `db8c9985081365e918a49876a7314c3a07f072cdcfc19415156721be0b86804f`
- Reviewer session: `reviewer-spec-001-20261102`

**Headline score (core): 93.0**  |  Diagnostic score (all): 90.7

> The auditor verifies evidence presence, section placement, and score arithmetic. It does not judge whether evidence is sufficient or true.

> Normalization covers full-width ASCII, NBSP, and ideographic space; it is not full Unicode NFKC.

| Dimension | Headline | Diagnostic | Weight |
|---|---:|---:|---:|
| 背景落地 | 87.5 | 87.5 | 2 |
| 痛点证据 | 50.0 | 50.0 | 3 |
| 目标可验收 | 100.0 | 87.5 | 3 |
| 非目标边界 | 100.0 | 100.0 | 2 |
| 架构自洽 | 100.0 | 90.0 | 2 |
| 详细设计完整 | 100.0 | 100.0 | 2 |
| NFR 可度量 | 100.0 | 100.0 | 2 |
| 上线可回滚 | 100.0 | 100.0 | 2 |
| 替代方案诚实 | 100.0 | 100.0 | 2 |
| 跨章自洽 | 100.0 | 100.0 | 2 |
| 证据卫生 | 100.0 | 100.0 | 3 |

## 背景落地 (background)

- **bg.current_state** [core, level 1] — 1.1 是否写清了模块职责、核心流程与当前技术栈？
  - evidence `1.1 系统现状`: "1. Agent 生成或修改一篇 Spec，并在 'depends_on' 中写下一个或多个文档路径。"
  - normalized: "1. Agent 生成或修改一篇 Spec,并在 'depends_on' 中写下一个或多个文档路径。"
  - evidence `1.1 系统现状`: "这条流程上没有任何自动校验：没有工具检查目标文档是否存在，也没有工具检查目标文档是否已批准。"
  - normalized: "这条流程上没有任何自动校验:没有工具检查目标文档是否存在,也没有工具检查目标文档是否已批准。"
  - rationale: 写了三步核心流程与 Agent/评审者职责，也点明当前没有自动校验；但没有交代当前技术栈或工具链。
- **bg.failure_chain** [core, level 2] — 1.2 是否有从触发到后果的分步失败链？
  - evidence `1.2 核心技术痛点`: "1. 一份 Spec 声明 'depends_on' 指向某篇架构文档。"
  - normalized: "1. 一份 Spec 声明 'depends_on' 指向某篇架构文档。"
  - evidence `1.2 核心技术痛点`: "4. 构建全绿、页面都在、检索正常，代码却照着过期或未批准的架构写完了。"
  - normalized: "4. 构建全绿、页面都在、检索正常,代码却照着过期或未批准的架构写完了。"
  - rationale: 1.2 用编号 1–4 串起声明依赖、下游采信、无人校验、代码照过期或未批准架构写完的分步失败链。
- **bg.trigger** [core, level 2] — 1.3 是否解释了为什么现在做？
  - evidence `1.3 重构/设计动机`: "触发点：office-hours 设计已于 2026-10-30 批准，v1 进入实现与接入阶段；如果此时仍只靠人工抽查，引用图缺陷会在接入后继续静默进入下游。"
  - normalized: "触发点:office-hours 设计已于 2026-10-30 批准,v1 进入实现与接入阶段;如果此时仍只靠人工抽查,引用图缺陷会在接入后继续静默进入下游。"
  - rationale: 1.3 给出 2026-10-30 office-hours 批准、v1 进入实现与接入的触发点，并说明此时不做会让缺陷继续静默进入下游。
- **bg.non_goals** [core, level 2] — 1.4 是否写了不做的后果与本次非目标？
  - evidence `1.4 非目标`: "不把闸门做出来的后果见 1.2：失败继续静默，代码继续照着过期或未批准的架构写。"
  - normalized: "不把闸门做出来的后果见 1.2:失败继续静默,代码继续照着过期或未批准的架构写。"
  - evidence `1.4 非目标`: "- 循环引用检测（Deferred to Phase 2；与 PDS-STD-004 §8.5.6 的偏离见第 7 章）。"
  - normalized: "- 循环引用检测(Deferred to Phase 2;与 PDS-STD-004 §8.5.6 的偏离见第 7 章)。"
  - rationale: 1.4 先说明不建闸门的后果（失败继续静默），再列出本次明确不做或延后的非目标清单。

## 痛点证据 (pain_evidence)

- **bg.pain_evidence** [core, level 1] — 1.2 是否有至少一条可核实的证据（事故、度量或引述）？
  - evidence `1.2 核心技术痛点`: "证据：截至 2026-11-01，在 '/Users/alang/Workspace'（排除本仓库、'.claude' 技能目录与 'crawlers/' 抓取语料）的 411 篇 Markdown 中，只有 1 篇声明了 'depends_on'，且为空数组；真实语料主要用表格维护依赖（'meta_standards'、'spec_upstream_refs'）。"
  - normalized: "证据:截至 2026-11-01,在 '/Users/alang/Workspace'(排除本仓库、'.claude' 技能目录与 'crawlers/' 抓取语料)的 411 篇 Markdown 中,只有 1 篇声明了 'depends_on',且为空数组;真实语料主要用表格维护依赖('meta_standards'、'spec_upstream_refs')。"
  - rationale: 有带日期、路径范围与数量的 411 篇/1 篇度量，可核实；但它证明的是 depends_on 几乎为零，而非缺失或未批准目标造成的事故，文档也自认没有事故证据。

## 目标可验收 (goal_verifiability)

- **goal.verifiable** [core, level 2] — 2.1 每个 Goal 是否都有验收方式？
  - evidence `2.1 目标`: "- **G1 引用存在**：每条 'depends_on' 的目标必须在受管语料中；缺失产出 'refs.missing_target'。验收：'TestMissingTargetBlocks'。"
  - normalized: "- **G1 引用存在**:每条 'depends_on' 的目标必须在受管语料中;缺失产出 'refs.missing_target'。验收:'TestMissingTargetBlocks'。"
  - evidence `2.1 目标`: "- **G5 稳定输出**：JSON 含 'schema_version'、'findings'、'extract_errors'、'stats'；findings 按 '(rule, document, target)' 排序，字节稳定。验收：'TestJSONSchemaAndStats'（仅覆盖 schema/stats）；排序与字节稳定性无独立测试。"
  - normalized: "- **G5 稳定输出**:JSON 含 'schema_version'、'findings'、'extract_errors'、'stats';findings 按 '(rule, document, target)' 排序,字节稳定。验收:'TestJSONSchemaAndStats'(仅覆盖 schema/stats);排序与字节稳定性无独立测试。"
  - rationale: G1–G5 每个 Goal 都带验收并指向具体测试；G5 还如实标注排序与字节稳定性无独立测试。
- **goal.quantified** [core, level 2] — 2.1 是否至少有一个 Goal 有数字/单位/明确阈值？
  - evidence `2.1 目标`: "- **G3 Fail-closed**：零受管文档在任一 phase 都退出 2；零解析文档与零依赖边在 phase ≥ 1 退出 2，策略降级不绕过这些覆盖断言；解析错误在 phase ≥ 1 且未触发策略降级时退出 2，触发降级时只报告并退出 0，不得静默通过。验收：'TestZeroGovernedDocumentsFailsClosed'、'TestZeroEdgesFailsClosed'、'TestExtractErrorFailsClosedInPhase1'（降级路径未覆盖）。"
  - normalized: "- **G3 Fail-closed**:零受管文档在任一 phase 都退出 2;零解析文档与零依赖边在 phase ≥ 1 退出 2,策略降级不绕过这些覆盖断言;解析错误在 phase ≥ 1 且未触发策略降级时退出 2,触发降级时只报告并退出 0,不得静默通过。验收:'TestZeroGovernedDocumentsFailsClosed'、'TestZeroEdgesFailsClosed'、'TestExtractErrorFailsClosedInPhase1'(降级路径未覆盖)。"
  - rationale: G3 以任一 phase、phase ≥ 1、退出 2 给出明确阈值，G4 也用 phase 0/1 与退出 1 的判定。
- **goal.count** [extended, level 1] — 2.1 Goals 是否控制在 3–6 个？
  - evidence `2.1 目标`: "- **G1 引用存在**：每条 'depends_on' 的目标必须在受管语料中；缺失产出 'refs.missing_target'。验收：'TestMissingTargetBlocks'。"
  - normalized: "- **G1 引用存在**:每条 'depends_on' 的目标必须在受管语料中;缺失产出 'refs.missing_target'。验收:'TestMissingTargetBlocks'。"
  - evidence `2.1 目标`: "- **G5 稳定输出**：JSON 含 'schema_version'、'findings'、'extract_errors'、'stats'；findings 按 '(rule, document, target)' 排序，字节稳定。验收：'TestJSONSchemaAndStats'（仅覆盖 schema/stats）；排序与字节稳定性无独立测试。"
  - normalized: "- **G5 稳定输出**:JSON 含 'schema_version'、'findings'、'extract_errors'、'stats';findings 按 '(rule, document, target)' 排序,字节稳定。验收:'TestJSONSchemaAndStats'(仅覆盖 schema/stats);排序与字节稳定性无独立测试。"
  - rationale: 2.1 共 G1–G5 五个 Goal，落在 3–6 个区间；但未按架构、性能或可靠性分类。

## 非目标边界 (non_goal_boundary)

- **nongoal.labeled** [core, level 2] — 2.2 每条 Non-Goal 是否标了本次不做或 Deferred to Phase 2？
  - evidence `2.2 非目标`: "- **NG1** 必填字段、status 枚举、状态流转、TTL、审批来源校验（Deferred to Phase 2）。"
  - normalized: "- **NG1** 必填字段、status 枚举、状态流转、TTL、审批来源校验(Deferred to Phase 2)。"
  - evidence `2.2 非目标`: "- **NG3** pre-commit 执行点（本次不做）。"
  - normalized: "- **NG3** pre-commit 执行点(本次不做)。"
  - rationale: NG1–NG7 每条都带 Deferred to Phase 2 或本次不做的标注。
- **nongoal.not_dumping_ground** [extended, level 2] — 2.2 非目标是否能追溯到痛点或边界？
  - evidence `2.2 非目标`: "- **NG1** 必填字段、status 枚举、状态流转、TTL、审批来源校验（Deferred to Phase 2）。"
  - normalized: "- **NG1** 必填字段、status 枚举、状态流转、TTL、审批来源校验(Deferred to Phase 2)。"
  - evidence `2.2 非目标`: "- **NG3** pre-commit 执行点（本次不做）。"
  - normalized: "- **NG3** pre-commit 执行点(本次不做)。"
  - evidence `2.2 非目标`: "- **NG7** Rego 非本地内置函数未禁用（本次不做；仅靠 10 秒求值上限兜底）。"
  - normalized: "- **NG7** Rego 非本地内置函数未禁用(本次不做;仅靠 10 秒求值上限兜底)。"
  - rationale: NG1/NG2/NG6/NG7 可追溯到状态可信与引用图、身份、策略安全边界，7.6 与 3.4 给出理由；NG3–NG5 是执行点、分发与语料范围的明确交付边界。

## 架构自洽 (architecture)

- **arch.solution_shape** [core, level 2] — 第 3 章是否有拓扑或等价的模块/数据流描述？
  - evidence `3. 架构设计与核心方案`: "CLI[CLI flags] --> Load[load standard.json]"
  - normalized: "CLI[CLI flags] --> Load[load standard.json]"
  - evidence `3. 架构设计与核心方案`: "7. 渲染报告并返回退出码。"
  - normalized: "7. 渲染报告并返回退出码。"
  - rationale: 第 3 章同时给出 3.1 的 mermaid 拓扑图、3.2 的模块职责和 3.3 的七步数据流。
- **arch.topology** [extended, level 2] — 3.1 拓扑图是否让每个模块有一句职责？
  - evidence `3.1 总体架构`: "CLI[CLI flags] --> Load[load standard.json]"
  - normalized: "CLI[CLI flags] --> Load[load standard.json]"
  - evidence `3.1 总体架构`: "单二进制内嵌 OPA 引擎；策略包从仓库 'policy/' 读取，'policy/data/standard.json' 是受管范围与可依赖状态的唯一事实来源。"
  - normalized: "单二进制内嵌 OPA 引擎;策略包从仓库 'policy/' 读取,'policy/data/standard.json' 是受管范围与可依赖状态的唯一事实来源。"
  - rationale: 3.1 的 flowchart 画出了 CLI 到退出码的完整拓扑；3.2 为每个模块各写一句职责。
- **arch.change_set** [extended, level 2] — 第 3 章是否写清新增/改造/保持不变？
  - evidence `3. 架构设计与核心方案`: "- **新增**：单二进制 'mdstandards'、'policy/' 策略包、CI 与 Agent hook 两个执行点。"
  - normalized: "- **新增**:单二进制 'mdstandards'、'policy/' 策略包、CI 与 Agent hook 两个执行点。"
  - evidence `3. 架构设计与核心方案`: "- **改造**：无。v1 不修改既有文档的 frontmatter，也不改既有 PDS 校验器。"
  - normalized: "- **改造**:无。v1 不修改既有文档的 frontmatter,也不改既有 PDS 校验器。"
  - evidence `3. 架构设计与核心方案`: "- **保持不变**：文档目录结构、'depends_on' 语义、人工评审流程。"
  - normalized: "- **保持不变**:文档目录结构、'depends_on' 语义、人工评审流程。"
  - rationale: 3.1 用新增、改造、保持不变三段完整描述变更集。
- **arch.rationale** [extended, level 1] — 关键选型是否有理由，详细对比指向第 7 章？
  - evidence `3.4 关键选型与约束`: "- **单 'main' 包**：v1 只有几百行、没有第二个调用方；拆 'cmd/' + 'internal/' 属于过早抽象（决策 D2）。"
  - normalized: "- **单 'main' 包**:v1 只有几百行、没有第二个调用方;拆 'cmd/' + 'internal/' 属于过早抽象(决策 D2)。"
  - rationale: 3.4 对策略即数据、单 main 包、路径身份、有界归一化等关键选型给出了理由，但没有把详细对比指向第 7 章。

## 详细设计完整 (detailed_design)

- **design.contract** [core, level 2] — 第 4 章是否至少有一份数据模型或接口契约？
  - evidence `4. 详细设计`: "JSON 输出契约：'schema_version'、'findings'、'extract_errors'、'stats'。"
  - normalized: "JSON 输出契约:'schema_version'、'findings'、'extract_errors'、'stats'。"
  - rationale: 第 4 章有 4.3 的实体字段表和 JSON 输出契约，至少包含一份数据模型与接口契约。
- **design.data_model** [extended, level 2] — 4.3 数据模型是否有字段类型、必填与校验？
  - evidence `4.3 数据模型与契约`: "\| Standard \| governed_globs \| string[] \| 是 \| 受管路径模式，非空 \|"
  - normalized: "\| Standard \| governed_globs \| string[] \| 是 \| 受管路径模式,非空 \|"
  - evidence `4.3 数据模型与契约`: "\| Standard \| identity \| string \| 是 \| 固定为 'repo_relative_path' \|"
  - normalized: "\| Standard \| identity \| string \| 是 \| 固定为 'repo_relative_path' \|"
  - rationale: 4.3 表按实体、字段、类型、必填、说明列出字段，并在说明列给出非空、固定值等校验约束。
- **design.interface** [extended, level 2] — 4.4 接口是否有请求/响应/状态码或方法签名？
  - evidence `4.4 接口与协议`: "\| '-phase' \| '1' \| '0' 只报告、'1' 基线 + 阻断、'2' 完全阻断 \|"
  - normalized: "\| '-phase' \| '1' \| '0' 只报告、'1' 基线 + 阻断、'2' 完全阻断 \|"
  - evidence `4.4 接口与协议`: "- 'github'：每条 finding 与 extract_error 各输出一行 '::error file=<path>::[<rule>] <message>'，另有一行 '::notice::' 统计。"
  - normalized: "- 'github':每条 finding 与 extract_error 各输出一行 '::error file=<path>::[<rule>] <message>',另有一行 '::notice::' 统计。"
  - rationale: 4.4 以 CLI 参数表给出请求侧参数与默认值，并定义 table/json/github 三种响应格式；退出码在 4.6 补齐。
- **design.sequence** [extended, level 2] — 4.5 复杂流程是否有时序图或等价分步说明？
  - evidence `4.5 核心流程与时序`: "alt 阻断阶段且未降级"
  - normalized: "alt 阻断阶段且未降级"
  - evidence `4.5 核心流程与时序`: "P-->>M: 标准或 policy.load_error"
  - normalized: "P-->>M: 标准或 policy.load_error"
  - rationale: 4.5 的 sequenceDiagram 覆盖正常调用，并用 alt 区分阻断与降级；policy.load_error 与 extract_errors 作为异常返回也画在图中。
- **design.errors** [extended, level 2] — 4.6 每个错误码是否有触发场景与处理策略？
  - evidence `4.6 异常处理与错误码`: "\| 'refs.status_not_dependable' \| 目标状态缺失或不在 'dependable_statuses' 中 \| 计为 finding；阻断阶段退出 1 \|"
  - normalized: "\| 'refs.status_not_dependable' \| 目标状态缺失或不在 'dependable_statuses' 中 \| 计为 finding;阻断阶段退出 1 \|"
  - evidence `4.6 异常处理与错误码`: "\| 'extract.parse_error' \| 无 frontmatter、YAML 非法、非映射等 \| 计入 'extract_errors'；未降级的阻断阶段退出 2，降级时只报告；不可被基线抑制 \|"
  - normalized: "\| 'extract.parse_error' \| 无 frontmatter、YAML 非法、非映射等 \| 计入 'extract_errors';未降级的阻断阶段退出 2,降级时只报告;不可被基线抑制 \|"
  - rationale: 4.6 表为四个规则 ID 分别给出触发场景与处理策略。

## NFR 可度量 (nfr)

- **nfr.measurable** [core, level 2] — 第 5 章每条 NFR 是否有数字、单位或验收方式？
  - evidence `5. 非功能性需求`: "- **NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒（'[estimated]'，待实测）。"
  - normalized: "- **NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒('[estimated]',待实测)。"
  - evidence `5. 非功能性需求`: "- **NFR-007** stdout 输出结构化报告；stderr 输出 'policy.load_error'、'coverage error'、'baseline error' 等可定位的错误原因；退出码固定为 0/1/2。"
  - normalized: "- **NFR-007** stdout 输出结构化报告;stderr 输出 'policy.load_error'、'coverage error'、'baseline error' 等可定位的错误原因;退出码固定为 0/1/2。"
  - rationale: 每条 NFR 都有数字、单位或可核实验收条件：性能项有秒数，可靠性、可观测性与安全项用退出码和具体行为断言。
- **nfr.coverage** [core, level 2] — 第 5 章是否覆盖性能、可靠性、可观测性、安全，或说明不适用？
  - evidence `5. 非功能性需求`: "- 长时间运行的指标、Trace 与告警不适用：本工具是短生命周期 CLI，由 CI 与 Agent hook 同步调用。"
  - normalized: "- 长时间运行的指标、Trace 与告警不适用:本工具是短生命周期 CLI,由 CI 与 Agent hook 同步调用。"
  - evidence `5. 非功能性需求`: "- 传输加密、身份认证与多租户权限不适用：工具只读取本地仓库，不对外提供服务。"
  - normalized: "- 传输加密、身份认证与多租户权限不适用:工具只读取本地仓库,不对外提供服务。"
  - rationale: 第 5 章分性能与资源、可靠性与容灾、可观测性、安全与合规四类，并对不适用的长时指标和认证写明理由。
- **nfr.honesty** [core, level 2] — 未实测的指标是否标了 [estimated] 或 未实现？
  - evidence `5. 非功能性需求`: "- **NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒（'[estimated]'，待实测）。"
  - normalized: "- **NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒('[estimated]',待实测)。"
  - evidence `5. 非功能性需求`: "- **NFR-002** Agent hook 场景下，≤200 篇文档的 p95 应 < 2 秒，超预算时降级为只警告（尚未实现）。"
  - normalized: "- **NFR-002** Agent hook 场景下,≤200 篇文档的 p95 应 < 2 秒,超预算时降级为只警告(尚未实现)。"
  - rationale: NFR-001 标 '[estimated]' 并注待实测，NFR-002 明确标尚未实现。

## 上线可回滚 (rollout)

- **rollout.release** [core, level 2] — 6.1 是否有发布策略与灰度节奏？
  - evidence `6.1 发布与灰度`: "- **阶段 0（只报告）**：不因违规或解析错误阻断（退出 0；零受管文档、策略加载失败等工具级错误仍退出 2），上传 'findings.json' 和注解，收集存量。"
  - normalized: "- **阶段 0(只报告)**:不因违规或解析错误阻断(退出 0;零受管文档、策略加载失败等工具级错误仍退出 2),上传 'findings.json' 和注解,收集存量。"
  - evidence `6.1 发布与灰度`: "晋级条件：连续两周的新增违规数为 0，才允许 0 → 1；基线为空才允许 1 → 2。"
  - normalized: "晋级条件:连续两周的新增违规数为 0,才允许 0 → 1;基线为空才允许 1 → 2。"
  - rationale: 6.1 给出阶段 0/1/2 的发布策略，以连续两周零新增、基线为空作为晋级观察条件。
- **rollout.migration** [core, level 2] — 6.2 是否有迁移/兼容方案，或说明无需迁移？
  - evidence `6.2 数据迁移与兼容性`: "- 不迁移历史数据；存量违规用 baseline 记录，不修改旧文档来凑合规。"
  - normalized: "- 不迁移历史数据;存量违规用 baseline 记录,不修改旧文档来凑合规。"
  - evidence `6.2 数据迁移与兼容性`: "- 真实语料几乎没有 'depends_on'（唯一一处为空数组），以 phase 1 直接运行会因“零依赖边”退出 2；回填 'depends_on' 是独立的 v2 任务。"
  - normalized: "- 真实语料几乎没有 'depends_on'(唯一一处为空数组),以 phase 1 直接运行会因“零依赖边”退出 2;回填 'depends_on' 是独立的 v2 任务。"
  - rationale: 6.2 明确不迁移历史数据，存量违规用 baseline 记录，并说明回填 depends_on 属独立 v2 任务。
- **rollout.smoke** [core, level 2] — 6.3 是否有冒烟清单与观察指标？
  - evidence `6.3 上线验证与冒烟`: "1. 'make check' 在本仓库退出 0；示例语料用 './mdstandards -root ./testdata/example-repo -policy ./policy -phase 1 -format json' 单独冒烟。"
  - normalized: "1. 'make check' 在本仓库退出 0;示例语料用 './mdstandards -root ./testdata/example-repo -policy ./policy -phase 1 -format json' 单独冒烟。"
  - evidence `6.3 上线验证与冒烟`: "观察指标：'findings' 数量、'extract_errors' 数量、'documents_matched'、'edges_evaluated'、退出码分布。"
  - normalized: "观察指标:'findings' 数量、'extract_errors' 数量、'documents_matched'、'edges_evaluated'、退出码分布。"
  - rationale: 6.3 有四项冒烟步骤，并列出 findings、extract_errors、documents_matched、edges_evaluated 与退出码分布作为观察指标。
- **rollout.rollback_trigger** [core, level 2] — 6.4 回滚触发条件是否是可观测数字？
  - evidence `6.4 回滚与熔断`: "- **误报**：任一 phase 1 运行中，经人工复核确认为误报的阻断 finding ≥ 1，或连续 2 个工作日误报率 > 5%。"
  - normalized: "- **误报**:任一 phase 1 运行中,经人工复核确认为误报的阻断 finding ≥ 1,或连续 2 个工作日误报率 > 5%。"
  - evidence `6.4 回滚与熔断`: "- **性能**：全量检查 p95 > 10 秒连续 3 次。"
  - normalized: "- **性能**:全量检查 p95 > 10 秒连续 3 次。"
  - rationale: 6.4 用误报 ≥1、误报率 >5%、工具错误出现 1 次、p95 >10 秒连续 3 次等可观测数字定义触发条件。
- **rollout.rollback_steps** [core, level 2] — 6.4 回滚步骤是否可执行？
  - evidence `6.4 回滚与熔断`: "1. 把 CI 的 '-phase' 从 1 调回 0，保留 'findings.json' 作为证据。"
  - normalized: "1. 把 CI 的 '-phase' 从 1 调回 0,保留 'findings.json' 作为证据。"
  - rationale: 6.4 给出五步可执行回滚步骤，从调低 phase 到规则修复与恢复阻断。

## 替代方案诚实 (alternatives)

- **alts.status_quo** [core, level 2] — 第 7 章是否有 Status Quo 及其成本/风险？
  - evidence `7. 替代方案与权衡`: "- **淘汰原因**：Agent 产出是指数的，人工吞吐固定；失败是静默的，抽查只能覆盖抽样，无法在合并前拦住引用图缺陷。"
  - normalized: "- **淘汰原因**:Agent 产出是指数的,人工吞吐固定;失败是静默的,抽查只能覆盖抽样,无法在合并前拦住引用图缺陷。"
  - rationale: 7.1 把人工抽查作为 Status Quo 方案 A，给出零成本优点，并在淘汰原因中说明吞吐固定与静默失败风险。
- **alts.rejected** [core, level 2] — 第 7 章是否至少有一个被否决方案，且淘汰原因对应当前约束？
  - evidence `7. 替代方案与权衡`: "- **淘汰原因**：失败模式是文档引用图，不是单个文件的字段格式；单文件检查抓不到“目标不存在或未批准”。"
  - normalized: "- **淘汰原因**:失败模式是文档引用图,不是单个文件的字段格式;单文件检查抓不到“目标不存在或未批准”。"
  - evidence `7. 替代方案与权衡`: "- **淘汰原因**：规则写在代码里，标准变更需要改代码；策略无法作为独立事实来源评审。"
  - normalized: "- **淘汰原因**:规则写在代码里,标准变更需要改代码;策略无法作为独立事实来源评审。"
  - rationale: 7.2 与 7.3 分别否决单文件检查与 Schema+薄校验，淘汰原因直接对应当前约束：引用图缺陷与策略不可独立评审。
- **alts.matrix** [core, level 2] — 第 7 章是否有与约束对应的权衡矩阵？
  - evidence `7. 替代方案与权衡`: "\| 维度 \| A 人工抽查 \| B 单文件检查 \| C Schema + 薄校验 \| D OPA/Rego（当前） \|"
  - normalized: "\| 维度 \| A 人工抽查 \| B 单文件检查 \| C Schema + 薄校验 \| D OPA/Rego(当前) \|"
  - evidence `7. 替代方案与权衡`: "\| 能发现引用图缺陷 \| 否 \| 否 \| 是 \| 是 \|"
  - normalized: "\| 能发现引用图缺陷 \| 否 \| 否 \| 是 \| 是 \|"
  - rationale: 7.5 矩阵以引用图缺陷、策略可独立评审、新增规则成本、运行时依赖、契合 v1 范围为维度，对应痛点与 3.4 约束。
- **alts.consistency** [core, level 2] — 第 7 章结论是否与 Goals、架构选型一致？
  - evidence `7. 替代方案与权衡`: "- **理由**：规则与程序分离，策略可独立评审；单二进制无额外运行时；基线与分阶段可复用同一套求值结果。"
  - normalized: "- **理由**:规则与程序分离,策略可独立评审;单二进制无额外运行时;基线与分阶段可复用同一套求值结果。"
  - rationale: 7.4 选择 D（OPA/Rego），与 3.4 的策略即数据以及 G1/G2 的引用、状态规则一致，矩阵结论也指向 D。

## 跨章自洽 (cross_consistency)

- **x.goals_trace** [core, level 2] — 后续章节是否逐条回应 Goals？
  - evidence `2.1 目标`: "- **G1 引用存在**：每条 'depends_on' 的目标必须在受管语料中；缺失产出 'refs.missing_target'。验收：'TestMissingTargetBlocks'。"
  - normalized: "- **G1 引用存在**:每条 'depends_on' 的目标必须在受管语料中;缺失产出 'refs.missing_target'。验收:'TestMissingTargetBlocks'。"
  - evidence `4.7.2 规则`: "- **REQ-010** 对每条依赖边，若目标路径不在语料中，应产出 'refs.missing_target'。"
  - normalized: "- **REQ-010** 对每条依赖边,若目标路径不在语料中,应产出 'refs.missing_target'。"
  - rationale: G1→REQ-010、G2→REQ-011/012、G3→REQ-020–022、G4→REQ-030–034、G5→REQ-040/041，每个 Goal 都能在后续章节找到对应设计。
- **x.nfr_arch** [core, level 2] — NFR 与架构、选型是否一致？
  - evidence `5. 非功能性需求`: "- **NFR-005** 闸门必须 fail-closed：零受管文档、策略加载失败在任何 phase 都不得返回成功；零解析文档、零依赖边在 'phase >= 1' 都不得返回成功；解析错误在阻断阶段且未触发策略降级时不得返回成功。覆盖断言不受策略降级影响。"
  - normalized: "- **NFR-005** 闸门必须 fail-closed:零受管文档、策略加载失败在任何 phase 都不得返回成功;零解析文档、零依赖边在 'phase >= 1' 都不得返回成功;解析错误在阻断阶段且未触发策略降级时不得返回成功。覆盖断言不受策略降级影响。"
  - evidence `4.7.3 覆盖断言`: "- **REQ-020** 受管匹配数为 0 时，进程应退出 2。"
  - normalized: "- **REQ-020** 受管匹配数为 0 时,进程应退出 2。"
  - rationale: NFR-005 的 fail-closed 对应 4.7.3 覆盖断言，NFR-006 的字节稳定对应 REQ-041 排序，NFR-008 的 1 MiB/10 秒对应 REQ-007 与 NFR-003。
- **x.rollback_release** [core, level 2] — 回滚与发布策略是否一致？
  - evidence `6.4 回滚与熔断`: "1. 把 CI 的 '-phase' 从 1 调回 0，保留 'findings.json' 作为证据。"
  - normalized: "1. 把 CI 的 '-phase' 从 1 调回 0,保留 'findings.json' 作为证据。"
  - evidence `6.1 发布与灰度`: "- **阶段 0（只报告）**：不因违规或解析错误阻断（退出 0；零受管文档、策略加载失败等工具级错误仍退出 2），上传 'findings.json' 和注解，收集存量。"
  - normalized: "- **阶段 0(只报告)**:不因违规或解析错误阻断(退出 0;零受管文档、策略加载失败等工具级错误仍退出 2),上传 'findings.json' 和注解,收集存量。"
  - rationale: 回滚第 1 步把 phase 从 1 调回 0，与 6.1 的阶段 0 只报告策略一致；修复后再晋级也与 6.1 的晋级条件衔接。

## 证据卫生 (evidence_hygiene)

- **x.no_placeholders** [core, level 2] — 是否没有把 无 / ... / 待补充 / TBD 当作实质回答？
  - evidence `1.2 核心技术痛点`: "缺口：本次没有拿到一次具体事故的时间、涉及文件数与投入工时。这一点记入附录 A。"
  - normalized: "缺口:本次没有拿到一次具体事故的时间、涉及文件数与投入工时。这一点记入附录 A。"
  - evidence `6.3 上线验证与冒烟`: "\| REQ-043 \| 无 \| 已实现，无测试 \|"
  - normalized: "\| REQ-043 \| 无 \| 已实现,无测试 \|"
  - rationale: 全文没有 TBD、待补充或省略号式占位；6.3 中的无是附带的明确无测试状态，不是用来掩盖内容缺失。
- **x.no_code** [core, level 2] — 是否没有堆砌实现代码、测试日志或偶发报错？
  - evidence `4.7.1 抽取`: "- **REQ-001** 工具应遍历每个 'governed_globs' 模式，把匹配的文件解析为语料中的一条 document。"
  - normalized: "- **REQ-001** 工具应遍历每个 'governed_globs' 模式,把匹配的文件解析为语料中的一条 document。"
  - rationale: 全文只有 mermaid 图、参数与字段表格、需求条目，没有实现代码、测试日志或偶发报错堆砌。
- **x.numbers_sourced** [core, level 2] — 关键数字是否有出处，或标了 [estimated]？
  - evidence `1.2 核心技术痛点`: "证据：截至 2026-11-01，在 '/Users/alang/Workspace'（排除本仓库、'.claude' 技能目录与 'crawlers/' 抓取语料）的 411 篇 Markdown 中，只有 1 篇声明了 'depends_on'，且为空数组；真实语料主要用表格维护依赖（'meta_standards'、'spec_upstream_refs'）。"
  - normalized: "证据:截至 2026-11-01,在 '/Users/alang/Workspace'(排除本仓库、'.claude' 技能目录与 'crawlers/' 抓取语料)的 411 篇 Markdown 中,只有 1 篇声明了 'depends_on',且为空数组;真实语料主要用表格维护依赖('meta_standards'、'spec_upstream_refs')。"
  - evidence `5. 非功能性需求`: "- **NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒（'[estimated]'，待实测）。"
  - normalized: "- **NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒('[estimated]',待实测)。"
  - rationale: 411 篇/1 篇度量带日期与语料范围，NFR-001 的 10 秒标 '[estimated]' 并注待实测，NFR-002 标尚未实现；其余数字是设计阈值。
- **x.no_fabricated_refs** [core, level 2] — 是否没有捏造 SPEC 引用、depends_on 目标或 M<N>？
  - evidence `7.6 已知限制与偏离`: "3. **循环引用未检测**：与 PDS-STD-004 §8.5.6 存在偏离。"
  - normalized: "3. **循环引用未检测**:与 PDS-STD-004 §8.5.6 存在偏离。"
  - rationale: frontmatter 的唯一 depends_on 是同时出现在关联资源的架构文档，正文引用 PDS-STD-004/001 的具体条款，全文没有 M<N> 占位。

## Audit failures

None.
