---
name: idea-writer
description: 记录和检查想法库里的 idea 文档。用于「记一下」「新想法」等场景：按模板建档、跑结构闸门、再做一次独立内容评分。不要用于 Spec、README、ADR 或其他非 idea 文档。
---

# idea-writer

把想法记录流程编译成一条可执行路径：先建档，再过结构闸门，最后做独立内容评分。结构错误阻断记录；内容分只诊断，绝不阻止想法入库。

## 适用范围

目标仓库是一个 idea 库：每个想法一个 Markdown 文件，路径为 `YYYY/YYYY-MM-DD/<slug>.md`，frontmatter 有 `id`、`title`、`date`、`tags`、`status`，正文有四个必需二级标题：一句话描述、想法、为什么有意思、下一步。模板另有一个可选二级标题「外部的声音」，用于外部反馈、同类尝试、反例或行业标准。

Spec、README、ADR 不在范围内。用户要写这些时，明确说明不在范围内并停止。

## 流程

1. **确认仓库**：仓库根应有 `.mdstandards.json`，内容为 `{"type":"idea"}`。没有就创建；这是仓库声明自己是 idea 库的唯一入口。
2. **建档**：优先用仓库自带的 `./scripts/new "<标题>" --tags <标签>`。脚本会生成 ID 和路径；不要手写文件名。
3. **填模板**：按模板填四个必需章节。用户口述零散时保留原意，不要替他拔高或加没说过的东西。写不出来时允许写「目前说不清，因为……」，也允许空正文或 `暂无`；这不违反结构闸门。「外部的声音」是可选章节：用户提供了外部反馈、同类尝试、反例或行业标准就写进去；没有就删掉整节，不要为了填满而编造。保留模板注释会被结构闸门判为占位符残留。
4. **结构闸门（同步，阻断）**：

   ```bash
   ${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool mdstandards \
     -root . -policy "${MDSTANDARDS_POLICY_DIR:-${CLAUDE_PLUGIN_ROOT}/policy}" \
     -phase 1 -format json
   ```

   退出 1 表示有 findings；逐条修到退出 0。退出 2 是配置、解析或覆盖错误，必须停下来读 stderr，不要当成通过。
5. **内容评分（落盘后，不阻断回复）**：
   - 计算 rubric 和目标的 sha256，准备 findings 头部：`schema_version`、`rubric_version`、`rubric_sha256`、`target`、`target_sha256`、`reviewer`。
   - 另起独立评审会话（不同 session id，且不接收写作对话），按 `idea/quality/reviewer.md` 评审，产出 `quality-findings.json`。先写 `claim.standalone` 的 `guess`，再读正文。
   - 运行审计：

     ```bash
     ${CLAUDE_PLUGIN_ROOT}/bin/doc-writer-tool qualitycheck -audit \
       -rubric "${CLAUDE_PLUGIN_ROOT}/skills/idea-writer/idea/quality/rubric.v1.yaml" \
       -findings quality-findings.json -doc "<idea 文件>" \
       -writer-session "<写作会话 id>" \
       -title "Idea Quality Report" -target-label Idea \
       -out "docs/quality/ideas/<id>/quality-report.md"
     ```

   - 报告写到 `docs/quality/ideas/<id>/`，不进受管文档。后台评审不可用或 5 分钟未完成时，退回同步评审并说明；评分失败就把报告标 `unavailable`，不要回滚已经建好的想法文件。
6. **告诉用户路径**：想法文件、结构闸门结果、内容报告路径。内容低分只提示「可以补什么」，不改状态、不阻止入库。

## 校准

- rubric：`idea/quality/rubric.v1.yaml`，唯一评分标准。
- 校准标签：`idea/quality/calibration/ideas.v1.yaml`，五篇想法的人工分诊标签（推进/观察/搁置），其中一篇留出。
- 先用其余四篇调 rubric，再对留出篇跑一次；至少三篇的「推进就绪度」档位与人工标签一致，才算有信号。不一致就改锚点，不写闸门。

## 禁止

- 不要在结构闸门失败时跳过或改 policy 绕过。
- 不要为了让内容分变高而替用户补没说过的东西。
- 不要把内容分当成阻断门槛；结构才阻断。
- 不要手写 idea 文件名或 ID；用 `scripts/new`。
