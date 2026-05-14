# gitviz

A terminal UI for exploring Git history, branches, and commits.

## Status

Early prototype / scaffold phase.

## Vision

`gitviz` aims to make Git history easier to understand in the terminal through:

- commit browsing
- branch exploration
- commit inspection
- visual history rendering
- interactive operations like checkout and rebase helpers

## MVP

The first version focuses on:

- viewing recent commits
- keyboard navigation
- commit detail panel
- branch list and checkout

## UI Direction

The dashboard is structured into four layers:

1. header bar
2. main content area
3. status/context strip
4. footer help bar

Supported first-class states include:

- no repository detected
- loading repository data
- no commits found
- command error

## Tech Stack

- Go
- Bubble Tea
- Bubbles
- Lip Gloss

## Run

```bash
go run ./cmd/gitviz
```

## Project Principles

- build small, useful increments
- keep the UI responsive
- prefer simple Git CLI integration first
- avoid overengineering the first versions

## Governance

- contribution guide: `docs/contributing.md`
- versioning policy: `docs/versioning.md`
- architecture: `docs/architecture.md`
- MVP scope: `docs/mvp.md`

## Roadmap

- [ ] commit list
- [ ] commit inspector
- [ ] branch switcher
- [ ] graph rendering
- [ ] diff viewer
- [ ] rebase helper

## Release Track

- `v0.0.1-dev`: scaffold + dashboard shell
- `v0.0.1-alpha`: basic git detection + commit flow
- `v0.0.1-beta`: MVP complete, broader testing
- `v0.0.1`: first public MVP release
