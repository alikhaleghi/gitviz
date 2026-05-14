# Architecture

## App Structure

- `cmd/gitviz`: application entrypoint
- `internal/`: project-private domain and UI logic
- `pkg/`: optional reusable packages if needed later

## Git Command Strategy

MVP shells out to `git` commands (`git log`, `git show`, `git branch`) for fast iteration and broad compatibility.

## UI Model

TUI is split into focused panes:

- commit list pane
- details pane
- footer/help row

## Why CLI over libgit2 for MVP

Using Git CLI first keeps complexity low, avoids CGO dependencies, and speeds up initial delivery.
