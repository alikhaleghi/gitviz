# GitHub Operations

This file defines the day-1 operating model for `gitviz`.

## Repository

Target repository: `alikhaleghi/gitviz`

## Branch Rules

- `main` should remain mergeable and stable
- all work happens in short-lived feature branches
- all merges to `main` go through pull requests
- prefer squash merge

## Labels

### Type

- `feature`
- `bug`
- `task`
- `docs`
- `refactor`

### Priority

- `priority:high`
- `priority:medium`
- `priority:low`

### Status

- `blocked`
- `in-progress`
- `good first issue`

### Area

- `ui`
- `git`
- `core`
- `docs`

## Milestones

- `v0.0.1-dev`
- `v0.0.1-alpha`
- `v0.0.1-beta`
- `v0.0.1`

## Project Board Columns

- Backlog
- Ready
- In Progress
- Review
- Done

## Issue Workflow

- each issue should have a type label and area label
- assign milestone when issue contributes to a release target
- branches and PRs should reference issue IDs

## Initial Issue Source

Use `docs/issues/*.md` as source drafts, then create matching GitHub issues.
