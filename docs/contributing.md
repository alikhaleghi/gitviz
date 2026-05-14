# Contributing

Thanks for your interest in contributing to `gitviz`.

## Workflow

- open or link an issue before starting non-trivial work
- create a short-lived branch from `main`
- open a pull request early and keep it focused
- use squash merge into `main`

## Branch Naming

Use one of:

- `feat/<issue>-<short-topic>`
- `fix/<issue>-<short-topic>`
- `docs/<issue>-<short-topic>`
- `refactor/<issue>-<short-topic>`

Examples:

- `feat/4-load-commits`
- `docs/8-mvp-rules`

## Commit Format

`<type>(<scope>): <summary>`

Examples:

- `feat(tui): add status strip to dashboard`
- `fix(git): handle non-repository directories`
- `docs(readme): clarify MVP scope`
- `chore(repo): add GitHub templates`

### Allowed types

- `feat`
- `fix`
- `docs`
- `refactor`
- `test`
- `ci`
- `chore`

## Issue References

Use:

- `[refs #<id>]` for partial work
- `[closes #<id>]` only when the issue is fully completed

You may also put `Closes #<id>` in the PR body for automatic closure.

## Pull Request Checklist

- scope matches one issue or one clear concern
- tests/checks pass locally
- UI changes include terminal screenshots when useful
- docs updated for behavior or workflow changes

## Code Guidelines

- MVP first, polish second
- keyboard-first interactions
- clear empty/error states on every screen
- use Git CLI integration before deeper abstractions
