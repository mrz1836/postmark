# Release Workflow & Versioning

> Structured releases ensure predictable deployments and clear communication of changes.

<br><br>

## 🚀 Semantic Versioning

We follow **Semantic Versioning (✧ SemVer)**:
`MAJOR.MINOR.PATCH` → `1.2.3`

| Segment   | Bumps When...                         | Examples        |
|-----------|---------------------------------------|-----------------|
| **MAJOR** | Breaking API change                   | `1.0.0 → 2.0.0` |
| **MINOR** | Back‑compatible feature / enhancement | `1.2.0 → 1.3.0` |
| **PATCH** | Back‑compatible bug fix / docs        | `1.2.3 → 1.2.4` |

<br><br>

## 📦 Release Tooling

* Releases are driven by **[goreleaser]** and configured in `.goreleaser.yml`.
* Install locally with Homebrew (Mac):
```bash
  brew install goreleaser
```

<br><br>

## 🔄 Release Workflow

| Step | Command                                   | Purpose                                                                                            |
|------|-------------------------------------------|----------------------------------------------------------------------------------------------------|
| 1    | `magex release:snapshot`                   | Build & upload a **snapshot** (pre‑release) for quick CI validation.                               |
| 2    | `magex version:bump push=true bump=patch` | Create and push a signed Git tag. Triggers GitHub Actions to package the release                   |
| 3    | GitHub Actions                            | CI runs `goreleaser release` on the tag; artifacts and changelog are published to GitHub Releases. |

> **Note for AI Agents:** Do not create or push tags automatically. Only the repository [codeowners](../CODEOWNERS) are authorized to tag and publish official releases.

[goreleaser]: https://github.com/goreleaser/goreleaser

<br><br>

## 📝 Changelog Management

### Automatic Generation
GoReleaser automatically generates changelogs from commit messages:
* Groups commits by type (`feat`, `fix`, `docs`, etc.)
* Excludes certain commit types (configured in `.goreleaser.yml`)
* Links to PRs and issues mentioned in commits

### Manual Additions
For significant releases, you may want to add a manual summary:
1. Create a draft release on GitHub
2. Edit the auto-generated changelog
3. Add a "Highlights" section at the top
4. Call out breaking changes prominently

<br><br>

## 🏷️ Version Tags

### Tag Format
* Release tags: `v1.2.3` (always prefix with `v`)
* Pre-release tags: `v1.2.3-rc.1`, `v1.2.3-beta.2`
* Development snapshots: Generated automatically, not tagged

<br><br>

## 📦 Release Artifacts

GoReleaser produces:
* **Binaries** for multiple platforms (Linux, macOS, Windows)
* **Docker images** (if configured)
* **Checksums** for verification
* **Release notes** from commits
* **Source archives** (tar.gz, zip)

All artifacts are automatically uploaded to GitHub Releases.

<br><br>

## 🔍 Pre-Release Checklist

Before tagging a release:

- [ ] All linters passing (`magex lint`)
- [ ] All tests passing (`magex test`)
- [ ] No security vulnerabilities (`magex deps:audit`)
- [ ] Documentation updated
- [ ] Version bumped if needed
- [ ] PR merged to main branch

<br><br>

## 🚨 Hotfix Process

For critical production fixes:

1. Create branch from the release tag: `git checkout -b hotfix/security-fix v1.2.3`
2. Apply minimal fix
3. Test thoroughly
4. Tag as patch: `v1.2.4`
5. Cherry-pick to main branch
6. Document in release notes

<br><br>

## 📊 Version History

Check release history:
```bash
# List all tags
git tag -l

# Show specific release info
git show v1.2.3

# Compare versions
git log v1.2.2..v1.2.3 --oneline
```

<br><br>

## 🤖 Automation

The release process is largely automated via GitHub Actions:
* **Trigger**: Push of a tag matching `v*`
* **Workflow**: `.github/workflows/fortress-release.yml`
* **Configuration**: `.goreleaser.yml`
* **Permissions**: Requires `GITHUB_TOKEN` with release permissions

> Manual intervention should rarely be needed. If issues arise, check the Actions tab.
