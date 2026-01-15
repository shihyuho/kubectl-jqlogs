---
name: makefile-executor
description: Analyzes Makefiles to understand their targets and executes the appropriate make command based on user intent. Specializes in parsing self-documenting Makefiles that use the `##` comment convention.
allowed-tools:
  - Read
  - Bash(make:*)
---

# Makefile

本 skill 旨在分析專案根目錄下的 `Makefile`，理解其中定義的 `targets`，並在適當的時機執行它們。

## 核心能力

- **解析 `Makefile`**: 讀取 `Makefile` 內容，識別可用的 `targets`、變數和它們之間的依賴關係。
- **意圖映射**: 將使用者的指令 (例如「安裝相依套件」、「用 Java 17 跑測試」) 映射到最符合的 `make` 指令，包含必要的參數。
- **執行指令**: 安全地執行 `make` 指令。

## 工作流程

1.  **探索 (Discovery)**: 檢查 `Makefile` 是否存在。若存在，優先執行 `make help` 來取得專案作者提供的指令清單。
2.  **分析 (Analysis)**: 如果 `make help` 不足，則直接讀取 `Makefile` 內容。本 skill 會：
    -   尋找 `target: ## comment` 格式來理解 `target` 的用途。
    -   尋找 `VARIABLE ?= default_value` 格式來識別可被覆寫的參數。
    -   分析 `target: dependency1 dependency2` 格式來理解執行順序。
3.  **映射 (Mapping)**: 根據分析結果，將使用者的需求轉換為一個具體的 `make` 指令，包含必要的 `OPTION=value` 前綴。
4.  **執行 (Execution)**: 在執行前向使用者確認即將運行的完整指令及其預期作用，然後透過 `Bash` tool 執行。

## `Makefile` 寫法分析與慣例

本 skill 的設計基於對良好實踐 `Makefile` 的分析，特別是你提供的範例：

-   **自文件化註解 (`##`)**:
    -   **格式**: `target-name: ## 說明文字...`
    -   **作用**: 讓 `Makefile` 變得可讀且易於維護的關鍵。本 skill 會利用此資訊來理解 `target` 的意圖。

-   **`help` Target**:
    -   **作用**: 一個好的 `Makefile` 應提供 `help` target，自動解析並印出說明選單。
    -   **本 skill 的應用**: skill 會優先嘗試執行 `make help`，這是最可靠的探索方式。

-   **分類註解 (`##@`)**:
    -   **格式**: `##@ Category Name`
    -   **作用**: 在 `help` 輸出中對 `targets` 進行分組，讓 skill 能更好地理解指令的領域（例如：開發、部署）。

-   **變數與參數 (`?=` and `$(call ...)` )**:
    -   **格式**: `VAR ?=` 或 `ifeq ($(strip $(VAR)),)`
    -   **作用**: `?=` 表示如果變數未被設定，則給予預設值，這暗示了此變數可由外部傳入。`ifeq` 檢查則常用於強制要求傳入變數。
    -   **本 skill 的應用**: skill 會識別這些模式，來判斷執行 `make` 指令時是否可以或必須提供 `VAR=value` 這樣的參數。

-   **目標依賴 (`target: dependency`)**:
    -   **格式**: `redeploy: undeploy deploy`
    -   **作用**: 定義了 `target` 的執行順序。在執行 `redeploy` 之前，`make` 會先確保 `undeploy` 和 `deploy` 都被執行。
    -   **本 skill 的應用**: 這幫助 skill 理解一個指令可能包含的多個步驟。

## 執行範例

**情境**: 使用者說：「我要用 Java 17 跑測試。」

**Skill 的思考與行動流程**:

1.  **意圖分析**: 使用者的意圖是「跑測試」，並帶有參數「Java 17」。
2.  **探索 `Makefile`**:
    -   偵測到 `Makefile`，執行 `make help`。
3.  **分析 `help` 輸出與檔案內容**:
    -   從 `help` 中找到 `test: ## Clean and test the compiled code.`，確認 `make test` 是正確指令。
    -   從 `help` 的 `Usage: OPTION=<value> make <target>` 推斷可以傳遞參數。
    -   進一步分析 `Makefile` 原始碼，找到 `test` target 的指令 `mvn test -e $(call java_version)`。
    -   追溯 `java_version` function，發現它使用了 `JAVA` 變數：`$(if $(filter-out "",$(JAVA)),-D'java.version=$(JAVA)',)`。
4.  **映射 Target 與參數**:
    -   確認 `make test` 是正確的 target。
    -   確認可以透過設定 `JAVA` 變數來指定 Java 版本。
    -   組合出最終指令：`JAVA=17 make test`。
5.  **準備執行**:
    -   (向使用者確認) "我將執行 `JAVA=17 make test` 來用 Java 17 執行測試，是否繼續？"
6.  **執行**:
    -   `run_shell_command("JAVA=17 make test")`
