# 第 0 章：标题与元数据

目标：让读者在 3 秒内知道这是什么、谁负责、能不能当依据。

## 标题

公式：`[系统/模块名] [核心动作/变更内容] 技术方案`。

- 避免："文档系统重构"（太笼统）。
- 可用："MDStandardsKit v1 元数据闸门 技术方案"。

## 元数据

YAML frontmatter 供闸门读取：`type`、`status`、`owner`、`depends_on`、`updated_at`。
正文表格供人读：

```markdown
| 字段 | 内容 |
| --- | --- |
| 标题 | <标题> |
| 主笔 | <author> |
| 评审人 | <reviewers> |
| 创建日期 | <YYYY-MM-DD> |
| 目标版本 | <version 或 无> |
| 关联资源 | <links 或 无> |
```

- 评审人必须包含下游依赖方负责人，不能只写自己。
- 状态与版本要显式：`status: In Review` 表示还不能作为编码依据。
- 修改核心内容后同步 `updated_at` 与 `changelog.md`。

## 检查

- 标题符合公式，不是一句话描述。
- frontmatter 五个字段齐全；`depends_on` 目标真实存在且状态可依赖。
- 正文表格有标题、主笔、评审人、创建日期。
