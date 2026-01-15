---
name: test-driven-development
description: Use when writing new functions, adding features, fixing bugs, or refactoring by applying TDD principles - write failing tests before implementation code, make them pass, then refactor.
allowed-tools:
  - Write
  - Read
  - Edit
  - Bash
  - Grep
---

# 測試驅動開發 (TDD)

對所有程式碼變更，遵循「紅 → 綠 → 重構」的循環。

## TDD 循環

1.  **紅 (RED)**: 撰寫一個會失敗的測試。
2.  **綠 (GREEN)**: 撰寫最精簡的程式碼使其通過。
3.  **重構 (REFACTOR)**: 改善程式碼品質。

### 對每個需求重複此循環

## 何時應用 TDD

✅ **永遠對以下情況使用 TDD：**

- 新的函式/方法
- 新功能
- Bug 修復 (先重現 bug)
- 重構現有程式碼
- API 變更

❌ **對以下情況可跳過 TDD：**

- UI 樣式微調
- 設定檔變更
- 文件更新

## 流程 (Java 範例)

### 1. 先寫失敗的測試 (RED)

```java
// src/test/java/com/example/CalculatorTest.java
import static org.junit.jupiter.api.Assertions.assertEquals;
import org.junit.jupiter.api.Test;

class CalculatorTest {
    @Test
    void calculatesTotalWithTax() {
        Calculator calculator = new Calculator();
        double[] prices = {100.0, 200.0};
        // 預期 300 * 1.08 = 324
        double expected = 324.0; 
        
        // 此時 calculateTotal 方法甚至還不存在
        double result = calculator.calculateTotal(prices);
        
        assertEquals(expected, result, 0.001);
    }
}

// 執行測試 - 編譯失敗或執行失敗 (RED)
```

### 2. 撰寫最精簡的程式碼使其通過 (GREEN)

```java
// src/main/java/com/example/Calculator.java
public class Calculator {
    public double calculateTotal(double[] prices) {
        double sum = 0;
        for (double price : prices) {
            sum += price;
        }
        return sum * 1.08; // 先用魔法數字讓測試通過
    }
}

// 再次執行測試 - 應該會通過 (GREEN)
```

### 3. 重構 (REFACTOR)

現在測試通過了，你可以安全地改善程式碼，例如提取常數、改善命名等。

```java
// Refactored Calculator.java
public class Calculator {
    private static final double TAX_RATE = 1.08;

    public double calculateTotal(double[] prices) {
        double sum = 0;
        for (double price : prices) {
            sum += price;
        }
        return sum * TAX_RATE;
    }
}

// 每次重構後都要再次執行測試，確保依然是 GREEN
```

## 關鍵規則

- 測試**必須**先失敗 (以驗證測試本身是有效的)
- 一個測試對應一個需求
- 測試行為，而非實作細節
- Commit 前務必執行**完整**的測試套件
- **絕不**跳過失敗的測試

## 常見陷阱

- 在寫測試前就先寫實作
- 沒寫實作，測試卻通過了 (假陽性, false positive)
- 測試實作細節，而非行為
- 沒有先執行測試以驗證它會失敗

## 驗證指令 (Maven/Gradle)

- 若在根目錄有找到 `Makefile` 檔案，優先讀取並試著找到可以執行測試的方法，例如 `make test`
- 若沒有找到 `Makefile` 檔案，或其中沒有找到可以執行測試的 target，則使用 `mvn test` 或 `gradle test` 來執行測試

## 與現有 Skills 的整合

可搭配使用

- `makefile`: 了解 Makefile 的使用方法

