# Auto-merge with a GitHub App

The [dependabot-auto-merge.yml](../workflows/dependabot-auto-merge.yml) and [auto-merge-on-approval.yml](../workflows/auto-merge-on-approval.yml) workflows approve PRs and enable auto-merge with a **GitHub App installation token**. They do not use personal access tokens.

Until the App is configured, both workflows still run their analysis, labels and alerts, and skip only the approve/merge step.

<br>

## Why a GitHub App

- **Short-lived, narrow tokens.** Each run mints a token for the current repository only, with only `contents: write` and `pull-requests: write`. It expires within an hour and is revoked when the job ends.
- **No personal account in the loop.** Approvals and merges show up as `<app-name>[bot]`, not as a person.
- **Merges trigger CI.** Merges made with the built-in `GITHUB_TOKEN` don't start new workflow runs, so push workflows on the base branch would silently not run. An App token avoids that.
- **Contained blast radius.** Use one App per account or organization, installed only on that account's repositories. A leaked key for one org can't reach another.

<br>

## How the key is protected

| Control | What it stops |
|---------|---------------|
| The private key lives in an **environment secret** (`auto-merge`) that only protected long-lived branches can use. | A workflow pushed to any other branch can't read it. Repository and organization secrets don't give you this. |
| `dependabot-auto-merge.yml` runs on `pull_request_target`, and `auto-merge-on-approval.yml` runs on `workflow_run`. | Both run the workflow file from a protected branch, never the copy inside the PR. Neither checks out PR code. |
| The job that holds the key runs only `actions/create-github-app-token` (pinned by SHA) and the `gh` CLI. | Any step in a job can read every secret that job uses. Keeping third-party code out of that job keeps it away from the key. |
| Approvals and auto-merge are pinned to the analyzed head commit. Dependabot approvals also require that Dependabot triggered the run and authored every commit, signed by GitHub. | Code pushed after the analysis can't ride an earlier approval. A push by anyone else revokes the bot approval. |
| The base branch must require status checks. | Auto-merge never merges before CI. |

<br>

## One-time setup

### 1. Create the App (once per account or organization)

Create it in the account's settings (*Developer settings → GitHub Apps → New GitHub App*):

| Setting | Value |
|---------|-------|
| Name | `<account>-automerge` (App names are unique across GitHub) |
| Webhook | **Inactive** (no events) |
| Repository permissions | **Contents: Read and write**, **Pull requests: Read and write** (Metadata is read-only and automatic) |
| Everything else | No access |
| Where can it be installed | **Only on this account** |

Then **install** it on the repositories that run these workflows (*Only select repositories*), and generate a private key.

### 2. Configure each repository

| Where | Name | Value |
|-------|------|-------|
| Settings → Secrets and variables → Actions → **Variables** | `AUTO_MERGE_APP_CLIENT_ID` | The App's Client ID. Setting this turns the approve/merge step on. |
| Settings → **Environments** → `auto-merge` → Deployment branches | *Selected branches and tags* | The default branch, plus any other long-lived branch that PRs target (for example `development`). Every branch listed must be protected by a ruleset that requires a pull request and status checks. |
| Settings → Environments → `auto-merge` → **Environment secrets** | `AUTO_MERGE_APP_PRIVATE_KEY` | The full contents of the App's `.pem` private key. |

Also make sure that:

- **Allow auto-merge** is enabled in the repository's general settings.
- The base branch has a ruleset or branch protection that **requires status checks**. Without one, the workflows refuse to auto-merge.

Delete the `.pem` file from your machine afterwards. To rotate the key, generate a new one, update the environment secret, then delete the old key from the App's settings.

<br>

## Troubleshooting

| Symptom | Likely cause |
|---------|--------------|
| The apply job fails at **Create GitHub App token** | The App isn't installed on this repository, or the environment secret is missing or wrong. |
| *Branch "…" is not allowed to deploy to auto-merge* | The PR's base branch isn't listed in the `auto-merge` environment's deployment branches. |
| The analysis ends in `skip-no-required-checks` | The base branch doesn't require any status checks. |
| The analysis ends in `skip-foreign-commits` or `revoke-foreign-push` | Someone other than Dependabot pushed to the PR, so it needs a human review. |
| The PR is approved, but auto-merge waits forever | The ruleset requires a review the App can't satisfy (for example, code owner review), or an extra approval rule applies. |

<br>

---

[← Back to Workflows](workflows.md)
