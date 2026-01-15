---
name: create-gh-pr
description: Guides the creation of high-quality Pull Requests that adhere to repository standards, covering template discovery, drafting descriptions, and executing PR creation commands.
allowed-tools:
  - Read
  - Bash(ls:*)
  - Bash(gh:*)
---

# Pull Request 建立助手

此技能旨在引導建立符合儲存庫（Repository）標準的高品質 Pull Request。

## 工作流程

請遵循以下步驟來建立 Pull Request：

### 1. 尋找範本

在儲存庫中搜尋 Pull Request 範本。

- 檢查 `.github/pull_request_template.md`
- 檢查 `.github/PULL_REQUEST_TEMPLATE.md`
- 如果存在多個範本（例如在 `.github/PULL_REQUEST_TEMPLATE/` 目錄下），請詢問使用者該使用哪一個，或根據情境選擇最合適的一個（例如 `bug_fix.md` 對比 `feature.md`）。

### 2. 閱讀範本

閱讀已識別範本檔案的內容。

### 3. 草擬描述

撰寫一份嚴格遵循範本結構的 PR 描述。

- **標題**：保留範本中的所有標題。
- **檢查清單**：檢視每個項目。若已完成請標記 `[x]`。若項目不適用，請保持未勾選或標記為 `[ ]`（視範本說明而定），或者如果範本允許彈性調整則移除它（但為了透明度，建議保留未勾選狀態）。
- **內容**：填寫各個區塊，清楚且簡潔地總結您的變更。
- **相關議題 (Issues)**：連結任何此 PR 修復或相關的議題 ([GitHub Keywords](references/github-keywords.md))

### 4. 建立 PR

- **標題**：如果儲存庫有採用，請確保標題遵循 Conventional Commits 格式（例如 `feat(ui): add new button`、`fix(core): resolve crash`）。
- **使用工具**
	- 優先使用 GitHub MCP Server 的 `create_pull_request` tool
	- 若無GitHub MCP Server，則使用 `gh` CLI 工具來建立 PR。
		```bash
		gh pr create --title "type(scope): succinct description" --body "..."
		```
	- 若以上皆無，則提供 Markdown 格式的 PR 內文，供使用者複製

### 5. 與現有 Skills 的整合

可搭配使用

- `create-git-commit`: 產生 git commit message，可參考 Conventional Commits 格式 

## 原則

- **合規性**：絕不要忽略 PR 範本。它的存在是有原因的。
- **完整性**：填寫所有相關的章節。
- **準確性**：不要勾選您尚未完成的任務項目。
