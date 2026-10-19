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

## How the workflows are triggered

None of the workflows use `pull_request_target`. GitHub's [workflow execution protections](https://github.blog/changelog/2026-09-17-workflow-execution-protections-in-github-actions-generally-available/) disable it by default on public repositories that have no Actions event policy allowing it, enforced from 2026-11-02. They use the `pull_request` → `workflow_run` pattern instead:

| Event | Workflow | Runs the copy from | Token and secrets |
|-------|----------|--------------------|-------------------|
| `pull_request`, `pull_request_review` | [auto-merge-trigger.yml](../workflows/auto-merge-trigger.yml) ("Auto-merge Trigger") | The PR (merge commit) | None (`permissions: {}`). It only echoes; each job is a wake-up hint for one consumer. |
| `workflow_run` of the trigger | [dependabot-auto-merge.yml](../workflows/dependabot-auto-merge.yml) | The default branch | `verify`: read-only. `analyze`: labels and comments. `apply`: GitHub App token from the `auto-merge` environment. |
| `workflow_run` of the trigger | [auto-merge-on-approval.yml](../workflows/auto-merge-on-approval.yml) | The default branch | `evaluate`: read-only. `apply`: GitHub App token from the `auto-merge` environment. |
| `workflow_run` of the trigger (fork PRs) | [pull-request-management.yml](../workflows/pull-request-management.yml), fork jobs | The default branch | `GITHUB_TOKEN`: a read-only lookup job finds the PR (skipping any ambiguous match), then labels, comments, cache cleanup. |
| `pull_request` (same-repo PRs) | [pull-request-management.yml](../workflows/pull-request-management.yml), same-repo job | The PR (merge commit) | `GITHUB_TOKEN` with the job's write scopes. Pushing a same-repo branch already takes write access. |

The consumers trust nothing the trigger's run produces: no artifacts, outputs, logs or files. Whether its job ran is only a hint that saves runs. They identify the PR from GitHub's event payload (`workflow_run.head_sha`, `head_repository`, `pull_requests[]`, or `GET /repos/{owner}/{repo}/commits/{sha}/pulls`) and re-read everything from the API.

For a Dependabot PR, the `verify` job checks that the PR is open, was opened by `dependabot[bot]`, isn't from a fork, is still on the commit the trigger ran for, and has only commits that Dependabot authored and GitHub signed. It then reads Dependabot's metadata itself, from the `updated-dependencies` block of that verified head commit's message. `dependabot/fetch-metadata` can't run under `workflow_run`: it only reads the PR from a `pull_request` event payload. The parser mirrors fetch-metadata's rules and fails closed: anything unexpected sends the PR to manual review.

<br>

## How the key is protected

| Control | What it stops |
|---------|---------------|
| The private key lives in an **environment secret** (`auto-merge`) that only protected long-lived branches can use. | A workflow pushed to any other branch can't read it. Repository and organization secrets don't give you this. |
| Both workflows that use the key run on `workflow_run`, whose ref is the default branch. | They always run the default branch's workflow file, never the copy inside the PR. Neither checks out PR code. |
| The job that holds the key runs only `actions/create-github-app-token` (pinned by SHA) and the `gh` CLI. | Any step in a job can read every secret that job uses. Keeping third-party code out of that job keeps it away from the key. |
| Approvals and auto-merge are pinned to the verified head commit. Dependabot approvals also require that Dependabot triggered the run (`workflow_run.actor`) and authored every commit, signed by GitHub. Every Dependabot evaluation that doesn't approve revokes the bot approval and the App's auto-merge. | Code pushed after the analysis can't ride an earlier approval: a commit by anyone else, a push someone else triggered, or a Dependabot update that now needs review dismisses the bot approval and disarms auto-merge. |
| The base branch must require status checks. | Auto-merge never merges before CI. |
| Dependabot approvals also require the base branch's rulesets to require an approval and void it on a new push. | A push that also edits the PR-side trigger can't ride an earlier bot approval (see below). |

### Required for Dependabot auto-merge: dismiss stale approvals on push

The revocation after someone else pushes to a Dependabot PR runs when the trigger runs for that push. The trigger is the PR's own copy, so anyone who can push to the PR and edit workflow files can stop it in the same push. The earlier bot approval would then still count, and the armed auto-merge would merge the new code once checks pass. That includes pushers who couldn't approve or merge on their own, such as a deploy key with write access.

So the Dependabot workflow only approves when a ruleset on the base branch:

- requires at least one approving review, **and**
- turns on **Dismiss stale pull request approvals when new commits are pushed**, or **Require approval of the most recent reviewable push**.

GitHub then drops the bot approval on any push, whatever the workflows do. The Dependabot workflow re-approves after Dependabot's own rebases. Without such a ruleset it leaves the PR for manual review (`skip-stale-approvals-kept`). The workflow can't read classic branch protection's review settings with its token, so the setting must come from a ruleset.

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
| Settings → **Environments** → `auto-merge` → Deployment branches | *Selected branches and tags* | The default branch: both workflows that use the key run on `workflow_run`, whose ref is always the default branch. Other long-lived branches (for example `development`) are no longer needed here, but listing them does no harm. Every branch listed must be protected by a ruleset that requires a pull request and status checks. |
| Settings → Environments → `auto-merge` → **Environment secrets** | `AUTO_MERGE_APP_PRIVATE_KEY` | The full contents of the App's `.pem` private key. |

Also make sure that:

- **Allow auto-merge** is enabled in the repository's general settings.
- The base branch has a ruleset or branch protection that **requires status checks**. Without one, the workflows refuse to auto-merge.
- The base branch's ruleset **requires an approval and dismisses stale approvals on push** (or requires approval of the most recent push). Without it, Dependabot PRs are left for manual review. See above.
- The three workflow files are on the **default branch**, and the trigger is on every branch that Dependabot PRs target. A `workflow_run` workflow only runs from the default branch, while the trigger runs from the PR's merge commit.
- If the repository has **Actions policies** (*Settings → Actions → Policies*), they allow these runs. Event rules must allow `pull_request`, `pull_request_review` and `workflow_run`. For actor rules: a Dependabot PR's runs, `workflow_run` included, are credited to `dependabot[bot]`. Runs caused by the App approving or arming are credited to the App.

Delete the `.pem` file from your machine afterwards. To rotate the key, generate a new one, update the environment secret, then delete the old key from the App's settings.

<br>

## Troubleshooting

| Symptom | Likely cause |
|---------|--------------|
| The apply job fails at **Create GitHub App token** | The App isn't installed on this repository, or the environment secret is missing or wrong. |
| *Branch "…" is not allowed to deploy to auto-merge* | The default branch isn't listed in the `auto-merge` environment's deployment branches. |
| Nothing runs for a new Dependabot PR | The PR has a merge conflict (GitHub runs no `pull_request` workflow until it's resolved), the trigger or the consumer isn't on the right branch (see above), or an Actions policy blocks the run (it ends in `startup_failure`). Check **Policy insights**. |
| **Verify Dependabot PR** ends in `skip-trigger-gated` | Normal for events that queue nothing (labels, edits, the App's own arming). Only opening, updating or reopening a Dependabot PR queues an evaluation. |
| **Verify Dependabot PR** ends in `skip-head-moved` | The PR changed after the trigger ran. The run for the new commit decides. |
| **Verify Dependabot PR** ends in `manual-no-metadata` | The head commit's message has no readable `updated-dependencies` block. It's left for manual review, and the apply job revokes any earlier bot approval. |
| **Verify Dependabot PR** ends in `revoke-unverified` | A GitHub API read failed (after retries) while verifying the PR, so the apply job revoked any earlier bot approval rather than leave it standing. Re-run the workflow to evaluate the PR again. |
| The analysis ends in `skip-no-required-checks` | The base branch doesn't require any status checks. |
| The analysis ends in `skip-stale-approvals-kept` | No ruleset on the base branch both requires an approval and dismisses stale approvals on push (or requires approval of the most recent push). See [above](#required-for-dependabot-auto-merge-dismiss-stale-approvals-on-push). |
| The run ends in `revoke-foreign-push` or `skip-foreign-commits` | Someone other than Dependabot pushed to the PR, so it needs a human review. The apply job revokes the bot approval and disarms auto-merge. |
| An **Auto-merge revoked** comment after Dependabot updated its own PR | The update changed what the PR does (for example, a grouped update now includes a major version), so the new decision doesn't auto-merge and the earlier bot approval was dismissed. |
| The apply job fails with **Revocation incomplete**, and the PR has an **Auto-merge revocation incomplete** comment | A GitHub API call failed while revoking. The job still disables auto-merge first and tries every dismissal, then lists what it couldn't undo. Remove any bot approval or auto-merge still on the PR by hand, or re-run the job. |
| The apply job reports **Nothing to revoke**, or that the App's approval of the current head stands | Normal. The evaluation didn't approve, but there was no bot approval or auto-merge to undo, or the PR's current commit already has an approval from an evaluation Dependabot triggered (for example, after a human reopened it, or after a newer run). |
| The PR is approved, but auto-merge waits forever | The ruleset requires a review the App can't satisfy (for example, code owner review), or an extra approval rule applies. |

<br>

---

[← Back to Workflows](workflows.md)
