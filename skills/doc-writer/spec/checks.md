# Spec 自检

逐条核对；每条都要给出文档中的位置或证据。`无`、`...`、`待补充` 不算通过。

## 元数据

- [ ] frontmatter 有 `type: spec`、`status`、`owner`、`depends_on`、`updated_at`。
- [ ] `depends_on` 的每个目标都真实存在，且状态是 `Approved` 或 `Published`。
- [ ] 正文元数据表有标题、主笔、评审人、创建日期；缺一不可。
- [ ] 标题符合 `[系统/模块名] [核心动作/变更内容] 技术方案`。

## 背景与痛点

- [ ] 1.1 写了现状与核心流程。
- [ ] 1.2 有失败链，以及一次具体事故、度量或引述。
- [ ] 1.3 写了触发点。
- [ ] 1.4 写了不做的后果与非目标。

## 方案

- [ ] 每个 Goal 都有验收方式。
- [ ] 每个 Non-Goal 标了 `本次不做` 或 `Deferred to Phase 2`。
- [ ] 第 3 章引用 `depends_on` 的架构文档承载全局结构，只画本次变更的拓扑，不复制全局拓扑与模块职责。
- [ ] 替代方案包含 Status Quo 和至少一个被否决方案。
- [ ] 第 6 章有回滚触发条件。
- [ ] 没有具体实现代码、测试日志或偶发报错堆砌。

## 执行拆解

- [ ] `Approved` / `Implemented` 的 Spec 同目录有 `tasks.md`；`Draft` / `In Review` 可选。
- [ ] frontmatter 有 `type: tasks`、`status`、`owner`、`depends_on: []`、`spec`、`spec_sha256`、`updated_at`。
- [ ] 每个 Goal 至少一个专用任务；每个任务至少覆盖一个 Goal/REQ。
- [ ] 依赖可解、无自环；两个无依赖路径的任务不写同一路径。
- [ ] 完成标准可观察；验证是命令、测试 ID 或产物路径。
- [ ] `taskscheck -spec <spec> -phase 1` 退出 0，或每条 finding 已修复。

## 规范

- [ ] 目录、标题、正文没有 `M<N>` 或里程碑名称。
- [ ] 没有捏造不存在的 SPEC 引用或 `depends_on` 目标。
- [ ] 修订时状态已回退 `In Review`，`updated_at` 与 `changelog.md` 已更新。

## 外部闸门

- [ ] `mdstandards -phase 1` 退出 0，或每条 finding 已进入人工 baseline。
- [ ] phase 0 的 `extract_errors` 已处理。
