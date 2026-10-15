# Go Essentials

> Non-negotiable practices that define professional Go development. Every function, package, and design decision should reflect these principles.

<br><br>

## 🌐 Context-First Design

Context should flow through your entire call stack—no exceptions.

* **Always pass `context.Context` as the first parameter** for any operation that could be canceled, timeout, or carry request-scoped values
* **Never store context in structs**—pass it explicitly through function calls
* **Use `context.Background()` only at the top level** (main or service initialization); tests use `t.Context()`
* **Derive child contexts** using `context.WithTimeout()`, `context.WithCancel()`, or `context.WithValue()`
* **Respect context cancellation** by checking `ctx.Done()` in long-running operations

```go
// ✅ Correct: Context as first parameter
func ProcessUserData(ctx context.Context, userID string) error {
    // Check for cancellation before expensive operations
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Pass context down the call chain
    return database.FetchUser(ctx, userID)
}

// 🚫 Incorrect: No context parameter
func ProcessUserData(userID string) error {
    return database.FetchUser(userID) // Can't be canceled or timeout
}
```

<br><br>

## 🔌 Interface Design Philosophy

Interfaces define contracts, not implementations. Keep them minimal and focused.

* **Accept interfaces, return concrete types** — caller decides what they need; you provide specific value
* **Keep interfaces small** — prefer single-method interfaces (`io.Reader`, `io.Writer`, `io.Closer`)
* **Define interfaces where they're used**, not where they're implemented (consumer-driven)
* **Use composition over large interfaces** — combine small interfaces when needed
* **Name single-method interfaces with `-er` suffix** (`Reader`, `Writer`, `Validator`)

```go
// ✅ Small, focused interface defined at point of use
type UserValidator interface {
    ValidateUser(ctx context.Context, user User) error
}

func ProcessSignup(ctx context.Context, validator UserValidator, user User) error {
    if err := validator.ValidateUser(ctx, user); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // Process signup...
    return nil
}

// 🚫 Large, monolithic interface
type UserService interface {
    ValidateUser(ctx context.Context, user User) error
    CreateUser(ctx context.Context, user User) error
    UpdateUser(ctx context.Context, user User) error
    DeleteUser(ctx context.Context, userID string) error
    ListUsers(ctx context.Context) ([]User, error)
    // ... 15 more methods
}
```

<br><br>

## ⚡ Goroutine Discipline

Goroutines are cheap to create but expensive to debug when mismanaged.

* **Always have a clear lifecycle** — know when goroutines start and how they terminate
* **Use context for cancellation** — never leave goroutines hanging
* **Avoid naked `go func()`** — wrap in functions that handle errors and cleanup
* **Use `sync.WaitGroup` or channels** for coordination and synchronization
* **Handle goroutine panics** — use `defer recover()` for background workers
* **Limit concurrency** — use worker pools or semaphores to prevent resource exhaustion

```go
// ✅ Well-managed goroutine with proper lifecycle
func ProcessBatch(ctx context.Context, items []Item) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(items))

    for _, item := range items {
        wg.Add(1)
        go func(item Item) {
            defer wg.Done()
            defer func() {
                if r := recover(); r != nil {
                    errCh <- fmt.Errorf("panic processing item %v: %v", item.ID, r)
                }
            }()

            select {
            case <-ctx.Done():
                errCh <- ctx.Err()
                return
            default:
            }

            if err := processItem(ctx, item); err != nil {
                errCh <- fmt.Errorf("failed to process item %v: %w", item.ID, err)
            }
        }(item)
    }

    wg.Wait()
    close(errCh)

    for err := range errCh {
        if err != nil {
            return err // Return first error encountered
        }
    }
    return nil
}

// 🚫 Unmanaged goroutine
func ProcessBatch(items []Item) {
    for _, item := range items {
        go func(item Item) {
            processItem(item) // No error handling, no cancellation, no lifecycle
        }(item)
    }
    // No way to know when processing is complete
}
```

<br><br>

## 🚫 No Global State

Global variables make code unpredictable, hard to test, and create hidden dependencies.

* **No package-level variables** that hold mutable state
* **Use dependency injection** — pass dependencies explicitly through constructors
* **Prefer `struct` fields over globals** for configuration and state
* **Use `context.Context`** for request-scoped values instead of globals
* **Constants and immutable data are acceptable** at package level

```go
// ✅ Dependency injection pattern
type UserService struct {
    db     Database
    logger Logger
    config Config
}

func NewUserService(db Database, logger Logger, config Config) *UserService {
    return &UserService{
        db:     db,
        logger: logger,
        config: config,
    }
}

func (s *UserService) CreateUser(ctx context.Context, user User) error {
    s.logger.Info("creating user", "userID", user.ID)
    return s.db.Insert(ctx, user)
}

// 🚫 Global state
var (
    globalDB     Database
    globalLogger Logger
    globalConfig Config
)

func CreateUser(ctx context.Context, user User) error {
    globalLogger.Info("creating user", "userID", user.ID) // Hidden dependency
    return globalDB.Insert(ctx, user)                     // Hard to test
}
```

<br><br>

## 🚫 No `init()` Functions

`init()` functions create unpredictable initialization order and hidden side effects.

* **Use explicit constructors** (`NewXxx()` functions) instead of `init()`
* **Make initialization lazy** — initialize on first use when possible
* **Use `sync.Once`** for one-time initialization that must happen exactly once
* **Prefer dependency injection** over package-level initialization
* **Exception**: Only use `init()` for registering with external systems (e.g., database drivers, but even then, prefer explicit registration)

```go
// ✅ Explicit initialization
type Cache struct {
    mu    sync.RWMutex
    data  map[string]any
    once  sync.Once
}

func NewCache() *Cache {
    return &Cache{
        data: make(map[string]any),
    }
}

func (c *Cache) ensureInitialized() {
    c.once.Do(func() {
        // One-time setup that's expensive
        c.data = loadInitialData()
    })
}

// 🚫 Hidden initialization
var globalCache map[string]any

func init() {
    globalCache = make(map[string]any)
    // This runs at import time - unpredictable order
    // Hard to test, hard to control
}
```

<br><br>

## ⚠️ Error Handling Excellence

* **Always check errors**
* **Use `if err != nil { return err }`** for early returns
* **Use `errors.Is()`** or `errors.As()` for error type checks
* **Use `fmt.Errorf` for wrapping errors** with context
* **Prefer `errors.New()`** over `fmt.Errorf`
* **Use custom error types sparingly**
* **Avoid returning ambiguous errors;** provide context
* **Avoid using `panic`** for expected errors; reserve it for unrecoverable situations
* **Use `errors.Unwrap()`** to access underlying errors when needed
* **Use `errors.Join()`** to combine multiple errors when appropriate
* **Wrap errors with context** using `fmt.Errorf("operation failed: %w", err)`
* **Return early on errors** — avoid deep nesting with guard clauses
* **Log errors at the boundary** — don't log the same error multiple times as it bubbles up

```go
// ✅ Comprehensive error handling
func ProcessPayment(ctx context.Context, payment Payment) error {
    if payment.Amount <= 0 {
        return errors.New("payment amount must be positive")
    }

    user, err := userRepo.GetUser(ctx, payment.UserID)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return fmt.Errorf("cannot process payment for unknown user %s", payment.UserID)
        }
        return fmt.Errorf("failed to fetch user %s: %w", payment.UserID, err)
    }

    if err := validatePaymentMethod(ctx, payment.Method); err != nil {
        return fmt.Errorf("invalid payment method: %w", err)
    }

    txn, err := chargePayment(ctx, payment)
    if err != nil {
        return fmt.Errorf("payment charge failed for user %s: %w", user.ID, err)
    }

    if err := auditRepo.LogTransaction(ctx, txn); err != nil {
        // Log but don't fail the payment
        log.Error("failed to audit transaction", "txnID", txn.ID, "error", err)
    }

    return nil
}

// 🚫 Poor error handling
func ProcessPayment(ctx context.Context, payment Payment) error {
    user, _ := userRepo.GetUser(ctx, payment.UserID) // Ignored error

    validatePaymentMethod(ctx, payment.Method) // Ignored return value

    txn, err := chargePayment(ctx, payment)
    if err != nil {
        return err // No context about what failed
    }

    auditRepo.LogTransaction(ctx, txn) // Ignored error
    return nil
}
```

<br><br>

## 📦 Module Hygiene

Go modules are the foundation of dependency management and reproducible builds.

* **Always use Go modules** — never develop outside a module
* **Pin dependencies to specific versions** — avoid `latest` or floating versions in production
* **Use `go mod tidy`** after any dependency changes to clean up unused dependencies
* **Prefer minimal module graphs** — avoid deep dependency trees when possible
* **Use `replace` directives sparingly** — only for development or emergency patches
* **Document major version upgrades** — breaking changes should be called out in PRs

```bash
# ✅ Proper module management workflow
go mod init github.com/company/project
go get github.com/stretchr/testify@v1.8.4  # Pin to specific version
go mod tidy                                  # Clean up go.mod and go.sum
go mod verify                               # Verify dependencies haven't been tampered with

# ✅ Check for vulnerabilities
govulncheck ./...

# ✅ Update dependencies (with care)
go get -u ./...  # Update to latest minor/patch versions
go mod tidy

# ✅ Use MAGE-X for common tasks
magex deps:update # updates dependencies and runs go mod tidy

# 🚫 Avoid using `replace` unless absolutely necessary
# go.mod
replace github.com/some/dependency => github.com/some/dependency v1.2.3

```

<br><br>

## 🔧 Performance & Profiling

Write code that performs well by default, and measure when optimization is needed.

* **Use `magex bench`** to establish performance baselines
* **Profile with `go tool pprof`** when investigating performance issues
* **Avoid premature optimization** — write clear code first, optimize bottlenecks later
* **Use benchmarks** to validate that optimizations actually improve performance
* **Pool expensive objects** with `sync.Pool` when allocation pressure is high
* **Consider memory allocation patterns** — prefer slices to maps for small datasets

```go
// ✅ Performance-conscious code with benchmarks
func BenchmarkUserProcessing(b *testing.B) {
    users := generateTestUsers(1000)

    for b.Loop() { // times only the loop, not the setup above
        processUsers(users)
    }
}

func processUsers(users []User) []ProcessedUser {
    // Pre-allocate slice to avoid repeated allocations
    result := make([]ProcessedUser, 0, len(users))

    for _, user := range users {
        processed := ProcessedUser{
            ID:   user.ID,
            Name: strings.ToUpper(user.Name), // Simple transformation
        }
        result = append(result, processed)
    }

    return result
}
```

<br><br>

## 🛠 Formatting & Linting

Code must be cleanly formatted and pass all linters before being committed.

```bash
magex format:fix
magex lint
```

> Refer to `.golangci.json` for the full set of enabled linters and formatters.

Editors should honor `.editorconfig` for indentation and whitespace rules, and
Git respects `.gitattributes` to enforce consistent line endings across
platforms.

<br><br>

## 💄 YAML Formatting

YAML files must be formatted consistently to ensure clean diffs and readable configuration files.

```bash
magex format:fix
```

> The `magex format:fix` command handles YAML formatting (via yamlfmt) along with Go, JSON, and other file types. CI automatically validates formatting using the same tools.
