# Testing Standards

> Comprehensive testing ensures code quality, reliability, and maintainability. These standards apply to all test code in the repository.

<br><br>

## 🧪 Testing Framework

We use the `testify` suite for unit tests. All tests must follow these conventions:

* Name tests using the pattern: `TestFunctionNameScenarioDescription` (no underscores) (PascalCase)
* Use `testify` when possible, do not use `testing` directly
* Use `testify/assert` for general assertions
* Use `testify/require` for:
	* All error or nil checks
	* Any test where failure should halt execution
	* Any test where a pointer or complex structure is required to be used after the check
* Use `require.InDelta` or `require.InEpsilon` for floating-point comparisons
* Prefer **table-driven tests** for clarity and reusability, always have a name for each test case
* Use subtests (`t.Run`) to isolate and describe scenarios
* If the test is in a test suite, always use the test suite instead of `t` directly
* **Optionally use** `t.Parallel()` , but try and avoid it unless testing for concurrency issues
* Avoid flaky, timing-sensitive, or non-deterministic tests
* Mock external dependencies — tests should be fast and deterministic
* Use descriptive test names that explain the scenario being tested
* Test error cases — ensure your error handling actually works
* Use the test's own helpers: `t.Context()` instead of `context.Background()`, `t.Setenv()` instead of `os.Setenv()`, `t.Chdir()` instead of `os.Chdir()`, and `t.TempDir()` instead of `os.MkdirTemp()`. `t.Setenv()` and `t.Chdir()` restore process state when the test ends; neither can be used in parallel tests or tests with parallel ancestors
* Handle all errors in tests properly:
	* `os.Unsetenv()` returns an error - use `require.NoError(t, err)`
	* `db.Close()` in defer statements - wrap in anonymous function: `defer func() { _ = db.Close() }()`

<br><br>

## 📝 Test Structure Example

```go
func TestUserServiceCreateUser(t *testing.T) {
    tests := []struct {
        name        string
        input       User
        setupMocks  func(*MockDB)
        wantErr     bool
        errContains string
    }{
        {
            name: "valid user creation",
            input: User{
                Name:  "Alice",
                Email: "alice@example.com",
            },
            setupMocks: func(db *MockDB) {
                db.On("Insert", mock.Anything, mock.Anything).Return(nil)
            },
            wantErr: false,
        },
        {
            name: "duplicate email error",
            input: User{
                Name:  "Bob",
                Email: "existing@example.com",
            },
            setupMocks: func(db *MockDB) {
                db.On("Insert", mock.Anything, mock.Anything).
                    Return(ErrDuplicateEmail)
            },
            wantErr:     true,
            errContains: "duplicate email",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            mockDB := new(MockDB)
            tt.setupMocks(mockDB)

            service := NewUserService(mockDB)

            // Execute
            err := service.CreateUser(t.Context(), tt.input)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
                if tt.errContains != "" {
                    assert.Contains(t, err.Error(), tt.errContains)
                }
            } else {
                require.NoError(t, err)
            }

            mockDB.AssertExpectations(t)
        })
    }
}
```

<br><br>

## 🔧 Running Tests

Run tests locally with:
```bash
magex test
```

> All tests must pass in CI prior to merge.

<br><br>

## 🔍 Fuzz Tests (Optional)

Fuzz tests help uncover unexpected edge cases by generating random inputs. While not required, they are encouraged for **small, self-contained functions**.

Best practices:
* Keep fuzz targets short and deterministic
* Seed the corpus with meaningful values
* Run fuzzers with `go test -fuzz=. -run=^$` when exploring edge cases
* Limit iterations for local runs to maintain speed

Example:
```go
func FuzzParseConfig(f *testing.F) {
    // Seed corpus with known inputs
    f.Add("valid: true")
    f.Add("count: 42")
    f.Add("")

    f.Fuzz(func(t *testing.T, input string) {
        // Function should not panic on any input
        cfg, err := ParseConfig(input)
        if err != nil {
            // Error is acceptable, panic is not
            return
        }

        // Validate parsed config is sensible
        require.NotNil(t, cfg)
    })
}
```

<br><br>

## 📈 Code Coverage

* Code coverage thresholds are configured in the GoFortress coverage system
* Aim to provide meaningful test coverage for all new logic and edge cases
* Cover every public function with at least one test
* Aim for >= 90% coverage across the codebase (ideally 100%)
* Use `go test -coverprofile=coverage.out ./...` to generate coverage reports

Generate coverage locally:
```bash
magex test:cover
```

View coverage in browser:
```bash
go tool cover -html=coverage.out
```

<br><br>

## 🎯 What to Test

### Unit Tests Should Cover:
* **Happy path** — normal, expected behavior
* **Edge cases** — boundary conditions, empty inputs, maximum values
* **Error paths** — all error returns should be tested
* **Concurrency** — race conditions, goroutine safety (when applicable)
* **State changes** — ensure mutations happen correctly
* **Resource cleanup** — verify resources are properly released

### What NOT to Test:
* **Third-party libraries** — assume they work correctly
* **Language features** — don't test Go itself
* **Private functions** — test through public API only
* **Generated code** — unless you wrote the generator

<br><br>

## 🏃 Test Performance

Keep tests fast and focused:
* **Use mocks** instead of real databases, APIs, or file systems
* **Parallelize** independent tests with `t.Parallel()` when beneficial
* **Skip slow tests** in short mode: `if testing.Short() { t.Skip("skipping in short mode") }`
* **Use test fixtures** sparingly — prefer generating test data in code
* **Avoid `time.Sleep`** — use channels, contexts, or synchronization primitives

<br><br>

## 🔨 Test Helpers

Create focused helper functions to reduce duplication:

```go
// Good: Focused helper with clear purpose
func requireUserEqual(t *testing.T, expected, actual User) {
    t.Helper()
    require.Equal(t, expected.ID, actual.ID)
    require.Equal(t, expected.Name, actual.Name)
    require.Equal(t, expected.Email, actual.Email)
    require.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
}

// Bad: Overly generic helper
func assertStuff(t *testing.T, a, b any) {
    // Too vague, hard to understand intent
}
```

<br><br>

## 🚫 Common Testing Mistakes

Avoid these patterns:

```go
// 🚫 Don't ignore errors in tests
result, _ := SomeFunction() // Bad!

// ✅ Always check errors
result, err := SomeFunction()
require.NoError(t, err)

// 🚫 Don't use global test state
var testDB *sql.DB // Bad! Makes tests interdependent

// ✅ Create fresh instances for each test
func setupTestDB(t *testing.T) *sql.DB {
    // Return new instance
}

// 🚫 Don't test implementation details
assert.Equal(t, 3, len(cache.internalMap)) // Bad!

// ✅ Test behavior through public API
assert.Equal(t, 3, cache.Size())
```

<br><br>

## 📊 Test Documentation

Document complex test scenarios:

```go
func TestComplexWorkflow(t *testing.T) {
    // This test verifies that the payment processor correctly handles
    // partial refunds when the original transaction was split across
    // multiple payment methods and one of them has expired.
    //
    // Setup: Create order with 2 payment methods
    // Action: Process partial refund
    // Expectation: Refund succeeds, expired method is skipped

    // ... test implementation
}
```
