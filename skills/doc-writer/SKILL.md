---
name: doc-writer
description: 按本仓库的 Spec 规范起草、修订和检查技术规格说明书。用于写新 Spec、修改已批准的 Spec、或在评审前检查背景、痛点与方案是否写清楚。不要用于里程碑纪要、README、ADR 或其他非 Spec 文档。
---

# doc-writer

把本仓库的 Spec 写作规范编译成一条可执行流程：先问清背景与痛点，再按固定骨架起草，最后用 mdstandards 过元数据闸门，并用 taskscheck 检查执行拆解。

## 适用范围

v1 只支持技术 Spec，目录为 `docs/specs/SPEC-<三位数字>-<kebab>/spec.md`，同目录有 `changelog.md`；`Approved` / `Implemented` 时同目录还有 `tasks.md`。

里程碑纪要、架构文档、README、ADR 不在 v1 范围内。用户要写这些时，明确说明不在范围内并停止，不要套用 Spec 流程。

## 章节路由

写哪一章就读哪一份，不要一次读全部：

| 章节 | 读 |
| --- | --- |
| 标题与元数据 | `spec/chapters/00-title-and-metadata.md` + `spec/standards.md` |
| 1 背景与上下文 | `spec/chapters/01-background.md` |
| 2 目标与非目标 | `spec/chapters/02-goals.md` |
| 3 架构设计与核心方案 | `spec/chapters/03-architecture.md` |
| 4 详细设计 | `spec/chapters/04-detailed-design.md` |
| 5 非功能性需求 | `spec/chapters/05-nfr.md` |
| 6 上线与迁移计划 | `spec/chapters/06-rollout.md` |
| 7 替代方案与权衡 | `spec/chapters/07-alternatives.md` |
| 8 执行拆解 | `spec/chapters/08-execution-plan.md` + `spec/tasks-template.md` |

目录布局、状态生命周期、changelog、引用与提交规范：`spec/standards.md`。

## 流程

1. **确认目标**：新建 Spec 还是修订已有 Spec。新建目录名只描述能力或领域，不出现 `M<N>`、里程碑编号或里程碑名称。
2. **读闸门**：读 `spec/gates.md`，一次只问一个问题。
3. **过闸门**：背景四问与痛点证据是硬性最小集。任一未答完，不写方案章节，只输出 Gate Report，不创建或修改 `spec.md`。
4. **起草**：读 `spec/template.md`，写每一章前读 `spec/chapters/` 下对应的指南。只写能被评审的内容；不贴实现代码、测试日志或偶发报错。
4b. **执行拆解**：Spec 进入 `Approved` / `Implemented`，或用户明确要执行计划时，读 `spec/chapters/08-execution-plan.md` 与 `spec/tasks-template.md`，在同目录写 `tasks.md`。任务从 Goals/REQs 派生，用依赖表达顺序，不复制架构。
5. **自检与质量评分**：
   - **5a 写作自检**：读 `spec/checks.md`，逐条核对。每条要求文档中的位置或证据；`无`、`...`、`待补充` 不算回答。
   - **5b 独立质量评审**：另起独立评审上下文，按 `spec/quality/reviewer.md` 产出 `quality-findings.json`；运行 `${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool qualitycheck -audit -writer-session <id>`，读 `quality-report.md`。写作修复后必须重新评审，只提交最终 findings/report。
6. **外部闸门**：先 `-phase 0` 看全貌，再 `-phase 1` 阻断。命令见下。
7. **修订**：修改已批准或已落地的 Spec 时，状态立即回退 `In Review`，更新 `updated_at`，在 `changelog.md` 追加记录；回退期间不改生产代码。

## 质量评分卡

- rubric：`spec/quality/rubric.v1.yaml`，唯一评分标准。
- 评审协议：`spec/quality/reviewer.md`。
- 输出契约：`spec/quality/findings.schema.json`。
- 产物：`docs/quality/SPEC-xxx/quality-report.md` 与 `quality-findings.json`，随 Spec 提交；`docs/quality/` 不受 `governed_globs` 管理。
- 评分卡只验证证据存在、章节位置与分数算术，不判断证据是否充分或真实。

```bash
${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool qualitycheck -validate-rubric -rubric ${CLAUDE_PLUGIN_ROOT}/skills/doc-writer/spec/quality/rubric.v1.yaml
${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool qualitycheck -audit -rubric <rubric> -findings <findings.json> -doc <spec.md> -writer-session <id> -out <report.md>
```

## 外部闸门

在目标仓库根目录运行。工具和 policy 随插件分发；`${CLAUDE_PLUGIN_ROOT}` 由 Claude Code 在技能正文里内联替换。项目级 policy 覆盖必须显式设置 `MDSTANDARDS_POLICY_DIR` 和 `MDSTANDARDS_ALLOW_POLICY_OVERRIDE=1`，否则 shim 会 exit 2。

```bash
${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool mdstandards -root . -policy "${MDSTANDARDS_POLICY_DIR:-${CLAUDE_PLUGIN_ROOT}/policy}" -phase 0 -format json
${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool mdstandards -root . -policy "${MDSTANDARDS_POLICY_DIR:-${CLAUDE_PLUGIN_ROOT}/policy}" -phase 1 -format json
${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool mdstandards --version
${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool taskscheck -spec <target>/spec.md -phase 0 -format json
${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool taskscheck -spec <target>/spec.md -phase 1 -format json
${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool taskscheck --version
```

- phase 0 退出 0 也可能包含 `extract_errors`；必须解析 JSON 并处理，不能只看退出码。
- phase 1 退出 0 通过；退出 1 是 findings；退出 2 是策略、解析或覆盖错误。
- 首次运行如果受管文档为 0 或依赖边为 0，phase 1 会退出 2。先建一份 `docs/architecture/` 下状态为 `Approved` 的架构文档，再建一份带真实 `depends_on` 的 Spec。
- `taskscheck` 只检查当前 Spec。`Draft` / `In Review` 没有 `tasks.md` 不报错；`Approved` / `Implemented` 必须有 `tasks.md`。

## 禁止

- 不要在背景未写清、痛点未给出证据时写方案章节。
- 不要编造 `depends_on` 目标或 SPEC 引用。
- 不要在 Spec 的目录、标题或正文里写里程碑编号或里程碑名称。
- 不要修改 `policy/` 或 baseline 来绕过闸门。
- 不要把 `spec/checks.md` 当作硬闸门；v1 的结构硬闸门是 mdstandards + taskscheck，`qualitycheck` 只诊断、不阻断。
