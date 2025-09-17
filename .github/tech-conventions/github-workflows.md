# GitHub Workflows Development

> GitHub Actions workflows are critical infrastructure that automate our CI/CD pipeline, testing, security scanning, and release processes. All workflows must be reliable, secure, and maintainable to ensure consistent developer experience and protect our codebase.

<br><br>

## 🎯 Workflow Development Priorities

* **Security First** — Use minimal permissions, pin all actions to full commit SHA, and never expose secrets unnecessarily
* **Reliability** — Workflows should be deterministic, idempotent, and handle failure gracefully
* **Performance** — Optimize for speed with appropriate caching, parallel execution, and efficient resource usage
* **Maintainability** — Write self-documenting workflows with clear structure, consistent naming, and comprehensive comments
* **Native GitHub Features** — Prefer built-in GitHub Actions and features over third-party alternatives when possible
* **Least Privilege** — Grant only the minimum permissions required for each job and step
* **Fail Fast** — Design workflows to catch issues early and provide clear, actionable error messages
* **Documentation** — Every workflow must include purpose, triggers, and maintainer information in header comments

<br><br>

## 📋 Workflow Development Guidelines

* **Action Pinning** — Pin all external actions to full commit SHA (e.g., `actions/checkout@2f3b4a2e0e471e13e2ea2bc2a350e888c9cf9b75`) for security and reproducibility
* **Permissions** — Use `permissions: read-all` as default, then grant specific write permissions only where needed
* **Concurrency Control** — Include concurrency groups to prevent resource conflicts and optimize runner usage
* **Environment Variables** — Use `env` blocks at appropriate levels (workflow, job, or step) to maintain clarity
* **Error Handling** — Use `continue-on-error` and `if` conditionals strategically to handle expected failures
* **Matrix Builds** — Leverage matrix strategies for testing across multiple versions, platforms, or configurations
* **Caching** — Implement intelligent caching for dependencies, build artifacts, and test results
* **Secrets Management** — Use GitHub Secrets for sensitive data; never hardcode credentials or tokens
* **Branch Protection** — Ensure critical workflows are required status checks before merging
* **Workflow Naming** — Use kebab-case names that clearly describe the workflow purpose

<br><br>

## 🏗️ Workflow Template

Use this template as the foundation for all new GitHub Actions workflows:

```yaml
# ------------------------------------------------------------------------------
#  [Workflow Name] Workflow
#
#  Purpose: [Brief description of what this workflow does and why it exists]
#
#  Triggers: [List the events that trigger this workflow]
#
#  Maintainer: @[github-username]
# ------------------------------------------------------------------------------

name: [workflow-name]

# --------------------------------------------------------------------
# Trigger Configuration
# --------------------------------------------------------------------
on:
  [trigger-events]

# --------------------------------------------------------------------
# Permissions
# --------------------------------------------------------------------
permissions: read-all

# --------------------------------------------------------------------
# Concurrency Control
# --------------------------------------------------------------------
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  [job-name]:
    runs-on: ubuntu-latest
    permissions:
      [specific-permissions]: [read|write]
    steps:
      # --------------------------------------------------------------------
      # [Step Category/Purpose]
      # --------------------------------------------------------------------
      - name: [Step Name]
        uses: [action@full-commit-sha]
        with:
          [parameters]

      # --------------------------------------------------------------------
      # [Next Step Category/Purpose]
      # --------------------------------------------------------------------
      - name: [Next Step Name]
        run: |
          [commands]
        env:
          [environment-variables]
```

<br><br>

## 🔍 Template Usage Notes

* Replace all bracketed placeholders `[...]` with actual values
* Use descriptive names for workflows, jobs, and steps
* Include section separators (dashes) to organize logical groups of steps
* Add environment variables in the appropriate scope (workflow, job, or step level)
* Document complex logic with inline comments using `#` within run blocks
* Test workflows thoroughly in draft PRs before merging to avoid CI/CD disruption

<br><br>

## 🚦 Workflow Validation Checklist

Before merging any workflow changes, verify:

* [ ] All external actions are pinned to full commit SHA
* [ ] Permissions follow least-privilege principle
* [ ] Concurrency control prevents resource conflicts
* [ ] Workflow header includes purpose, triggers, and maintainer
* [ ] Section separators organize steps logically
* [ ] Environment variables are properly scoped
* [ ] Error handling covers expected failure scenarios
* [ ] Workflow has been tested in a draft PR
* [ ] Documentation reflects any new workflow dependencies or requirements

<br><br>

## 💡 Best Practices Examples

### Secure Action Pinning
```yaml
# ✅ Good: Pinned to full commit SHA
- uses: actions/checkout@93ea575cb5d8a053eaa0ac8fa3b40d7e05a33cc8 # v3.1.0

# 🚫 Bad: Using tag (mutable)
- uses: actions/checkout@v3

# 🚫 Bad: Using branch (very mutable)
- uses: actions/checkout@main
```

### Proper Permissions
```yaml
# ✅ Good: Minimal permissions
permissions:
  contents: read
  pull-requests: write

# 🚫 Bad: Overly broad permissions
permissions: write-all
```

### Efficient Caching
```yaml
# ✅ Good: Restore and save cache with fallbacks
- name: Cache Go modules
  uses: actions/cache@88522ab9f39a2ea568f7027eddc7d8d8bc9d59c8 # v3.3.1
  with:
    path: |
      ~/go/pkg/mod
      ~/.cache/go-build
    key: ${{ runner.os }}-go-${{ hashFiles('${{ $GO_SUM_FILE }}') }}
```

<br><br>

## 🔧 Common Patterns

### Multi-Platform Testing
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
    go-version: ['1.21', '1.22']
runs-on: ${{ matrix.os }}
steps:
  - uses: actions/setup-go@93397bea11091df50f3d7e59dc26a7711a8bcfbe # v4.1.0
    with:
      go-version: ${{ matrix.go-version }}
```

### Conditional Execution
```yaml
- name: Deploy to staging
  if: github.ref == 'refs/heads/develop' && github.event_name == 'push'
  run: ./deploy.sh staging
```

### Job Dependencies
```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    # ... test steps

  deploy:
    needs: test
    if: success()
    runs-on: ubuntu-latest
    # ... deploy steps
```

<br><br>

## 🚨 Security Considerations

### Handling Secrets
```yaml
# ✅ Good: Use GitHub secrets
env:
  API_KEY: ${{ secrets.API_KEY }}

# 🚫 Bad: Hardcoded secret
#env:
#  API_KEY: "1234"
```

### Third-Party Actions
* Review source code before using
* Pin to commit SHA after review
* Prefer official actions (actions/*)
* Consider forking critical third-party actions

### Pull Request Safety
```yaml
# Limit permissions for PR workflows
on:
  pull_request_target:
    types: [opened, synchronize]

permissions:
  contents: read
  pull-requests: write

# Validate PR source
- name: Check PR safety
  if: github.event.pull_request.head.repo.fork == true
  run: echo "::warning::PR from fork - limited permissions applied"
```

<br><br>

## 📊 Monitoring & Debugging

### Workflow Insights
* Check Actions tab for run history
* Review timing to identify bottlenecks
* Monitor failure patterns
* Track runner usage and costs

### Debugging Techniques
```yaml
# Enable debug logging
- name: Enable debug mode
  run: echo "::debug::Debug mode enabled"

# Add step debugging
- name: Debug environment
  run: |
    echo "Event: ${{ github.event_name }}"
    echo "Ref: ${{ github.ref }}"
    echo "SHA: ${{ github.sha }}"
    env | sort

# Use tmate for interactive debugging (carefully!)
- name: Setup tmate session
  if: ${{ failure() && github.event_name == 'workflow_dispatch' }}
  uses: mxschmitt/action-tmate@v3
  timeout-minutes: 15
```

<br><br>

## 📚 Resources

### Internal
* [CI Validation](ci-validation.md) - CI/CD pipeline overview
* [Security Practices](security-practices.md) - Security guidelines

### External
* [GitHub Actions Documentation](https://docs.github.com/en/actions)
* [Security Hardening Guide](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)
* [Workflow Syntax Reference](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions)
