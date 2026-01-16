---
description: 協助使用者建立 GitHub Pull Request
---

1. **載入 Skill** (Optional): 
    - 檢查是否設定了以下 Skill:
        - `create-gh-pr`
    - 若有，請讀取，使用該 Skill 提供的 PR 模板或規範。
2. **檢查分支狀態**：
   - 執行 `git status` 與 `git log @{u}..` 檢查目前分支狀態與是否有未推送的 commit。
   - 如果有未推送的 commit，詢問使用者是否要執行 `git push`（如果是新分支，需包含 `-u origin <branch>`）。
3. **準備 PR 內容**：
   - 根據 commit messages 生成 PR 的 Title 和 Body。
   - **Title**: 應簡潔明瞭，符合 Conventional Commits (e.g., `feat: add new login page`)。
   - **Body**: 應總結變更內容，列出重點修改，並連結相關 Issue (如 `Closes #123`)。
4. **建立 Pull Request**：
   - 構建 `gh pr create` 指令，建議包含 `--title` 和 `--body` 參數。
   - 若使用者偏好或自動生成內容不完整，可建議執行 `gh pr create --web` 讓使用者在瀏覽器完成。
   - 詢問使用者是否執行。
