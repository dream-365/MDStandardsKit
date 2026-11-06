---
type: architecture
status: Approved
owner: alang
depends_on: []
---

# MDStandardsKit 架构

> 派生自 2026-10-30 批准的 office-hours 设计文档
> `~/.gstack/projects/MDStandardsKit/alang-unknown-design-20261030-124625.md`，
> 并经过 `/plan-eng-review` 的三轮规格评审与一轮外部独立评审。
> 本文只记录已批准的决策，不重复评审过程。

## 1. 问题与目标

文档由 Agent 生成，元数据随生成路径抖动。真正的失败不是字段格式难看，而是：
spec 声明依赖某篇架构文档，Agent 把这篇文档及其状态当作开发依据，却没有任何东西
验证它存在、也没有东西验证它已被批准。失败是静默的：构建全绿、页面都在、检索正常，
代码却照着过期或未批准的架构写完了。

目标：在文档进入下游之前，用机器闸门检查**文档引用图**的完整性，而不是检查单个文件的格式。

## 2. 范围

### v1 包含

- `spec` 类型的两条规则：引用存在、目标状态可依赖。
- `idea` 类型（v1.1）：仓库级类型声明、结构规则、独立内容评分卡；spec 的引用图规则不变。
- 一个二进制，内嵌 OPA 引擎，策略包从仓库读取。
- 两个执行点：CI 合并时、Agent 写完自检时。
- 三阶段 Rollout 与存量违规基线。

### v1 不包含

- `spec` 类型的必填字段校验、status 枚举校验、状态流转校验、复审有效期、审批来源校验（全部 v2）；`idea` 类型的字段、标题、路径校验见 §5.2。
- 循环引用检测（PDS-STD-004 §8.5.6 要求，v1 明确偏离，见 §9）。
- pre-commit 执行点（v1.1）。
- 包管理器分发。
- 多仓库语料合并。

## 3. 组件

```
                 ┌─────────────────────────────┐
   CLI flags ───▶│ mdstandards (single binary) │
                 │                             │
                 │  1. load standard + config  │──▶ standard.json + .mdstandards.json
                 │  2. resolve document type   │──▶ document_types[type].globs
                 │  3. extract corpus          │──▶ 遍历 globs；按类型抽字段
                 │  4. coverage assertions     │──▶ 零文档 / 需要图时零边 → exit 2
                 │  5. evaluate Rego           │──▶ 内嵌 OPA + policy/*.rego
                 │  6. baseline filter         │──▶ 精确三元组匹配
                 │  7. render report           │──▶ table | json | github
                 └──────────────┬──────────────┘
                                │
                                ▼
                      exit 0 / 1 / 2
```

代码布局（单 `main` 包，无内部公共接口）：

| 文件 | 职责 |
|---|---|
| `main.go` | 参数解析、流程编排、策略加载、退出码 |
| `extract.go` | 遍历受管文件、解析 frontmatter、构建语料 |
| `normalize.go` | 路径归一化（D6 + D12） |
| `report.go` | findings 结构、三种输出格式、基线过滤 |
| `policy/*.rego` | spec 两条图规则 + idea 结构规则 + 载入失败守卫 |
| `policy/data/standard.json` | 唯一事实来源；`document_types` 存放类型档案 |

### 3.1 taskscheck（独立二进制）

`taskscheck` 不并入 mdstandards。它读取 `docs/specs/**/spec.md` 与同目录 `tasks.md`，校验任务 ID、依赖无环、Goal 专用覆盖、写冲突、完成标准与 `spec_sha256`。mdstandards 的 0/1/2 与 `-dump-corpus` 契约不变。边界原因见 `docs/designs/tasks-boundary-spike.md`：`document_types` 是仓库级配置，tasks 与 spec 需要在同一仓库共存。

## 4. 数据契约

### 4.1 standard.json

```json
{
  "standard": {
    "standard_version": "1",
    "governed_globs": ["docs/specs/**/*.md", "docs/architecture/**/*.md", "adr/**/*.md"],
    "dependable_statuses": ["Approved", "Published"],
    "identity": "repo_relative_path"
  }
}
```

`document_types` 是可选的类型档案，键是仓库级 `.mdstandards.json` 声明的 `type`：

```json
{
  "document_types": {
    "idea": {
      "globs": ["20*/**/*.md"],
      "required_fields": ["id", "title", "date", "tags", "status"],
      "status_values": ["萌芽", "探索中", "已落地", "搁置"],
      "required_headings": ["一句话描述", "想法", "为什么有意思", "下一步"],
      "id_pattern": "^[A-Z]{8}$",
      "date_field": "date",
      "graph": false,
      "placeholder_markers": [
        "<!-- 一句话说清楚这个想法是什么。写不出来就说明还没想清楚。 -->",
        "<!-- 展开说说：它是什么样子的？怎么运作？ -->",
        "<!-- 它解决了什么问题？或者只是让你好奇？ -->",
        "optional:outside-voice",
        "<!-- 如果想推进，最小的一步是什么？ -->"
      ]
    }
  }
}
```

`spec` 仍用顶层 `governed_globs`；其他类型用档案里的 `globs`。`graph: false` 表示该类型不要求依赖边，
零依赖边覆盖断言不触发。宽 glob 让路径形状规则能看见两层/四层路径，而不是静默忽略。
`optional:outside-voice` 是可选章节「外部的声音」的模板哨兵：整节缺失不报错，哨兵残留则命中占位符规则。

外层的 `standard` 键是必须的：OPA 把数据文档按根合并，不用文件名做命名空间。
没有这一层，`data.standard` 会是 undefined。

### 4.2 corpus

```json
{
  "documents": [
    {"path": "specs/001.md", "type": "spec", "status": "Approved", "owner": "alang",
     "depends_on": ["docs/architecture/sys.md"]}
  ],
  "extract_errors": [
    {"path": "specs/002.md", "reason": "missing frontmatter"}
  ],
  "stats": {"documents_matched": 2, "documents_parsed": 2, "edges_evaluated": 1}
}
```

每个受管文件必须产出一条 document 或一条 extract_error，不会被静默丢弃。

idea 语料额外带 `document_type`、`id`、`title`、`date`、`tags`、`sections`、`missing`、
`date_invalid`、`tags_invalid`。这些字段都带 `omitempty`，只在 `docType == "idea"` 时填充，
所以 spec 的 `-dump-corpus` 输出逐字节不变。

### 4.3 findings

```json
{
  "schema_version": "1",
  "findings": [
    {"rule": "refs.missing_target", "document": "specs/001.md",
     "target": "docs/architecture/sys.md", "severity": "error", "message": "..."}
  ],
  "extract_errors": [],
  "stats": {"documents_matched": 2, "documents_parsed": 2, "edges_evaluated": 1}
}
```

排序固定为 `(rule, document, target)`，保证同一份输入产出字节稳定的输出。

固定规则 ID：

| ID | 含义 |
|---|---|
| `refs.missing_target` | `depends_on` 指向的目标不在受管语料中 |
| `refs.status_not_dependable` | 目标状态不在 `dependable_statuses` 中，或缺失 |
| `extract.parse_error` | 受管文件无法解析（无 frontmatter、YAML 非法等） |
| `policy.load_error` | 策略包加载失败或标准缺失 |

### 4.4 基线

```json
{"rule": "refs.missing_target", "document": "specs/001.md", "target": "docs/architecture/sys.md"}
```

匹配是 `(rule, document, target)` 的精确匹配。路径重命名会让旧条目失配，
从而作为新违规暴露——这是预期行为。解析错误永远不能进入基线。

### 4.5 退出码

| 码 | 含义 |
|---|---|
| 0 | 干净；或 `-phase 0` 的只报告模式；或 `-phase 1` 性能预算降级（≤200 篇受管文档单次超过 2 秒时只警告） |
| 1 | 存在违规且处于阻断状态 |
| 2 | 工具或抽取错误：策略/配置加载失败、版本不支持、零受管文档、需要依赖图的类型零依赖边、未降级时存在解析错误 |

## 5. 规则与失败模式

规则 1（`refs.rego`）：每个 `depends_on` 目标必须存在于语料索引中。
规则 2（`status.rego`）：被依赖文档的 `status` 必须逐字等于
`dependable_statuses` 中的某一项；缺失同样违规。

守卫（`status.rego`）：`data.standard.dependable_statuses` 不存在时产出
`policy.load_error`，防止策略读不到标准却静默全绿。

| 失败场景 | 处理 | 是否静默 |
|---|---|---|
| glob 写错，匹配零个文件 | 覆盖断言 → exit 2 | 否 |
| 语料里零条依赖边 | 覆盖断言 → exit 2（phase ≥ 1） | 否 |
| 策略包编译失败 / 标准缺失 | exit 2 | 否 |
| 单个文件解析失败 | 计入 `extract_errors`，phase ≥ 1 且未降级时 exit 2 | 否 |
| phase 1 单次门禁超过 2 秒且 ≤200 篇受管文档 | 打印警告并降级为只警告；覆盖断言与策略加载错误仍 exit 2；phase 2 不降级 | 否（有警告） |
| Agent 修改 `policy/` 或 `baseline.json` | `-policy-from-head` 从 HEAD 读策略；策略或 baseline 有未提交改动时降级为只报告 | 否（有警告） |
| Agent hook 根本没装 | 无检测手段 | **是，已知关键缺口** |

### 5.2 idea 类型规则

`policy/idea.rego` 只在 `input.document_type == "idea"` 时生效；`.mdstandards.json` 损坏、
type 未知或含未知键时，CLI 在求值前 exit 2。

| 规则 | 含义 |
|---|---|
| `idea.missing_field` | 缺少 `id` / `date` / `tags` / `status`（title 由 `empty_title` 处理） |
| `idea.empty_title` | `title` 为空或只有空白 |
| `idea.id_format` | `id` 不匹配 `^[A-Z]{8}$` |
| `idea.id_duplicate` | 同一语料里出现重复 `id` |
| `idea.date_format` | 日期不是 `YYYY-MM-DD` |
| `idea.date_invalid` | 日期格式正确但不是真实日历日期，或带时间 |
| `idea.path_shape` | 路径不是恰好三层，或年份/日期目录格式不对 |
| `idea.date_dir_mismatch` | frontmatter `date` 与日期目录不一致 |
| `idea.year_dir_mismatch` | 年份目录与 `date` 的年份不一致 |
| `idea.status_invalid` | 状态不在 `status_values` |
| `idea.tags_invalid` | `tags` 不是字符串列表 |
| `idea.heading_missing` | 缺少四个必需二级标题之一（可选的「外部的声音」不算） |
| `idea.placeholder_remaining` | 模板 HTML 注释或 `id: ABCDEFGH` 未替换 |

空正文、`暂无`、`目前说不清，因为……` 都是合法内容，不产生结构 finding。「外部的声音」是可选章节：整节缺失不产生 finding；章节存在时，模板注释仍按 `idea.placeholder_remaining` 检查。

### 5.3 tasks 类型规则

`taskscheck` 是独立二进制，规则不进入 Rego。`tasks.md` 用 `type: tasks`、`spec: ./spec.md` 和 `spec_sha256` 钉住同目录的 Spec；`Approved` / `Implemented` 的 Spec 必须有 tasks.md，`Draft` / `In Review` 可选。

| 规则组 | 规则 ID |
|---|---|
| 文件与生命周期 | `tasks.missing_file`、`tasks.orphan_tasks`、`tasks.invalid_frontmatter`、`tasks.parse_error` |
| Spec 关联 | `tasks.missing_spec`、`tasks.wrong_spec_path`、`tasks.stale_plan`、`spec.no_goals` |
| 任务图 | `tasks.invalid_task_id`、`tasks.duplicate_id`、`tasks.unknown_dependency`、`tasks.self_dependency`、`tasks.dependency_cycle` |
| 覆盖 | `tasks.empty_coverage`、`tasks.unknown_coverage`、`tasks.goal_not_dedicated`、`tasks.uncovered_requirement`（warning） |
| 完成标准 | `tasks.empty_done_criterion`、`tasks.banned_done_criterion`、`tasks.empty_verify`、`tasks.verify_not_executable` |
| 写入 | `tasks.invalid_write_path`、`tasks.write_conflict`、`tasks.case_collision`（warning） |
| 不适用 | `tasks.invalid_not_applicable` |

退出码：0 干净或 `-phase 0`；1 阻断 finding；2 工具/解析错误（包括 `tasks.no_specs`、`tasks.invalid_spec_flag`、`tasks.usage`、`spec.no_goals`）。

## 6. 安全与信任边界

闸门要拦的对象和闸门本身在同一个工作区里。因此：

- 二进制只内嵌 Rego **引擎**，策略包以仓库内的 bundle 为准，避免两套权威来源。
- `-policy-from-head` 从已提交的 `HEAD` 读取策略；工作区策略或 `baseline.json` 有未提交改动时
  打印警告并降级为只报告，不让 Agent 改完规则或基线立刻重跑。
- `policy/` 与 `baseline.json` 已由 `.github/CODEOWNERS` 保护。
- frontmatter 是 Agent 写出的不可信输入：单文件上限 1 MiB，疑似 YAML 别名展开
  （`*` 超过 100 个或 `<<` 超过 10 个）直接拒绝，策略求值 10 秒超时。
- 已知残留风险：`status: Approved` 由同一条写入路径设置，策略无法区分人和 Agent。
  v1 只降低可见性成本，不做来源校验。

## 7. 执行点

两条命令必须逐字相同，保证 CI 和本地看到同一套规则：

```bash
# CI：注解 + 机器可读工件
mdstandards -root . -policy ./policy -phase 1 -format github
mdstandards -root . -policy ./policy -phase 1 -format json -out findings.json

# Agent stop hook
mdstandards -root . -policy ./policy -phase 1 -format json
```

`-policy` 相对**当前工作目录**解析；`-root` 只决定哪些文件受管。
使用 `-policy-from-head` 时，策略目录必须在 `-root` 之内。

## 8. Rollout

- **阶段 0（只报告）**：退出码恒为 0，上传 `findings.json` 和注解，收集存量。
- **阶段 1（基线 + 阻断新增）**：存量写入 `baseline.json`，只对基线之外的新违规退出非零。
- **阶段 2（清空基线）**：基线清空后转为完全阻断。

晋级条件：连续两周的新增违规数为 0，才允许 0 → 1；基线为空才允许 1 → 2。

## 9. 与组织标准（PDS）的偏离

| 偏离 | 原因 | 何时收回 |
|---|---|---|
| `Draft` / `In Review` 引用一律违规，而 PDS-STD-004 §8.5.2 允许"注明状态"后引用 | v1 不解析引用处的注记 | v2 引入注记解析 |
| `Deprecated` / `Archived` 引用一律违规，而 §8.5.8 允许说明替代关系后引用 | 同上 | v2 |
| 不检测循环引用，而 §8.5.6 禁止 | 楔子只做两条规则 | v2 |
| 身份用路径，而 PDS-STD-001 §9.3.3 要求 `doc_id` 跨版本稳定 | 零迁移成本 | v2 评估 doc_id |

## 10. 未决问题

1. Agent 写入路径的具体形态（决定 hook 怎么装）。
2. 谁有权把 `status` 改成 `Approved`，v2 用什么来源校验机制。
3. 是否需要复审有效期（`reviewed_at` + TTL）。
4. 多仓库语料如何合并。

## 11. 已批准的决策

| ID | 决策 |
|---|---|
| D1 | v1 只保留 CI + Agent hook，pre-commit 推迟到 v1.1 |
| D2 | 单 `main` 包 + `policy/` bundle，不拆 `cmd/` 与 `extract/` |
| D3 | 只构建 darwin-arm64 与 linux-amd64 |
| D4 | tag 触发构建，产物挂 GitHub Release，不做包管理器 |
| D5 | Go + 官方 OPA `rego` 包，`GOOS/GOARCH` 交叉编译 |
| D6 | 有界路径归一化，大小写敏感 |
| D7 | CI 与 hook 分开算延迟预算 |
| D8 | 继续独立 OPA 方案，不扩展既有 PDS 校验器 |
| D9/D10 | 状态词汇取 `Approved` + `Published`，逐字匹配 |
| D11 | 仓库级 `.mdstandards.json` 只声明 type；globs 只在 policy 档案里 |
| D12 | `idea` 结构错误阻断，内容评分只诊断；空正文/`暂无` 合法 |
| D13 | 内容评分复用 `qualitycheck`，增加 `guess` 字段和报告标题 flag；评分引擎不重写 |
| D11 | 策略从 `HEAD` 读，策略文件受 CODEOWNERS 保护 |
| D12 | `..` 解析后仍在仓库内则接受，越界才违规 |
| D14 | tasks 用独立 `taskscheck`；Goal 硬覆盖、REQ 只报告；`tasks.md` 用 `spec_sha256` 钉住 `spec.md` |
