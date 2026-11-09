# Spec 标准

写任何一章之前先读这份；它约束目录、元数据、状态、changelog 与引用规范。

## 目录布局

- 一个 Spec 一个目录：`docs/specs/SPEC-<三位数字>-<kebab>/`。
- 目录名和标题只描述能力或领域，不出现 `M<N>`、里程碑编号或里程碑名称。
- 固定文件：`spec.md`、`changelog.md`；`Approved` / `Implemented` 时加 `tasks.md`。
- 附属资源：`assets/`（架构图、流程图、截图，禁止外部图床）、`schemas/` 或 `examples/`（JSON Schema、Protobuf、示例数据）。
- 质量产物：`docs/quality/SPEC-xxx/quality-report.md` 与 `quality-findings.json`，随 Spec 提交；`docs/quality/` 不在 `governed_globs` 内。
- 总索引：`docs/specs/README.md` 汇总编号、名称、状态与链接；新增 Spec 或状态流转时同步。
- 归档：`Deprecated` / `Superseded` 移到 `docs/specs/archive/`，保留原目录名。

## 元数据

YAML frontmatter（闸门读取）：

```yaml
type: spec
status: Draft
owner: example
depends_on:
  - docs/architecture/system.md
updated_at: 2026-11-01
```

正文元数据表（人读）：标题、主笔、评审人、创建日期、目标版本、关联资源。标题公式：`[系统/模块名] [核心动作/变更内容] 技术方案`。

`depends_on` 只写真实存在的受管文档；路径相对仓库根，大小写敏感。

## 状态生命周期

| 状态 | 含义 | 行为约束 |
| --- | --- | --- |
| Draft | 起草或征求意见 | 禁止据此写生产代码 |
| In Review | 正式评审 | 收集意见，修正争议 |
| Approved | 评审通过 | 准许编码，作为 SSOT；必须有 `tasks.md` |
| Implemented | 已落地 | 当前版本生效；必须有 `tasks.md` |
| Superseded | 被新方案替代 | 归档，仅作历史记录 |
| Deprecated | 废弃 | 归档，不再维护 |

强回退：任何核心内容修改（AC、架构、Goals/Non-Goals、边界规则）提交即回退 `In Review`；大幅重写回退 `Draft`。回退到 `Approved` 之前禁止修改或合并对应生产代码。

## 执行拆解

`Approved` / `Implemented` 的 Spec 必须在同目录提交 `tasks.md`；`Draft` / `In Review` 可选。

- `tasks.md` 从 `spec.md` 派生，不复制架构、数据模型或接口契约。
- 任务表固定七列：ID、任务、覆盖、依赖、写入、完成标准、验证。
- 每个 Goal 至少一个专用任务；每个任务至少覆盖一个 ID；依赖无环；无依赖路径的任务不写同一路径。
- `spec_sha256` 钉住 `spec.md`；Spec 改动后重算。
- 完整规则：`spec/chapters/08-execution-plan.md`；模板：`spec/tasks-template.md`；闸门：`taskscheck`。

## Changelog

进入 `In Review` / `Approved` / `Implemented` 后，架构调整、AC 变更、增删 Goal、边界修订必须写 `changelog.md`。

```markdown
# Changelog: SPEC-001 <名称>

## [YYYY-MM-DD] - <变更摘要>
- **修改类型**: Added / Changed / Fixed / Deprecated
- **影响范围**: 例如 AC-1.2, 3.2 Data Model
- **变更原因**: <为什么进行此修改>
- **具体内容**: <修改的具体细节>
```

修订时同步 frontmatter 的 `updated_at`。

## 引用与写作

- 跨 Spec 引用用 `[编号]()`；非必要不引用，避免 Spec 间过度交叉。
- 严禁捏造不存在的 Spec、Issue 或链接。
- 架构图里不写 Spec 编号。
- 第 3 章引用 `depends_on` 架构文档承载全局结构，不复制全局拓扑、模块职责与历史决策；只写本次变更范围内的拓扑、职责与数据流。
- LaTeX：行内 `$O(N)$`，块级 `$$ ... $$`。
- 严禁写具体实现代码；允许抽象类、接口、数据契约。
- 不把测试日志、偶发报错堆进背景；结论提炼后放进现状/差距或验收。
- 围绕业务主线叙事，不把"写 Spec"本身当叙事线。

## 质量评分

- rubric：`spec/quality/rubric.v1.yaml`；item、权重或锚点变化即升版本，旧分数不可比。
- 独立评审：另起会话，按 `reviewer.md` 产出 `quality-findings.json`；写作 Agent 不参与打分。
- 审计：`qualitycheck -audit -rubric <rubric> -findings <findings.json> -doc <spec.md> -writer-session <id>`；证据必须逐字存在，分数由工具复算。
- v1 只诊断、不阻断；校准与阈值在 follow-up。

## 提交

- 一个 commit 只涉及一个 Spec；`spec.md`、`changelog.md` 与附属资源可合并提交。
- 修改多个 Spec 时按 Spec 拆分提交。
