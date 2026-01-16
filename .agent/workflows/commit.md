---
description: 根據目前的 git diff 產生 commit message
---
1. **載入 Skill** (Optional): 
    - 檢查是否設定了以下 Skill:
        - `create-git-commit`
    - 若有，請讀取，並遵循其規範。
2. 執行 `git diff` 查看未暫存的變更。
3. 執行 `git diff --cached` 查看已暫存的變更。
4. 根據變更內容，分析修改的意圖（是新增功能、修補錯誤還是重構程式碼）。
5. 產生符合 Conventional Commits 規範以及上述 Skill 規範的 commit message。
6. 顯示結果給使用者，不需直接執行 commit 指令，除非使用者要求。
