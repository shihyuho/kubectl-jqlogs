---
name: create-git-commit
description: 根據 Git Repository 的風格建立 Commit。當使用者要求，或 Agent 需要建立 Commit 時，必須使用此技能。
allowed-tools:
  - Bash(git diff:*)
  - Bash(git log:*)
  - Bash(sed:*)
  - Bash(tr:*)
  - Bash(sort:*)
  - Bash(xargs:*)
---

# 建立 Git Commit

## 工作流程

### 1. 收集資訊

若有 Staged 的變更：

```bash
git diff --staged
git log --oneline -10
```

若沒有 Staged 的變更：

```bash
git diff
git log --oneline -10
```

### 2. 判斷

從過去的 Commit 歷史紀錄判斷以下項目。若無法判斷，請向使用者確認。

#### 語言判斷

| 語言 | 模式 |
| --- | --- |
| 繁體中文 | 包含繁體中文 |
| 英文 | 只包含英文 |
| 其他 | 上述以外 |

#### 風格判斷

| 風格 | 模式 |
| --- | --- |
| Conventional Commits | `feat:`, `fix:` 等 |
| 簡易格式 | 上述以外 |

#### Scope 判斷

判斷過去 10 筆 Commit 是否包含 Scope（例如：`feat(auth):`)。

若包含，使用以下指令取得 Scope 列表：

```bash
git log --oneline -100 | sed -n 's/^[a-f0-9]* [^(:]*(\([^)]*\)):.*/\1/p' | tr ',' '\n' | sed 's/^ *//' | sort -u | xargs | sed 's/ /, /g'
```

### 3. 產生 Commit Message

根據判斷結果產生 Commit Message。

- **標題（第一行）**: What（做了什麼）的簡潔說明，50 字元以內
- **內文（Body）**: Why（為什麼這麼做）的補充說明
- **程式碼參考**: 提及程式碼或路徑時，使用反引號（`）包起來
- **語言**: 配合判斷出的語言
- **風格**: 遵循判斷出的風格（[Conventional Commits](references/conventional-commits.md)）
- **Scope**: 從 Scope 列表中選擇合適的項目（若適用）

### 4. 呈現與確認

將產生的 Commit Message 呈現給使用者確認。

**格式範例：**

```
建議的 Commit Message:

feat(user): 新增使用者大頭貼上傳功能

- 實作上傳 API endpoint
- 新增圖片格式驗證
```

詢問使用者是否滿意，或是否需要修改。

### 5. 執行提交 (可選)

如果使用者確認滿意，可以使用以下指令協助提交：

```bash
git commit -m "feat(user): 新增使用者大頭貼上傳功能" -m "- 實作上傳 API endpoint" -m "- 新增圖片格式驗證"
```
