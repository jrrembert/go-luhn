# Contributing to go-luhn

Thank you for your interest in contributing!

## Prerequisites

- Go 1.19 or later

## Development setup

```bash
git clone https://github.com/jrrembert/go-luhn.git
cd go-luhn
go build ./...
go test ./...
```

## Test-Driven Development

This project follows TDD. For every feature or bug fix:

1. **Red** — write a failing test that defines the expected behavior
2. **Green** — write the minimum code to make the test pass
3. **Refactor** — clean up the code while keeping tests green

Run tests:

```bash
go test ./...
```

## Workflow

- Every change must have a GitHub issue. [Open one](https://github.com/jrrembert/go-luhn/issues/new/choose) before starting work.
- Create your branch from `rc` for features/fixes, or from `main` for chores/docs:
  ```bash
  git checkout rc
  git checkout -b feature/my-feature
  ```
- Use a git worktree to keep your working directory clean:
  ```bash
  git worktree add ../go-luhn-my-feature feature/my-feature
  ```
- Follow [Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`

## Pull requests

- Target `rc` for `feat:` and `fix:` changes; target `main` for `chore:`, `docs:`, `refactor:`, `test:`
- Fill in all sections of the PR template
- Include `Closes #N` in the PR body to auto-close the issue

## Code style

- `gofmt` for formatting
- `golangci-lint run` must pass before opening a PR

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
