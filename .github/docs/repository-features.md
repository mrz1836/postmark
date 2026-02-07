# Repository Features

A comprehensive list of built-in features that ship with this repository.

<br>

---

<br>

* **Continuous Integration on Autopilot** with [GitHub Actions](https://github.com/features/actions) – every push is built, tested, and reported in minutes.
* **Pull‑Request Flow That Merges Itself** thanks to [auto‑merge](../workflows/auto-merge-on-approval.yml) and hands‑free [Dependabot auto‑merge](../workflows/dependabot-auto-merge.yml).
* **One‑Command Builds** powered by battle‑tested [MAGE-X](https://github.com/mrz1836/mage-x) targets for linting, testing, releases, and more.
* **First‑Class Dependency Management** using native [Go Modules](https://github.com/golang/go/wiki/Modules).
* **Uniform Code Style** via [gofumpt](https://github.com/mvdan/gofumpt) plus zero‑noise linting with [golangci‑lint](https://github.com/golangci/golangci-lint).
* **Confidence‑Boosting Tests** with [testify](https://github.com/stretchr/testify), the Go [race detector](https://blog.golang.org/race-detector), crystal‑clear [HTML coverage](https://blog.golang.org/cover) snapshots, and automatic reporting via internal coverage system.
* **Hands‑Free Releases** delivered by [GoReleaser](https://github.com/goreleaser/goreleaser) whenever you create a [new Tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging).
* **Relentless Dependency & Vulnerability Scans** via [Dependabot](https://dependabot.com) (runs daily at 8am to ensure broadcast dependencies are always current), [Nancy](https://github.com/sonatype-nexus-community/nancy), and [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck).
* **Security Posture by Default** with [CodeQL](https://docs.github.com/en/github/finding-security-vulnerabilities-and-errors-in-your-code/about-code-scanning), [OpenSSF Scorecard](https://openssf.org), and secret‑leak detection via [gitleaks](https://github.com/gitleaks/gitleaks).
* **Automatic Syndication** to [pkg.go.dev](https://pkg.go.dev/) on every release for instant godoc visibility.
* **Polished Community Experience** using rich templates for [Issues & PRs](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-go-broadcastsitory).
* **All the Right Meta Files** (`LICENSE`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SUPPORT.md`, `SECURITY.md`) pre‑filled and ready.
* **Code Ownership** clarified through a [CODEOWNERS](../CODEOWNERS) file, keeping reviews fast and focused.
* **Zero‑Noise Dev Environments** with tuned editor settings ([`.editorconfig`](../../.editorconfig)) plus curated *ignore* files for [VS Code](../../.editorconfig), [Docker](../../.dockerignore), and [Git](../../.gitignore).
* **Label Sync Magic**: your repo labels stay in lock‑step with [.github/labels.yml](../labels.yml).
* **Friendly First PR Workflow** – newcomers get a warm welcome thanks to a dedicated [workflow](../workflows/pull-request-management.yml).
* **Standards‑Compliant Docs** adhering to the [standard‑readme](https://github.com/RichardLitt/standard-readme/blob/master/spec.md) spec.
* **Instant Cloud Workspaces** via [Gitpod](https://gitpod.io/) – spin up a fully configured dev environment with automatic linting and tests.
* **Out‑of‑the‑Box VS Code Happiness** with a preconfigured [Go](https://code.visualstudio.com/docs/languages/go) workspace and [`.vscode`](../../.vscode) folder with all the right settings.
* **Optional Release Broadcasts** to your community via [Slack](https://slack.com), [Discord](https://discord.com), or [Twitter](https://twitter.com) – plug in your webhook.
* **AI Playbook** – machine‑readable guidelines in [tech conventions](../tech-conventions/ai-compliance.md)
* **Go-Pre-commit System** - [High-performance Go-native pre-commit hooks](https://github.com/mrz1836/go-pre-commit) with 17x faster execution—run the same formatting, linting, and tests before every commit, just like CI.
* **Zero Python Dependencies** - Pure Go implementation with [modular environment-based configuration](../env/README.md).
* **DevContainers for Instant Onboarding** – Launch a ready-to-code environment in seconds with [VS Code DevContainers](https://containers.dev/) and the included [.devcontainer.json](../../.devcontainer.json) config.

<br>

---

[← Back to README](../../README.md)
