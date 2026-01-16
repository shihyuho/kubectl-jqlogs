# kubectl-jqlogs

[English](README.md) | [繁體中文](README-zh_TW.md)

**讓你的 JSON日誌再次易讀。**

`kubectl-jqlogs` 的運作方式與 `kubectl logs` 完全相同，但內建了 `jq` 引擎，可自動美化 JSON 輸出。無需管線 (pipes)，無需額外工具，開箱即由，提供乾淨、可查詢的日誌。

![License](https://img.shields.io/github/license/shihyuho/kubectl-jqlogs)
![Release](https://img.shields.io/github/v/release/shihyuho/kubectl-jqlogs)


## ✨ 特色

- **混合日誌處理**：無縫處理混合內容。JSON 日誌會被格式化，而標準文字日誌則照原樣列印。不再有 "parse error" 讓你的管線崩潰。
- **智慧查詢語法**：簡化的欄位選擇語法。使用 `.level .msg` 代替複雜的字串插值。
- **[gojq](https://github.com/itchyny/gojq) 的強大功能**：
  - **擴充功能**：支援 `--yaml-output` 將 JSON 日誌轉換為 YAML。
  - **精確度**：使用任意精確度處理大數字和繁重計算 (使用 `math/big`)。
  - **標準合規**：完全實作純 jq 語言，無外部 C 依賴。
- **原生體驗**：支援所有標準 `kubectl logs` 旗標和自動補全。
- **原始輸出**：內建支援原始輸出 (`-r`) 和彩色輸出。

## 安裝

### 使用 Krew (推薦)

外掛程式發布後，您可以透過 [Krew](https://krew.sigs.k8s.io/) 安裝：

```bash
kubectl krew install jqlogs
```

或者直接從 Release manifest 安裝：

```bash
kubectl krew install --manifest=kubectl-jqlogs.yaml
```

### 使用 Go

```bash
go install github.com/shihyuho/kubectl-jqlogs@latest
```

### 手動安裝

從 [Releases](https://github.com/shihyuho/kubectl-jqlogs/releases)頁面下載適合您平台的二進位檔案，並將其放置在您的 `$PATH` 中。

## 使用方法

語法與 `kubectl logs` 相同，最後可加上選用的 jq 查詢。

### 基本用法

檢視日誌並自動格式化 JSON：

```bash
kubectl jqlogs -n my-namespace my-pod
```

### 使用 jq 查詢

要使用 jq 查詢，請用 `--` 分隔並提供查詢。

**依層級過濾：**

```bash
kubectl jqlogs -n my-namespace my-pod -- .level
```

### 進階 JQ 旗標

`kubectl-jqlogs` 支援標準 `gojq` 旗標：

- `-r`, `--raw-output`：輸出原始字串，而非 JSON 文字。
- `-c`, `--compact-output`：緊湊輸出而非美化列印。
- `-C`, `--color-output`：彩色化 JSON 輸出。
- `--yaml-output`：輸出為 YAML。

#### 範例

**YAML 輸出：**
```bash
kubectl jqlogs --yaml-output -n my-ns my-pod
# level: info
# msg: hello
```

**簡易欄位選擇 (Simple Field Selection):**

這是 `kubectl-jqlogs` 新增的 **自定義功能 (custom capability)**，旨在簡化日誌檢查。您無需編寫冗長的 `jq` 結構 (如 `{msg: .msg}` 或 `"\(.msg)"`)，只需列出您想要的欄位 (例如 `.level .msg`)，外掛程式就會自動為您格式化。

> **注意**：這純粹是一項增強功能。**完全支援所有標準 `jq` 語法**，因此您隨時可以使用複雜的過濾器、`select()`、`map()` 和管線。

```bash
kubectl jqlogs -n my-namespace my-pod -- .level .message
# Output: "info An error occurred"

# 對於以 '@' 開頭的欄位，您可以直接使用它們：
kubectl jqlogs -n my-namespace my-pod -- .@timestamp .message
# Output: "2026-01-15... An error occurred"
```

**原始輸出 (可讀的堆疊追蹤)：**

使用 `-r` 輸出不帶引號的原始字串，這可以正確呈現換行符 (`\n`)。
```bash
kubectl jqlogs -r -n my-namespace my-pod -- .message
# Output:
# Error: ...
#   at com.example...
```

**直接選擇訊息：**

```bash
kubectl jqlogs -n my-namespace my-pod -- 'select(.level=="error") | .msg'
```

### 串流日誌

使用 `-f` 追蹤日誌：

```bash
kubectl jqlogs -f -n my-namespace my-pod
```

## Shell 別名 (Alias)

為了節省時間，建議使用 shell 別名。將 `kubectl logs` 替換為更短的指令，如 `klo`：

### Bash / Zsh

```bash
# Add to your .bashrc or .zshrc
alias klo='kubectl jqlogs'
```

現在您可以簡單地執行：

```bash
klo -n my-ns my-pod
```

## 運作原理

`kubectl-jqlogs` 充當 `kubectl logs` 的包裝器。它執行原生指令，擷取輸出串流，並處理每一行：

1. 如果該行是有效的 JSON，它會套用指定的 jq 查詢 (預設為 `.`) 並美化列印結果。
2. 如果該行不是 JSON，則照原樣列印。
3. 如果查詢對某一行失敗，會報告錯誤但串流繼續。

## License

MIT
