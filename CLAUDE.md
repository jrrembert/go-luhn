# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Go library implementing the Luhn algorithm for generating and validating checksums (commonly used for credit card validation). Published as `github.com/jrrembert/go-luhn`.

## Commands

```bash
go build ./...                        # Build the library
go test ./...                         # Run all tests
go test -run TestGenerate             # Run tests matching a pattern
go test -v ./...                      # Run all tests with verbose output
go vet ./...                          # Run Go vet
golangci-lint run                     # Lint with golangci-lint
```

## Releases

Releases are fully automated via [semantic-release](https://github.com/semantic-release/semantic-release). No manual versioning or publishing is needed.

- **Release candidates**: Push `feat:`/`fix:` commits to the `rc` branch → publishes pre-release versions (e.g., `1.0.0-rc.1`)
- **Stable releases**: Merge `rc` into `main` → publishes stable versions (e.g., `1.0.0`)
- Every stable release must be preceded by at least one release candidate
- After each release, the workflow automatically registers the new version with the [Go module proxy](https://proxy.golang.org) so it's immediately discoverable on [pkg.go.dev](https://pkg.go.dev/github.com/jrrembert/go-luhn) and available via `go get`

## Architecture

- `luhn.go` — All algorithm implementations; exports 6 public functions: `Generate`, `Validate`, `Random`, `GenerateModN`, `ValidateModN`, `ChecksumModN`
- `luhn_test.go` — Co-located tests (Go standard `testing` package)
- `go.mod` — Module definition (`github.com/jrrembert/go-luhn`, no external dependencies)
- `docs/SPEC.md` — Canonical algorithm specification with all 106 test vectors

## Git Workflow

- **Default PR target is `rc`** — all `feat:` and `fix:` branches target `rc`, not `main`
- `chore:`, `docs:`, `refactor:`, `test:` branches may target `main` directly (they don't trigger releases)
- To publish a stable release, merge `rc` → `main` via PR
- Never push directly to `main` or `rc`
- Never force push to protected branches (`main`, `rc`) without explicit approval from the repo admin
- Every change must have a GitHub issue. If the user provides an issue number, use it. Otherwise, create one before starting work
- Every PR must include `Closes #N` in the body to auto-close its issue on merge
- **Note**: `Closes #N` only auto-closes when merging into the default branch (`main`). PRs targeting `rc` won't auto-close issues — include `Closes #N` in the release PR (`rc` → `main`) or close manually
- Branch naming: use prefixes `feature/`, `fix/`, `chore/` (e.g., `feature/add-auth`, `fix/login-bug`)
- Commits, PR titles, and issue titles follow conventional commit format: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`, `perf:`, `ci:`
- Always merge PRs via GitHub UI or `gh pr merge` — never merge locally with `git merge` then push. Local merges break GitHub's `Closes #N` auto-close linking.
- **Merge strategy**:
  - `--rebase` for all regular PRs: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`, `perf:`, `ci:`
  - `--merge` (merge commit) for release PRs (`rc` → `main`) — keeps the branches in sync and avoids SHA divergence
- PRs use the template at `.github/PULL_REQUEST_TEMPLATE.md` — fill in all sections (Summary, Changes, Test plan)
- Use `/pr` or `/pr <issue-number>` to create pull requests with the standard format
- Always create the feature branch from the appropriate base branch (`rc` for features/fixes, `main` for chore/docs) **before** writing code, not at commit time
- Use a git worktree for implementation work (e.g., `git worktree add ../go-luhn-<name> <branch>`). This avoids issues with stale branch state in the main working directory. Do not remove the worktree until the PR has been merged

This project follows Test-Driven Development (TDD). For every feature or bug fix:

1. **Red**: Write a failing test first that defines the expected behavior
2. **Green**: Write the minimum code to make the test pass
3. **Refactor**: Clean up the code while keeping tests green

## Code Style

`gofmt` and `golangci-lint` enforce Go style. Key rules:
- Standard Go formatting (tabs, not spaces)
- All exported functions must have doc comments
- Error strings are lowercase, no trailing punctuation (Go convention)
- Function names use Go conventions: `GenerateModN`, `ValidateModN`, `ChecksumModN`

## PR Review Workflow

- To fetch line-level review comments, filter to only needed fields:
  ```bash
  gh api repos/{owner}/{repo}/pulls/{n}/comments --jq '.[] | {id, body, path, line}'
  ```
- Reply to a comment using its `id`:
  ```bash
  gh api repos/{owner}/{repo}/pulls/{n}/comments/{id}/replies -f body="..."
  ```

## Documentation

- When making changes that affect setup, build, or test commands, update both this file and `README.md`
- Project specification lives in `docs/SPEC.md` — update when requirements or architecture change
