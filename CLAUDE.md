# CLAUDE.md

## 军规

- 所有改动先改需求，再改代码。需求包括 Spec、设计文档、模板、标准与 rubric 等定义行为的文件；需求未更新，不改实现。

## 沟通

- 用户常在手机上操作：回复尽量简洁明了，先给结论，少铺垫。
- 决策问题保留必要选项与推荐，但不要重复大段背景。

## Testing

- Run: `go test ./...`
- Full check: `make check`（build + vet + test + mdstandards 自检）
