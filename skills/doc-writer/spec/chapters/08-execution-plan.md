# 第 8 章：执行拆解

目标：把 Spec 变成一张依赖有序、可独立验证的任务表，让执行者一眼看出下一步和并行车道。

## 什么时候写

- `Approved` / `Implemented` 的 Spec：必须有 `tasks.md`。
- `Draft` / `In Review`：可选；存在就检查。
- 单批可完成的小 Spec：仍写 `tasks.md`，用 `not_applicable: <理由>` 明示，不静默缺失。

## 固定 frontmatter

```yaml
---
type: tasks
status: Draft
owner: <owner>
depends_on: []
spec: ./spec.md
spec_sha256: "<spec.md 的 sha256>"
updated_at: <YYYY-MM-DD>
---
```

- `depends_on` 固定为空数组；Spec 关联走 `spec` 字段。
- `spec_sha256` 是 `spec.md` 全文的 sha256；Spec 一改就要重算。用 `taskscheck -print-spec-hash spec.md` 生成。
- 不适用时写 `not_applicable: <具体理由>`，不写任务表。

## 任务表

列名固定：

```markdown
| ID | 任务 | 覆盖 | 依赖 | 写入 | 完成标准 | 验证 |
|---|---|---|---|---|---|---|
| T1 | 抽取语料 | G1, REQ-001 | — | extract.go | dump 出三份语料 | go test ./... -run TestDumpCorpus |
| T2 | 路径归一化 | G1, REQ-005 | T1 | normalize.go | 越界路径产出 extract_error | go test ./... -run TestNormalize |
```

- `ID`：`T[A-Za-z0-9_-]*`，与 `G<n>` / `REQ-<三位数字>` 命名空间隔离。
- `覆盖`：逗号分隔的 Goal / REQ ID；每个任务至少一个。每个 Goal 至少有一个专用任务（该任务的 `覆盖` 含这个 Goal、不含其它 Goal；可以另含 REQ）。
- `依赖`：`—` 或逗号分隔的任务 ID；必须可解、无自指、无环。
- `写入`：`—` 或逗号分隔的仓库相对路径；以 `/` 结尾表示目录前缀。两个无依赖路径的任务不能写同一路径。
- `完成标准`：可观察，不是“完成 / works / tested / 功能可用 / 跑一下”。
- `验证`：命令、测试 ID 或产物路径；不能只写“手动看看”。
- 单元格里的字面竖线写成 `\|`。

## 避坑

- 不要复制架构、数据模型或接口契约；任务表只回答“先做什么、依赖谁、写完怎么验”。
- 不要写工时估算或进度状态；v1 不是项目管理工具。
- 不要用一个任务覆盖全部 Goals；Goal 覆盖会被专用任务规则挡住。
- 不要让并行任务写同一个文件；用 `写入` 标出来，闸门会算冲突。

## 检查

- `taskscheck -spec <target>/spec.md -phase 0` 看全貌，`-phase 1` 阻断。
- 每个 Goal 有专用任务；依赖无环；写入无冲突；完成标准与验证非空且可执行。
- Spec 改动后重算 `spec_sha256`，否则 `tasks.stale_plan` 会让闸门变红。
