# Modular Environment Configuration

Configuration is split into domain-specific files loaded in numeric order. Later files override earlier ones (last wins).

<br>

## File Structure

```
.github/env/
├── load-env.sh           # Shell loader script
├── README.md             # This file
│
├── 00-core.env           # Go versions, runners, feature flags, timeouts
├── 10-coverage.env       # go-coverage settings
├── 10-mage-x.env         # MAGE-X build system configuration
├── 10-pre-commit.env     # go-pre-commit settings
├── 10-security.env       # Gitleaks, Nancy, Govulncheck
├── 20-redis.env          # Redis service configuration
├── 20-workflows.env      # Stale, labels, dependabot, PR management
├── 90-project.env        # Project-specific overrides (not synced)
└── 99-local.env          # Local development (gitignored)
```

<br>

## Naming Convention

| Prefix | Purpose | Examples |
|--------|---------|----------|
| `00-` | Core / foundation | Go versions, runners, feature flags |
| `10-` | Tool configuration | mage-x, coverage, pre-commit, security |
| `20-` | Services & workflows | Redis, workflow automation |
| `90-` | Project overrides | Project-specific settings (not synced) |
| `99-` | Local development | Machine-specific (gitignored) |

<br>

## Override Behavior

Files load in sorted order. Variables in later files override earlier ones:

```bash
# 00-core.env
GO_PRIMARY_VERSION=1.24.x

# 90-project.env (loaded later, wins)
GO_PRIMARY_VERSION=1.23.x  # This value is used
```

<br>

## Usage

### In GitHub Actions

The loader is called automatically by the `load-env` composite action:

```yaml
- uses: ./.github/actions/load-env
  id: load-env
```

### Local Development

```bash
source .github/env/load-env.sh

# With verbose output
source .github/env/load-env.sh --verbose

# Or via environment variable
ENV_LOADER_VERBOSE=1 source .github/env/load-env.sh
```

### Verifying Configuration

```bash
# Check a specific variable
source .github/env/load-env.sh && echo $GO_PRIMARY_VERSION

# List all exported variables
source .github/env/load-env.sh && env | grep -E '^[A-Z_]+'
```

<br>

## Adding New Variables

1. **Identify the domain** — which file does this variable belong in?
2. **Add with a comment** — explain what the variable controls
3. **Test locally** — source the loader and verify
4. **Commit** — changes sync to other repos via go-broadcast (if configured)

```bash
# In 10-mage-x.env
# Maximum parallel workers for mage builds
MAGE_X_MAX_WORKERS=4
```

<br>

## Project Overrides (`90-project.env`)

Settings specific to this repository that should not be synced to other repos:

```bash
# 90-project.env - Project-specific overrides
# These settings are NOT synced via go-broadcast

GO_COVERAGE_THRESHOLD=80.0
ENABLE_REDIS_SERVICE=true
```

<br>

## Local Development (`99-local.env`)

Create `99-local.env` for machine-specific settings (gitignored):

```bash
# 99-local.env - Local development overrides

GO_COVERAGE_USE_LOCAL=true
GO_COVERAGE_LOCAL_PATH=/Users/me/projects/go-coverage/go-coverage
MAGE_X_VERBOSE=true
```

<br>

## CI Behavior

- In CI (`CI=true`), the loader skips `99-local.env`
- Variables are exported so downstream workflow steps can access them
- The composite action converts exported vars to JSON for workflow compatibility

<br>

## Troubleshooting

**Variables not available in workflow steps** — Ensure `set -a` is enabled in the loader (exports all sourced variables).

**Wrong value being used** — Check load order. Later files override earlier ones. Use `--verbose` to see which files are loaded.

**Local overrides not working** — Make sure `99-local.env` exists and you're not running in CI mode.
