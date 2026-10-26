# 4. Use mise and Go modules for tooling

Date: 2026-10-16

## Status

Superseded by [ADR-0023](0023-use-just-and-go-tool-directive-for-tooling.md).

## Context

The project needs a pinned Go version, pinned dependencies, and pinned
command-line tools (sqlc, goose, golangci-lint) so every machine gets the
same environment.

Go has dependency management built in. `go.mod` records the module path, the
Go version, and the required libraries. `go.sum` records a checksum for each
downloaded library, so a changed download is detected. Together they play
the role of `pyproject.toml` and `uv.lock` in the Python version. There is
no virtual environment: libraries are stored in a shared cache outside the
project, and each project builds against the versions in its own `go.mod`.

Options for installing command-line tools:

- **`go install` or the `tool` directive in `go.mod`**: builds the tool from
  source. golangci-lint's documentation discourages both for its own tool,
  because the result depends on the local Go version and tool dependencies
  can mix with the project's dependencies.
- **mise**: installs prebuilt release binaries and pins their versions in
  `mise.toml`, alongside Go itself.

## Decision

- **mise** pins Go 1.27 in `mise.toml`.
- **Go modules** manage libraries in `go.mod` and `go.sum`. The module path
  is `github.com/hbabb/tiebreaker-go`.
- **mise** also pins the command-line tools in `mise.toml`: sqlc,
  golangci-lint, and the Doppler CLI from the mise registry, and goose from
  its GitHub releases (`github:pressly/goose`), since goose is not in the
  mise registry.
- Common commands are mise tasks, the same as the Python version
  (`mise run test`, `mise run lint`, `mise run dev`).

## Consequences

- Anyone with mise can reproduce the environment from the repository:
  `mise install` installs Go and the tools, and `go` downloads libraries on
  the first build.
- Tool versions do not mix with the app's library versions in `go.mod`.
- The compiled app is a single binary file with no runtime to install, which
  simplifies deployment later. Docker may still be added for production; that
  is a separate decision.
