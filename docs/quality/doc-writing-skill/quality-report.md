# Spec Quality Report

- Target: `docs/designs/doc-writing-skill.md`
- Rubric: 1.0.0 (`acb806c85c7a357f1cdfeaea80a8dc951333aefe5b630fe8576e0e0ce7f25947`)
- Target SHA-256: `38be1d04747b3edfe78415a0a7ce69898753fa0bd6bdd423543eab97756233ce`
- Reviewer session: `reviewer-design-20261102`

**Headline score (core): 15.0**  |  Diagnostic score (all): 15.0

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
| 跨章自洽 | 50.0 | 50.0 | 2 |
| 证据卫生 | 91.7 | 91.7 | 3 |

## 背景落地 (background)

- **bg.current_state** [core, level 0] — 1.1 是否写清了模块职责、核心流程与当前技术栈？
  - rationale: 文档没有「1.1 系统现状」章节；现状只散落在 Problem Statement 与 What already exists，未集中写模块职责、核心流程与当前技术栈。
  - missing: 缺少「1.1 系统现状」：应集中说明现有规范/工具链的模块职责、核心调用链与数据流，以及 Go、mdstandards 等技术栈。
- **bg.failure_chain** [core, level 0] — 1.2 是否有从触发到后果的分步失败链？
  - rationale: 没有「1.2 核心技术痛点」章节；Problem Statement 只描述了产出不一致和评审顺序，没有从触发到后果的分步失败链。
  - missing: 缺少「1.2 核心技术痛点」：应给出编号的失败链，从 Agent 读取规范不稳定、关键问题未回答到产出解决错问题的方案。
- **bg.trigger** [core, level 0] — 1.3 是否解释了为什么现在做？
  - rationale: 文档没有「1.3 重构/设计动机」章节；虽然 Problem Statement 说明现有规范不稳定，但没有独立的触发点与为什么现在做的时间理由。
  - missing: 缺少「1.3 重构/设计动机」：应写出具体触发事件或时间点，并解释为什么现在迁移规范而不是以后。
- **bg.non_goals** [core, level 0] — 1.4 是否写了不做的后果与本次非目标？
  - rationale: 文档没有「1.4 非目标」章节；非目标散落在 NOT in scope 与 Constraints，且未说明不做这些项的后果。
  - missing: 缺少「1.4 非目标」：应集中列出本次不做的范围，并说明不做会带来的后果。

## 痛点证据 (pain_evidence)

- **bg.pain_evidence** [core, level 0] — 1.2 是否有至少一条可核实的证据（事故、度量或引述）？
  - rationale: 没有「1.2 核心技术痛点」章节；用户引述只出现在 Problem Statement，未作为痛点章节中与失败链对应的证据。
  - missing: 缺少「1.2 核心技术痛点」中的可核实证据：应把事故、度量或带出处的引述放入痛点章节，并与失败链逐步对应。

## 目标可验收 (goal_verifiability)

- **goal.verifiable** [core, level 0] — 2.1 每个 Goal 是否都有验收方式？
  - rationale: 文档没有「2.1 目标」章节，也没有逐条 Goal；Success Criteria 是验收项，不是带验收方式的 Goal 清单。
  - missing: 缺少「2.1 目标」：应列出每个 Goal，并为每个 Goal 附上测试、指标或人工检查等验收方式。
- **goal.quantified** [core, level 0] — 2.1 是否至少有一个 Goal 有数字/单位/明确阈值？
  - rationale: 没有「2.1 目标」章节，因此不存在带数字、单位或明确阈值的 Goal。
  - missing: 缺少「2.1 目标」中的量化 Goal：至少一个 Goal 应有数字、单位或明确阈值，如延迟、行数或覆盖率。
- **goal.count** [extended, level 0] — 2.1 Goals 是否控制在 3–6 个？
  - rationale: 文档没有 Goals 列表，无法按 3–6 个计数。
  - missing: 缺少「2.1 目标」的 Goal 列表：应给出 3–6 个 Goal，并按架构、性能、可靠性等分类。

## 非目标边界 (non_goal_boundary)

- **nongoal.labeled** [core, level 0] — 2.2 每条 Non-Goal 是否标了本次不做或 Deferred to Phase 2？
  - rationale: 没有「2.2 非目标」章节；NOT in scope 的条目虽有 deferred/rejected 标注，但未组织成该章节内逐条标注的 Non-Goal 列表。
  - missing: 缺少「2.2 非目标」：应逐条列出 Non-Goal，并标注「本次不做」「Deferred to Phase 2」或「永远不做」。
- **nongoal.not_dumping_ground** [extended, level 0] — 2.2 非目标是否能追溯到痛点或边界？
  - rationale: 没有「2.2 非目标」章节，无法在该章节内逐条追溯非目标与痛点或边界的对应关系。
  - missing: 缺少「2.2 非目标」：应逐条说明该非目标对应的痛点或架构边界。

## 架构自洽 (architecture)

- **arch.solution_shape** [core, level 0] — 第 3 章是否有拓扑或等价的模块/数据流描述？
  - rationale: 文档没有「第 3 章 架构设计与核心方案」；目录结构与 Flow diagram 散落在推荐方案和 Eng Review，未构成该章拓扑或数据流。
  - missing: 缺少「3. 架构设计与核心方案」：应给出技能、检查器与数据流的拓扑图或等价模块描述。
- **arch.topology** [extended, level 0] — 3.1 拓扑图是否让每个模块有一句职责？
  - rationale: 没有「3.1 总体架构」章节；目录树只列文件，Flow diagram 只列阶段，未给每个模块一句职责。
  - missing: 缺少「3.1 总体架构」：应提供拓扑图，并为每个模块写一句职责。
- **arch.change_set** [extended, level 0] — 第 3 章是否写清新增/改造/保持不变？
  - rationale: 没有第 3 章，未在该章集中写清新增、改造与保持不变的部分。
  - missing: 缺少第 3 章的新增/改造/不变清单：应分别列出新增文件、改造文件和保持不变的部分。
- **arch.rationale** [extended, level 0] — 关键选型是否有理由，详细对比指向第 7 章？
  - rationale: 没有「3.4 关键选型与约束」章节；Approaches Considered 有否决理由，但不在该章节，也未把详细对比指向第 7 章。
  - missing: 缺少「3.4 关键选型与约束」：应写关键选型理由，并把详细对比指向第 7 章。

## 详细设计完整 (detailed_design)

- **design.contract** [core, level 0] — 第 4 章是否至少有一份数据模型或接口契约？
  - rationale: 文档没有「第 4 章 详细设计」；routes.json 契约与 check schema 在推荐方案和检查器章节，未作为详细设计契约呈现。
  - missing: 缺少「4. 详细设计」：应至少给出数据模型或接口契约，如 routes.json、check JSON 的字段与语义。
- **design.data_model** [extended, level 0] — 4.3 数据模型是否有字段类型、必填与校验？
  - rationale: 没有「4.3 数据模型与契约」章节；routes.json 只有示例字段，没有字段类型、必填与校验规则。
  - missing: 缺少「4.3 数据模型与契约」：应给出字段表或 schema，包含类型、必填与校验。
- **design.interface** [extended, level 0] — 4.4 接口是否有请求/响应/状态码或方法签名？
  - rationale: 没有「4.4 接口与协议」章节；doccheck CLI 模式与退出码散落在检查器章节，未在接口章节整理请求、响应与状态码。
  - missing: 缺少「4.4 接口与协议」：应集中给出接口签名、请求/响应结构或状态码。
- **design.sequence** [extended, level 0] — 4.5 复杂流程是否有时序图或等价分步说明？
  - rationale: 没有「4.5 核心流程与时序」章节；路由算法与 Flow diagram 分散，未在该章节给出含异常路径的时序或分步说明。
  - missing: 缺少「4.5 核心流程与时序」：应给出复杂流程的时序图或编号步骤，并覆盖异常路径。
- **design.errors** [extended, level 0] — 4.6 每个错误码是否有触发场景与处理策略？
  - rationale: 没有「4.6 异常处理与错误码」章节；exit 0/1/2 在外部闸门章节有触发与处理，但未作为该章节的错误码表。
  - missing: 缺少「4.6 异常处理与错误码」：应集中列出每个错误码或退出码的触发场景与处理策略。

## NFR 可度量 (nfr)

- **nfr.measurable** [core, level 0] — 第 5 章每条 NFR 是否有数字、单位或验收方式？
  - rationale: 文档没有「第 5 章 非功能性需求」；只有 Performance review 中的 <10s 扫描预算，未形成逐条 NFR。
  - missing: 缺少「5. 非功能性需求」：应逐条给出性能、可靠性、可观测性、安全等 NFR 的数字、单位或验收方式。
- **nfr.coverage** [core, level 0] — 第 5 章是否覆盖性能、可靠性、可观测性、安全，或说明不适用？
  - rationale: 没有第 5 章；Performance review 只覆盖性能，未覆盖可靠性、可观测性与安全，也未说明不适用。
  - missing: 缺少「5. 非功能性需求」：应覆盖性能、可靠性、可观测性、安全，或明确说明哪类不适用。
- **nfr.honesty** [core, level 0] — 未实测的指标是否标了 [estimated] 或 未实现？
  - rationale: 没有第 5 章；<10s/1000 docs 的扫描预算是未实测数字，且未标 [estimated] 或 未实现。
  - missing: 缺少「5. 非功能性需求」的诚实标注：未实测指标应标 [estimated]，未实现项应标 未实现。

## 上线可回滚 (rollout)

- **rollout.release** [core, level 0] — 6.1 是否有发布策略与灰度节奏？
  - rationale: 没有「6.1 发布与灰度」章节；分发与 release workflow 描述了发布方式，但没有灰度节奏。
  - missing: 缺少「6.1 发布与灰度」：应给出发布策略与灰度、观察节奏。
- **rollout.migration** [core, level 0] — 6.2 是否有迁移/兼容方案，或说明无需迁移？
  - rationale: 没有「6.2 数据迁移与兼容性」章节；迁移步骤在分发与迁移章节，未作为该节方案。
  - missing: 缺少「6.2 数据迁移与兼容性」：应给出迁移/兼容方案，或明确说明无需迁移及理由。
- **rollout.smoke** [core, level 0] — 6.3 是否有冒烟清单与观察指标？
  - rationale: 没有「6.3 上线验证与冒烟」章节；CI 检查步骤散落，未形成冒烟清单与观察指标。
  - missing: 缺少「6.3 上线验证与冒烟」：应给出冒烟清单与观察指标。
- **rollout.rollback_trigger** [core, level 0] — 6.4 回滚触发条件是否是可观测数字？
  - rationale: 没有「6.4 回滚与熔断」章节；文档只要求生成的 Spec 写回滚触发条件，没有本设计自身的可观测回滚阈值。
  - missing: 缺少「6.4 回滚与熔断」：应给出可观测的回滚触发数字，如错误率、延迟或资源阈值。
- **rollout.rollback_steps** [core, level 0] — 6.4 回滚步骤是否可执行？
  - rationale: 没有「6.4 回滚与熔断」章节；没有发布或分发失败时的可执行回滚步骤，只有 Spec 状态回退规则。
  - missing: 缺少「6.4 回滚与熔断」：应给出编号或命令形式的可执行回滚步骤，不能只写「可以回滚」。

## 替代方案诚实 (alternatives)

- **alts.status_quo** [core, level 0] — 第 7 章是否有 Status Quo 及其成本/风险？
  - rationale: 没有「第 7 章 替代方案与权衡」；Problem Statement 描述了现状与代价，但未把 Status Quo 作为替代方案列出。
  - missing: 缺少「7. 替代方案与权衡」中的 Status Quo：应列出维持现状的成本与风险。
- **alts.rejected** [core, level 0] — 第 7 章是否至少有一个被否决方案，且淘汰原因对应当前约束？
  - rationale: 没有第 7 章；Approach A/C 有否决理由，但不在该章节，无法按该节要求核对淘汰原因与约束的对应。
  - missing: 缺少「7. 替代方案与权衡」：应至少给出一个被否决方案，并说明淘汰原因如何对应当前约束。
- **alts.matrix** [core, level 0] — 第 7 章是否有与约束对应的权衡矩阵？
  - rationale: 没有第 7 章，也没有与约束对应的权衡矩阵；Approach A 只散落 Effort/Risk 两项。
  - missing: 缺少「7. 替代方案与权衡」的权衡矩阵：应给出与约束对应的维度矩阵。
- **alts.consistency** [core, level 0] — 第 7 章结论是否与 Goals、架构选型一致？
  - rationale: 没有第 7 章；推荐方案 B 与后续 Eng Review 的 Spec-only v1 范围相互矛盾，无法形成一致的结论。
  - missing: 缺少「7. 替代方案与权衡」：应给出与 Goals、架构选型一致的结论，并解决 B 与 Spec-only v1 的冲突。

## 跨章自洽 (cross_consistency)

- **x.goals_trace** [core, level 1] — 后续章节是否逐条回应 Goals？
  - evidence `Success Criteria`: "- 生成一份真实的里程碑纪要，通过 'checks/milestone.json'：五章齐全、决策带理由、In/Out Scope 成对、Backlog 表有 Issue 列且非空、DoD 可观测。In Scope 与 Backlog 行的逐项对应由评审确认。"
  - normalized: "- 生成一份真实的里程碑纪要,通过 'checks/milestone.json':五章齐全、决策带理由、In/Out Scope 成对、Backlog 表有 Issue 列且非空、DoD 可观测。In Scope 与 Backlog 行的逐项对应由评审确认。"
  - evidence `NOT in scope`: "- Milestone route, 'routes.json', project profile, and five-mode 'tools/doccheck': deferred to v2, because the milestone corpus is not mapped and the router has no deterministic executor (S4, S5)."
  - normalized: "- Milestone route, 'routes.json', project profile, and five-mode 'tools/doccheck': deferred to v2, because the milestone corpus is not mapped and the router has no deterministic executor (S4, S5)."
  - rationale: Success Criteria 中的里程碑与 router/doccheck 验收项在后续章节没有对应交付，且被 NOT in scope 延后；仅有部分验收项能追溯到设计。
- **x.nfr_arch** [core, level 2] — NFR 与架构、选型是否一致？
  - evidence `4. Performance review`: "One full 'mdstandards' scan per generated spec. At 1000 docs the NFR budget is <10s per scan; the Spec-only flow runs phase 0 optionally and phase 1 once."
  - normalized: "One full 'mdstandards' scan per generated spec. At 1000 docs the NFR budget is <10s per scan; the Spec-only flow runs phase 0 optionally and phase 1 once."
  - evidence `Constraints`: "- 不引入新的运行时依赖；检查器用 Go 编写，复用已有工具链。"
  - normalized: "- 不引入新的运行时依赖;检查器用 Go 编写,复用已有工具链。"
  - rationale: 性能预算按每次生成一次 mdstandards 扫描计算，与复用 Go 工具链、不新增运行时依赖的选型一致。
- **x.rollback_release** [core, level 0] — 回滚与发布策略是否一致？
  - rationale: 文档有 Distribution Plan 与 release workflow 的分发安排，但没有任何发布回滚步骤或触发条件，二者无法对应。
  - missing: 缺少与发布策略配套的回滚方案：应说明发布失败或安装后发现问题时如何撤回、回退，并让回滚步骤对应 release 步骤。

## 证据卫生 (evidence_hygiene)

- **x.no_placeholders** [core, level 2] — 是否没有把 无 / ... / 待补充 / TBD 当作实质回答？
  - evidence `检查器 tools/doccheck`: "5. 'doccheck -placeholders <path>...'：检查模板与生成文档里的 '{{...}}'、整行 '<...>'、整行 '...' 残留。"
  - normalized: "5. 'doccheck -placeholders <path>...':检查模板与生成文档里的 '{{...}}'、整行 '<...>'、整行 '...' 残留。"
  - rationale: 全文没有把 无、待补充、TBD 当作实质回答；'...' 和 '<...>' 只作为占位符语法出现，并被列为残留检查对象。
- **x.no_code** [core, level 2] — 是否没有堆砌实现代码、测试日志或偶发报错？
  - evidence `routes.json 契约`: "      "metadata": "yaml_frontmatter""
  - normalized: ""metadata": "yaml_frontmatter""
  - evidence `检查器 tools/doccheck`: "{"items":[{"id":"pain_evidence","status":"pass\|fail","locator":"1.2","expected":"...","actual":"..."}],"summary":{"pass":0,"fail":1,"advisory":0},"exit_code":1}"
  - normalized: "{"items":[{"id":"pain_evidence","status":"pass\|fail","locator":"1.2","expected":"...","actual":"..."}],"summary":{"pass":0,"fail":1,"advisory":0},"exit_code":1}"
  - rationale: 代码块都是目录树、JSON/YAML 契约与检查器输出 schema，没有实现代码、测试日志或偶发报错堆砌。
- **x.numbers_sourced** [core, level 1] — 关键数字是否有出处，或标了 [estimated]？
  - evidence `Reality rows`: "REALITY: numbers pass 820 guide lines measured with 'wc -l'; 14 files counted; 2 components counted; 23 mapping rows counted 'docs/designs/doc-writing-skill.md:17'"
  - normalized: "REALITY: numbers pass 820 guide lines measured with 'wc -l'; 14 files counted; 2 components counted; 23 mapping rows counted 'docs/designs/doc-writing-skill.md:17'"
  - evidence `4. Performance review`: "At 1000 docs the NFR budget is <10s per scan"
  - normalized: "At 1000 docs the NFR budget is <10s per scan"
  - rationale: 820 行、23 条映射等数字在 Reality rows 标了测量方式，但 1000 docs 下 <10s 的扫描预算是未实测数字，未标 [estimated]。
- **x.no_fabricated_refs** [core, level 2] — 是否没有捏造 SPEC 引用、depends_on 目标或 M<N>？
  - evidence `Reality rows`: "\| P7 \| target-repo has no policy bundle \| 'ls <workspace>/target-repo/policy' (absent) \| verified \|"
  - normalized: "\| P7 \| target-repo has no policy bundle \| 'ls <workspace>/target-repo/policy' (absent) \| verified \|"
  - evidence `Spec 模块（references/spec.md）`: "- 'depends_on' 只写真实存在的受管文档；写入前确认目标文件存在，并且状态可依赖。"
  - normalized: "- 'depends_on' 只写真实存在的受管文档;写入前确认目标文件存在,并且状态可依赖。"
  - rationale: 引用的 main.go、policy 路径由 Reality rows 逐条核实；depends_on 规则要求只写真实存在且可依赖的文档，全文没有具体 SPEC 编号或捏造的 M<N> 引用。

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
