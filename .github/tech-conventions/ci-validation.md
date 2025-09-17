# CI & Validation

> Continuous Integration ensures code quality, security, and consistency across all contributions.

<br><br>

## 🧩 Automated Checks

CI automatically runs on every PR to verify:

* Formatting (use `magex format:fix` which runs: `go fmt` and `goimports` and `gofumpt`)
* Linting (`magex lint`)
* Tests (`magex test`)
* Fuzz tests (if applicable) (`magex test:fuzz`)
* This codebase uses GitHub Actions; test workflows reside in `.github/workflows/fortress.yml` and `.github/workflows/fortress-test-suite.yml`.
* Pin each external GitHub Action to a **full commit SHA** (e.g., `actions/checkout@2f3b4a2e0e471e13e2ea2bc2a350e888c9cf9b75`) as recommended by GitHub's [security hardening guidance](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions#using-pinned-actions). Dependabot will track and update these pinned versions automatically.

Failing PRs will be blocked. AI agents should iterate until CI passes.

<br><br>

## 🔄 CI Pipeline Overview

### PR Validation Pipeline
1. **Code Quality Checks**
   - Format validation (gofmt, goimports, gofumpt)
   - Linting with 60+ linters via golangci-lint
   - YAML formatting with Prettier

2. **Testing Suite**
   - Unit tests with coverage reporting
   - Race condition detection
   - Integration tests (if applicable)
   - Fuzz tests for critical paths

3. **Security Scanning**
   - Vulnerability scanning with govulncheck
   - Static analysis with CodeQL
   - OSSAR security checks
   - Dependency vulnerability alerts

4. **Build Validation**
   - Multi-platform builds (Linux, macOS, Windows)
   - Binary artifact generation
   - Docker image builds (if configured)

<br><br>

## ✅ Required Status Checks

All PRs must pass these checks before merge:

| Check Name      | Description                 | Required |
|-----------------|-----------------------------|----------|
| `lint`          | Code formatting and linting | ✅ Yes    |
| `test`          | Unit tests with coverage    | ✅ Yes    |
| `test-race`     | Race condition detection    | ✅ Yes    |
| `security-scan` | Vulnerability scanning      | ✅ Yes    |
| `build`         | Multi-platform builds       | ✅ Yes    |
| `codeql`        | Static security analysis    | ✅ Yes    |

<br><br>

## 🚀 CI Performance

### Optimization Strategies
* **Parallel execution** of independent jobs
* **Caching** for Go modules and build artifacts
* **Matrix builds** for multi-platform testing
* **Conditional runs** to skip unnecessary checks

### Cache Management
```yaml
- uses: actions/cache@v3
  with:
    path: |
      ~/go/pkg/mod
      ~/.cache/go-build
    key: ${{ runner.os }}-go-${{ hashFiles('${{ $GO_SUM_FILE }}') }}
```

<br><br>

## 🔧 Local CI Validation

Run CI checks locally before pushing:

```bash
# Individual checks
magex lint         # Linting only
magex test         # Tests only
magex test:race    # Race detection
magex deps:audit   # Security scan
```

<br><br>

## 📊 Coverage Requirements

### Thresholds
* Minimum coverage: 80% (enforced)
* Target coverage: 90%+ (recommended)
* New code coverage: 95%+ (best practice)

### Coverage Reports
* Generated on every PR
* Visible in PR comments
* Historical tracking via GitHub Pages
* Badge updates automatically

<br><br>

## 🛡️ Security Checks

### Automated Security Scanning
1. **govulncheck** - Go vulnerability database
2. **CodeQL** - Semantic code analysis
3. **OSSAR** - Open Source Security Analysis
4. **Dependabot** - Dependency updates

### Security Workflow
```bash
# Local security check
magex deps:audit

# Components checked:
# - Known CVEs in dependencies
# - Hard-coded secrets
# - SQL injection vulnerabilities
# - Path traversal risks
```

<br><br>

## 🔄 CI Configuration

### Workflow Files
* `.github/workflows/fortress.yml` - Main CI pipeline
* `.github/workflows/fortress-test-suite.yml` - Extended test suite
* `.github/workflows/codeql-analysis.yml` - Security analysis

### Configuration Best Practices
* Use environment variables for configuration
* Pin action versions to commit SHAs
* Set appropriate timeouts
* Use concurrency controls
* Implement proper error handling

<br><br>

## 🚨 Troubleshooting CI Failures

### Common Issues

1. **Linting Failures**
   ```bash
   # Fix formatting
   magex format:fix

   # Check specific linter
   golangci-lint run --enable-only <linter-name>
   ```

2. **Test Failures**
   ```bash
   # Run failing test locally
   go test -v -run TestName ./package

   # Debug with more output
   go test -v -race -count=1 ./...
   ```

3. **Coverage Drops**
   ```bash
   # Generate coverage report
   magex test:cover

   # Find uncovered lines
   go tool cover -html=coverage.out
   ```

4. **Security Vulnerabilities**
   ```bash
   # Check vulnerabilities
   magex deps:audit

   # Update dependencies
   magex deps:update
   ```

<br><br>

## 🔀 CI for Different Branches

### Branch Protection Rules
* **main/master**: All checks required
* **release/***: Additional release checks
* **feature/***: Standard checks
* **dependabot/***: Security checks prioritized

### Merge Requirements
1. All CI checks must pass
2. Code review approval required
3. Branch must be up to date
4. No merge conflicts

<br><br>

## 📈 CI Metrics

Monitor CI health through:
* Build duration trends
* Failure rate by check type
* Flaky test detection
* Coverage trends over time

Access metrics via:
* GitHub Actions insights
* Repository analytics
* Custom dashboards (if configured)
