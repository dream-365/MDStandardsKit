# Spec Quality Report

- Target: `docs/specs/SPEC-002-idea-type/spec.md`
- Rubric: 1.0.0 (`acb806c85c7a357f1cdfeaea80a8dc951333aefe5b630fe8576e0e0ce7f25947`)
- Target SHA-256: `52181d85687c1b9cc525c35fd0ca1799b695ec284f158af541e4146e3865fcdf`
- Reviewer session: `reviewer-spec-002-20261103b`

**Headline score (core): 93.3**  |  Diagnostic score (all): 87.7

> The auditor verifies evidence presence, section placement, and score arithmetic. It does not judge whether evidence is sufficient or true.

> Normalization covers full-width ASCII, NBSP, and ideographic space; it is not full Unicode NFKC.

| Dimension | Headline | Diagnostic | Weight |
|---|---:|---:|---:|
| 背景落地 | 87.5 | 87.5 | 2 |
| 痛点证据 | 100.0 | 100.0 | 3 |
| 目标可验收 | 100.0 | 87.5 | 3 |
| 非目标边界 | 100.0 | 83.3 | 2 |
| 架构自洽 | 100.0 | 90.0 | 2 |
| 详细设计完整 | 100.0 | 75.0 | 2 |
| NFR 可度量 | 75.0 | 75.0 | 2 |
| 上线可回滚 | 83.3 | 83.3 | 2 |
| 替代方案诚实 | 87.5 | 87.5 | 2 |
| 跨章自洽 | 83.3 | 83.3 | 2 |
| 证据卫生 | 100.0 | 100.0 | 3 |

## 背景落地 (background)

- **bg.current_state** [core, level 1] — 1.1 是否写清了模块职责、核心流程与当前技术栈？
  - evidence `1.1 系统现状`: "'scripts/new' 负责建文件、生成 8 位大写 ID、校验标题非空、校验日期格式与 ID 唯一。"
  - normalized: "'scripts/new' 负责建文件、生成 8 位大写 ID、校验标题非空、校验日期格式与 ID 唯一。"
  - rationale: 写清了想法库结构、scripts/new 的职责与流程，但未交代当前技术栈（语言/框架/版本），不足以完整理解现状。
- **bg.failure_chain** [core, level 2] — 1.2 是否有从触发到后果的分步失败链？
  - evidence `1.2 核心技术痛点`: "2. 章节缺失、模板注释还在、状态写错、日期与目录对不上，没有任何工具发现。"
  - normalized: "2. 章节缺失、模板注释还在、状态写错、日期与目录对不上,没有任何工具发现。"
  - rationale: 1.2 以编号 1–3 给出从建文件到问题进入库后才被发现的完整失败链，符合分步因果要求。
- **bg.trigger** [core, level 2] — 1.3 是否解释了为什么现在做？
  - evidence `1.3 重构/设计动机`: "触发点是想法库开始承载越来越多的 Agent 生成内容，而现有 'scripts/new' 只管文件名和 ID。"
  - normalized: "触发点是想法库开始承载越来越多的 Agent 生成内容,而现有 'scripts/new' 只管文件名和 ID。"
  - rationale: 明确给出触发点（Agent 生成内容增多）与为什么现在做（现有工具只管文件名和 ID）。
- **bg.non_goals** [core, level 2] — 1.4 是否写了不做的后果与本次非目标？
  - evidence `1.4 非目标`: "不做结构闸门的后果见 1.2：空壳想法继续入库，内容低分的想法和人判断混在一起。"
  - normalized: "不做结构闸门的后果见 1.2:空壳想法继续入库,内容低分的想法和人判断混在一起。"
  - rationale: 先说明不做的后果，再以 5 条清单列出本次非目标。

## 痛点证据 (pain_evidence)

- **bg.pain_evidence** [core, level 2] — 1.2 是否有至少一条可核实的证据（事故、度量或引述）？
  - evidence `1.2 核心技术痛点`: "证据：截至 2026-11-02，'~/Workspace/ideas' 有 5 篇想法；在没有结构闸门之前，'scripts/new' 建出的文件只保证 title/date/id 三项。"
  - normalized: "证据:截至 2026-11-02,'~/Workspace/ideas' 有 5 篇想法;在没有结构闸门之前,'scripts/new' 建出的文件只保证 title/date/id 三项。"
  - evidence `1.2 核心技术痛点`: "2026-11-02 的冒烟测试对「万花筒」跑完 idea rubric，headline 25.0、diagnostic 35.0，证明内容评分链路能产出可审计的报告。"
  - normalized: "2026-11-02 的冒烟测试对「万花筒」跑完 idea rubric,headline 25.0、diagnostic 35.0,证明内容评分链路能产出可审计的报告。"
  - rationale: 证据带日期、仓库与具体覆盖范围，并直接对应失败链中「结构无人检查、内容无标准」的后果。

## 目标可验收 (goal_verifiability)

- **goal.verifiable** [core, level 2] — 2.1 每个 Goal 是否都有验收方式？
  - evidence `2.1 目标`: "验收：'TestIdeaConfigFailsClosed'、'TestIdeaValidRepo'。"
  - normalized: "验收:'TestIdeaConfigFailsClosed'、'TestIdeaValidRepo'。"
  - rationale: G1–G5 每个 Goal 后均附「验收：」，给出测试、实测或 make check 等验收方式。
- **goal.quantified** [core, level 2] — 2.1 是否至少有一个 Goal 有数字/单位/明确阈值？
  - evidence `2.1 目标`: "缺失、损坏、未知 type、未知键都 exit 2。"
  - normalized: "缺失、损坏、未知 type、未知键都 exit 2。"
  - rationale: G1 用退出码 2 作为明确阈值，满足至少一个 Goal 有数字/明确阈值的条件。
- **goal.count** [extended, level 1] — 2.1 Goals 是否控制在 3–6 个？
  - evidence `2.1 目标`: "**G5 spec 不回归**"
  - normalized: "**G5 spec 不回归**"
  - rationale: 共 G1–G5 五个 Goal，数量落在 3–6 内；但未按架构/性能/可靠性分类，故未达 level 2。

## 非目标边界 (non_goal_boundary)

- **nongoal.labeled** [core, level 2] — 2.2 每条 Non-Goal 是否标了本次不做或 Deferred to Phase 2？
  - evidence `2.2 非目标`: "**NG1** 想法状态流转校验（Deferred to Phase 2）。"
  - normalized: "**NG1** 想法状态流转校验(Deferred to Phase 2)。"
  - rationale: NG1–NG5 每条都带「Deferred to Phase 2」或「本次不做」标注。
- **nongoal.not_dumping_ground** [extended, level 1] — 2.2 非目标是否能追溯到痛点或边界？
  - evidence `2.2 非目标`: "**NG5** 内容分作为阻断门槛（本次不做；v0.2.0 只诊断）。"
  - normalized: "**NG5** 内容分作为阻断门槛(本次不做;v0.2.0 只诊断)。"
  - rationale: NG1/NG3/NG5 可追溯到现状或动机，但 NG2、NG4 与非目标边界关系弱，未做到每条可追溯。

## 架构自洽 (architecture)

- **arch.solution_shape** [core, level 2] — 第 3 章是否有拓扑或等价的模块/数据流描述？
  - evidence `3. 架构设计与核心方案`: "CLI[mdstandards] --> Config[.mdstandards.json]"
  - normalized: "CLI[mdstandards] --> Config[.mdstandards.json]"
  - rationale: 第 3 章含 mermaid 拓扑图，并在 3.2 给出各模块的新增/改造职责，构成等价的模块与数据流描述。
- **arch.topology** [extended, level 1] — 3.1 拓扑图是否让每个模块有一句职责？
  - evidence `3.1 总体架构`: "CLI[mdstandards] --> Config[.mdstandards.json]"
  - normalized: "CLI[mdstandards] --> Config[.mdstandards.json]"
  - rationale: 有可读的拓扑图，但图中只有模块名，未为每个模块写一句职责说明。
- **arch.change_set** [extended, level 2] — 第 3 章是否写清新增/改造/保持不变？
  - evidence `3. 架构设计与核心方案`: "- 单 'main' 包，不拆 'cmd/' + 'internal/'。"
  - normalized: "- 单 'main' 包,不拆 'cmd/' + 'internal/'。"
  - rationale: 3.2 分「新增 / 改造 / 保持不变」三节逐项列出变更集。
- **arch.rationale** [extended, level 2] — 关键选型是否有理由，详细对比指向第 7 章？
  - evidence `3.4 关键选型与约束`: "方案对比见第 7 章。"
  - normalized: "方案对比见第 7 章。"
  - rationale: 3.4 每项关键选型均给出理由（类型配置+专用 Rego、fail-closed 等），并指向第 7 章对比。

## 详细设计完整 (detailed_design)

- **design.contract** [core, level 2] — 第 4 章是否至少有一份数据模型或接口契约？
  - evidence `4. 详细设计`: "\| 'document_type' \| 语料级类型，idea 时为 'idea' \|"
  - normalized: "\| 'document_type' \| 语料级类型,idea 时为 'idea' \|"
  - rationale: 4.3 提供 idea 文档在 corpus 中的数据模型字段表，满足至少一份数据模型契约。
- **design.data_model** [extended, level 0] — 4.3 数据模型是否有字段类型、必填与校验？
  - rationale: 字段表只有「字段 / 含义」两列，没有字段类型、必填标记或校验规则。
  - missing: 4.3 的字段表缺少类型、必填与校验信息；需为每个字段补类型（如 string/array）、是否必填及校验约束后才算达标。
- **design.interface** [extended, level 2] — 4.4 接口是否有请求/响应/状态码或方法签名？
  - evidence `4.4 接口与协议`: "'-policy' 相对当前工作目录解析；idea 仓库必须显式传插件的 policy 目录，或设置 'MDSTANDARDS_POLICY_DIR'。输出沿用 'table\|json\|github'；退出码见 4.6。"
  - normalized: "'-policy' 相对当前工作目录解析;idea 仓库必须显式传插件的 policy 目录,或设置 'MDSTANDARDS_POLICY_DIR'。输出沿用 'table\|json\|github';退出码见 4.6。"
  - rationale: 给出 CLI 调用签名、输出格式与退出码引用，构成接口与协议说明。
- **design.sequence** [extended, level 2] — 4.5 复杂流程是否有时序图或等价分步说明？
  - evidence `4.5 核心流程与时序`: "4. 覆盖断言：零受管文档、'phase >= 1' 的零解析文档都退出 2；idea 类型不要求依赖边。"
  - normalized: "4. 覆盖断言:零受管文档、'phase >= 1' 的零解析文档都退出 2;idea 类型不要求依赖边。"
  - rationale: 4.5 用 1–6 编号步骤描述核心流程，并在第 4 步覆盖零文档/零解析等异常路径。
- **design.errors** [extended, level 1] — 4.6 每个错误码是否有触发场景与处理策略？
  - evidence `4.6 异常处理与错误码`: "\| 2 \| 配置/策略加载失败、未知 type、零受管文档、零解析文档、解析错误 \|"
  - normalized: "\| 2 \| 配置/策略加载失败、未知 type、零受管文档、零解析文档、解析错误 \|"
  - rationale: 4.6 列出退出码及触发场景，但未给出每个码的处理策略。

## NFR 可度量 (nfr)

- **nfr.measurable** [core, level 1] — 第 5 章每条 NFR 是否有数字、单位或验收方式？
  - evidence `5. 非功能性需求`: "**NFR-2 性能**：9 篇 idea 的 phase 1 在 1 秒内完成（2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s）；单文件上限沿用 1 MiB。"
  - normalized: "**NFR-2 性能**:9 篇 idea 的 phase 1 在 1 秒内完成(2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s);单文件上限沿用 1 MiB。"
  - rationale: 仅性能项有数字与单位；NFR-1/3/4/5 未给出数字或明确验收方式，故为部分达标。
- **nfr.coverage** [core, level 2] — 第 5 章是否覆盖性能、可靠性、可观测性、安全，或说明不适用？
  - evidence `5. 非功能性需求`: "**NFR-4 可观测性**：每条 finding 带 'rule'、'document'、'target'、'message'；内容报告带 rubric 版本、目标 SHA-256、reviewer 会话。"
  - normalized: "**NFR-4 可观测性**:每条 finding 带 'rule'、'document'、'target'、'message';内容报告带 rubric 版本、目标 SHA-256、reviewer 会话。"
  - rationale: 覆盖性能（NFR-2）、可靠性（NFR-1 确定性）、可观测性（NFR-4）与安全（NFR-3）。
- **nfr.honesty** [core, level 2] — 未实测的指标是否标了 [estimated] 或 未实现？
  - evidence `5. 非功能性需求`: "**NFR-2 性能**：9 篇 idea 的 phase 1 在 1 秒内完成（2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s）；单文件上限沿用 1 MiB。"
  - normalized: "**NFR-2 性能**:9 篇 idea 的 phase 1 在 1 秒内完成(2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s);单文件上限沿用 1 MiB。"
  - rationale: 量化指标标注为「实测」，全文未把未实测指标写成已达成，也无需 [estimated] 标注。

## 上线可回滚 (rollout)

- **rollout.release** [core, level 2] — 6.1 是否有发布策略与灰度节奏？
  - evidence `6.1 发布与灰度`: "3. 插件版本升到 '0.2.0'，release tag 必须与 'plugin.json' 一致；先对 'testdata/idea-repo' 冒烟，再对真实 ideas 仓库灰度。"
  - normalized: "3. 插件版本升到 '0.2.0',release tag 必须与 'plugin.json' 一致;先对 'testdata/idea-repo' 冒烟,再对真实 ideas 仓库灰度。"
  - rationale: 给出发布版本/tag 要求，以及先在 testdata 冒烟、再对真实仓库灰度的节奏。
- **rollout.migration** [core, level 2] — 6.2 是否有迁移/兼容方案，或说明无需迁移？
  - evidence `6.2 数据迁移与兼容性`: "- 不迁移既有想法文件；新增配置即可启用结构闸门。"
  - normalized: "- 不迁移既有想法文件;新增配置即可启用结构闸门。"
  - rationale: 明确说明不迁移既有文件及启用方式，属于「无需迁移 + 理由」的兼容方案。
- **rollout.smoke** [core, level 2] — 6.3 是否有冒烟清单与观察指标？
  - evidence `6.3 上线验证与冒烟`: "1. 'testdata/idea-repo' phase 1 退出 0。"
  - normalized: "1. 'testdata/idea-repo' phase 1 退出 0。"
  - rationale: 6.3 给出 3 条冒烟清单，并以退出 0、make check 全绿作为观察结果。
- **rollout.rollback_trigger** [core, level 1] — 6.4 回滚触发条件是否是可观测数字？
  - evidence `6.4 回滚与熔断`: "- 熔断触发：真实 ideas 仓库出现结构误报、插件 shim 或 release 冒烟失败、spec golden 回归。"
  - normalized: "- 熔断触发:真实 ideas 仓库出现结构误报、插件 shim 或 release 冒烟失败、spec golden 回归。"
  - rationale: 触发条件是可观测事件，但不是错误率/延迟等可观测数字阈值，未达 level 2。
- **rollout.rollback_steps** [core, level 2] — 6.4 回滚步骤是否可执行？
  - evidence `6.4 回滚与熔断`: "- 回滚：移除 '.mdstandards.json' 或回退插件版本；spec 行为不受影响。"
  - normalized: "- 回滚:移除 '.mdstandards.json' 或回退插件版本;spec 行为不受影响。"
  - rationale: 给出可执行的具体回滚动作（移除配置或回退版本），不止「可以回滚」。

## 替代方案诚实 (alternatives)

- **alts.status_quo** [core, level 2] — 第 7 章是否有 Status Quo 及其成本/风险？
  - evidence `7. 替代方案与权衡`: "\| Status Quo \| 零改动 \| 空壳想法继续入库，内容无标准 \| 否决 \|"
  - normalized: "\| Status Quo \| 零改动 \| 空壳想法继续入库,内容无标准 \| 否决 \|"
  - rationale: Status Quo 行同时给出成本/风险（空壳入库、内容无标准）与结论。
- **alts.rejected** [core, level 2] — 第 7 章是否至少有一个被否决方案，且淘汰原因对应当前约束？
  - evidence `7. 替代方案与权衡`: "\| 方案 C：独立 'mdstandards idea' 命令 \| 不碰 spec policy \| 解析重复、类型支持名义化 \| 否决 \|"
  - normalized: "\| 方案 C:独立 'mdstandards idea' 命令 \| 不碰 spec policy \| 解析重复、类型支持名义化 \| 否决 \|"
  - rationale: 方案 C 被否决，淘汰原因（解析重复、类型支持名义化）对应当前「复用评分引擎、不重复解析」的约束。
- **alts.matrix** [core, level 1] — 第 7 章是否有与约束对应的权衡矩阵？
  - evidence `7. 替代方案与权衡`: "\| 方案 A：类型配置 + 专用 Rego \| 最小改动、spec 不动、规则集中在 policy \| 第三种类型要再补代码 \| **采用** \|"
  - normalized: "\| 方案 A:类型配置 + 专用 Rego \| 最小改动、spec 不动、规则集中在 policy \| 第三种类型要再补代码 \| **采用** \|"
  - rationale: 存在方案对比表，但列维度是通用「优点 / 缺点」，未按具体约束组织成权衡矩阵维度。
- **alts.consistency** [core, level 2] — 第 7 章结论是否与 Goals、架构选型一致？
  - evidence `7. 替代方案与权衡`: "\| 方案 A：类型配置 + 专用 Rego \| 最小改动、spec 不动、规则集中在 policy \| 第三种类型要再补代码 \| **采用** \|"
  - normalized: "\| 方案 A:类型配置 + 专用 Rego \| 最小改动、spec 不动、规则集中在 policy \| 第三种类型要再补代码 \| **采用** \|"
  - rationale: 采用方案 A 与 3.4 的关键选型、Goals 中 spec 不回归/类型配置方向一致。

## 跨章自洽 (cross_consistency)

- **x.goals_trace** [core, level 1] — 后续章节是否逐条回应 Goals？
  - evidence `全文`: "规则见架构文档 §5.2。每条规则独立命中并带稳定的 rule ID：ID 格式与唯一性、日期格式与真实性、路径形状与目录一致性、状态枚举、tags 类型、标题非空、占位符残留、四个必需章节。"
  - normalized: "规则见架构文档 §5.2。每条规则独立命中并带稳定的 rule ID:ID 格式与唯一性、日期格式与真实性、路径形状与目录一致性、状态枚举、tags 类型、标题非空、占位符残留、四个必需章节。"
  - rationale: G1/G2/G4/G5 可在第 4、6 章找到对应设计，但 G3「空正文/暂无 不误伤」在后续章节缺少明确对应说明。
- **x.nfr_arch** [core, level 2] — NFR 与架构、选型是否一致？
  - evidence `全文`: "- **单二进制、单主包**：不拆 'cmd/' + 'internal/'，不新增运行时依赖。"
  - normalized: "- **单二进制、单主包**:不拆 'cmd/' + 'internal/',不新增运行时依赖。"
  - rationale: NFR-3 fail-closed、NFR-5 零依赖等与 3.4 选型、3.1 架构一致，未见矛盾。
- **x.rollback_release** [core, level 2] — 回滚与发布策略是否一致？
  - evidence `全文`: "- 触发后停止灰度，回退到上一插件版本；内容评分产物不阻断入库，不参与回滚判定。"
  - normalized: "- 触发后停止灰度,回退到上一插件版本;内容评分产物不阻断入库,不参与回滚判定。"
  - rationale: 回滚动作与 6.1 的冒烟→灰度发布节奏对应，回退到上一插件版本与版本发布策略一致。

## 证据卫生 (evidence_hygiene)

- **x.no_placeholders** [core, level 2] — 是否没有把 无 / ... / 待补充 / TBD 当作实质回答？
  - evidence `全文`: "空正文、'暂无'、'目前说不清，因为……' 不产生结构 finding。"
  - normalized: "空正文、'暂无'、'目前说不清,因为……' 不产生结构 finding。"
  - rationale: 全文未出现 待补充/TBD 等占位作为实质回答；文中「暂无」是合法内容示例而非占位。
- **x.no_code** [core, level 2] — 是否没有堆砌实现代码、测试日志或偶发报错？
  - evidence `全文`: "mdstandards -root ~/Workspace/ideas -policy <plugin>/policy -phase 0\|1 -format table\|json"
  - normalized: "mdstandards -root ~/Workspace/ideas -policy <plugin>/policy -phase 0\|1 -format table\|json"
  - rationale: 仅有一个 mermaid 拓扑图与一条 CLI 用法示例，没有实现代码、测试日志或偶发报错堆砌。
- **x.numbers_sourced** [core, level 2] — 关键数字是否有出处，或标了 [estimated]？
  - evidence `全文`: "**NFR-2 性能**：9 篇 idea 的 phase 1 在 1 秒内完成（2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s）；单文件上限沿用 1 MiB。"
  - normalized: "**NFR-2 性能**:9 篇 idea 的 phase 1 在 1 秒内完成(2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s);单文件上限沿用 1 MiB。"
  - rationale: 关键数字均带日期/实测上下文（5 篇、25.0/35.0、0.01s 等），未见未实测数字冒充已达成。
- **x.no_fabricated_refs** [core, level 2] — 是否没有捏造 SPEC 引用、depends_on 目标或 M<N>？
  - evidence `全文`: "  - docs/architecture/mdstandards-architecture.md"
  - normalized: "- docs/architecture/mdstandards-architecture.md"
  - rationale: 仅引用 depends_on 的架构文档与正文相对链接，未见伪造的 SPEC 引用或里程碑 M<N>（M2 为芯片型号）。

## Audit failures

None.
