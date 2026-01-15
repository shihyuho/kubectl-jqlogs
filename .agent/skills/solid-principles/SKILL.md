---
name: solid-principles
description: Use during implementation when designing modules, functions, and components requiring SOLID principles for maintainable, flexible architecture.
allowed-tools:
  - Read
  - Edit
  - Grep
  - Glob
---


# SOLID 原則

應用 SOLID 設計原則，打造可維護、有彈性的程式碼架構。

## 五大原則

### 1. 單一職責原則 (Single Responsibility Principle, SRP)

### 一個模組應該只有一個變更的理由

### Java 模式

```java
// BAD - 多重職責
public class UserService {
    public void createUser(String username) {
        // 1. 驗證使用者名稱
        if (username == null || username.isEmpty()) {
            // ...
        }

        // 2. 將使用者存入資料庫
        // DataSource ds = ...;
        // ds.saveUser(username);

        // 3. 寄送歡迎郵件
        // EmailClient email = new EmailClient();
        // email.send("welcome@corp.com", "Welcome!");
    }
}

// GOOD - 單一職責
public class UserValidator {
    public void validate(String username) {
        // 驗證邏輯
    }
}

public class UserRepository {
    public void save(String username) {
        // 資料庫儲存邏輯
    }
}

public class UserNotifier {
    public void sendWelcomeEmail(String username) {
        // 寄送郵件邏輯
    }
}
```

**問問自己：**「這個模組做的『那一件』事情是什麼？」

### 2. 開放封閉原則 (Open/Closed Principle, OCP)

**軟體實體應對擴充開放，對修改封閉。**

### Java 模式

```java
// BAD - 新增支付方式時需要修改
public class PaymentProcessor {
    public void process(PaymentInfo info) {
        if ("creditcard".equals(info.getType())) {
            // 處理信用卡
        } else if ("paypal".equals(info.getType())) {
            // 處理 PayPal
        }
        // 新增支付方式時必須修改此函式
    }
}

// GOOD - 透過 interface 進行擴充
public interface PaymentProvider {
    void processPayment(PaymentInfo info);
}

public class CreditCardProvider implements PaymentProvider {
    public void processPayment(PaymentInfo info) {
        // Stripe 邏輯
    }
}

public class PayPalProvider implements PaymentProvider {
    public void processPayment(PaymentInfo info) {
        // PayPal 邏輯
    }
}

// 使用 - 新增 provider 時不需修改此處程式碼
public class PaymentProcessor {
    public void process(PaymentProvider provider, PaymentInfo info) {
        provider.processPayment(info);
    }
}
```

**問問自己：**「我能否在不修改現有程式碼的情況下增加新功能？」

### 3. 里氏替換原則 (Liskov Substitution Principle, LSP)

### 子型別必須可以替換其基礎型別

### Java 模式

```java
// BAD - 違反 LSP
public class Rectangle {
    protected int width, height;

    public void setWidth(int width) { this.width = width; }
    public void setHeight(int height) { this.height = height; }
    public int getArea() { return width * height; }
}

public class Square extends Rectangle {
    // 設定寬度時，高度也跟著改變，破壞了 Rectangle 的契約
    @Override
    public void setWidth(int width) {
        this.width = width;
        this.height = width;
    }

    @Override
    public void setHeight(int height) {
        this.width = height;
        this.height = height;
    }
}

// GOOD - 正確的抽象
public interface Shape {
    int getArea();
}

public class Rectangle implements Shape {
    private final int width, height;
    public Rectangle(int width, int height) {
        this.width = width;
        this.height = height;
    }
    public int getArea() { return width * height; }
}

public class Square implements Shape {
    private final int side;
    public Square(int side) { this.side = side; }
    public int getArea() { return side * side; }
}
```

**問問自己：**「我能否用它的父類別/interface 替換它，而行為不會被破壞？」

### 4. 介面隔離原則 (Interface Segregation Principle, ISP)

**用戶端不應被迫依賴於它們不使用的介面。**

### Java 模式

```java
// BAD - 臃腫的 interface
public interface Worker {
    void work();
    void eat();
    void takeBreak();
    // Robot被迫要實作 eat() 和 takeBreak()
}

public class Human implements Worker {
    public void work() { /* ... */ }
    public void eat() { /* ... */ }
    public void takeBreak() { /* ... */ }
}

public class Robot implements Worker {
    public void work() { /* ... */ }
    public void eat() { 
        // 無法實作，違反原則
        throw new UnsupportedOperationException();
    }
    public void takeBreak() {
        // 無法實作，違反原則
        throw new UnsupportedOperationException();
    }
}


// GOOD - 分離的 interfaces
public interface Workable {
    void work();
}

public interface Feedable {
    void eat();
}

public class Human implements Workable, Feedable {
    public void work() { /* ... */ }
    public void eat() { /* ... */ }
}

public class Robot implements Workable {
    public void work() { /* ... */ }
}
```

**問問自己：**「這個 interface 是否強迫實作者去定義未被使用的方法？」

### 5. 依賴反轉原則 (Dependency Inversion Principle, DIP)

### 依賴抽象，而非具體實作

### Java 模式

```java
// BAD - 直接依賴具體實作
public class UserManager {
    private final StripeApi stripeApi;

    public UserManager() {
        this.stripeApi = new StripeApi(); // 高度耦合
    }

    public void processPayment(double amount) {
        stripeApi.charge(amount);
    }
}

// GOOD - 依賴抽象 (Dependency Injection)
public interface PaymentGateway {
    void charge(double amount);
}

public class StripeGateway implements PaymentGateway {
    public void charge(double amount) {
        // Stripe 邏輯
    }
}

public class UserManager {
    private final PaymentGateway gateway;

    // 依賴被注入，而非自行建立
    public UserManager(PaymentGateway gateway) {
        this.gateway = gateway;
    }

    public void processPayment(double amount) {
        gateway.charge(amount);
    }
}

// 使用
// PaymentGateway stripe = new StripeGateway();
// UserManager manager = new UserManager(stripe);
```

**問問自己：**「我能否在不修改相依程式碼的情況下，替換掉實作？」

## 應用檢查清單

### 撰寫新程式碼之前

- [ ] 識別單一職責
- [ ] 設計擴充點 (behaviours, interfaces)
- [ ] 在實作之前先定義抽象
- [ ] 保持 interface 最小且專注

### 實作期間

- [ ] 每個模組只有一個變更的理由 (SRP)
- [ ] 新功能是透過擴充，而非修改 (OCP)
- [ ] 實作遵守契約 (LSP)
- [ ] Interface 保持最小化 (ISP)
- [ ] 依賴是可被注入/可設定的 (DIP)

### Code Review 期間

- [ ] 職責是否被清楚地分離？
- [ ] 我們能否在不修改現有程式碼的情況下新增功能？
- [ ] 所有實作是否都履行了它們的契約？
- [ ] Interface 是否專注且最小化？
- [ ] 依賴是否被抽象化了？

## 程式碼庫中的常見違規

### SRP 違規

- GraphQL resolvers 同時也包含商業邏輯 (應使用 command handlers)
- Components 同時抓取資料並進行渲染 (應使用 hooks + 展示型 components)

### OCP 違規

- 針對型別使用過長的 if/else 或 case 陳述式 (應使用 behaviours/多型)
- 寫死的 provider 邏輯 (應使用依賴注入)

### LSP 違規

- 在基礎型別會回傳 nil/error tuple 的情況下，實作卻 raise exception
- 不同實作之間的回傳型別不一致

### ISP 違規

- 臃腫的 GraphQL types 要求所有欄位 (應使用 fragments)
- 巨石般的 component props (應分割成專注的 interfaces)

### DIP 違規

- 直接呼叫外部服務 (應使用 behaviours 包裝)
- 寫死的 Repo 呼叫 (應注入 repository)

## 與現有 Skills 的整合

可搭配使用

- `boy-scout-rule`: 在改善程式碼時應用 SOLID
- `test-driven-development`: 為每個職責撰寫測試

## 記住

**SOLID 的核心在於管理依賴與職責，而非創造更多程式碼。**

好的設計源於務實地應用這些原則，而非教條式地遵守。
