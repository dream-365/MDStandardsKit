# MDStandardsKit

给 Agent 产出的 Markdown 文档装一道元数据闸门。

它检查的不是单个文件的格式，而是**文档引用图**：一篇 spec 声明它依赖某篇架构文档，
那篇文档必须存在，并且处于可被依赖的状态（`Approved` 或 `Published`）。构建全绿、
页面都在，不代表 Agent 没有照着过期的架构写代码。

## 快速开始

```bash
make build

# 只报告，永不阻断
./mdstandards -root . -policy ./policy -phase 0 -format table

# 阻断：有违规 exit 1，工具或抽取错误 exit 2
./mdstandards -root . -policy ./policy -phase 1 -format table

# 给 CI 的机器可读输出
./mdstandards -root . -policy ./policy -phase 1 -format json -out findings.json
```

试一下自带示例：

```bash
./mdstandards -root ./testdata/example-repo -policy ./policy -phase 1 -format table
```

## 它检查什么

v1 只有两条规则：

| 规则 ID | 含义 |
|---|---|
| `refs.missing_target` | `depends_on` 指向的目标不是受管文档 |
| `refs.status_not_dependable` | 目标状态不是 `Approved` / `Published`，或缺失 |

外加三条工程保障：策略加载失败即停、零受管文档 / 零依赖边即停、解析失败即停。
静默通过是这套工具最想消灭的东西，所以它自己也不允许静默。

## 文档格式

受管文档需要 YAML frontmatter：

```yaml
---
type: spec
status: Approved
owner: alang
depends_on:
  - docs/architecture/system.md
---
```

`depends_on` 里的路径相对仓库根，允许 `./`、`..`（不得越出仓库）、反斜杠和重复斜杠；
大小写敏感。

## 退出码

| 码 | 含义 |
|---|---|
| 0 | 干净，或 `-phase 0` 的只报告模式 |
| 1 | 存在违规且处于阻断状态 |
| 2 | 策略加载失败、版本不支持、零受管文档、零依赖边、存在解析错误 |

## 目录

```
.
├── main.go / extract.go / normalize.go / report.go   单 main 包
├── *_test.go
├── policy/                  Rego 规则与标准（唯一事实来源）
├── testdata/example-repo/   可运行的示例语料
├── docs/
│   ├── architecture/        架构文档
│   └── specs/               v1 规格
├── Makefile
└── .github/workflows/       CI 与发布
```

代码刻意保持单包：v1 只有几百行、没有第二个调用方，拆成 `cmd/` + `internal/`
属于过早抽象（决策 D2）。

## 设计文档

- [架构](docs/architecture/mdstandards-architecture.md)
- [v1 规格](docs/specs/mdstandards-v1.md)

## 已知限制

- 真实语料还没有 `depends_on`，直接对既有仓库运行会因"零依赖边"退出 2，这是设计行为。
- `status: Approved` 由同一条 Agent 写入路径设置，v1 不校验审批来源。
- 不检测循环引用，与 PDS-STD-004 §8.5.6 存在偏离。
- Agent hook 未安装时是静默无操作，当前无检测手段。

## 命令

```bash
make build   # 构建
make test    # go test ./...
make vet     # go vet ./...
make check   # vet + test + 对自身仓库自检
make clean   # 删除二进制
```
