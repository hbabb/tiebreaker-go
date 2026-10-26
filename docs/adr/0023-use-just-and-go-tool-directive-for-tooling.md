# 23. Use Just and the go.mod tool directive for tooling

Date: 2026-10-25

## Status

Accepted

Supersedes [ADR-0004](0004-use-mise-and-go-modules-for-tooling.md).

## Context

ADR-0004 used mise for three jobs: installing Go, installing the project's
command-line tools, and running tasks. Each job now has a better home.

- **Go itself** is installed and kept up to date by the operating system's
  package manager. A second copy managed by mise is not needed.
- **Command-line tools** can be tracked in `go.mod`. Go 1.24 added the
  `tool` directive: `go get -tool <package>@<version>` pins a tool in
  `go.mod` the same way a library is pinned, and `go tool <name>` builds and
  runs it, caching the binary after the first run. The tool versions are
  then covered by `go.sum` and by Dependabot, with no second version file.
- **Tasks** move to Just. A `justfile` holds named recipes, each with a
  comment that `just --list` shows as its description. A recipe can list
  other recipes as dependencies, which run in order before it.

golangci-lint's documentation discourages building it from source with
`go install` or the `tool` directive (ADR-0004), so it stays a prebuilt
binary installed outside the project.

## Decision

- **Go** comes from the system package manager. `go.mod` sets the minimum
  version with its `go` line.
- **Go modules** manage libraries in `go.mod` and `go.sum`, unchanged from
  ADR-0004.
- **Command-line tools** built in Go are pinned with the `tool` directive in
  `go.mod` and run with `go tool`: templ, goose, sqlc, and govulncheck.
- **golangci-lint** is installed as a prebuilt binary and called directly as
  `golangci-lint`. CI installs it with `golangci/golangci-lint-action`.
- **The Doppler CLI** is installed separately. It is not a Go tool that
  `go get -tool` can pin.
- **Tasks** are Just recipes in `justfile` at the project root, replacing
  the mise tasks. `mise.toml` is removed.

## Consequences

- Setting up the project needs Go, Just, golangci-lint, and the Doppler CLI.
  Every other tool is fetched by `go` on first use.
- Tool versions now share `go.mod` with the app's libraries. A tool's own
  dependencies appear in `go.mod` and `go.sum` as indirect requirements,
  though they are not compiled into the app.
- Dependabot opens update pull requests for the pinned tools, as it does
  for libraries.
- The first `go tool <name>` run builds the tool, so it is slower than later
  runs.
- Mentions of mise and `mise.toml` in ADR-0005, ADR-0011, ADR-0017,
  ADR-0019, and ADR-0022 describe the setup at the time of those decisions.
  This ADR replaces them; those records are not edited.
