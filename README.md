A keyboard-driven terminal UI for exploring Git history, inspecting commits, and navigating branches — all without leaving your terminal.

```bash
go install github.com/alikhaleghi/gitviz/cmd/gitviz@latest
```

---

## Features

- **Commit browser** — scroll through recent commits with real-time detail loading
- **Commit inspector** — author, date, full message, and changed files at a glance
- **Dual-pane layout** — responsive split adapts to your terminal width
- **Fully keyboard-driven** — never reach for the mouse
- **Contextual help** — footer hints change based on what you're doing
- **Git CLI-native** — no CGO, no libgit2, just `git` commands

## Demo

```
┌───────────────────────────────────────────────────────┐
│ gitviz  Repo: detected  |  View: commits              │
│ Path: /home/user/project                              │
│───────────────────────────────────────────────────────│
│ ┌─ Commits ───────────┐  ┌─ Details ───────────────┐ │
│ │ ▸ a94bc2c  feat(git) │  │ a94bc2c                 │ │
│ │   10f55bb  feat(tui) │  │ Author:  John           │ │
│ │   eb87d9c  Merge PR  │  │ Date:    2026-05-15     │ │
│ │   136ffbe  feat(git) │  │                         │ │
│ │                      │  │ added commit detail     │ │
│ │                      │  │ inspector with...       │ │
│ │                      │  │                         │ │
│ │                      │  │ ── Files ──             │ │
│ │                      │  │ M  model.go             │ │
│ │                      │  │ M  details.go           │ │
│ └──────────────────────┘  └────────────────────────┘ │
│ Status: inspected a94bc2c                             │
│───────────────────────────────────────────────────────│
│ q quit  r refresh  ↑↓/w s move  enter/e inspect      │
└───────────────────────────────────────────────────────┘
```

## Keyboard Controls

| Key | Action |
|---|---|
| `↑` `↓` / `k` `j` / `w` `s` | Navigate commit list |
| `Enter` / `e` | Inspect selected commit |
| `Tab` / `Shift+Tab` | Switch focus between panes |
| `Esc` | Return to commit list |
| `r` | Refresh commit log |
| `b` | Branch view *(coming soon)* |
| `q` / `Ctrl+C` | Quit |

## Getting Started

### Install

```bash
# Install globally
go install github.com/alikhaleghi/gitviz/cmd/gitviz@latest

# Then run from any Git repo
cd my-project
gitviz
```

> **Note:** Binary installs to `$(go env GOPATH)/bin/gitviz`. Make sure that directory is in your `$PATH`, or run `$(go env GOPATH)/bin/gitviz` instead.

### From source

```bash
git clone https://github.com/alikhaleghi/gitviz.git
cd gitviz
make run        # run directly
# or
make build      # build to ./bin/gitviz
./bin/gitviz
```

### Prerequisites

- Go 1.24+
- Works best inside a Git repository (run it outside one too — you'll get a helpful hint)

## Project Structure

```
cmd/gitviz/          # Application entrypoint
internal/
├── app/             # Domain logic (coming soon)
└── tui/             # TUI components
    ├── model.go     # Core state, update loop, git operations
    ├── header.go    # Header bar renderer
    ├── commitlist.go# Commit list pane
    ├── details.go   # Commit detail inspector
    ├── footer.go    # Status strip and help bar
    └── styles.go    # Lip Gloss styles and helpers
pkg/                 # Reusable packages (future)
```

Every UI state — no repository, no commits, commit selected, load error — has a dedicated rendering path. No blank screens.

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — Go TUI framework (Model-View-Update)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — terminal styling

## Development

```bash
make run      # Run from source
make build    # Build binary to bin/gitviz
make test     # Run all tests
make fmt      # Format code
make lint     # Run go vet
```

### Design Principles

- **MVP first** — ship the smallest useful thing, then iterate
- **Keyboard-first** — every action must be reachable without a mouse
- **CLI-native** — shell out to `git` before reaching for complex integrations
- **No blank screens** — every pane has a meaningful empty, error, or loading state
- **Under a second** — startup must feel instant for normal-sized repos

## Roadmap

- [x] Project scaffold and TUI layout
- [x] Component architecture
- [x] Keyboard navigation
- [x] Git log integration
- [x] Commit detail inspector
- [ ] Branch list and switcher
- [ ] Branch checkout
- [ ] Graph rendering
- [ ] Diff viewer

## Versioning

[Semantic versioning](https://semver.org/) with pre-release stages: `vMAJOR.MINOR.PATCH[-STAGE.N]`

| Stage | Description |
|---|---|
| `dev` | Internal development snapshots |
| `alpha` | Early preview, core workflow exists |
| `beta` | Feature-complete for milestone, broader testing |
| *(none)* | Stable release |

Current: `v0.0.1-dev.3`

## Contributing

See [CONTRIBUTING.md](docs/contributing.md) for branch naming, commit format, and PR workflow.

## License

MIT — see [LICENSE](LICENSE).
