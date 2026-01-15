---
name: boy-scout-rule
description: Use when modifying existing files, refactoring, improving code quality, or touching legacy code by applying the Boy Scout Rule to leave code better than you found it.
allowed-tools:
  - Read
  - Edit
  - Grep
  - Bash
---

# 童子軍原則 (Boy Scout Rule)

> 「讓營地比你發現時更乾淨。」

**永遠要讓程式碼比你發現時更好。** 當你接觸一個檔案時，就做一些漸進式的改善。

## 該改善什麼

**程式碼品質**:

- 移除無用程式碼 (被註解的區塊、未使用的函式)
- 修復你所接觸檔案中的 linting 問題
- 改善不清晰的命名 (`x`, `temp`, `data` → 改為描述性名稱)
- 增加型別註記 (Type annotations)
- 將魔法數字 (magic numbers) 提取為具名常數
- 簡化複雜的邏輯
- 增加缺失的錯誤處理
- 更新過時的註解
- 修正格式
- 移除未使用的 imports/變數
- 整合重複的程式碼

## 不該做什麼

- ❌ 大規模且無關的重構
- ❌ 在沒有測試的情況下改變行為
- ❌ 修復檔案中的所有問題 (保持專注)
- ❌ 在沒有測試的情況下進行破壞性變更 (Breaking changes)
- ❌ 過早優化 (Premature optimization)
- ❌ 變更不相關的部分

## 流程

1.  **變更前**: 閱讀檔案，記錄下明顯的問題，執行 linter
2.  **進行主要變更**: 實作功能/修復問題，撰寫測試
3.  **進行改善**: 修復 linting，改善命名，增加型別，提取常數，移除無用程式碼
4.  **執行驗證**: `mvn clean verify` 或 `gradle check`
5.  **文件化**: 在 commit message 中包含童子軍原則的改善項目

## Coding Style

若有發現 `makefile` skill，則檢查是否有程式碼風格驗證的 target，若有則執行，並驗證 git status 是否為 clean。

## Commit Message

若有發現 `create-git-commit` skill，則檢查是否符合該 skill 的要求。

## 改善範例

**Before**:

```java
// Item.java
// public class Item {
//     public double price;
// }

public double calculateTotal(List<Item> items) { // 命名不佳
    double t = 0; // 不好的變數名
    for (int i = 0; i < items.size(); i++) { // 傳統 for 迴圈
        t += items.get(i).price * 1.08; // 魔法數字
    }
    return t;
}
```

**After**:

```java
// Item.java
// public class Item {
//     public double getPrice() { return price; }
// }

private static final double TAX_RATE = 1.08;

public double calculateTotal(List<Item> items) {
    return items.stream()
                .mapToDouble(item -> item.getPrice() * TAX_RATE)
                .sum();
}
```

**改善**: 使用 final 常數、改善命名、使用現代的 Stream API、簡化邏輯。

## 關鍵規則

- 將改善與主要變更放在同一個 commit 中
- 專注於「爆炸半徑」(靠近你變更的程式碼)
- 優先考慮可讀性，而非炫技
- 改善後務必執行完整的測試套件
- 當不確定商業邏輯時，尋求協助
- 小的改善 > 完美的重構

## 檢查清單

- [ ] 我已移除所接觸檔案中的無用程式碼
- [ ] 我已修復 linting 問題
- [ ] 我已改善變更處的命名
- [ ] 我已在需要的地方增加型別註記
- [ ] 我已提取魔法數字
- [ ] 我已更新過時的註解
- [ ] 我已移除未使用的 imports
- [ ] 我已簡化複雜的邏輯
- [ ] 我已增加錯誤處理
- [ ] 所有測試都通過
- [ ] 沒有新的 linting 錯誤
- [ ] 已在 commit message 中記錄

## 整合

**實作期間**: 進行變更，應用童子軍原則，驗證，然後一起 commit

**Code Review 期間**: 尋找應用童- Scout 原則的機會，並表揚好的實踐

**修復 Bug 期間**: 修復 bug，改善周邊程式碼，增加測試，進行清理

## 記住

### 漸進式改善，而非追求完美

- 小的改善會隨著時間累積
- 每個接觸到的檔案都是一個機會
- 當一個好的程式碼庫守護者
- 當對較大的改善不確定時，尋求 review

## 與現有 Skills 的整合

可搭配使用

- `makefile`: 了解 Makefile 的使用方法
- `create-git-commit`: 建立 git commit message  
