# Spec Quality Report

- Target: `testdata/example-repo/docs/specs/example-spec.md`
- Rubric: 1.0.0 (`acb806c85c7a357f1cdfeaea80a8dc951333aefe5b630fe8576e0e0ce7f25947`)
- Target SHA-256: `eb17265c5d53537c5cb7b062a2914939238c85072ef80940f1db456731389db8`
- Reviewer session: `reviewer-legacy-20261102`

**Headline score (core): 12.0**  |  Diagnostic score (all): 12.0

> The auditor verifies evidence presence, section placement, and score arithmetic. It does not judge whether evidence is sufficient or true.

> Normalization covers full-width ASCII, NBSP, and ideographic space; it is not full Unicode NFKC.

| Dimension | Headline | Diagnostic | Weight |
|---|---:|---:|---:|
| 背景落地 | 0.0 | 0.0 | 2 |
| 痛点证据 | 0.0 | 0.0 | 3 |
| 目标可验收 | 0.0 | 0.0 | 3 |
| 非目标边界 | 0.0 | 0.0 | 2 |
| 架构自洽 | 0.0 | 0.0 | 2 |
| 详细设计完整 | 0.0 | 0.0 | 2 |
| NFR 可度量 | 0.0 | 0.0 | 2 |
| 上线可回滚 | 0.0 | 0.0 | 2 |
| 替代方案诚实 | 0.0 | 0.0 | 2 |
| 跨章自洽 | 0.0 | 0.0 | 2 |
| 证据卫生 | 100.0 | 100.0 | 3 |

## 背景落地 (background)

- **bg.current_state** [core, level 0] — 1.1 是否写清了模块职责、核心流程与当前技术栈？
  - rationale: 文档正文只有一句依赖说明，没有系统现状章节，未写模块职责、核心流程或当前技术栈。
  - missing: 缺少「1.1 系统现状」：应说明模块职责、核心调用链与数据流，以及当前技术栈。
- **bg.failure_chain** [core, level 0] — 1.2 是否有从触发到后果的分步失败链？
  - rationale: 没有痛点章节，全文不存在从触发到后果的分步失败链。
  - missing: 缺少「1.2 核心技术痛点」：应给出编号的触发、中间失败到用户可见后果的因果链。
- **bg.trigger** [core, level 0] — 1.3 是否解释了为什么现在做？
  - rationale: 全文只有一句“应当通过检查”的说明，没有解释为什么现在做这项工作。
  - missing: 缺少「1.3 重构/设计动机」：应写明业务或技术触发点，以及为什么现在而不是以后做。
- **bg.non_goals** [core, level 0] — 1.4 是否写了不做的后果与本次非目标？
  - rationale: 文档没有非目标章节，既未列非目标，也未说明不做的后果。
  - missing: 缺少「1.4 非目标」：应给出本次明确不做的范围，并说明不做会带来什么后果。

## 痛点证据 (pain_evidence)

- **bg.pain_evidence** [core, level 0] — 1.2 是否有至少一条可核实的证据（事故、度量或引述）？
  - rationale: 没有痛点章节，也没有事故、度量或带出处的引述等可核实证据。
  - missing: 缺少「1.2 核心技术痛点」中的可核实证据：应给出事故时间/文件/后果、度量数字或带出处的引述，并与失败链对应。

## 目标可验收 (goal_verifiability)

- **goal.verifiable** [core, level 0] — 2.1 每个 Goal 是否都有验收方式？
  - rationale: 全文没有目标章节，没有任何 Goal，更无验收方式。
  - missing: 缺少「2.1 目标」：应为每个 Goal 附上测试、指标或人工检查等验收方式。
- **goal.quantified** [core, level 0] — 2.1 是否至少有一个 Goal 有数字/单位/明确阈值？
  - rationale: 文档没有任何 Goal，因此不存在带数字、单位或阈值的量化目标。
  - missing: 缺少「2.1 目标」中的量化目标：至少一个 Goal 应有数字、单位或明确阈值，如 P95 延迟或覆盖率。
- **goal.count** [extended, level 0] — 2.1 Goals 是否控制在 3–6 个？
  - rationale: 全文没有 Goal 列表，数量为 0，不在 3–6 个范围内。
  - missing: 缺少「2.1 目标」的 Goal 列表：应给出 3–6 个 Goal，并按架构/性能/可靠性等分类。

## 非目标边界 (non_goal_boundary)

- **nongoal.labeled** [core, level 0] — 2.2 每条 Non-Goal 是否标了本次不做或 Deferred to Phase 2？
  - rationale: 文档没有非目标章节，不存在任何带「本次不做」或「Deferred to Phase 2」标注的条目。
  - missing: 缺少「2.2 非目标」：每条 Non-Goal 都应标注「本次不做」或「Deferred to Phase 2」。
- **nongoal.not_dumping_ground** [extended, level 0] — 2.2 非目标是否能追溯到痛点或边界？
  - rationale: 没有非目标条目，无法判断其是否可追溯到痛点或边界。
  - missing: 缺少「2.2 非目标」：应逐条说明该非目标对应的痛点或架构边界。

## 架构自洽 (architecture)

- **arch.solution_shape** [core, level 0] — 第 3 章是否有拓扑或等价的模块/数据流描述？
  - rationale: 文档没有第 3 章，正文无拓扑图，也无模块职责或数据流描述。
  - missing: 缺少「3. 架构设计与核心方案」：应给出拓扑图，或等价的模块职责与数据流描述。
- **arch.topology** [extended, level 0] — 3.1 拓扑图是否让每个模块有一句职责？
  - rationale: 没有总体架构章节，既无拓扑图也无模块职责说明。
  - missing: 缺少「3.1 总体架构」：应提供可读的拓扑图，并为每个模块写一句职责。
- **arch.change_set** [extended, level 0] — 第 3 章是否写清新增/改造/保持不变？
  - rationale: 第 3 章不存在，文档未说明任何新增、改造或保持不变的部分。
  - missing: 缺少第 3 章的新增/改造/不变清单：应分别列出本次新增、改造与保持不变的部分。
- **arch.rationale** [extended, level 0] — 关键选型是否有理由，详细对比指向第 7 章？
  - rationale: 文档没有关键选型章节，没有任何选型理由，也未指向第 7 章。
  - missing: 缺少「3.4 关键选型与约束」：应说明关键选型理由，并把详细对比指向第 7 章。

## 详细设计完整 (detailed_design)

- **design.contract** [core, level 0] — 第 4 章是否至少有一份数据模型或接口契约？
  - rationale: 文档没有第 4 章，未提供任何数据模型、字段表或接口契约。
  - missing: 缺少「4. 详细设计」：至少应有一份数据模型（字段表/schema）或接口契约（请求/响应或方法签名）。
- **design.data_model** [extended, level 0] — 4.3 数据模型是否有字段类型、必填与校验？
  - rationale: 没有数据模型章节，不存在带类型、必填和校验的字段定义。
  - missing: 缺少「4.3 数据模型与契约」：应给出字段表或 schema，标明类型、必填与校验规则。
- **design.interface** [extended, level 0] — 4.4 接口是否有请求/响应/状态码或方法签名？
  - rationale: 没有接口章节，未定义请求、响应、状态码或方法签名。
  - missing: 缺少「4.4 接口与协议」：应定义每个接口的请求/响应结构、状态码或方法签名。
- **design.sequence** [extended, level 0] — 4.5 复杂流程是否有时序图或等价分步说明？
  - rationale: 文档没有第 4 章，也就不存在任何核心流程的时序图或分步说明。
  - missing: 缺少「4.5 核心流程与时序」：应给出时序图或编号分步说明，并覆盖异常路径。
- **design.errors** [extended, level 0] — 4.6 每个错误码是否有触发场景与处理策略？
  - rationale: 没有异常处理章节，未定义任何错误码、触发场景或处理策略。
  - missing: 缺少「4.6 异常处理与错误码」：应为每个错误码写明触发场景与处理策略。

## NFR 可度量 (nfr)

- **nfr.measurable** [core, level 0] — 第 5 章每条 NFR 是否有数字、单位或验收方式？
  - rationale: 文档没有第 5 章，没有任何 NFR，更无数字、单位或验收方式。
  - missing: 缺少「5. 非功能性需求」：每条 NFR 都应有数字、单位或明确验收方式，如 P95、QPS、SLA。
- **nfr.coverage** [core, level 0] — 第 5 章是否覆盖性能、可靠性、可观测性、安全，或说明不适用？
  - rationale: 第 5 章缺失，性能、可靠性、可观测性、安全四类均未覆盖，也未说明不适用。
  - missing: 缺少「5. 非功能性需求」：应覆盖性能、可靠性、可观测性、安全，或逐类说明不适用及理由。
- **nfr.honesty** [core, level 0] — 未实测的指标是否标了 [estimated] 或 未实现？
  - rationale: 文档没有任何指标，因此不存在已实测或未实测的标注，也无 [estimated] 或「未实现」标记。
  - missing: 缺少「5. 非功能性需求」：应给出指标并明确标注 [estimated] 或「未实现」，避免把未实测值写成已达成。

## 上线可回滚 (rollout)

- **rollout.release** [core, level 0] — 6.1 是否有发布策略与灰度节奏？
  - rationale: 文档没有发布与灰度章节，未写发布方式，也没有灰度节奏。
  - missing: 缺少「6.1 发布与灰度」：应写明蓝绿/金丝雀/Feature Flag 等发布策略及观察期与放量节奏。
- **rollout.migration** [core, level 0] — 6.2 是否有迁移/兼容方案，或说明无需迁移？
  - rationale: 没有迁移与兼容性章节，既无迁移方案，也未说明无需迁移。
  - missing: 缺少「6.2 数据迁移与兼容性」：应给出 Backfill/双写/影子表等方案及兼容策略，或明确说明无需迁移及理由。
- **rollout.smoke** [core, level 0] — 6.3 是否有冒烟清单与观察指标？
  - rationale: 没有上线验证章节，无冒烟步骤，也无观察指标。
  - missing: 缺少「6.3 上线验证与冒烟」：应给出可执行的冒烟清单与对应的观察指标。
- **rollout.rollback_trigger** [core, level 0] — 6.4 回滚触发条件是否是可观测数字？
  - rationale: 没有回滚与熔断章节，不存在任何可观测的回滚触发数字。
  - missing: 缺少「6.4 回滚与熔断」：应把回滚触发条件写成可观测数字，如错误率、延迟或资源阈值。
- **rollout.rollback_steps** [core, level 0] — 6.4 回滚步骤是否可执行？
  - rationale: 没有回滚与熔断章节，文档未给出任何回滚步骤或命令。
  - missing: 缺少「6.4 回滚与熔断」：应给出可执行的编号回滚步骤或命令，而不只是「可以回滚」。

## 替代方案诚实 (alternatives)

- **alts.status_quo** [core, level 0] — 第 7 章是否有 Status Quo 及其成本/风险？
  - rationale: 文档没有第 7 章，未描述 Status Quo，也未给出其成本或风险。
  - missing: 缺少「7. 替代方案与权衡」：应描述 Status Quo 及其成本、风险或触发阈值。
- **alts.rejected** [core, level 0] — 第 7 章是否至少有一个被否决方案，且淘汰原因对应当前约束？
  - rationale: 第 7 章缺失，没有任何被否决方案，也没有淘汰原因。
  - missing: 缺少「7. 替代方案与权衡」：至少应有一个被否决方案的思路、优点及对应当前约束的淘汰原因。
- **alts.matrix** [core, level 0] — 第 7 章是否有与约束对应的权衡矩阵？
  - rationale: 文档没有第 7 章，不存在任何权衡矩阵。
  - missing: 缺少「7. 替代方案与权衡」：应给出与当前约束维度对应的权衡矩阵。
- **alts.consistency** [core, level 0] — 第 7 章结论是否与 Goals、架构选型一致？
  - rationale: 第 7 章缺失，且全文没有 Goals 或架构选型可供对照，结论一致性无从谈起。
  - missing: 缺少「7. 替代方案与权衡」的结论：应明确所选方案，并使其与第 2 章 Goals、第 3 章架构选型一致。

## 跨章自洽 (cross_consistency)

- **x.goals_trace** [core, level 0] — 后续章节是否逐条回应 Goals？
  - rationale: 文档没有 Goals，也没有第 3/4/5 章，未回应任何目标。
  - missing: 缺少 Goals 与后续设计章节：应为每个 Goal 在第 3/4/5 章找到对应设计与验收。
- **x.nfr_arch** [core, level 0] — NFR 与架构、选型是否一致？
  - rationale: 文档既无 NFR 也无架构/选型，无法建立二者的一致性对应。
  - missing: 缺少第 3、5 章：应写出 NFR，并在架构与选型中说明如何满足这些指标。
- **x.rollback_release** [core, level 0] — 回滚与发布策略是否一致？
  - rationale: 文档既无发布/灰度策略也无回滚方案，二者无法对应。
  - missing: 缺少第 6 章发布与回滚内容：回滚步骤应与发布/灰度策略相互对应。

## 证据卫生 (evidence_hygiene)

- **x.no_placeholders** [core, level 2] — 是否没有把 无 / ... / 待补充 / TBD 当作实质回答？
  - evidence `全文`: "这份 spec 依赖同目录树上的一份已批准架构文档，因此应当通过检查。"
  - normalized: "这份 spec 依赖同目录树上的一份已批准架构文档,因此应当通过检查。"
  - rationale: 全文没有把 无、...、待补充、TBD 当作实质回答，正文仅有这一句完整陈述。
- **x.no_code** [core, level 2] — 是否没有堆砌实现代码、测试日志或偶发报错？
  - evidence `全文`: "这份 spec 依赖同目录树上的一份已批准架构文档，因此应当通过检查。"
  - normalized: "这份 spec 依赖同目录树上的一份已批准架构文档,因此应当通过检查。"
  - rationale: 全文没有代码块、测试日志或偶发报错，只有一段依赖说明文字。
- **x.numbers_sourced** [core, level 2] — 关键数字是否有出处，或标了 [estimated]？
  - evidence `全文`: "这份 spec 依赖同目录树上的一份已批准架构文档，因此应当通过检查。"
  - normalized: "这份 spec 依赖同目录树上的一份已批准架构文档,因此应当通过检查。"
  - rationale: not applicable: 文档没有提出任何性能、资源或覆盖率数字，因此没有需要标注出处或 [estimated] 的关键数字。
- **x.no_fabricated_refs** [core, level 2] — 是否没有捏造 SPEC 引用、depends_on 目标或 M<N>？
  - evidence `全文`: "docs/architecture/example-architecture.md"
  - normalized: "docs/architecture/example-architecture.md"
  - evidence `全文`: "这份 spec 依赖同目录树上的一份已批准架构文档，因此应当通过检查。"
  - normalized: "这份 spec 依赖同目录树上的一份已批准架构文档,因此应当通过检查。"
  - rationale: depends_on 指向的 docs/architecture/example-architecture.md 实际存在且 status 为 Approved；全文没有 SPEC 引用或 M<N>。

## Chapter notes

- section "1.1 系统现状" not found
- section "1.2 核心技术痛点" not found
- section "1.3 重构/设计动机" not found
- section "1.4 非目标" not found
- section "2.1 目标" not found
- section "2.2 非目标" not found
- section "3. 架构设计与核心方案" not found
- section "3.1 总体架构" not found
- section "3.4 关键选型与约束" not found
- section "4. 详细设计" not found
- section "4.3 数据模型与契约" not found
- section "4.4 接口与协议" not found
- section "4.5 核心流程与时序" not found
- section "4.6 异常处理与错误码" not found
- section "5. 非功能性需求" not found
- section "6.1 发布与灰度" not found
- section "6.2 数据迁移与兼容性" not found
- section "6.3 上线验证与冒烟" not found
- section "6.4 回滚与熔断" not found
- section "7. 替代方案与权衡" not found

## Audit failures

None.
