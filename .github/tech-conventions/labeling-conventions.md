# Labeling Conventions (GitHub)

> Labels serve as shared vocabulary for categorizing issues, pull requests, and discussions. Proper labeling improves triage, prioritization, automation, and clarity across the engineering lifecycle.

Current labels are located in `.github/labels.yml` and automatically synced into GitHub upon updating the `main/master` branch.

<br><br>

## 🎨 Standard Labels & Usage

| Label Name               | Color     | Description                                                         | When to Use                                                                 |
|--------------------------|-----------|---------------------------------------------------------------------|-----------------------------------------------------------------------------|
| `automerge`              | `#fef2c0` | Safe to merge automatically (e.g., from CI or bot)                  | Label added by automation or trusted reviewers                              |
| `automated-sync`         | `#fef2c0` | Used for referencing the automated sync process                     | PRs created by automated file sync tools or bots                            |
| `bug-P1`                 | `#b23128` | **Critical bug**, highest priority, impacts all users               | Regressions, major system outages, critical service bugs                    |
| `bug-P2`                 | `#de3d32` | **Moderate bug**, medium priority, affects a subset                 | Broken functionality with known workaround or scoped impact                 |
| `bug-P3`                 | `#f44336` | **Minor bug**, lowest priority, limited user impact                 | Edge case issues, cosmetic UI glitches, legacy bugs                         |
| `chore`                  | `#006b75` | Low-impact, internal tasks                                          | Dependency bumps, code formatting, comment fixes                            |
| `coverage-override`      | `#ff8c00` | Override coverage requirements for this PR                          | Special cases where coverage drop is acceptable (e.g., removing dead code)  |
| `dependabot`             | `#1a70c7` | PR or issue created by or related to Dependabot                     | Automated PRs/issues from Dependabot or related dependency bots             |
| `dependencies`           | `#006b75` | Dependency updates, version bumps, etc.                             | Any PR/issue related to dependency management or upgrades                   |
| `devcontainer`           | `#006b75` | Used for referencing DevContainers                                  | Changes to `.devcontainer` or container-based dev environments              |
| `docker`                 | `#006b75` | Used for referencing Docker related issues                          | Dockerfile, docker-compose, containerization changes                        |
| `documentation`          | `#0075ca` | Improvements or additions to project docs                           | Updates to `README`, onboarding docs, usage guides, code comments           |
| `feature`                | `#0e8a16` | Any new **major feature or capability**                             | Adding new API, CLI command, UI section, or module                          |
| `github-actions`         | `#006b75` | Used for referencing GitHub Actions                                 | Workflow, action, or CI/CD pipeline changes                                 |
| `gomod`                  | `#006b75` | Used for referencing Go Modules                                     | Changes to `go.mod`, `go.sum`, or module management                         |
| `hot-fix`                | `#b60205` | Time-sensitive or production-impacting fix                          | Used with `bug-P1` or urgent code/config changes that must ship immediately |
| `idea`                   | `#cccccc` | Suggestions or brainstorming candidates                             | Feature ideas, process improvements, early-stage thoughts                   |
| `npm`                    | `#006b75` | Used for referencing npm packages and dependencies                  | Changes to package.json, npm scripts, or Node.js tooling                    |
| `performance`            | `#8bc34a` | Performance improvements or optimizations                           | Optimizing code, reducing latency, improving resource usage                 |
| `prototype`              | `#d4c5f9` | Experimental work that may be unstable or incomplete                | Spike branches, POCs, proof-of-concept work                                 |
| `question`               | `#cc317c` | A request for clarification or feedback                             | Use for technical questions, code understanding, process queries            |
| `refactor`               | `#ffa500` | Non-functional changes to improve structure or readability          | Code cleanups, abstraction, splitting monoliths                             |
| `requires-manual-review` | `#ffd700` | PR or issue requires manual review by a maintainer or security team | Use when automation cannot determine safety or correctness                  |
| `security`               | `#d73a4a` | Security-related issue, vulnerability, or fix                       | Vulnerabilities, security patches, or sensitive changes                     |
| `size/L`                 | `#01579b` | Large change (201–500 lines)                                        | Large PRs, major refactors                                                  |
| `size/M`                 | `#0288d1` | Medium change (51–200 lines)                                        | Moderate PRs, new features, refactors                                       |
| `size/S`                 | `#4fc3f7` | Small change (11–50 lines)                                          | Small, focused PRs                                                          |
| `size/XL`                | `#002f6c` | Very large change (>500 lines)                                      | Massive PRs, sweeping changes                                               |
| `size/XS`                | `#b3e5fc` | Very small change (≤10 lines)                                       | Tiny PRs, typo fixes, single-line changes                                   |
| `stale`                  | `#c2e0c6` | Inactive, obsolete, or no longer relevant                           | Used for automated cleanup or manual archiving of old PRs/issues            |
| `test`                   | `#c2e0c6` | Changes to tests or test infrastructure                             | Unit tests, mocks, integration tests, CI coverage enhancements              |
| `ui-ux`                  | `#fbca04` | Frontend or user experience-related changes                         | CSS/HTML/JS updates, UI behavior tweaks, design consistency                 |
| `update`                 | `#006b75` | General updates not tied to a specific bug or feature               | Routine code changes, small improvements, silent enhancements               |
| `work-in-progress`       | `#fbca04` | Not ready to merge, actively under development                      | Blocks `automerge`, signals in-progress discussion or implementation        |

<br><br>

## 🧠 Labeling Best Practices

* Apply labels at the time of PR/issue creation, or during triage.
* Use **only one priority label** (`bug-P1`, `P2`, `P3`) per item.
* Combine labels as needed (e.g., `feature` + `ui-ux` + `test`).
* Don't forget to remove outdated labels (e.g., `work-in-progress` → after merge readiness).
* Use `requires-manual-review` for PRs that need human eyes before merging, especially if they involve security or critical changes.
* Use `automerge` only when the PR is fully approved, tested, and ready for production.
* Use `stale` for issues or PRs that have been inactive for 30+ days, to signal they may need attention or closure.
* Use `size/*` labels to indicate the complexity of the change, which helps reviewers understand the effort involved.
* Use `security` for any PR that addresses vulnerabilities, security patches, or sensitive changes that require extra scrutiny.

<br><br>

## 🤖 Automated Labeling

Several labels are applied automatically:

### Size Labels
GitHub Actions automatically calculate and apply size labels based on:
* Lines of code changed
* Number of files modified
* Complexity of changes

### Type Labels
Branch prefixes automatically map to labels:
* `feat/*` → `feature`
* `fix/*` → `bug-P3` (upgrade manually if higher priority)
* `docs/*` → `documentation`
* `test/*` → `test`
* `chore/*` → `chore`

### Bot Labels
* Dependabot PRs automatically receive `dependabot` and `dependencies`
* Stale bot applies `stale` after inactivity period

<br><br>

## 📋 Label Groups

Labels are conceptually grouped:

### Priority/Severity
* `bug-P1`, `bug-P2`, `bug-P3`
* `hot-fix`
* `security`

### Type of Change
* `feature`, `bug-*`, `documentation`, `test`
* `refactor`, `chore`, `update`
* `performance`

### Status/State
* `work-in-progress`
* `requires-manual-review`
* `stale`
* `automerge`
* `automated-sync`
* `coverage-override`

### Size/Scope
* `size/XS`, `size/S`, `size/M`, `size/L`, `size/XL`

### Technology/Area
* `docker`, `github-actions`, `gomod`
* `ui-ux`, `devcontainer`
* `dependencies`, `npm`

<br><br>

## 🔧 Managing Labels

### Adding New Labels
1. Update `.github/labels.yml`
2. Submit PR with justification
3. Labels sync automatically on merge

### Renaming Labels
1. Update in `labels.yml`
2. GitHub preserves label associations
3. Update any automation that references old names

### Removing Labels
1. Remove from `labels.yml`
2. Consider impact on existing issues/PRs
3. Manually clean up if needed

> Label changes should be discussed with maintainers to ensure consistency across the project.
