# MDStandardsKit

给 Agent 产出的 Markdown 文档装一道元数据闸门。

它检查的不是单个文件的格式，而是**文档引用图**：一篇 spec 声明它依赖某篇架构文档，
那篇文档必须存在，并且处于可被依赖的状态（`Approved` 或 `Published`）。构建全绿、
页面都在，不代表 Agent 没有照着过期的架构写代码。

一个仓库可以同时有多种文档类型：在 `.mdstandards.json` 里列出 `types`，每个文件按自己的
frontmatter `type:` 套用不同规则。

## 作为 Claude Code 插件安装

doc-writer 和 idea-writer 技能随本仓库作为插件发布。安装固定版本：

```bash
claude plugin marketplace add dream-365/MDStandardsKit#v0.3.0 --scope user
claude plugin install doc-writer@mdstandards --scope user
```

安装后在任意项目运行 `/doc-writer:doc-writer` 或 `/doc-writer:idea-writer`。升级或回滚：

```bash
claude plugin marketplace remove mdstandards
claude plugin marketplace add dream-365/MDStandardsKit#v0.3.0 --scope user
claude plugin install doc-writer@mdstandards --scope user
```

插件自带 `mdstandards`、`qualitycheck` 和 `taskscheck` 的下载与构建兜底：支持平台首次运行会从 Release 下载并校验 sha256，失败时用 Go 1.27.1+ 本地构建。

## 快速开始

需要 Go 1.27.1+。可以直接安装，也可以本地构建：

```bash
# 安装到 $GOBIN
go install github.com/dream-365/MDStandardsKit@latest
go install github.com/dream-365/MDStandardsKit/tools/qualitycheck@latest
go install github.com/dream-365/MDStandardsKit/tools/taskscheck@latest

# 或者本地构建
make build

# 只报告，永不阻断
./mdstandards -root . -policy ./policy -phase 0 -format table

# 阻断：有违规 exit 1，工具或抽取错误 exit 2
./mdstandards -root . -policy ./policy -phase 1 -format table

# 给 CI 的机器可读输出
./mdstandards -root . -policy ./policy -phase 1 -format json -out findings.json
```

`v0.2.2` 及更早的 tag 仍声明旧的 module 名；`v0.3.0` 起使用仓库地址，发布后 `@latest` 可用。

试一下自带示例：

```bash
./mdstandards -root ./testdata/example-repo -policy ./policy -phase 1 -format table
```

## 它检查什么

`spec` 类型有两条引用图规则：

| 规则 ID | 含义 |
|---|---|
| `refs.missing_target` | `depends_on` 指向的目标不是受管文档 |
| `refs.status_not_dependable` | 目标状态不是 `Approved` / `Published`，或缺失 |

`idea` 类型有一组结构规则：

| 规则 ID | 含义 |
|---|---|
| `idea.missing_field` | 缺少 `id` / `title` / `date` / `tags` / `status` |
| `idea.empty_title` | `title` 为空或只有空白 |
| `idea.id_format` / `idea.id_duplicate` | `id` 不是 8 位大写字母，或在语料里重复 |
| `idea.date_format` / `idea.date_invalid` | 日期不是 `YYYY-MM-DD`，或不是真实日历日期 |
| `idea.path_shape` / `idea.date_dir_mismatch` / `idea.year_dir_mismatch` | 路径不是 `YYYY/YYYY-MM-DD/<slug>.md`，或日期与目录对不上 |
| `idea.status_invalid` | 状态不在 `萌芽 / 探索中 / 已落地 / 搁置` |
| `idea.tags_invalid` | `tags` 不是字符串列表 |
| `idea.heading_missing` | 缺少四个必需二级标题之一（可选的「外部的声音」不算） |
| `idea.placeholder_remaining` | 模板 HTML 注释或 `id: ABCDEFGH` 未替换 |

外加工程保障：策略加载失败即停、零受管文档即停、需要依赖图的类型零依赖边即停、未降级时解析失败即停。
静默通过是这套工具最想消灭的东西，所以它自己也不允许静默。

## idea 类型

在想法库根目录放 `.mdstandards.json`：

```json
{"type":"idea"}
```

文件路径固定为 `YYYY/YYYY-MM-DD/<slug>.md`。结构错误在记录时阻断；内容评分只诊断，不阻止入库。
空正文、`暂无`、`目前说不清，因为……` 都是允许的。内容评分用 `skills/idea-writer/` 的 rubric，
报告写到 `docs/quality/ideas/<id>/`，不进入受管文档。

正文有四个必需章节（一句话描述、想法、为什么有意思、下一步），以及一个可选章节「外部的声音」。
外部的声音可以整节省略；如果保留章节，模板注释必须替换，否则命中 `idea.placeholder_remaining`。
它当前不进入 rubric，只作为评审上下文。

## 多类型仓库

在仓库根放 `.mdstandards.json`，列出仓库实际出现的全部类型：

```json
{"types":["spec","architecture","changelog","tasks","index"]}
```

每个受管文件按 frontmatter `type:` 路由；类型必须在 `policy/data/standard.json` 的
`document_types` 里有档案，并出现在 `types` 中。缺失、未知、未启用或与路径不兼容的 `type`
都会产生 routing error：`-phase 0` 只报告，`-phase 1` 未降级时 exit 2。

`rule_set` 决定套用哪组规则：

| `rule_set` | 规则 | 说明 |
|---|---|---|
| `spec` | `refs.*` | 引用图检查；`graph: true` 时参与零边断言 |
| `idea` | `idea.*` | 结构规则 |
| `none` | 无 | 只解析、可作为 `depends_on` 目标 |

发现范围始终包含顶层 `governed_globs`；缩小 `types` 不会让顶层 globs 覆盖的文件静默消失。
迁移时先跑 phase 0，修完 routing error 再切 phase 1。legacy `{"type":"idea"}` / `{"type":"spec"}`
继续有效；新 policy 使用 `standard_version: 2`，multi 模式要求 2。

## 执行拆解 tasks.md

`Approved` / `Implemented` 的 Spec 需要在同目录提交 `tasks.md`；`Draft` / `In Review` 可选。它从 `spec.md` 派生，用依赖表达顺序，不复制架构。

```yaml
---
type: tasks
status: In Progress
owner: example
depends_on: []
spec: ./spec.md
spec_sha256: "<spec.md 的 sha256>"
updated_at: 2026-11-02
---
```

任务表固定七列：ID、任务、覆盖、依赖、写入、完成标准、验证。`taskscheck` 检查 ID 唯一、依赖可解无环、每个 Goal 有专用任务、写冲突、完成标准可执行、spec 哈希新鲜；REQ 覆盖只报告。

```bash
./taskscheck -root . -phase 1 -format table
./taskscheck -spec docs/specs/SPEC-002-idea-type/spec.md -phase 1
```

## 文档格式

受管文档需要 YAML frontmatter：

```yaml
---
type: spec
status: Approved
owner: example
depends_on:
  - docs/architecture/system.md
---
```

`depends_on` 里的路径相对仓库根，允许 `./`、`..`（不得越出仓库）、反斜杠和重复斜杠；
大小写敏感。`type` 在多类型模式下是 load-bearing 字段：它决定套用哪组规则，修改它需要评审保护。

## 退出码

| 码 | 含义 |
|---|---|
| 0 | 干净，或 `-phase 0` 的只报告模式，或 `-phase 1` 性能预算降级（≤200 篇受管文档单次超过 2 秒时只警告） |
| 1 | 存在违规且处于阻断状态 |
| 2 | 策略/配置加载失败、版本不支持、零受管文档、phase ≥ 1 零解析、需要依赖图的类型零依赖边、未降级时存在解析或 routing 错误 |

`-phase 2` 不参与性能预算降级；覆盖断言与策略加载错误在降级时仍退出 2。

## 目录

```
.
├── main.go / extract.go / normalize.go / report.go   单 main 包
├── *_test.go
├── .claude-plugin/          Claude Code 插件与 marketplace manifest
├── bin/doc-writer-tool      技能调用的工具 dispatcher
├── skills/doc-writer/       Spec 写作技能
├── skills/idea-writer/      idea 写作与评分技能
├── scripts/shim-test.sh     dispatcher 行为测试
├── policy/                  Rego 规则与标准（唯一事实来源）
├── testdata/example-repo/   可运行的 Spec 示例语料
├── testdata/idea-repo/      可运行的 idea 示例语料
├── testdata/mixed-repo/     可运行的多类型示例语料
├── tools/qualitycheck/      Spec / idea 内容质量评分卡审计器
├── tools/taskscheck/         tasks.md 执行拆解闸门
├── docs/
│   ├── architecture/        架构文档
│   ├── designs/             设计记录
│   ├── quality/             质量评分卡报告
│   └── specs/               v1 规格
├── Makefile
└── .github/workflows/       CI 与发布
```

`mdstandards` 代码刻意保持单包：v1 只有几百行、没有第二个调用方，拆成 `cmd/` + `internal/`
属于过早抽象（决策 D2）。

## 设计文档

- [架构](docs/architecture/mdstandards-architecture.md)
- [v1 规格](docs/specs/SPEC-001-mdstandards/spec.md)
- [doc-writer 技能设计](docs/designs/doc-writing-skill.md)
- [doc-writer 插件分发设计](docs/designs/doc-writer-plugin-distribution.md)
- [Spec 内容质量评分卡设计](docs/designs/spec-quality-scorecard.md)
- [idea 类型设计](docs/designs/idea-document-type.md)
- [执行拆解设计](docs/designs/spec-execution-tasks.md)
- [tasks 边界 spike](docs/designs/tasks-boundary-spike.md)
- [质量报告索引](docs/quality/README.md)

## 已知限制

- 真实语料还没有 `depends_on`，直接对既有仓库运行会因"零依赖边"退出 2，这是设计行为。
- `status: Approved` 由同一条 Agent 写入路径设置，v1 不校验审批来源。
- 不检测循环引用，与 PDS-STD-004 §8.5.6 存在偏离。
- Agent hook 未安装时是静默无操作，当前无检测手段。
- 插件分发只支持 darwin/arm64、darwin/amd64、linux/amd64、linux/arm64；Windows 不在 v0.1 范围内。
- `taskscheck` v1 只硬覆盖 Goal，REQ 覆盖只报告；没有 baseline；`Draft` / `In Review` 的 Spec 不强制 `tasks.md`。
- 多类型模式要求每个受管文件带可识别的 `type:`；`type:` 与 `.mdstandards.json` 的 `types` 都是 load-bearing，需要 CODEOWNERS/评审保护。
- tasks 由独立二进制检查，不并入 mdstandards 的规则系统；`type: tasks` 只让 mdstandards 解析该文档。

## 命令

```bash
make build       # 构建 mdstandards
make taskscheck  # 构建 taskscheck
make test        # go test ./...
make vet         # go vet ./...
make check       # vet + test + mdstandards 与 taskscheck 自检
make clean       # 删除二进制
```

## 许可证

[MIT](LICENSE)
