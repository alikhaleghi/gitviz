# gitviz — AGENTS.md

These rules are loaded into context every session. Follow them strictly.

## Engineering Principles

- MVP first, polish second
- No premature abstraction
- Shell out to `git` before deeper integrations
- Every feature must be keyboard-first
- Every screen must have clear empty/error states
- Keep startup under one second for normal repos
- Avoid hidden side effects for destructive git actions

## Branch Naming

```
feat/<issue>-<short-topic>
fix/<issue>-<short-topic>
docs/<issue>-<short-topic>
refactor/<issue>-<short-topic>
```

## Commit Format

```
<type>(<scope>): <summary> [refs #<id>]
```

Allowed types: `feat`, `fix`, `docs`, `refactor`, `test`, `ci`, `chore`

**Every commit MUST include an issue reference** — `[refs #<id>]` for partial work, `[closes #<id>]` only when fully completed. This links commits to the GitHub issue page.

## PR Workflow

- Short-lived branches from `main`
- Squash merge into `main`
- PRs must reference issue IDs
- PR body must follow `.github/PULL_REQUEST_TEMPLATE.md` exactly — fill every applicable section, leave or remove the rest

## Versioning

- Semver with stages: `vMAJOR.MINOR.PATCH[-STAGE.N]`
- Pre-1.0: breaking changes allowed, bumps must reflect scope

## Auth & Identity

- `gh` commands run as the Ahur System user — use for issue, PR, and release creation
- Git remote is set to `https://alikhaleghi:<token>@github.com/alikhaleghi/gitviz.git` — commits appear as Ali Khaleghi
- Workflow: Ahur System creates issues and infrastructure; Ali Khaleghi authors all solution commits
- Closed issue drafts are archived in `docs/issues/closed/` with filenames matching the GitHub issue number — only open drafts stay in `docs/issues/open/`

## Issue Manager

When creating or updating issues via `gh`:

1. **Determine type** from the draft or request:
   - `[Feature]` → use `.github/ISSUE_TEMPLATE/feature.md` sections (Summary, Problem, Proposed solution, Acceptance criteria, Notes)
   - `[Task]` → use `.github/ISSUE_TEMPLATE/task.md` sections (Task, Checklist, Notes)
   - `[Bug]` → use `.github/ISSUE_TEMPLATE/bug.md` sections (Description, Steps, Expected, Environment, Notes)
2. **Construct the body** using the exact section headers from the matching template
3. **Set title** with the type prefix: `[Feature]`, `[Task]`, or `[Bug]`
4. **Labels**: use `--add-label` with existing labels (`enhancement` for features, `bug` for bugs). Tasks may have no label if `task` doesn't exist. If label add fails (Ahur System permissions), note it for the user.
5. **Save draft** in `docs/issues/open/<number>-<short-topic>.md` with the returned issue number

When the user says "sync issue drafts":
- Move closed issue `.md` files into `docs/issues/closed/`
- Rename files to match the real GitHub issue number
- Keep only open draft files in `docs/issues/open/`

## GitHub Ops

- Every issue needs a type label + area label
- Milestones: `v0.0.1-dev` → `v0.0.1-alpha` → `v0.0.1-beta` → `v0.0.1`

## Version Bumping

When bumping the version (tag, milestone, stage), update EVERY location that references the old version:
- `README.md` — "Current:" line
- `docs/versioning.md` — example tags
- `docs/github-operations.md` — milestone references
- Any other file with the old version string

Search the entire project with `grep` before confirming the bump is complete.

## Law of Committing

- Break feature work into atomic, buildable commits — never commit everything at once
- After each commit, push to remote so GitHub links appear on the issue page
- After each commit, provide a structured comment message for the user to post on the issue. Format:

  ```
  Commit: [`<hash>`]

  <2-4 sentence summary of what changed and why. Include bullet points of key files changed.>

  Changes:
  - `<file>`: <change description>
  - `<file>`: <change description>
  ```
- Do not proceed to the next commit until the user confirms or acknowledges
- Each commit must compile (`go build ./...`) and tests must pass (`go test ./...`)
- Every commit MUST include `[refs #<id>]` or `[closes #<id>]` so GitHub links it to the issue — if the issue ID is unknown, ask the user before committing
- If multiple commits are pushed together, provide a single comment summarizing all of them

## Architecture

- `cmd/gitviz` — entrypoint
- `internal/` — project-private domain and UI logic
- `pkg/` — optional reusable packages
- TUI layout: commit list pane + details pane + footer/help row
