# 22. Format and lint with golangci-lint and templ

Date: 2026-10-19

## Status

Accepted

Supersedes [ADR-0011](0011-use-gofmt-and-golangci-lint.md).

## Context

ADR-0011 chose `gofumpt` for formatting, run with `go fmt ./...`, and
golangci-lint v2 with its default linter set. Two gaps have appeared since:

- **`.templ` files were not formatted.** `go fmt` only reads `.go` files.
  The templ CLI ships its own formatter, `templ fmt`, which also keeps the
  generated `_templ.go` files in step with their templates (ADR-0019).
- **Imports were not grouped or cleaned up.** `gofumpt` sorts imports but does
  not add missing ones, remove unused ones, or separate the project's own
  packages from third-party ones. `goimports` does all three, and applies
  `gofumpt` as well.

golangci-lint v2 includes formatters as well as linters. `golangci-lint fmt`
rewrites files, and `golangci-lint run` reports any file that is not
formatted, so one configuration file covers both.

The default linter set does not check for security problems, such as an
`http.Server` with no timeouts or `math/rand` used for tokens, or for
mistakes specific to this project's libraries: database rows left open,
testify assertions used incorrectly, and inconsistent `log/slog` calls.

## Decision

- **Formatting:**
  - `.templ` files: `templ fmt .`
  - Go files: `golangci-lint fmt`, with the `gofumpt` and `goimports`
    formatters. `goimports` puts imports starting with
    `github.com/hbabb/tiebreaker-go` in their own group.
- **Linting:** `golangci-lint run`, starting from the `standard` set, with
  these added:
  - `gosec`: security checks.
  - `bodyclose`, `noctx`: HTTP response bodies are closed, and outgoing
    requests carry a context.
  - `rowserrcheck`, `sqlclosecheck`: database rows and statements are closed,
    and row errors are checked (ADR-0013).
  - `errorlint`: errors are compared with `errors.Is` and `errors.As`.
  - `testifylint`: testify is used correctly (ADR-0020).
  - `sloglint`: `log/slog` calls follow one style.
  - `usestdlibvars`: standard library constants, such as `http.MethodGet`,
    are used instead of raw strings.
  - `misspell`: common English misspellings.
  - `nolintlint`: every `//nolint` comment names the linter it silences and
    gives a reason.
- `noctx` is turned off for `_test.go` files. Test requests built with
  `httptest` never leave the process.
- Generated files, including `_templ.go`, are skipped by both linters and
  formatters. This is golangci-lint's default for files marked
  `Code generated ... DO NOT EDIT.`
- The mise `format` task runs `templ fmt .`, then `golangci-lint fmt`. The
  `generate` task runs `templ generate`, and runs after `format`, because
  formatting a template changes the positions recorded in its generated file.

## Consequences

- One command formats every file type in the project.
- `golangci-lint run` fails on unformatted code, so CI catches it without a
  separate formatting check.
- `go fmt ./...` is no longer the project's format command. Its output is
  still valid, but `goimports` may regroup imports afterwards.
- New linters can flag existing code. Each one added later is recorded in a
  new ADR or in the commit that adds it.
