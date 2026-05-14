# Versioning Policy

This project uses semantic versioning with pre-release stages.

Format:

`vMAJOR.MINOR.PATCH[-STAGE.N]`

Examples:

- `v0.0.1-dev.1`
- `v0.0.1-alpha.1`
- `v0.0.1-beta.1`
- `v0.0.1`

## Stages

### dev

Internal development snapshots for active milestone implementation.

### alpha

Early preview release. Core workflow exists but is incomplete or unstable.

### beta

Feature-complete for the milestone scope and ready for broader testing.

### stable

No stage suffix. Release is considered usable for its declared scope.

## Bump Rules

### Patch

Increment patch for:

- bug fixes
- documentation updates
- small UX improvements
- internal cleanup without major user-facing scope changes

### Minor

Increment minor for:

- meaningful user-facing features
- significant workflow additions
- milestone-level scope expansion

### Major

Increment major when the project reaches a stable public contract, starting with `v1.0.0`.

## Pre-1.0 Policy

Before `v1.0.0`, breaking changes are allowed, but version bumps must still reflect visible scope and release maturity rather than development effort.
