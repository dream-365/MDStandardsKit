# Spec Quality Report

- Target: `docs/specs/SPEC-001-mdstandards/spec.md`
- Rubric: 1.0.0 (`acb806c85c7a357f1cdfeaea80a8dc951333aefe5b630fe8576e0e0ce7f25947`)
- Target SHA-256: `bb8d4081ca8b4738feeeea3a8c22353dabf47fdeb466803ef4593b8b4d0f1d56`
- Reviewer session: `reviewer-spec-001-20261103`

**Headline score (core): 91.0**  |  Diagnostic score (all): 87.4

> The auditor verifies evidence presence, section placement, and score arithmetic. It does not judge whether evidence is sufficient or true.

> Normalization covers full-width ASCII, NBSP, and ideographic space; it is not full Unicode NFKC.

| Dimension | Headline | Diagnostic | Weight |
|---|---:|---:|---:|
| 背景落地 | 87.5 | 87.5 | 2 |
| 痛点证据 | 50.0 | 50.0 | 3 |
| 目标可验收 | 100.0 | 87.5 | 3 |
| 非目标边界 | 100.0 | 83.3 | 2 |
| 架构自洽 | 100.0 | 90.0 | 2 |
| 详细设计完整 | 100.0 | 100.0 | 2 |
| NFR 可度量 | 75.0 | 75.0 | 2 |
| 上线可回滚 | 100.0 | 100.0 | 2 |
| 替代方案诚实 | 100.0 | 100.0 | 2 |
| 跨章自洽 | 100.0 | 100.0 | 2 |
| 证据卫生 | 100.0 | 100.0 | 3 |

## 背景落地 (background)

- **bg.current_state** [core, level 1] — 1.1 是否写清了模块职责、核心流程与当前技术栈？
  - evidence `1.1 系统现状`: "当前的核心流程如下："
  - normalized: "当前的核心流程如下:"
  - evidence `1.1 系统现状`: "1. Agent 生成或修改一篇 Spec，并在 'depends_on' 中写下一个或多个文档路径。"
  - normalized: "1. Agent 生成或修改一篇 Spec,并在 'depends_on' 中写下一个或多个文档路径。"
  - rationale: 写清了文档 frontmatter 的角色与三步核心流程，但未交代当前技术栈（工具、语言、运行方式），不足以完整还原现状。
- **bg.failure_chain** [core, level 2] — 1.2 是否有从触发到后果的分步失败链？
  - evidence `1.2 核心技术痛点`: "失败链如下："
  - normalized: "失败链如下:"
  - evidence `1.2 核心技术痛点`: "4. 构建全绿、页面都在、检索正常，代码却照着过期或未批准的架构写完了。"
  - normalized: "4. 构建全绿、页面都在、检索正常,代码却照着过期或未批准的架构写完了。"
  - rationale: 1.2 以编号步骤给出从 Spec 声明 depends_on 到代码照过期架构写完的完整失败链，触发与后果可直接核对。
- **bg.trigger** [core, level 2] — 1.3 是否解释了为什么现在做？
  - evidence `1.3 重构/设计动机`: "触发点：office-hours 设计已于 2026-10-30 批准，v1 进入实现与接入阶段；如果此时仍只靠人工抽查，引用图缺陷会在接入后继续静默进入下游。"
  - normalized: "触发点:office-hours 设计已于 2026-10-30 批准,v1 进入实现与接入阶段;如果此时仍只靠人工抽查,引用图缺陷会在接入后继续静默进入下游。"
  - rationale: 明确给出设计批准时间与接入阶段，并说明继续人工抽查的具体时间性风险，回答了为什么现在做。
- **bg.non_goals** [core, level 2] — 1.4 是否写了不做的后果与本次非目标？
  - evidence `1.4 非目标`: "不把闸门做出来的后果见 1.2：失败继续静默，代码继续照着过期或未批准的架构写。"
  - normalized: "不把闸门做出来的后果见 1.2:失败继续静默,代码继续照着过期或未批准的架构写。"
  - evidence `1.4 非目标`: "- 必填字段、status 枚举、状态流转、复审有效期、审批来源校验（Deferred to Phase 2）。"
  - normalized: "- 必填字段、status 枚举、状态流转、复审有效期、审批来源校验(Deferred to Phase 2)。"
  - rationale: 先说明不做的后果，再列出明确延后或本次不做的条目，两者齐备。

## 痛点证据 (pain_evidence)

- **bg.pain_evidence** [core, level 1] — 1.2 是否有至少一条可核实的证据（事故、度量或引述）？
  - evidence `1.2 核心技术痛点`: "只有 1 篇声明了 'depends_on'，且为空数组"
  - normalized: "只有 1 篇声明了 'depends_on',且为空数组"
  - rationale: 有带日期、路径范围与排除条件的度量，属可核实证据；但该度量说明依赖边几乎不存在，并不能直接证明失败链已发生，且正文自述缺少具体事故。

## 目标可验收 (goal_verifiability)

- **goal.verifiable** [core, level 2] — 2.1 每个 Goal 是否都有验收方式？
  - evidence `2.1 目标`: "验收：'TestMissingTargetBlocks'。"
  - normalized: "验收:'TestMissingTargetBlocks'。"
  - evidence `2.1 目标`: "验收：'TestJSONSchemaAndStats'、'TestJSONStatsCounts'、'TestSortFindingsOrder'、'TestJSONByteStableAcrossRuns'。"
  - normalized: "验收:'TestJSONSchemaAndStats'、'TestJSONStatsCounts'、'TestSortFindingsOrder'、'TestJSONByteStableAcrossRuns'。"
  - rationale: G1–G5 每条目标后都写明验收测试，抽查 G1 与 G5 可见验收方式明确。
- **goal.quantified** [core, level 2] — 2.1 是否至少有一个 Goal 有数字/单位/明确阈值？
  - evidence `2.1 目标`: "零受管文档在任一 phase 都退出 2"
  - normalized: "零受管文档在任一 phase 都退出 2"
  - evidence `2.1 目标`: "phase 1 对基线外的新违规退出 1"
  - normalized: "phase 1 对基线外的新违规退出 1"
  - rationale: G3、G4 以 phase 与退出码给出明确阈值，满足至少一个 Goal 有数字或明确阈值。
- **goal.count** [extended, level 1] — 2.1 Goals 是否控制在 3–6 个？
  - evidence `2.1 目标`: "**G1 引用存在**"
  - normalized: "**G1 引用存在**"
  - evidence `2.1 目标`: "**G5 稳定输出**"
  - normalized: "**G5 稳定输出**"
  - rationale: 2.1 共 5 个 Goal，处于 3–6 个区间，但未按架构/性能/可靠性分类，故为 1。

## 非目标边界 (non_goal_boundary)

- **nongoal.labeled** [core, level 2] — 2.2 每条 Non-Goal 是否标了本次不做或 Deferred to Phase 2？
  - evidence `2.2 非目标`: "**NG1** 必填字段、status 枚举、状态流转、TTL、审批来源校验（Deferred to Phase 2）。"
  - normalized: "**NG1** 必填字段、status 枚举、状态流转、TTL、审批来源校验(Deferred to Phase 2)。"
  - rationale: NG1–NG7 每条都以「本次不做」或「Deferred to Phase 2」标注。
- **nongoal.not_dumping_ground** [extended, level 1] — 2.2 非目标是否能追溯到痛点或边界？
  - evidence `2.2 非目标`: "**NG2** 循环引用检测（Deferred to Phase 2）。"
  - normalized: "**NG2** 循环引用检测(Deferred to Phase 2)。"
  - rationale: 部分非目标（如循环引用、审批来源）能对应引用图痛点或 v1 边界，但条目未逐条说明追溯到哪条痛点/边界，另有分发、多仓库等属范围外事项。

## 架构自洽 (architecture)

- **arch.solution_shape** [core, level 2] — 第 3 章是否有拓扑或等价的模块/数据流描述？
  - evidence `3. 架构设计与核心方案`: "CLI[CLI flags] --> Load[load standard.json]"
  - normalized: "CLI[CLI flags] --> Load[load standard.json]"
  - evidence `3. 架构设计与核心方案`: "- **新增**：单二进制 'mdstandards'、'policy/' 策略包、CI 与 Agent hook 两个执行点。"
  - normalized: "- **新增**:单二进制 'mdstandards'、'policy/' 策略包、CI 与 Agent hook 两个执行点。"
  - rationale: 第 3 章给出 mermaid 拓扑，并用新增/改造/不变清单与模块职责补充，具备模块与数据流描述。
- **arch.topology** [extended, level 2] — 3.1 拓扑图是否让每个模块有一句职责？
  - evidence `3.1 总体架构`: "CLI[CLI flags] --> Load[load standard.json]"
  - normalized: "CLI[CLI flags] --> Load[load standard.json]"
  - evidence `3.1 总体架构`: "Cover --> Eval[evaluate Rego]"
  - normalized: "Cover --> Eval[evaluate Rego]"
  - rationale: 3.1 的拓扑图节点以职责短语标注每一步，读者可据此理解模块分工。
- **arch.change_set** [extended, level 2] — 第 3 章是否写清新增/改造/保持不变？
  - evidence `3. 架构设计与核心方案`: "- **新增**：单二进制 'mdstandards'、'policy/' 策略包、CI 与 Agent hook 两个执行点。"
  - normalized: "- **新增**:单二进制 'mdstandards'、'policy/' 策略包、CI 与 Agent hook 两个执行点。"
  - evidence `3. 架构设计与核心方案`: "- **改造**：无。v1 不修改既有文档的 frontmatter，也不改既有 PDS 校验器。"
  - normalized: "- **改造**:无。v1 不修改既有文档的 frontmatter,也不改既有 PDS 校验器。"
  - evidence `3. 架构设计与核心方案`: "- **保持不变**：文档目录结构、'depends_on' 语义、人工评审流程。"
  - normalized: "- **保持不变**:文档目录结构、'depends_on' 语义、人工评审流程。"
  - rationale: 3.1 明确分列新增、改造（无）与保持不变三类，覆盖要求。
- **arch.rationale** [extended, level 1] — 关键选型是否有理由，详细对比指向第 7 章？
  - evidence `3.4 关键选型与约束`: "- **策略即数据**：受管范围与可依赖状态写在 'standard.json'，不写死在程序里。"
  - normalized: "- **策略即数据**:受管范围与可依赖状态写在 'standard.json',不写死在程序里。"
  - rationale: 3.4 为关键选型给出理由，但没有把详细对比指向第 7 章，故为 1。

## 详细设计完整 (detailed_design)

- **design.contract** [core, level 2] — 第 4 章是否至少有一份数据模型或接口契约？
  - evidence `4. 详细设计`: "JSON 输出契约：'schema_version'、'findings'、'extract_errors'、'stats'。"
  - normalized: "JSON 输出契约:'schema_version'、'findings'、'extract_errors'、'stats'。"
  - rationale: 第 4 章 4.3 提供数据模型字段表，并给出 JSON 输出契约，满足至少一份数据模型或接口契约。
- **design.data_model** [extended, level 2] — 4.3 数据模型是否有字段类型、必填与校验？
  - evidence `4.3 数据模型与契约`: "\| Standard \| standard_version \| string \| 是 \| 当前只接受 '"1"' \|"
  - normalized: "\| Standard \| standard_version \| string \| 是 \| 当前只接受 '"1"' \|"
  - evidence `4.3 数据模型与契约`: "\| Document \| status \| string \| 否 \| 缺失按不可依赖处理 \|"
  - normalized: "\| Document \| status \| string \| 否 \| 缺失按不可依赖处理 \|"
  - rationale: 字段表逐项给出类型与必填，并在说明列给出取值约束或缺失处理。
- **design.interface** [extended, level 2] — 4.4 接口是否有请求/响应/状态码或方法签名？
  - evidence `4.4 接口与协议`: "\| '-phase' \| '1' \| '0' 只报告、'1' 基线 + 阻断、'2' 完全阻断 \|"
  - normalized: "\| '-phase' \| '1' \| '0' 只报告、'1' 基线 + 阻断、'2' 完全阻断 \|"
  - evidence `4.4 接口与协议`: "- 'json'：机器可读报告，供 CI 归档与下游消费。"
  - normalized: "- 'json':机器可读报告,供 CI 归档与下游消费。"
  - rationale: 4.4 给出 CLI 参数（请求契约）与三种输出协议（响应格式），接口输入输出可核对。
- **design.sequence** [extended, level 2] — 4.5 复杂流程是否有时序图或等价分步说明？
  - evidence `4.5 核心流程与时序`: "sequenceDiagram"
  - normalized: "sequenceDiagram"
  - evidence `4.5 核心流程与时序`: "alt 阻断阶段且未降级"
  - normalized: "alt 阻断阶段且未降级"
  - rationale: 4.5 提供时序图，并用 alt 分支覆盖阻断阶段与降级两种路径。
- **design.errors** [extended, level 2] — 4.6 每个错误码是否有触发场景与处理策略？
  - evidence `4.6 异常处理与错误码`: "\| 'refs.missing_target' \| 'depends_on' 目标不在受管语料中 \| 计为 finding；阻断阶段退出 1 \|"
  - normalized: "\| 'refs.missing_target' \| 'depends_on' 目标不在受管语料中 \| 计为 finding;阻断阶段退出 1 \|"
  - rationale: 4.6 规则表逐条列出规则 ID、触发场景与处理策略，并补充退出码表。

## NFR 可度量 (nfr)

- **nfr.measurable** [core, level 1] — 第 5 章每条 NFR 是否有数字、单位或验收方式？
  - evidence `5. 非功能性需求`: "- **NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒。"
  - normalized: "- **NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒。"
  - evidence `5. 非功能性需求`: "- **NFR-004** 构建目标为 darwin-arm64 与 linux-amd64。"
  - normalized: "- **NFR-004** 构建目标为 darwin-arm64 与 linux-amd64。"
  - rationale: 部分 NFR 有明确数字或单位，但 NFR-004/005/006/009 在第 5 章内没有数字、单位或明确验收方式，未达到每条都可度量。
- **nfr.coverage** [core, level 2] — 第 5 章是否覆盖性能、可靠性、可观测性、安全，或说明不适用？
  - evidence `5. 非功能性需求`: "**NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒。"
  - normalized: "**NFR-001** 单次全量检查在 1000 篇受管文档上应 < 10 秒。"
  - evidence `5. 非功能性需求`: "**NFR-005** 闸门必须 fail-closed"
  - normalized: "**NFR-005** 闸门必须 fail-closed"
  - evidence `5. 非功能性需求`: "**NFR-007** stdout 输出结构化报告"
  - normalized: "**NFR-007** stdout 输出结构化报告"
  - evidence `5. 非功能性需求`: "**NFR-008** frontmatter 是不可信输入"
  - normalized: "**NFR-008** frontmatter 是不可信输入"
  - evidence `5. 非功能性需求`: "长时间运行的指标、Trace 与告警不适用"
  - normalized: "长时间运行的指标、Trace 与告警不适用"
  - rationale: 第 5 章分别覆盖性能、可靠性、可观测性与安全，并对不适用项给出理由。
- **nfr.honesty** [core, level 2] — 未实测的指标是否标了 [estimated] 或 未实现？
  - evidence `5. 非功能性需求`: "2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 38ms"
  - normalized: "2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 38ms"
  - evidence `5. 非功能性需求`: "2026-11-03 实测 p95 9ms（20 次，Apple M2 / macOS 27.0.1 / arm64）"
  - normalized: "2026-11-03 实测 p95 9ms(20 次,Apple M2 / macOS 27.0.1 / arm64)"
  - rationale: 性能数字均标注实测日期与环境，未见把未实测指标写成已达成；超时路径也说明用注入 deadline 验证。

## 上线可回滚 (rollout)

- **rollout.release** [core, level 2] — 6.1 是否有发布策略与灰度节奏？
  - evidence `6.1 发布与灰度`: "晋级条件：连续两周的新增违规数为 0，才允许 0 → 1；基线为空才允许 1 → 2。"
  - normalized: "晋级条件:连续两周的新增违规数为 0,才允许 0 → 1;基线为空才允许 1 → 2。"
  - rationale: 6.1 给出阶段 0/1/2 的发布策略与可量化的晋级条件，构成灰度节奏。
- **rollout.migration** [core, level 2] — 6.2 是否有迁移/兼容方案，或说明无需迁移？
  - evidence `6.2 数据迁移与兼容性`: "- 不迁移历史数据；存量违规用 baseline 记录，不修改旧文档来凑合规。"
  - normalized: "- 不迁移历史数据;存量违规用 baseline 记录,不修改旧文档来凑合规。"
  - rationale: 明确说明无需迁移历史数据，并给出存量违规用 baseline 记录、版本变更需显式升级的兼容策略。
- **rollout.smoke** [core, level 2] — 6.3 是否有冒烟清单与观察指标？
  - evidence `6.3 上线验证与冒烟`: "观察指标：'findings' 数量、'extract_errors' 数量、'documents_matched'、'edges_evaluated'、退出码分布。"
  - normalized: "观察指标:'findings' 数量、'extract_errors' 数量、'documents_matched'、'edges_evaluated'、退出码分布。"
  - rationale: 6.3 有编号冒烟清单，并另列观察指标，两者齐备。
- **rollout.rollback_trigger** [core, level 2] — 6.4 回滚触发条件是否是可观测数字？
  - evidence `6.4 回滚与熔断`: "- **性能**：全量检查 p95 > 10 秒连续 3 次。"
  - normalized: "- **性能**:全量检查 p95 > 10 秒连续 3 次。"
  - rationale: 6.4 的触发条件含误报率、错误次数、p95 与时序阈值，均可观测。
- **rollout.rollback_steps** [core, level 2] — 6.4 回滚步骤是否可执行？
  - evidence `6.4 回滚与熔断`: "1. 把 CI 的 '-phase' 从 1 调回 0，保留 'findings.json' 作为证据。"
  - normalized: "1. 把 CI 的 '-phase' 从 1 调回 0,保留 'findings.json' 作为证据。"
  - rationale: 回滚步骤为编号的可执行操作，包含调 phase、停用规则、修 policy、清基线等具体动作。

## 替代方案诚实 (alternatives)

- **alts.status_quo** [core, level 2] — 第 7 章是否有 Status Quo 及其成本/风险？
  - evidence `7. 替代方案与权衡`: "- **淘汰原因**：Agent 产出是指数的，人工吞吐固定；失败是静默的，抽查只能覆盖抽样，无法在合并前拦住引用图缺陷。"
  - normalized: "- **淘汰原因**:Agent 产出是指数的,人工吞吐固定;失败是静默的,抽查只能覆盖抽样,无法在合并前拦住引用图缺陷。"
  - rationale: 7.1 以 Status Quo 为方案 A，并说明人工抽查的吞吐与静默失败成本/风险。
- **alts.rejected** [core, level 2] — 第 7 章是否至少有一个被否决方案，且淘汰原因对应当前约束？
  - evidence `7. 替代方案与权衡`: "- **淘汰原因**：失败模式是文档引用图，不是单个文件的字段格式；单文件检查抓不到“目标不存在或未批准”。"
  - normalized: "- **淘汰原因**:失败模式是文档引用图,不是单个文件的字段格式;单文件检查抓不到“目标不存在或未批准”。"
  - rationale: 7.2/7.3 给出被否决方案及理由，淘汰原因直接对应引用图这一核心约束。
- **alts.matrix** [core, level 2] — 第 7 章是否有与约束对应的权衡矩阵？
  - evidence `7. 替代方案与权衡`: "\| 能发现引用图缺陷 \| 否 \| 否 \| 是 \| 是 \|"
  - normalized: "\| 能发现引用图缺陷 \| 否 \| 否 \| 是 \| 是 \|"
  - rationale: 7.5 权衡矩阵维度包含引用图发现能力、策略独立评审、规则成本与运行时依赖，均与 v1 约束对应。
- **alts.consistency** [core, level 2] — 第 7 章结论是否与 Goals、架构选型一致？
  - evidence `7. 替代方案与权衡`: "- **理由**：规则与程序分离，策略可独立评审；单二进制无额外运行时；基线与分阶段可复用同一套求值结果。"
  - normalized: "- **理由**:规则与程序分离,策略可独立评审;单二进制无额外运行时;基线与分阶段可复用同一套求值结果。"
  - rationale: 第 7 章选择方案 D，与第 3 章「策略即数据/单二进制」及 G1–G5 的分阶段、稳定输出目标一致。

## 跨章自洽 (cross_consistency)

- **x.goals_trace** [core, level 2] — 后续章节是否逐条回应 Goals？
  - evidence `全文`: "- **REQ-010** 对每条依赖边，若目标路径不在语料中，应产出 'refs.missing_target'。"
  - normalized: "- **REQ-010** 对每条依赖边,若目标路径不在语料中,应产出 'refs.missing_target'。"
  - evidence `全文`: "- **REQ-030** 'phase = 0' 时，进程不因 finding 或 extract_error 阻断（退出 0），但 findings 与 extract_errors 仍应完整输出；零受管文档、策略加载失败等工具级错误仍按 REQ-020 与前置条件退出 2。"
  - normalized: "- **REQ-030** 'phase = 0' 时,进程不因 finding 或 extract_error 阻断(退出 0),但 findings 与 extract_errors 仍应完整输出;零受管文档、策略加载失败等工具级错误仍按 REQ-020 与前置条件退出 2。"
  - evidence `全文`: "- **REQ-041** findings 应按 '(rule, document, target)' 排序，保证字节稳定。"
  - normalized: "- **REQ-041** findings 应按 '(rule, document, target)' 排序,保证字节稳定。"
  - rationale: G1/G2 落到 REQ-010～012，G3/G4 落到 REQ-020～033，G5 落到 REQ-040～042/NFR-006，逐条有后续设计回应。
- **x.nfr_arch** [core, level 2] — NFR 与架构、选型是否一致？
  - evidence `全文`: "- **NFR-005** 闸门必须 fail-closed：零受管文档、策略加载失败在任何 phase 都不得返回成功；零解析文档、零依赖边在 'phase >= 1' 都不得返回成功；解析错误在阻断阶段且未触发降级时不得返回成功。覆盖断言不受降级影响。"
  - normalized: "- **NFR-005** 闸门必须 fail-closed:零受管文档、策略加载失败在任何 phase 都不得返回成功;零解析文档、零依赖边在 'phase >= 1' 都不得返回成功;解析错误在阻断阶段且未触发降级时不得返回成功。覆盖断言不受降级影响。"
  - evidence `全文`: "单二进制内嵌 OPA 引擎；策略包从仓库 'policy/' 读取，'policy/data/standard.json' 是受管范围与可依赖状态的唯一事实来源。"
  - normalized: "单二进制内嵌 OPA 引擎;策略包从仓库 'policy/' 读取,'policy/data/standard.json' 是受管范围与可依赖状态的唯一事实来源。"
  - rationale: NFR-005/006/009 与第 3 章的覆盖断言、稳定排序、HEAD 策略完整性设计一一对应，未见冲突。
- **x.rollback_release** [core, level 2] — 回滚与发布策略是否一致？
  - evidence `全文`: "1. 把 CI 的 '-phase' 从 1 调回 0，保留 'findings.json' 作为证据。"
  - normalized: "1. 把 CI 的 '-phase' 从 1 调回 0,保留 'findings.json' 作为证据。"
  - evidence `全文`: "晋级条件：连续两周的新增违规数为 0，才允许 0 → 1；基线为空才允许 1 → 2。"
  - normalized: "晋级条件:连续两周的新增违规数为 0,才允许 0 → 1;基线为空才允许 1 → 2。"
  - rationale: 回滚第一步即把 phase 从 1 调回 0，与 6.1 的分阶段晋级策略方向一致。

## 证据卫生 (evidence_hygiene)

- **x.no_placeholders** [core, level 2] — 是否没有把 无 / ... / 待补充 / TBD 当作实质回答？
  - evidence `全文`: "- **改造**：无。v1 不修改既有文档的 frontmatter，也不改既有 PDS 校验器。"
  - normalized: "- **改造**:无。v1 不修改既有文档的 frontmatter,也不改既有 PDS 校验器。"
  - rationale: 全文未把「无」「...」「待补充」「TBD」当作无解释的实质回答；少量「无」均带具体说明（如改造为无并解释原因）。
- **x.no_code** [core, level 2] — 是否没有堆砌实现代码、测试日志或偶发报错？
  - evidence `全文`: "'''mermaid"
  - normalized: "'''mermaid"
  - evidence `全文`: "- **REQ-010** 对每条依赖边，若目标路径不在语料中，应产出 'refs.missing_target'。"
  - normalized: "- **REQ-010** 对每条依赖边,若目标路径不在语料中,应产出 'refs.missing_target'。"
  - rationale: 代码块仅为 mermaid 图，其余为实现要求、表格与测试名，没有 Go 实现代码、测试日志或偶发报错堆砌。
- **x.numbers_sourced** [core, level 2] — 关键数字是否有出处，或标了 [estimated]？
  - evidence `全文`: "411 篇 Markdown"
  - normalized: "411 篇 Markdown"
  - evidence `全文`: "2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 38ms"
  - normalized: "2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 38ms"
  - evidence `全文`: "2026-11-03 实测 p95 9ms（20 次，Apple M2 / macOS 27.0.1 / arm64）"
  - normalized: "2026-11-03 实测 p95 9ms(20 次,Apple M2 / macOS 27.0.1 / arm64)"
  - rationale: 关键度量注明日期、环境、样本或命令；设计阈值属规格值，无需外部出处。
- **x.no_fabricated_refs** [core, level 2] — 是否没有捏造 SPEC 引用、depends_on 目标或 M<N>？
  - evidence `全文`: "docs/architecture/mdstandards-architecture.md"
  - normalized: "docs/architecture/mdstandards-architecture.md"
  - evidence `全文`: "3. **循环引用未检测**：与 PDS-STD-004 §8.5.6 存在偏离。"
  - normalized: "3. **循环引用未检测**:与 PDS-STD-004 §8.5.6 存在偏离。"
  - rationale: 引用均为具体文档路径或 PDS 条款，前后一致，未发现 M<N> 或明显捏造的 depends_on 目标。

## Audit failures

None.
