# Dependency Management

> Dependency hygiene is critical for security, reproducibility, and developer experience. Follow these practices to ensure our module stays stable, up to date, and secure.

<br><br>

## 📦 Module Management

* All dependencies must be managed via **Go Modules** (`go.mod`, `go.sum`)

* After adding, updating, or removing imports, run:

  ```bash
  go mod tidy
  ```

* Periodically refresh dependencies with:

  ```bash
  go get -u ./...
  ```

> Avoid unnecessary upgrades near release windows—review major version bumps carefully for breaking changes.

<br><br>

## 🛡️ Security Scanning

* Use [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) to identify known vulnerabilities:
```bash
  govulncheck ./...
```

* Run via magex command:
```bash
  magex deps:audit
```

* Run [gitleaks](https://github.com/gitleaks/gitleaks) before committing code to detect hardcoded secrets or sensitive data in the repository:
```bash
brew install gitleaks
gitleaks detect --source . --log-opts="--all" --verbose
```

* Address critical advisories before merging changes into `main/master`

* Document any intentionally ignored vulnerabilities with clear justification and issue tracking

* We follow the [OpenSSF](https://openssf.org) best practices to ensure this repository remains compliant with industry‑standard open source security guidelines

<br><br>

## 📁 Version Control

* Never manually edit `go.sum`
* Do not vendor dependencies; we rely on modules for reproducibility
* Lockstep upgrades across repos (when applicable) should be coordinated and noted in PRs

> Changes to dependencies should be explained in the PR description and ideally linked to the reason (e.g., bug fix, security advisory, feature requirement).

<br><br>

## 🔄 Dependency Update Workflow

### Regular Updates
1. **Check for updates**
   ```bash
   go list -u -m all
   ```

2. **Update minor/patch versions**
   ```bash
   go get -u ./...
   go mod tidy
   ```

3. **Test thoroughly**
   ```bash
   magex test
   magex test:race
   magex bench
   ```

4. **Security scan**
   ```bash
   magex deps:audit
   ```

### Major Version Updates
1. **Review breaking changes** in release notes
2. **Update import paths** if required
3. **Fix compilation errors**
4. **Update tests** for new behavior
5. **Document in PR** what changed and why

<br><br>

## 🤖 Automated Dependency Management

### Dependabot Configuration
* Configured in `.github/dependabot.yml`
* Checks for updates weekly
* Groups minor/patch updates
* Creates separate PRs for major versions

### Auto-merge Rules
* Minor/patch updates with passing CI can auto-merge
* Major updates require manual review
* Security updates prioritized for review

<br><br>

## 📊 Dependency Analysis

### Check dependency graph
```bash
go mod graph
```

### Identify unused dependencies
```bash
go mod tidy -v
```

### Analyze module size impact
```bash
go mod download -json | jq '.Dir' | xargs du -sh | sort -h
```

<br><br>

## 🚫 Dependency Guidelines

### DO:
* **Pin to specific versions** in production
* **Review licenses** before adding dependencies
* **Prefer standard library** when possible
* **Use minimal dependencies** for core functionality
* **Document unusual dependencies** in code comments

### DON'T:
* **Use `latest` tags** in production
* **Import unused packages**
* **Use replace directives** except for emergencies
* **Add dependencies for trivial functionality**
* **Ignore security advisories**

<br><br>

## 🔍 Evaluating New Dependencies

Before adding a new dependency, consider:

1. **Necessity**: Can we implement this ourselves simply?
2. **Maintenance**: Is the project actively maintained?
3. **Security**: Any known vulnerabilities?
4. **License**: Compatible with our project?
5. **Size**: How much does it increase binary size?
6. **Quality**: Well-tested? Good documentation?
7. **Dependencies**: Does it bring many transitive dependencies?

<br><br>

## 📝 Replace Directives

Use `replace` only when absolutely necessary:

```go
// Temporary fix for critical bug until upstream releases
replace github.com/broken/package v1.2.3 => github.com/fork/package v1.2.4-fixed

// Local development only - remove before committing
replace github.com/company/module => ../local-module
```

Document why the replacement is needed and track removal in an issue.

<br><br>

## 🔐 Private Dependencies

For private modules:

1. **Configure authentication**
   ```bash
   git config --global url."git@github.com:company/".insteadOf "https://github.com/company/"
   ```

2. **Set GOPRIVATE**
   ```bash
   export GOPRIVATE=github.com/company/*
   ```

3. **Document setup** in README for team members

<br><br>

## 📈 Monitoring Dependencies

### Track outdated dependencies
```bash
# Show available updates
go list -u -m all | grep '\['

# Count total dependencies
go mod graph | wc -l
```

### Review dependency changes
```bash
# See what changed in go.mod
git diff go.mod

# Detailed view of go.sum changes
git diff go.sum | grep '^[+-]' | sort
```

> Regular dependency maintenance prevents security issues and reduces upgrade complexity.
