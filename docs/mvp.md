# MVP Scope

## In Scope

- commit list
- commit details panel
- branch list and branch switcher

## Out of Scope

- rebase helper
- advanced graph rendering
- GitHub integration
- diff syntax highlighting

## Engineering Rules

- MVP first, polish second
- no premature abstraction
- shell out to `git` before deeper integrations
- every feature must be keyboard-first
- every screen must have clear empty/error states
- keep startup under one second for normal repos
- avoid hidden side effects for destructive git actions
