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

- 两条规则：引用存在、目标状态可依赖。
- 一个二进制，内嵌 OPA 引擎，策略包从仓库读取。
- 两个执行点：CI 合并时、Agent 写完自检时。
- 三阶段 Rollout 与存量违规基线。

### v1 不包含

- 必填字段校验、status 枚举校验、状态流转校验、复审有效期、审批来源校验（全部 v2）。
- 循环引用检测（PDS-STD-004 §8.5.6 要求，v1 明确偏离，见 §9）。
- pre-commit 执行点（v1.1）。
- 包管理器分发。
- 多仓库语料合并。

## 3. 组件

```
                 ┌─────────────────────────────┐
   CLI flags ───▶│ mdstandards (single binary) │
                 │                             │
                 │  1. load standard.json      │──▶ 读 policy/data/standard.json
                 │  2. extract corpus          │──▶ 遍历 governed_globs
                 │  3. coverage assertions     │──▶ 零文档 / 零边 → exit 2
                 │  4. evaluate Rego           │──▶ 内嵌 OPA + policy/*.rego
                 │  5. baseline filter         │──▶ 精确三元组匹配
                 │  6. render report           │──▶ table | json | github
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
| `policy/*.rego` | 两条规则 + 载入失败守卫 |
| `policy/data/standard.json` | 唯一事实来源 |

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
| 0 | 干净；或 `-phase 0` 的只报告模式 |
| 1 | 存在违规且处于阻断状态 |
| 2 | 工具或抽取错误：策略加载失败、版本不支持、零受管文档、零依赖边、存在解析错误 |

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
| 单个文件解析失败 | 计入 `extract_errors`，phase ≥ 1 时 exit 2 | 否 |
| Agent 修改 `policy/` 或 `baseline.json` | `-policy-from-head` 从 HEAD 读策略；脏工作区降级为只报告 | 否（有警告） |
| Agent hook 根本没装 | 无检测手段 | **是，已知关键缺口** |

## 6. 安全与信任边界

闸门要拦的对象和闸门本身在同一个工作区里。因此：

- 二进制只内嵌 Rego **引擎**，策略包以仓库内的 bundle 为准，避免两套权威来源。
- `-policy-from-head` 从已提交的 `HEAD` 读取策略；工作区策略有未提交改动时
  打印警告并降级为只报告，不让 Agent 改完规则立刻重跑。
- `policy/`、`standard.json`、`baseline.json` 应受 CODEOWNERS 保护（**尚未在本仓库落地**）。
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
| D11 | 策略从 `HEAD` 读，策略文件受 CODEOWNERS 保护 |
| D12 | `..` 解析后仍在仓库内则接受，越界才违规 |
