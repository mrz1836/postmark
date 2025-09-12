# Pull Request Guidelines

> Pull Requests—whether authored by humans or AI agents—must follow a consistent structure to ensure clarity, accountability, and ease of review.

<br><br>

## 🔖 Title Format

```
[Subsystem] Imperative and concise summary of change
```

Examples:

* `[API] Add pagination to client search endpoint`
* `[DB] Migrate legacy rate table schema`
* `[CI] Remove deprecated GitHub Action for testing`

> Use the imperative mood ("Add", "Fix", "Update") to match the style of commit messages and changelogs.

<br><br>

## 📝 Pull Request Description

Every PR must include the following **four** sections in the description:

### 1. **What Changed**

> A clear, bullet‑pointed or paragraph‑level summary of the technical changes.

### 2. **Why It Was Necessary**

> Context or motivation behind the change. Reference related issues, discussions, or bugs if applicable.

### 3. **Testing Performed**

> Document:
>
> * Test suites run (e.g., `TestCreateOriginationAccount`)
> * Edge cases covered
> * Manual steps that were taken (if any)

### 4. **Impact / Risk**

> Call out:
>
> * Breaking changes
> * Regression risk
> * Performance implications
> * Changes in developer experience (e.g., local dev setup, CI time)

<br><br>

## 💡 Additional PR Guidelines

* Link related issues with keywords like `Closes #123` or `Fixes #456` if there is a known issue.
* Keep PRs focused and minimal. Prefer multiple small PRs over large ones when possible.
* Use draft PRs early for feedback on in progress work.
* Releases are deployed using **goreleaser**.
* Rules for the release build are located in `.goreleaser.yml` and executed via `.github/workflows/release.yml`.

<br><br>

## 📋 PR Template Example

```markdown
## What Changed

* Added pagination support to the `/api/v1/users` endpoint
* Implemented cursor-based pagination using user IDs
* Added `limit` and `cursor` query parameters
* Updated OpenAPI specification

## Why It Was Necessary

Users were experiencing timeouts when fetching large user lists. This change implements
pagination to improve performance and reduce memory usage.

Closes #123

## Testing Performed

* Added unit tests in `TestUsersPagination`
* Tested edge cases:
  * Empty result sets
  * Invalid cursor values
  * Limits exceeding maximum
* Manual testing with 10k+ user dataset
* Verified backwards compatibility with existing clients

## Impact / Risk

* **Breaking Change**: None - pagination is optional
* **Performance**: 80% reduction in response time for large datasets
* **Risk**: Low - feature is behind query parameters
* **Migration**: None required
```

<br><br>

## 🔍 PR Review Checklist

Before requesting review, ensure:

- [ ] Code follows project conventions (see go-essentials.md)
- [ ] Tests are included and passing
- [ ] Documentation is updated (if applicable)
- [ ] Commits are clean and follow conventions
- [ ] PR description includes all four required sections
- [ ] No sensitive data or secrets are exposed
- [ ] CI checks are passing

<br><br>

## 🤝 Review Etiquette

### For Authors:
* **Respond to all comments** — Even if just to acknowledge
* **Push fixes as new commits** — Don't force-push during review
* **Be open to feedback** — Reviews improve code quality
* **Provide context** — Help reviewers understand your decisions

### For Reviewers:
* **Be constructive** — Suggest improvements, don't just criticize
* **Review promptly** — Aim for same-day initial review
* **Focus on important issues** — Nitpicks can be marked as optional
* **Approve explicitly** — Use GitHub's review approval feature

<br><br>

## 🚀 Merging

* **Squash and merge** for feature branches with messy history
* **Rebase and merge** for clean, logical commit sequences
* **Never force-push** to main/master branches
* **Delete branches** after merging to keep the repository clean

> The merge strategy may vary by project. Check with maintainers if unsure..
