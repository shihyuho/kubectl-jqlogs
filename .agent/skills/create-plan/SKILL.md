---
name: create-plan
description: Create a concise plan. Use when a user explicitly asks for a plan related to a coding task.

---

# 建立計畫

## 目標

將使用者的提示轉化為一個**單一、可執行的計畫**，計劃中的每一個項目都應該是一個清晰、簡潔的待辦事項，方便追蹤與執行。並在最終的助理訊息中呈現。

## 最小工作流程

在整個工作流程中，以唯讀模式操作。不要寫入或更新檔案。

1.  **快速掃描上下文**
  *   閱讀 `README.md` 和任何明顯的文件 (`docs/`, `CONTRIBUTING.md`, `ARCHITECTURE.md`)。
  *   瀏覽相關檔案（最有可能被修改的檔案）。
  *   識別限制（語言、框架、CI/測試指令、部署形式）。

2.  **僅在受阻時提出後續問題**
  *   **最多提出 1-2 個問題**。
  *   只有在沒有答案就無法負責任地制定計畫時才提出；優先選擇多選題。
  *   如果不確定但沒有受阻，則做出合理假設並繼續。

3.  **使用以下模板建立計畫**
  *   首先用**1 個簡短的段落**描述意圖和方法。
  *   清楚指出**在範圍內 (in scope)** 和**不在範圍內 (not in scope)** 的項目，越簡短越好。
  *   然後提供一個**小的檢查清單**（預設 6-10 個項目）。
    *   每個清單項目都應該是一個具體的行動，並且在有幫助時提及檔案/指令。
    *   **讓項目原子化且有順序**：探索 → 變更 → 測試 → 推出。
    *   **動詞優先**：「新增…」、「重構…」、「驗證…」、「部署…」。
  *   在適用時，至少包含一個用於**測試/驗證**的項目和一個用於**邊緣案例/風險**的項目。
  *   如果存在未知數，包含一個簡短的**未解決問題 (Open questions)** 部分（最多 3 個）。

4.  **不要在計畫前面加上元解釋；僅按照模板輸出計畫**

## 計畫模板（請嚴格遵循）

```markdown
# Plan

<1–3 sentences: what we're doing, why, and the high-level approach.>

## Scope
- In:
- Out:

## Action items
[ ] <Step 1>
[ ] <Step 2>
[ ] <Step 3>
[ ] <Step 4>
[ ] <Step 5>
[ ] <Step 6>

## Open questions
- <Question 1>
- <Question 2>
- <Question 3>
```

## 語系指南

以繁體中文輸出

## 清單項目指南
好的清單項目：
- 指向可能的檔案/模組：src/..., app/..., services/...
- 命名具體的驗證：「執行 npm test」、「為 X 新增單元測試」
- 在相關時包含安全的推出方式：功能開關 (feature flag)、遷移計畫 (migration plan)、回溯筆記 (rollback note)

避免：
- 模糊的步驟（「處理後端」、「處理身份驗證」）
- 過多的微小步驟
- 編寫程式碼片段（讓計畫保持與實作無關）

## 與現有 Skills 的整合

可搭配使用

- `makefile`: 了解 Makefile 的使用方法
