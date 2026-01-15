# Conventional Commits

## 格式

```
<type>[可選的 scope]: <description>

[可選的 body]

[可選的 footer]
```

## Type 前綴

| Type | 說明 |
| --- | --- |
| `feat` | 新增功能 |
| `fix` | 修正 Bug |
| `docs` | 只修改文件 |
| `style` | 程式碼風格（排版、分號等） |
| `refactor` | 重構（沒有新增功能或修正 Bug） |
| `perf` | 改善效能 |
| `test` | 新增或更新測試 |
| `build` | 修改建置系統或依賴項目 |
| `ci` | 修改 CI 設定 |
| `chore` | 其他變更（工具、設定等） |
| `revert` | 還原先前的 Commit |

### Type 選擇的判斷標準

Conventional Commits 是以搭配 Semantic Versioning 為前提的規範。

因此，Type 的選擇標準是 **「從使用者的角度來看，是否有變化」**。

**使用者的定義：**

- 如果是應用程式：使用該 App 的使用者
- 如果是函式庫：使用該函式庫的開發者

**判斷標準：**

- 從使用者角度看有變化 → `feat`, `fix`, `perf`
- 從使用者角度看沒變化 → `docs`, `style`, `refactor`, `test`, `build`, `ci`, `chore`

**常見錯誤：**

1.  因為「修正了」所以用 `fix`
    - 範例：`fix: 修正測試設定`
    - 錯在哪裡：修改測試設定並不會修好使用者的問題，所以不該用 `fix`。
    - 應該怎麼想：這個 Commit 的主旨是測試，所以應該用 `test`。

2.  因為「改了 Style」所以用 `style`
    - 範例：`style: 修改按鈕樣式`
    - 錯在哪裡：`style` 是指程式碼風格（排版、分號等），而不是指設計或 CSS 樣式表，所以不該用 `style`。
    - 應該怎麼想：這對使用者來說是個變化，所以應該用 `feat` 或 `fix`。

3.  因為是「小變更」所以用 `chore`
    - 範例：`chore: 修改文字`
    - 錯在哪裡：這對使用者來說是個變化，所以不該用 `chore`。
    - 應該怎麼想：這對使用者來說是個變化，所以應該用 `feat` 或 `fix`。

## Footer

加在 Body 後面，並以一個空行分隔。格式為 `token: value` 或 `token #value`：

```
feat(auth): 新增登入功能

實作 OAuth2 認證

Reviewed-by: Alice
Refs #123
```

## 範例

```
feat(auth): 新增 OAuth2 登入功能
```

```
fix: 解決 user service 的 null pointer 問題
```

```
docs: 在 README 中追加設定步驟
```

```
refactor(api): 簡化 response 處理
```

## 使用 SemVer 的情況

### Breaking Change

在 Type/Scope 後面加上 `!`，或是在 Footer 加上 `BREAKING CHANGE:`：

```
feat!: 移除已棄用的 API endpoint
feat(api)!: 變更 response 格式

BREAKING CHANGE: API response 從物件變為陣列
```

### 與 SemVer 的關聯

| Type | SemVer |
| --- | --- |
| `fix` | PATCH（例: 1.0.0 → 1.0.1） |
| `feat` | MINOR（例: 1.0.0 → 1.1.0） |
| `BREAKING CHANGE` | MAJOR（例: 1.0.0 → 2.0.0） |