---
type: changelog
status: In Review
owner: alang
depends_on: []
updated_at: 2026-11-02
---

# Changelog: SPEC-002 idea 类型

## [2026-11-02] - 初稿

- **修改类型**: Added
- **影响范围**: 新增 `docs/specs/SPEC-002-idea-type/spec.md`；`policy/idea.rego`、`extract.go`、`main.go`、`tools/qualitycheck`、`skills/idea-writer/`、`.claude-plugin/`。
- **变更原因**: 想法库有模板但没有结构验证和内容评估；需要给 mdstandards 增加 `idea` 类型。
- **具体内容**:
  - 仓库级 `.mdstandards.json` 声明 `{"type":"idea"}`；globs 只在 policy 的 `document_types.idea` 里。
  - 结构规则覆盖 ID、日期、路径、状态、tags、标题、占位符、必需章节。
  - `date` 同时兼容 YAML 字符串和 `time.Time`；非零时钟的 datetime 报错。
  - 内容评分复用 `qualitycheck`，增加 `guess` 字段和 `-title` / `-target-label`。
  - 结构错误阻断；内容分只诊断，不阻止入库。
- **状态变更**: 新建，`In Review`，等待独立评审。
