# 🪝 go-pre-commit: Pure Go Pre-commit Hooks

## Overview

**go-pre-commit** is a production-ready, high-performance Go-native pre-commit framework that delivers **17x faster execution** than traditional Python-based solutions. It provides automated code quality and consistency checks that run before every commit, ensuring your codebase maintains high standards without slowing down development.

🔗 **Repository**: [github.com/mrz1836/go-pre-commit](https://github.com/mrz1836/go-pre-commit)

### Why go-pre-commit?

**Performance Benefits:**
- ⚡ **17x faster execution** - Complete pre-commit checks in <2 seconds
- 🔄 **Parallel processing** - All checks run simultaneously for maximum speed
- 📦 **Zero Python dependencies** - Pure Go binary, no runtime requirements

**Developer Experience:**
- 🛠️ **Make/Mage integration** - Seamlessly wraps existing Makefile/mage targets
- ⚙️ **Environment-driven configuration** - All settings via environment variables
- 🎯 **Production ready** - Comprehensive test coverage and validation
- 🔧 **Cross-platform** - Works identically on macOS, Linux, and Windows

**Team Benefits:**
- 🚀 **Easy onboarding** - New developers productive in minutes
- 📋 **Consistent standards** - Same checks across all environments
- 🔒 **Quality gates** - Prevent common issues before they reach CI/CD

## Quick Start

### Installation

```bash
# Install the go-pre-commit tool
go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest

# Install hooks in your repository
cd your-go-project
go-pre-commit install

# Verify installation
go-pre-commit --version
```

### First Usage

```bash
# Make some changes to your code
echo "package main" > example.go

# Commit normally - hooks run automatically
git add .
git commit -m "feat: add example file"
# ✅ All checks passed in <2s
```

That's it! The pre-commit hooks are now active and will run on every commit.

## Configuration

go-pre-commit uses environment variables for all configuration, typically stored in `.github/.env.base` for team-wide settings and optionally `.github/.env.custom` for local overrides.

### Basic Configuration

```bash
# Enable the system
ENABLE_GO_PRE_COMMIT=true

# Tool version (recommended to pin for consistency)
GO_PRE_COMMIT_VERSION=v1.1.11

# Individual check control
GO_PRE_COMMIT_ENABLE_FUMPT=true      # Code formatting
GO_PRE_COMMIT_ENABLE_LINT=true       # Linting
GO_PRE_COMMIT_ENABLE_MOD_TIDY=true   # Module tidying
GO_PRE_COMMIT_ENABLE_WHITESPACE=true # Trailing whitespace removal
GO_PRE_COMMIT_ENABLE_EOF=true        # End-of-file newlines
```

### Performance Tuning

```bash
# System performance settings
GO_PRE_COMMIT_PARALLEL_WORKERS=2     # Number of parallel workers
GO_PRE_COMMIT_TIMEOUT_SECONDS=120    # Total timeout for all checks
GO_PRE_COMMIT_FAIL_FAST=false        # Stop on first failure
GO_PRE_COMMIT_LOG_LEVEL=debug        # Logging verbosity

# File handling
GO_PRE_COMMIT_MAX_FILE_SIZE_MB=10    # Skip files larger than this
GO_PRE_COMMIT_MAX_FILES_OPEN=100     # Maximum concurrent file handles

# CI/CD optimization
GO_PRE_COMMIT_ALL_FILES=true         # Check all files (vs changed only)
```

### Tool Versions

Pin specific tool versions for consistency across environments:

```bash
# External tool versions
GO_PRE_COMMIT_FUMPT_VERSION=latest
GO_PRE_COMMIT_GOIMPORTS_VERSION=latest
```

## Available Checks

go-pre-commit provides 5 core checks that cover the most important code quality aspects:

| Check          | Description                     | Integration    | Speed              |
|----------------|---------------------------------|----------------|--------------------|
| **fumpt**      | Advanced Go formatting          | `magex format` | 6ms (37% faster)   |
| **lint**       | Comprehensive linting           | `magex lint`   | 68ms (94% faster)  |
| **mod-tidy**   | Module dependency cleanup       | `magex tidy`   | 110ms (53% faster) |
| **whitespace** | Trailing whitespace removal     | Built-in       | <1ms               |
| **eof**        | End-of-file newline enforcement | Built-in       | <1ms               |

### Check Details

**Code Formatting (fumpt):**
- Uses [gofumpt](https://github.com/mvdan/gofumpt) for stricter formatting than `gofmt`
- Integrates with your existing `magex format` target
- Automatically fixes formatting issues when possible

**Linting (lint):**
- Runs comprehensive linting via `magex lint`
- Supports 60+ linters through golangci-lint integration
- Fails fast on critical issues, warns on style violations

**Module Management (mod-tidy):**
- Ensures `go.mod` and `go.sum` are clean and minimal
- Removes unused dependencies automatically
- Verifies module integrity

**Built-in Text Processing:**
- **Whitespace**: Removes trailing whitespace from all text files
- **EOF**: Ensures all files end with a single newline character

## Developer Workflow

### Daily Development

```bash
# Normal workflow - hooks run automatically
git add .
git commit -m "feat: implement new feature"
# ✅ fumpt: 6ms
# ✅ lint: 68ms
# ✅ mod-tidy: 110ms
# ✅ whitespace: <1ms
# ✅ eof: <1ms
# ✅ Total: 185ms
```

### Manual Execution

```bash
# Run all checks on staged files
go-pre-commit run

# Run all checks on all files
go-pre-commit run --all-files

# Run specific checks only
go-pre-commit run --checks fumpt,lint

# Debug with verbose output
go-pre-commit run --verbose

# Preview what would run without executing
go-pre-commit run --dry-run
```

### Skipping Checks

Sometimes you need to bypass checks for work-in-progress commits:

```bash
# Skip specific checks
SKIP=lint git commit -m "wip: work in progress"

# Skip all pre-commit checks
GO_PRE_COMMIT_SKIP=all git commit -m "hotfix: critical fix"

# Skip checks using Git's built-in flag
git commit --no-verify -m "emergency: bypass all hooks"
```

### Status and Management

```bash
# Check installation status
go-pre-commit status
go-pre-commit status --verbose

# Uninstall hooks
go-pre-commit uninstall

# Reinstall hooks (useful after updates)
go-pre-commit install --force
```

## CI/CD Integration

go-pre-commit integrates seamlessly with GitHub Actions and other CI/CD systems:

### GitHub Actions

The tool automatically integrates with the fortress workflow system:

```yaml
# .github/workflows/fortress-pre-commit.yml (automatically managed)
pre-commit:
  name: 🪝 Pre-commit Checks
  if: needs.setup.outputs.pre-commit-enabled == 'true'
  uses: ./.github/workflows/fortress-pre-commit.yml
```

**CI Optimization:**
- Binary caching reduces installation time
- Tool caching speeds up repeated runs
- Smart file detection (all files vs changed files)
- Parallel execution with configurable workers

### Performance in CI

```bash
# CI performance (typical times)
Total pipeline: <2s (17x faster than Python pre-commit)
├── fumpt: 6ms
├── lint: 68ms
├── mod-tidy: 110ms
├── whitespace: <1ms
└── eof: <1ms
```

## Migration from Embedded GoFortress System

### Key Changes

The project has migrated from the embedded GoFortress pre-commit system (`.github/pre-commit/`) to the external go-pre-commit tool:

| Aspect           | Old (Embedded)                             | New (External)                                                         |
|------------------|--------------------------------------------|------------------------------------------------------------------------|
| **Location**     | `.github/pre-commit/gofortress-pre-commit` | `go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest` |
| **Installation** | `cd .github/pre-commit && magex build`     | `go-pre-commit install`                                                |
| **Maintenance**  | Part of repository                         | External tool, versioned independently                                 |
| **Updates**      | Manual code updates                        | `go install` latest version                                            |
| **Distribution** | Repository-specific                        | Reusable across Go projects                                            |

### Migration Steps

1. **Uninstall old system** (if previously installed):
   ```bash
   # Remove old hooks (if they exist)
   rm -f .git/hooks/pre-commit
   ```

2. **Install new system**:
   ```bash
   # Install external tool
   go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest

   # Install hooks
   go-pre-commit install
   ```

3. **Update configuration** (already done in `.github/.env.base`):
   ```bash
   # Configuration now uses GO_PRE_COMMIT_ prefixes
   ENABLE_GO_PRE_COMMIT=true
   GO_PRE_COMMIT_VERSION=v1.1.11
   # ... other settings
   ```

4. **Verify migration**:
   ```bash
   # Test the new system
   go-pre-commit run --all-files
   git add .
   git commit -m "feat: migrate to external go-pre-commit"
   ```

### Benefits of External Tool

- **Reusability**: Same tool works across multiple Go projects
- **Maintenance**: Independent versioning and updates
- **Community**: Shared improvements across projects
- **Simplicity**: No embedded code to maintain

## Troubleshooting

### Common Issues

**Installation Problems:**

```bash
# "go-pre-commit not found"
# Fix: Ensure GOPATH/bin is in your PATH
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# Alternative: Install to a directory in your PATH
go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest
```

**Hook Issues:**

```bash
# "Hook already exists"
go-pre-commit install --force

# "Permission denied"
chmod +x .git/hooks/pre-commit

# "Hooks not running"
go-pre-commit status --verbose  # Check installation status
```

**Performance Issues:**

```bash
# Increase timeout for slower systems
GO_PRE_COMMIT_TIMEOUT_SECONDS=180 git commit -m "test"

# Reduce parallel workers on resource-constrained systems
GO_PRE_COMMIT_PARALLEL_WORKERS=1 git commit -m "test"

# Enable fail-fast to stop on first error
GO_PRE_COMMIT_FAIL_FAST=true git commit -m "test"
```

### Environment Verification

```bash
# Verify installation
go-pre-commit --version
go env GOPATH
echo $PATH | grep "$(go env GOPATH)/bin"

# Test individual components
magex format    # Test fumpt integration
magex lint      # Test lint integration
magex tidy      # Test mod-tidy integration
```

## Performance Benchmarks

### Execution Time Comparison

| Operation           | Python pre-commit | go-pre-commit | Improvement |
|---------------------|-------------------|---------------|-------------|
| **Total Pipeline**  | ~15,000ms         | ~900ms        | 17x faster  |
| **Code Formatting** | ~200ms            | ~6ms          | 33x faster  |
| **Linting**         | ~3,000ms          | ~68ms         | 44x faster  |
| **Module Tidying**  | ~400ms            | ~110ms        | 4x faster   |
| **Text Processing** | ~100ms            | ~1ms          | 100x faster |

### Real-World Performance

**Development Workflow:**
- Small commits (1-5 files): <1s total
- Medium commits (10-20 files): <2s total
- Large commits (50+ files): <5s total

**CI/CD Pipeline:**
- Full repository scan: <10s
- Changed files only: <2s
- Cache hit rate: >90%

## Best Practices

### Team Setup

**Repository Configuration:**
1. Pin tool version in `.github/.env.base`: `GO_PRE_COMMIT_VERSION=v1.1.11`
2. Document installation in README or onboarding guides
3. Set up CI/CD integration for consistent enforcement
4. Use `.github/.env.custom` for developer-specific overrides

**Developer Onboarding:**
```bash
# Include in developer setup scripts
go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest
go-pre-commit install
```

### Maintenance

**Regular Updates:**
```bash
# Update go-pre-commit itself
go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest

# Update tool versions in .env.base
GO_PRE_COMMIT_GOLANGCI_LINT_VERSION=v2.6.0
GO_PRE_COMMIT_FUMPT_VERSION=v0.9.0
```

**Performance Monitoring:**
- Track commit times to ensure performance remains optimal
- Adjust `PARALLEL_WORKERS` based on team hardware
- Monitor CI/CD pipeline times for regressions

## Summary

go-pre-commit transforms the pre-commit experience for Go projects by providing:

- **17x faster execution** than Python alternatives
- **Zero external dependencies** beyond Go toolchain
- **Seamless integration** with existing build systems (magex)
- **Production-ready performance** for teams of any size
- **External tool maintenance** - no embedded code to maintain

The migration from the embedded GoFortress system to the external tool provides better maintainability and reusability while preserving all performance benefits.

**Next Steps:**
1. Install go-pre-commit: `go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest`
2. Set up hooks: `go-pre-commit install`
3. Verify with: `go-pre-commit run --all-files`
4. Start enjoying faster, more reliable pre-commit checks

---

📚 **Complete Documentation**: [github.com/mrz1836/go-pre-commit](https://github.com/mrz1836/go-pre-commit)
🚀 **Performance Reports**: See benchmarks section for detailed metrics
🤝 **Community Support**: Open issues for questions and feature requests
