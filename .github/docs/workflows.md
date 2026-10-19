# GitHub Workflows

All GitHub Actions workflows in this repository are powered by modular configuration files – your one-stop shop for tweaking CI/CD behavior without touching a single YAML file!

<br>

## The Workflow Control Center

**Configuration Files:**
- **[`.github/env/`](../env/README.md)** – Modular environment configuration split into domain-specific files loaded in numeric order

These configuration files control everything from:
- **Go version matrix** (test on multiple versions or just one)
- **Runner selection** (Ubuntu or macOS, your wallet decides)
- **Feature toggles** (coverage, fuzzing, linting, race detection, benchmarks)
- **Security tool versions** (gitleaks, nancy, govulncheck)
- **Auto-merge behaviors** (how aggressive should the bots be?)
- **PR management rules** (size labels, auto-assignment, welcome messages)

<br>

## Workflows

| Workflow | Description |
|----------|-------------|
| [auto-merge-on-approval.yml](../workflows/auto-merge-on-approval.yml) | Automatically merges PRs after approval and all required checks, following strict rules. Runs from the default branch via `workflow_run` of the trigger. |
| [auto-merge-trigger.yml](../workflows/auto-merge-trigger.yml) | Permission-less `pull_request`/`pull_request_review` entry point (no secrets, no checkout). It wakes up auto-merge-on-approval on reviews and PR changes, dependabot-auto-merge on new or updated Dependabot PRs, and PR management for fork PRs. |
| [codeql-analysis.yml](../workflows/codeql-analysis.yml) | Analyzes code for security vulnerabilities using [GitHub CodeQL](https://codeql.github.com/). |
| [dependabot-auto-merge.yml](../workflows/dependabot-auto-merge.yml) | Automatically merges [Dependabot](https://github.com/dependabot) PRs that meet all requirements. Runs from the default branch via `workflow_run` of the trigger, verifies the PR and its commits through the API, and reads Dependabot's metadata from the verified commit. |
| [fortress.yml](../workflows/fortress.yml) | Runs the GoFortress security and testing workflow, including linting, testing, releasing, and vulnerability checks. |
| [pull-request-management.yml](../workflows/pull-request-management.yml) | Labels PRs by branch prefix, assigns a default user if none is assigned, and welcomes new contributors with a comment. Same-repo PRs run on `pull_request`; fork PRs run from the default branch via `workflow_run` of the trigger. |
| [scorecard.yml](../workflows/scorecard.yml) | Runs [OpenSSF](https://openssf.org/) Scorecard to assess supply chain security. |
| [stale.yml](../workflows/stale-check.yml) | Warns about (and optionally closes) inactive issues and PRs on a schedule or manual trigger. |
| [sync-labels.yml](../workflows/sync-labels.yml) | Keeps GitHub labels in sync with the declarative manifest at [`.github/labels.yml`](../labels.yml). |

The two auto-merge workflows approve and merge with a GitHub App instead of a personal access token. None of these workflows use `pull_request_target`. See **[Auto-merge with a GitHub App](auto-merge-github-app.md)** for how they are triggered and for the one-time setup.

<br>

---

[← Back to README](../../README.md)
