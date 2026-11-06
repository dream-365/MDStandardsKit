# Spec Quality Report

- Spec: `docs/specs/SPEC-002-idea-type/spec.md`
- Rubric: 1.1.0 (`1e6c4740ed85512b1657e6798e481c893fe015f4ce2a140665ff6ecb2ef7b29e`)
- Target SHA-256: `b22b25a5f6f2784fe3fb95cd115d216b462b4451192d919c5613b910606b6ea4`
- Reviewer session: `reviewer-spec-002-20261106`

**Headline score (core): 96.7**  |  Diagnostic score (all): 91.8

> The auditor verifies evidence presence, section placement, and score arithmetic. It does not judge whether evidence is sufficient or true.

> Normalization covers full-width ASCII, NBSP, and ideographic space; it is not full Unicode NFKC.

| Dimension | Headline | Diagnostic | Weight |
|---|---:|---:|---:|
| 背景落地 | 87.5 | 87.5 | 2 |
| 痛点证据 | 100.0 | 100.0 | 3 |
| 目标可验收 | 100.0 | 87.5 | 3 |
| 非目标边界 | 100.0 | 83.3 | 2 |
| 架构自洽 | 100.0 | 91.7 | 2 |
| 详细设计完整 | 100.0 | 83.3 | 2 |
| NFR 可度量 | 100.0 | 100.0 | 2 |
| 上线可回滚 | 83.3 | 83.3 | 2 |
| 替代方案诚实 | 100.0 | 100.0 | 2 |
| 跨章自洽 | 100.0 | 100.0 | 2 |
| 证据卫生 | 91.7 | 91.7 | 3 |

## 背景落地 (background)

- **bg.current_state** [core, level 1] — 1.1 是否写清了模块职责、核心流程与当前技术栈？
  - evidence `1.1 系统现状`: "'scripts/new' 负责建文件、生成 8 位大写 ID、校验标题非空、校验日期格式与 ID 唯一。它不检查章节是否存在、占位符是否替换、状态是否合法、日期与目录是否一致，也不评价内容。"
  - normalized: "'scripts/new' 负责建文件、生成 8 位大写 ID、校验标题非空、校验日期格式与 ID 唯一。它不检查章节是否存在、占位符是否替换、状态是否合法、日期与目录是否一致,也不评价内容。"
  - rationale: 写了 idea 库的文件组织、章节模板与 'scripts/new' 的职责，但未交代当前技术栈（语言、工具或版本）。
- **bg.failure_chain** [core, level 2] — 1.2 是否有从触发到后果的分步失败链？
  - evidence `1.2 核心技术痛点`: "1. Agent 按模板建了一个想法文件。 2. 章节缺失、模板注释还在、状态写错、日期与目录对不上，没有任何工具发现。 3. 想法进入库，之后靠人翻文件时才发现问题；内容层面也没有统一标准判断「值不值得现在推进」。"
  - normalized: "1. Agent 按模板建了一个想法文件。 2. 章节缺失、模板注释还在、状态写错、日期与目录对不上,没有任何工具发现。 3. 想法进入库,之后靠人翻文件时才发现问题;内容层面也没有统一标准判断「值不值得现在推进」。"
  - rationale: 1.2 用 1→2→3 的编号链说明从建文件、无人发现到入库后靠人翻文件的后果。
- **bg.trigger** [core, level 2] — 1.3 是否解释了为什么现在做？
  - evidence `1.3 重构/设计动机`: "触发点是想法库开始承载越来越多的 Agent 生成内容，而现有 'scripts/new' 只管文件名和 ID。"
  - normalized: "触发点是想法库开始承载越来越多的 Agent 生成内容,而现有 'scripts/new' 只管文件名和 ID。"
  - rationale: 1.3 点明触发点（Agent 生成内容增多）并解释为何现在做。
- **bg.non_goals** [core, level 2] — 1.4 是否写了不做的后果与本次非目标？
  - evidence `1.4 非目标`: "不做结构闸门的后果见 1.2：空壳想法继续入库，内容低分的想法和人判断混在一起。"
  - normalized: "不做结构闸门的后果见 1.2:空壳想法继续入库,内容低分的想法和人判断混在一起。"
  - rationale: 1.4 先写不做结构闸门的后果，再列本次非目标。

## 痛点证据 (pain_evidence)

- **bg.pain_evidence** [core, level 2] — 1.2 是否有至少一条可核实的证据（事故、度量或引述）？
  - evidence `1.2 核心技术痛点`: "证据：截至 2026-11-02，'~/Workspace/ideas' 有 5 篇想法；在没有结构闸门之前，'scripts/new' 建出的文件只保证 title/date/id 三项。"
  - normalized: "证据:截至 2026-11-02,'~/Workspace/ideas' 有 5 篇想法;在没有结构闸门之前,'scripts/new' 建出的文件只保证 title/date/id 三项。"
  - rationale: 给出截至 2026-11-02 的篇数与 'scripts/new' 仅保证三项的事实，对应该失败链第 2 步。

## 目标可验收 (goal_verifiability)

- **goal.verifiable** [core, level 2] — 2.1 每个 Goal 是否都有验收方式？
  - evidence `2.1 目标`: "- **G3 记录不误伤**：空正文、'暂无'、'目前说不清，因为……' 不产生结构 finding。验收：'TestIdeaValidRepo' 与真实 ideas 语料（2026-11-03 实测 9 篇）phase 1 通过。"
  - normalized: "- **G3 记录不误伤**:空正文、'暂无'、'目前说不清,因为……' 不产生结构 finding。验收:'TestIdeaValidRepo' 与真实 ideas 语料(2026-11-03 实测 9 篇)phase 1 通过。"
  - evidence `2.1 目标`: "- **G4 内容评分**：rubric 四个 core 项计算「推进就绪度」，成熟度和「非显而易见」只做诊断；独立 reviewer 先盲读再评分；'qualitycheck' 复算分数并验证逐字证据。验收：'-validate-rubric'、万花筒冒烟报告。"
  - normalized: "- **G4 内容评分**:rubric 四个 core 项计算「推进就绪度」,成熟度和「非显而易见」只做诊断;独立 reviewer 先盲读再评分;'qualitycheck' 复算分数并验证逐字证据。验收:'-validate-rubric'、万花筒冒烟报告。"
  - evidence `2.1 目标`: "- **G6 可选外部声音**：模板在「为什么有意思」与「下一步」之间提供可选章节「外部的声音」；整节缺失不产生 finding；空正文允许（与其他章节一致）；保留含 'optional:outside-voice' 哨兵的模板注释则命中 'idea.placeholder_remaining'。验收：'TestIdeaOptionalVoiceSection'。"
  - normalized: "- **G6 可选外部声音**:模板在「为什么有意思」与「下一步」之间提供可选章节「外部的声音」;整节缺失不产生 finding;空正文允许(与其他章节一致);保留含 'optional:outside-voice' 哨兵的模板注释则命中 'idea.placeholder_remaining'。验收:'TestIdeaOptionalVoiceSection'。"
  - rationale: G1–G6 每条都以测试、'make check' 或冒烟报告给出验收方式。
- **goal.quantified** [core, level 2] — 2.1 是否至少有一个 Goal 有数字/单位/明确阈值？
  - evidence `2.1 目标`: "验收：'TestIdeaValidRepo' 与真实 ideas 语料（2026-11-03 实测 9 篇）phase 1 通过。"
  - normalized: "验收:'TestIdeaValidRepo' 与真实 ideas 语料(2026-11-03 实测 9 篇)phase 1 通过。"
  - rationale: G3 给出 2026-11-03 实测 9 篇的样本数量（数字+单位）。
- **goal.count** [extended, level 1] — 2.1 Goals 是否控制在 3–6 个？
  - evidence `2.1 目标`: "- **G5 spec 不回归**：spec 的 '-dump-corpus' 输出逐字节不变；'make check' 全绿。验收：'TestSpecDumpCorpusGolden'、'make check'。 - **G6 可选外部声音**：模板在「为什么有意思」与「下一步」之间提供可选章节「外部的声音」；整节缺失不产生 finding；空正文允许（与其他章节一致）；保留含 'optional:outside-voice' 哨兵的模板注释则命中 'idea.placeholder_remaining'。验收：'TestIdeaOptionalVoiceSection'。"
  - normalized: "- **G5 spec 不回归**:spec 的 '-dump-corpus' 输出逐字节不变;'make check' 全绿。验收:'TestSpecDumpCorpusGolden'、'make check'。 - **G6 可选外部声音**:模板在「为什么有意思」与「下一步」之间提供可选章节「外部的声音」;整节缺失不产生 finding;空正文允许(与其他章节一致);保留含 'optional:outside-voice' 哨兵的模板注释则命中 'idea.placeholder_remaining'。验收:'TestIdeaOptionalVoiceSection'。"
  - rationale: G1–G6 共 6 个，在 3–6 范围内，但未按架构/性能/可靠性分类。

## 非目标边界 (non_goal_boundary)

- **nongoal.labeled** [core, level 2] — 2.2 每条 Non-Goal 是否标了本次不做或 Deferred to Phase 2？
  - evidence `2.2 非目标`: "- **NG1** 想法状态流转校验（Deferred to Phase 2）。"
  - normalized: "- **NG1** 想法状态流转校验(Deferred to Phase 2)。"
  - evidence `2.2 非目标`: "- **NG4** Windows 分发（本次不做）。"
  - normalized: "- **NG4** Windows 分发(本次不做)。"
  - rationale: NG1–NG6 每条都带「Deferred to Phase 2」或「本次不做」。
- **nongoal.not_dumping_ground** [extended, level 1] — 2.2 非目标是否能追溯到痛点或边界？
  - evidence `2.2 非目标`: "- **NG4** Windows 分发（本次不做）。"
  - normalized: "- **NG4** Windows 分发(本次不做)。"
  - rationale: 多数非目标可追溯到 1.2/1.3 或 v0.2.0 边界，NG4（Windows 分发）与痛点无对应。

## 架构自洽 (architecture)

- **arch.solution_shape** [core, level 2] — 第 3 章是否有本次变更的拓扑或等价的模块/数据流描述？
  - evidence `3. 架构设计与核心方案`: "'''mermaid flowchart LR     CLI[mdstandards] --> Config[.mdstandards.json]     CLI --> Policy[policy/data/standard.json]     Config --> Type[document type]     Policy --> Profile[document_types.idea]     Type --> Profile     Profile --> Extract[Extract root globs docType profile]     Extract --> Corpus[corpus + document_type]     Corpus --> Rego[policy/idea.rego]     Rego --> Findings[findings]     Findings --> Report[table / json / github] '''"
  - normalized: "'''mermaid flowchart LR CLI[mdstandards] --> Config[.mdstandards.json] CLI --> Policy[policy/data/standard.json] Config --> Type[document type] Policy --> Profile[document_types.idea] Type --> Profile Profile --> Extract[Extract root globs docType profile] Extract --> Corpus[corpus + document_type] Corpus --> Rego[policy/idea.rego] Rego --> Findings[findings] Findings --> Report[table / json / github] '''"
  - rationale: 3.1 的流程图给出本次变更从配置、解析、Rego 到报告的拓扑/数据流。
- **arch.topology** [extended, level 1] — 3.1 拓扑图是否让每个模块有一句职责？
  - evidence `3.1 总体架构`: "    CLI[mdstandards] --> Config[.mdstandards.json]"
  - normalized: "CLI[mdstandards] --> Config[.mdstandards.json]"
  - rationale: 有流程图，但节点只有名称，未给每个模块一句职责；职责散见 3.2/3.4。
- **arch.change_set** [extended, level 2] — 第 3 章是否写清新增/改造/保持不变？
  - evidence `3. 架构设计与核心方案`: "**新增**：  - 'policy/idea.rego'：idea 结构规则。 - 'skills/idea-writer/'：SKILL、rubric、reviewer、calibration、findings schema。 - '.mdstandards.json' 读取与类型解析。 - 'docs/quality/ideas/' 报告目录约定。"
  - normalized: "**新增**: - 'policy/idea.rego':idea 结构规则。 - 'skills/idea-writer/':SKILL、rubric、reviewer、calibration、findings schema。 - '.mdstandards.json' 读取与类型解析。 - 'docs/quality/ideas/' 报告目录约定。"
  - evidence `3. 架构设计与核心方案`: "**改造**：  - 'extract.go'：'Extract(root, globs, docType, profile)'；idea 字段与章节解析；'date' 同时接受字符串和 'time.Time'。 - 'main.go'：类型→globs 解析；配置失败 fail closed；零依赖边断言按 'graph' 生效。 - 'policy/refs.rego'、'policy/status.rego'：显式 spec 类型守卫。 - 'tools/qualitycheck'：'-title' / '-target-label'；findings schema 增加 'guess'。"
  - normalized: "**改造**: - 'extract.go':'Extract(root, globs, docType, profile)';idea 字段与章节解析;'date' 同时接受字符串和 'time.Time'。 - 'main.go':类型→globs 解析;配置失败 fail closed;零依赖边断言按 'graph' 生效。 - 'policy/refs.rego'、'policy/status.rego':显式 spec 类型守卫。 - 'tools/qualitycheck':'-title' / '-target-label';findings schema 增加 'guess'。"
  - evidence `3. 架构设计与核心方案`: "**保持不变**：  - spec 的规则、退出码和 '-dump-corpus' 字节输出。 - 单 'main' 包，不拆 'cmd/' + 'internal/'。 - 不新增运行时依赖。"
  - normalized: "**保持不变**: - spec 的规则、退出码和 '-dump-corpus' 字节输出。 - 单 'main' 包,不拆 'cmd/' + 'internal/'。 - 不新增运行时依赖。"
  - rationale: 3.2 分「新增／改造／保持不变」三节列出具体文件与契约。
- **arch.rationale** [extended, level 2] — 关键选型是否有理由，详细对比指向第 7 章？
  - evidence `3.4 关键选型与约束`: "- **类型配置 + 专用 Rego**：类型由仓库级 '.mdstandards.json' 声明，规则集中在 'policy/idea.rego'；spec 的规则、退出码与 '-dump-corpus' 契约不变。方案对比见第 7 章。"
  - normalized: "- **类型配置 + 专用 Rego**:类型由仓库级 '.mdstandards.json' 声明,规则集中在 'policy/idea.rego';spec 的规则、退出码与 '-dump-corpus' 契约不变。方案对比见第 7 章。"
  - rationale: 3.4 给出类型配置+专用 Rego 的理由，并明确「方案对比见第 7 章」。
- **arch.references_source** [extended, level 2] — 第 3 章是否引用 depends_on 架构文档，而不是自行展开其全局结构？
  - evidence `3. 架构设计与核心方案`: "全局结构与模块职责见 [架构文档](../../architecture/mdstandards-architecture.md)；本节只描述 idea 类型变更范围内的数据流。"
  - normalized: "全局结构与模块职责见 [架构文档](../../architecture/mdstandards-architecture.md);本节只描述 idea 类型变更范围内的数据流。"
  - rationale: 3.1 明确把全局结构与模块职责交给 depends_on 架构文档，本章只描述 idea 变更范围内的数据流。

## 详细设计完整 (detailed_design)

- **design.contract** [core, level 2] — 第 4 章是否至少有一份数据模型或接口契约？
  - evidence `4. 详细设计`: "\| 'document_type' \| 语料级类型，idea 时为 'idea' \|"
  - normalized: "\| 'document_type' \| 语料级类型,idea 时为 'idea' \|"
  - rationale: 4.3 给出 corpus 字段数据模型，4.4 给出 CLI 接口，至少一份契约成立。
- **design.data_model** [extended, level 1] — 4.3 数据模型是否有字段类型、必填与校验？
  - evidence `4.3 数据模型与契约`: "\| 字段 \| 含义 \| \|---\|---\| \| 'document_type' \| 语料级类型，idea 时为 'idea' \|"
  - normalized: "\| 字段 \| 含义 \| \|---\|---\| \| 'document_type' \| 语料级类型,idea 时为 'idea' \|"
  - rationale: 字段表只有「含义」列，没有类型、必填与校验规则。
- **design.interface** [extended, level 2] — 4.4 接口是否有请求/响应/状态码或方法签名？
  - evidence `4.4 接口与协议`: "'''bash mdstandards -root ~/Workspace/ideas -policy <plugin>/policy -phase 0\|1 -format table\|json '''"
  - normalized: "'''bash mdstandards -root ~/Workspace/ideas -policy <plugin>/policy -phase 0\|1 -format table\|json '''"
  - evidence `4.4 接口与协议`: "'-policy' 相对当前工作目录解析；idea 仓库必须显式传插件的 policy 目录，或设置 'MDSTANDARDS_POLICY_DIR'。输出沿用 'table\|json\|github'；退出码见 4.6。"
  - normalized: "'-policy' 相对当前工作目录解析;idea 仓库必须显式传插件的 policy 目录,或设置 'MDSTANDARDS_POLICY_DIR'。输出沿用 'table\|json\|github';退出码见 4.6。"
  - rationale: 4.4 给出 CLI 方法签名/参数、输出格式，并指向 4.6 的退出码。
- **design.sequence** [extended, level 2] — 4.5 复杂流程是否有时序图或等价分步说明？
  - evidence `4.5 核心流程与时序`: "1. 读取仓库根 '.mdstandards.json'，解析 'type=idea'。 2. 从 'standard.json' 取出 'document_types.idea' 的 globs 与结构配置。 3. 遍历 globs，解析 frontmatter 与正文标题，产出 document 或 extract_error。 4. 覆盖断言：零受管文档、'phase >= 1' 的零解析文档都退出 2；idea 类型不要求依赖边。 5. 'policy/idea.rego' 求值，产出结构 findings。 6. 'phase 1' 有 finding 退出 1；否则退出 0。内容评分在记录流程之外异步进行。"
  - normalized: "1. 读取仓库根 '.mdstandards.json',解析 'type=idea'。 2. 从 'standard.json' 取出 'document_types.idea' 的 globs 与结构配置。 3. 遍历 globs,解析 frontmatter 与正文标题,产出 document 或 extract_error。 4. 覆盖断言:零受管文档、'phase >= 1' 的零解析文档都退出 2;idea 类型不要求依赖边。 5. 'policy/idea.rego' 求值,产出结构 findings。 6. 'phase 1' 有 finding 退出 1;否则退出 0。内容评分在记录流程之外异步进行。"
  - rationale: 4.5 用 1–6 编号步骤描述解析到退出的流程，含零文档/解析失败分支。
- **design.errors** [extended, level 1] — 4.6 每个错误码是否有触发场景与处理策略？
  - evidence `4.6 异常处理与错误码`: "\| 码 \| 含义 \| \|---\|---\| \| 0 \| 干净，或 phase 0 \| \| 1 \| phase 1 有结构 finding \| \| 2 \| 配置/策略加载失败、未知 type、零受管文档、零解析文档、解析错误 \|"
  - normalized: "\| 码 \| 含义 \| \|---\|---\| \| 0 \| 干净,或 phase 0 \| \| 1 \| phase 1 有结构 finding \| \| 2 \| 配置/策略加载失败、未知 type、零受管文档、零解析文档、解析错误 \|"
  - rationale: 退出码各有含义与触发场景，但没有写调用方的恢复/处理策略。

## NFR 可度量 (nfr)

- **nfr.measurable** [core, level 2] — 第 5 章每条 NFR 是否有数字、单位或验收方式？
  - evidence `5. 非功能性需求`: "- **NFR-2 性能**：9 篇 idea 的 phase 1 在 1 秒内完成（2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s）；单文件上限沿用 1 MiB。"
  - normalized: "- **NFR-2 性能**:9 篇 idea 的 phase 1 在 1 秒内完成(2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s);单文件上限沿用 1 MiB。"
  - rationale: NFR-2 给出 1 秒阈值、0.01s 实测与 1 MiB 上限；其余 NFR 以可验证行为验收。
- **nfr.coverage** [core, level 2] — 第 5 章是否覆盖性能、可靠性、可观测性、安全，或说明不适用？
  - evidence `5. 非功能性需求`: "- **NFR-3 安全**：'.mdstandards.json' 使用 'DisallowUnknownFields'；配置错误 fail closed；frontmatter 仍按不可信输入处理。"
  - normalized: "- **NFR-3 安全**:'.mdstandards.json' 使用 'DisallowUnknownFields';配置错误 fail closed;frontmatter 仍按不可信输入处理。"
  - evidence `5. 非功能性需求`: "- **NFR-4 可观测性**：每条 finding 带 'rule'、'document'、'target'、'message'；内容报告带 rubric 版本、目标 SHA-256、reviewer 会话。"
  - normalized: "- **NFR-4 可观测性**:每条 finding 带 'rule'、'document'、'target'、'message';内容报告带 rubric 版本、目标 SHA-256、reviewer 会话。"
  - rationale: 确定性/性能、安全、可观测性分别由 NFR-1/2、NFR-3、NFR-4 覆盖。
- **nfr.honesty** [core, level 2] — 未实测的指标是否标了 [estimated] 或 未实现？
  - evidence `5. 非功能性需求`: "- **NFR-2 性能**：9 篇 idea 的 phase 1 在 1 秒内完成（2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s）；单文件上限沿用 1 MiB。"
  - normalized: "- **NFR-2 性能**:9 篇 idea 的 phase 1 在 1 秒内完成(2026-11-03 在 Apple M2 / macOS 27.0.1 / arm64 实测 0.01s);单文件上限沿用 1 MiB。"
  - rationale: 性能数字标明 2026-11-03 的实测环境与结果，未把未实测项写成已达成。

## 上线可回滚 (rollout)

- **rollout.release** [core, level 1] — 6.1 是否有发布策略与灰度节奏？
  - evidence `6.1 发布与灰度`: "1. 在 '~/Workspace/ideas' 增加 '.mdstandards.json'。 2. 'scripts/new' 在工具可用时跑 phase 0 检查；不可用只打印提示。 3. 插件版本升到 '0.2.0'，release tag 必须与 'plugin.json' 一致；先对 'testdata/idea-repo' 冒烟，再对真实 ideas 仓库灰度。"
  - normalized: "1. 在 '~/Workspace/ideas' 增加 '.mdstandards.json'。 2. 'scripts/new' 在工具可用时跑 phase 0 检查;不可用只打印提示。 3. 插件版本升到 '0.2.0',release tag 必须与 'plugin.json' 一致;先对 'testdata/idea-repo' 冒烟,再对真实 ideas 仓库灰度。"
  - rationale: 有发布步骤与 testdata→真实仓库的顺序，但没有灰度批次或观察期。
- **rollout.migration** [core, level 2] — 6.2 是否有迁移/兼容方案，或说明无需迁移？
  - evidence `6.2 数据迁移与兼容性`: "- 不迁移既有想法文件；新增配置即可启用结构闸门。"
  - normalized: "- 不迁移既有想法文件;新增配置即可启用结构闸门。"
  - rationale: 明确说明不迁移既有文件，新增配置即可启用。
- **rollout.smoke** [core, level 1] — 6.3 是否有冒烟清单与观察指标？
  - evidence `6.3 上线验证与冒烟`: "1. 'testdata/idea-repo' phase 1 退出 0。 2. 真实 ideas 仓库 phase 1 退出 0（2026-11-03 为 9 篇）。 3. 'make check' 全绿，含 spec golden 与插件 shim 测试。"
  - normalized: "1. 'testdata/idea-repo' phase 1 退出 0。 2. 真实 ideas 仓库 phase 1 退出 0(2026-11-03 为 9 篇)。 3. 'make check' 全绿,含 spec golden 与插件 shim 测试。"
  - rationale: 有 3 条冒烟步骤，但只写退出 0/全绿，没有错误率、延迟等观察指标。
- **rollout.rollback_trigger** [core, level 2] — 6.4 回滚触发条件是否是可观测数字？
  - evidence `6.4 回滚与熔断`: "- 熔断触发：真实 ideas 仓库出现结构误报、插件 shim 或 release 冒烟失败、spec golden 回归；真实语料 phase 1 耗时超过 1 秒（对应 NFR-2）也触发。"
  - normalized: "- 熔断触发:真实 ideas 仓库出现结构误报、插件 shim 或 release 冒烟失败、spec golden 回归;真实语料 phase 1 耗时超过 1 秒(对应 NFR-2)也触发。"
  - rationale: 触发条件新增「phase 1 耗时超过 1 秒」的可观测数字，并对应 NFR-2。
- **rollout.rollback_steps** [core, level 2] — 6.4 回滚步骤是否可执行？
  - evidence `6.4 回滚与熔断`: "- 回滚：移除 '.mdstandards.json' 或回退插件版本；spec 行为不受影响。"
  - normalized: "- 回滚:移除 '.mdstandards.json' 或回退插件版本;spec 行为不受影响。"
  - evidence `6.4 回滚与熔断`: "- 触发后停止灰度，回退到上一插件版本；内容评分产物不阻断入库，不参与回滚判定。"
  - normalized: "- 触发后停止灰度,回退到上一插件版本;内容评分产物不阻断入库,不参与回滚判定。"
  - rationale: 给出移除配置或回退插件版本的具体动作，并说明触发后停止灰度。

## 替代方案诚实 (alternatives)

- **alts.status_quo** [core, level 2] — 第 7 章是否有 Status Quo 及其成本/风险？
  - evidence `7. 替代方案与权衡`: "\| Status Quo \| 零改动 \| 空壳想法继续入库，内容无标准 \| 否决 \|"
  - normalized: "\| Status Quo \| 零改动 \| 空壳想法继续入库,内容无标准 \| 否决 \|"
  - rationale: Status Quo 行写明零改动与「空壳想法继续入库，内容无标准」的成本。
- **alts.rejected** [core, level 2] — 第 7 章是否至少有一个被否决方案，且淘汰原因对应当前约束？
  - evidence `7. 替代方案与权衡`: "\| 方案 B：声明式类型 schema \| 类型系统更通用 \| 重构大、测试迁移、过度抽象风险 \| 延后 \|"
  - normalized: "\| 方案 B:声明式类型 schema \| 类型系统更通用 \| 重构大、测试迁移、过度抽象风险 \| 延后 \|"
  - evidence `7. 替代方案与权衡`: "\| 方案 C：独立 'mdstandards idea' 命令 \| 不碰 spec policy \| 解析重复、类型支持名义化 \| 否决 \|"
  - normalized: "\| 方案 C:独立 'mdstandards idea' 命令 \| 不碰 spec policy \| 解析重复、类型支持名义化 \| 否决 \|"
  - rationale: 方案 B/C 被否决，理由（重构大、解析重复）对应当前最小改动与 spec 不动约束。
- **alts.matrix** [core, level 2] — 第 7 章是否有与约束对应的权衡矩阵？
  - evidence `7. 替代方案与权衡`: "\| 方案 \| 优点 \| 缺点 \| 结论 \| \|---\|---\|---\|---\| \| Status Quo \| 零改动 \| 空壳想法继续入库，内容无标准 \| 否决 \| \| 方案 A：类型配置 + 专用 Rego \| 最小改动、spec 不动、规则集中在 policy \| 第三种类型要再补代码 \| **采用** \| \| 方案 B：声明式类型 schema \| 类型系统更通用 \| 重构大、测试迁移、过度抽象风险 \| 延后 \| \| 方案 C：独立 'mdstandards idea' 命令 \| 不碰 spec policy \| 解析重复、类型支持名义化 \| 否决 \|"
  - normalized: "\| 方案 \| 优点 \| 缺点 \| 结论 \| \|---\|---\|---\|---\| \| Status Quo \| 零改动 \| 空壳想法继续入库,内容无标准 \| 否决 \| \| 方案 A:类型配置 + 专用 Rego \| 最小改动、spec 不动、规则集中在 policy \| 第三种类型要再补代码 \| **采用** \| \| 方案 B:声明式类型 schema \| 类型系统更通用 \| 重构大、测试迁移、过度抽象风险 \| 延后 \| \| 方案 C:独立 'mdstandards idea' 命令 \| 不碰 spec policy \| 解析重复、类型支持名义化 \| 否决 \|"
  - rationale: 表格按方案列出优点/缺点/结论，覆盖改动量、spec 影响与抽象风险。
- **alts.consistency** [core, level 2] — 第 7 章结论是否与 Goals、架构选型一致？
  - evidence `7. 替代方案与权衡`: "\| 方案 A：类型配置 + 专用 Rego \| 最小改动、spec 不动、规则集中在 policy \| 第三种类型要再补代码 \| **采用** \|"
  - normalized: "\| 方案 A:类型配置 + 专用 Rego \| 最小改动、spec 不动、规则集中在 policy \| 第三种类型要再补代码 \| **采用** \|"
  - rationale: 结论采用方案 A，与 3.4 选型及 G1/G2 的类型配置+专用 Rego 一致。

## 跨章自洽 (cross_consistency)

- **x.goals_trace** [core, level 2] — 后续章节是否逐条回应 Goals？
  - evidence `4.8 内容评分`: "core 项：'claim.standalone'(w2)、'problem.trigger'(w2)、'next.testable'(w3)、'next.bounded'(w1)。"
  - normalized: "core 项:'claim.standalone'(w2)、'problem.trigger'(w2)、'next.testable'(w3)、'next.bounded'(w1)。"
  - evidence `4.2 结构规则`: "「外部的声音」是可选章节：整节缺失不报错；空正文也允许（与其他章节一致）；章节存在时，含 'optional:outside-voice' 哨兵的模板注释按 'idea.placeholder_remaining' 处理。"
  - normalized: "「外部的声音」是可选章节:整节缺失不报错;空正文也允许(与其他章节一致);章节存在时,含 'optional:outside-voice' 哨兵的模板注释按 'idea.placeholder_remaining' 处理。"
  - evidence `6.2 数据迁移与兼容性`: "spec 的规则、退出码与 '-dump-corpus' 字节输出保持不变；idea 规则只在 'type=idea' 仓库生效。"
  - normalized: "spec 的规则、退出码与 '-dump-corpus' 字节输出保持不变;idea 规则只在 'type=idea' 仓库生效。"
  - rationale: G1/G2 见 4.5/4.2，G3/G5 见 6.3/6.2，G4 见 4.8，G6 见 4.2，均可追溯。
- **x.nfr_arch** [core, level 2] — NFR 与架构、选型是否一致？
  - evidence `5. 非功能性需求`: "- **NFR-1 确定性**：同一语料、同一 policy 的 JSON 输出字节稳定；spec '-dump-corpus' 与改造前逐字节一致。"
  - normalized: "- **NFR-1 确定性**:同一语料、同一 policy 的 JSON 输出字节稳定;spec '-dump-corpus' 与改造前逐字节一致。"
  - evidence `3.4 关键选型与约束`: "- **配置 fail-closed**：'.mdstandards.json' 使用 'DisallowUnknownFields'；缺失、损坏、未知 type 或未知键都退出 2。"
  - normalized: "- **配置 fail-closed**:'.mdstandards.json' 使用 'DisallowUnknownFields';缺失、损坏、未知 type 或未知键都退出 2。"
  - rationale: NFR-1/3/5 分别对应 6.2 的字节不变、3.4 的 fail-closed 与不新增依赖。
- **x.rollback_release** [core, level 2] — 回滚与发布策略是否一致？
  - evidence `6.1 发布与灰度`: "3. 插件版本升到 '0.2.0'，release tag 必须与 'plugin.json' 一致；先对 'testdata/idea-repo' 冒烟，再对真实 ideas 仓库灰度。"
  - normalized: "3. 插件版本升到 '0.2.0',release tag 必须与 'plugin.json' 一致;先对 'testdata/idea-repo' 冒烟,再对真实 ideas 仓库灰度。"
  - evidence `6.4 回滚与熔断`: "- 回滚：移除 '.mdstandards.json' 或回退插件版本；spec 行为不受影响。"
  - normalized: "- 回滚:移除 '.mdstandards.json' 或回退插件版本;spec 行为不受影响。"
  - rationale: 回滚即撤销 6.1 的配置启用与插件升级，并停止灰度，二者一致。

## 证据卫生 (evidence_hygiene)

- **x.no_placeholders** [core, level 2] — 是否没有把 无 / ... / 待补充 / TBD 当作实质回答？
  - evidence `2.1 目标`: "- **G3 记录不误伤**：空正文、'暂无'、'目前说不清，因为……' 不产生结构 finding。验收：'TestIdeaValidRepo' 与真实 ideas 语料（2026-11-03 实测 9 篇）phase 1 通过。"
  - normalized: "- **G3 记录不误伤**:空正文、'暂无'、'目前说不清,因为……' 不产生结构 finding。验收:'TestIdeaValidRepo' 与真实 ideas 语料(2026-11-03 实测 9 篇)phase 1 通过。"
  - rationale: 「暂无」「……」仅作为允许内容的示例出现，全文未把 待补充/TBD 当实质回答。
- **x.no_code** [core, level 2] — 是否没有堆砌实现代码、测试日志或偶发报错？
  - evidence `4.4 接口与协议`: "'''bash mdstandards -root ~/Workspace/ideas -policy <plugin>/policy -phase 0\|1 -format table\|json '''"
  - normalized: "'''bash mdstandards -root ~/Workspace/ideas -policy <plugin>/policy -phase 0\|1 -format table\|json '''"
  - rationale: 代码块只有 mermaid 架构图与 CLI 用法，没有实现代码、测试日志或报错堆砌。
- **x.numbers_sourced** [core, level 1] — 关键数字是否有出处，或标了 [estimated]？
  - evidence `1.2 核心技术痛点`: "证据：截至 2026-11-02，'~/Workspace/ideas' 有 5 篇想法；在没有结构闸门之前，'scripts/new' 建出的文件只保证 title/date/id 三项。2026-11-02 的冒烟测试对「万花筒」跑完 idea rubric，headline 25.0、diagnostic 35.0，证明内容评分链路能产出可审计的报告。"
  - normalized: "证据:截至 2026-11-02,'~/Workspace/ideas' 有 5 篇想法;在没有结构闸门之前,'scripts/new' 建出的文件只保证 title/date/id 三项。2026-11-02 的冒烟测试对「万花筒」跑完 idea rubric,headline 25.0、diagnostic 35.0,证明内容评分链路能产出可审计的报告。"
  - rationale: 冒烟分数与篇数只有日期，没有报告路径/链接；仅部分实测数字带环境说明。
- **x.no_fabricated_refs** [core, level 2] — 是否没有捏造 SPEC 引用、depends_on 目标或 M<N>？
  - evidence `全文`: "depends_on:   - docs/architecture/mdstandards-architecture.md"
  - normalized: "depends_on: - docs/architecture/mdstandards-architecture.md"
  - rationale: 只引用 depends_on 架构文档与设计记录，未出现 SPEC 编号或 M<N>。

## Audit failures

None.
