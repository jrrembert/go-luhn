# Release Process

This document describes the release process for `github.com/jrrembert/go-luhn`.

## Overview

Releases are fully automated via [semantic-release](https://github.com/semantic-release/semantic-release). The project uses a two-phase release process:

1. **Release candidates** — published from the `rc` branch for testing
2. **Stable releases** — published from `main` after RC validation

No manual versioning or tagging is needed. Go modules are published automatically via git tags.

## Release Candidate Workflow

Every stable release is preceded by one or more release candidates:

1. Push `feat:` or `fix:` commits to `rc` (via PR)
2. semantic-release creates prerelease tags (e.g., `v1.0.0-rc.1`, `v1.0.0-rc.2`)
3. A GitHub prerelease is created with auto-generated release notes
4. Test the RC: `go get github.com/jrrembert/go-luhn@v1.0.0-rc.1`
5. When satisfied, merge `rc` into `main` via PR

## Stable Release Workflow

When `rc` is merged into `main`, semantic-release:

1. Analyzes commit messages to determine the version bump
2. Creates a git tag (e.g., `v1.0.0`)
3. Creates a GitHub Release with auto-generated release notes

**Important**: Release PRs (`rc` → `main`) must use **merge commits** (`gh pr merge --merge`), not rebase. Rebase replays commits with new SHAs, causing `rc` and `main` to diverge and requiring a sync step after every release. Merge commits preserve the shared history between branches.

This is enforced automatically by `.github/workflows/auto-merge-release.yml` — when a PR from `rc` to `main` is opened, the workflow enables GitHub auto-merge with the merge commit strategy. Once all required status checks pass, the PR merges automatically.


Go module consumers can then install the stable version:

```bash
go get github.com/jrrembert/go-luhn@v1.0.0
```

## Commit Message → Version Bump

semantic-release uses [conventional commits](https://www.conventionalcommits.org/) to determine the version bump:

| Commit prefix | Version bump | Example |
|---|---|---|
| `fix:` | Patch (1.0.0 → 1.0.1) | `fix: handle empty string input` |
| `feat:` | Minor (1.0.0 → 1.1.0) | `feat: add batch validation` |
| `feat!:` or `BREAKING CHANGE:` footer | Major (1.0.0 → 2.0.0) | `feat!: rename Generate to Compute` |
| `chore:`, `docs:`, `refactor:`, `test:` | No release | `chore: update CI workflow` |

## Configuration

### Workflow

The GitHub Actions workflow (`.github/workflows/release.yml`) runs on every push to `main` or `rc`:

1. **Checkout** with full git history (`fetch-depth: 0`)
2. **Setup Node.js** (semantic-release runtime)
3. **Run** `npx semantic-release`

The workflow includes guardrails:
- **Fork protection**: only runs on `github.com/jrrembert/go-luhn`
- **Concurrency control**: prevents parallel releases on the same branch
- **Minimal permissions**: `contents: write` for tags, `issues: write` and `pull-requests: write` for release comments

### Semantic-release

Configured in `.releaserc.json` with three plugins:

1. `@semantic-release/commit-analyzer` — determines version bump from commits
2. `@semantic-release/release-notes-generator` — generates changelog
3. `@semantic-release/github` — creates GitHub Releases

No npm or Go-specific publish plugins are needed. Go modules are distributed via git tags and the Go module proxy automatically indexes tagged versions.

### Environment

The only required secret is `GITHUB_TOKEN`, which is automatically provided by GitHub Actions.

## Verification

After a release, verify it succeeded:

```bash
# Check git tags
git fetch --tags
git tag -l 'v*'

# Check GitHub Releases
gh release list

# Verify Go module proxy (may take a few minutes)
curl https://proxy.golang.org/github.com/jrrembert/go-luhn/@v/v1.0.0.info

# Install and test
go get github.com/jrrembert/go-luhn@v1.0.0
```

You can also check the [GitHub Releases page](https://github.com/jrrembert/go-luhn/releases) and the [Actions tab](https://github.com/jrrembert/go-luhn/actions/workflows/release.yml) for workflow run status.

### Dry Run

To test what semantic-release would do without creating a release:

```bash
npx semantic-release --dry-run
```

## Post-release sync

After each stable release, `rc` must be synced with `main` to pick up any `chore:`/`docs:` commits that targeted `main` directly. This is automated by `.github/workflows/sync-rc.yml`.

### How it works

1. A stable release tag (e.g., `v1.1.0`) triggers the `sync` job
2. The workflow sets a `sync/rc-up-to-date` commit status to **pending** on `rc`'s HEAD, blocking all PRs to `rc` via branch protection
3. It attempts to merge `main` into `rc`
4. **If clean merge**: pushes directly to `rc`. The push triggers the `clear-sync-status` job, which sets the status to **success**, unblocking PRs
5. **If conflicts**: creates a sync PR (`chore/sync-rc-after-release` → `rc`) for manual resolution. PRs to `rc` remain blocked until the sync PR is merged

### Manual conflict resolution

If the sync PR has conflicts:

1. Check out the `chore/sync-rc-after-release` branch locally
2. Resolve conflicts, commit, and push
3. Merge the sync PR — this pushes to `rc`, triggering the `clear-sync-status` job
4. Verify the `sync/rc-up-to-date` status is **success** on `rc`'s HEAD

## Troubleshooting

### No release was created

- Check commit messages — only `feat:` and `fix:` prefixes trigger releases
- Commits with `chore:`, `docs:`, `refactor:`, `test:` prefixes do not trigger releases
- Run `npx semantic-release --dry-run` locally to debug

### Wrong version bump

- Review commit messages for correct conventional commit format
- A `feat:` commit triggers minor bump, `fix:` triggers patch
- Use `BREAKING CHANGE:` footer or `!` suffix for major bumps

### Release workflow failed

- Check the [Actions tab](https://github.com/jrrembert/go-luhn/actions/workflows/release.yml) for error details
- Ensure `GITHUB_TOKEN` has write permissions (configured in workflow `permissions`)
- Verify the branch is `main` or `rc` (other branches don't trigger releases)
