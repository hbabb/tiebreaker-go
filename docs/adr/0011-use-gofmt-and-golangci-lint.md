# 11. Use gofmt and golangci-lint for formatting and linting

Date: 2026-10-16

## Status

Superseded by [ADR-0022](0022-format-and-lint-with-golangci-lint-and-templ.md).

## Context

Go code needs a formatter (for consistent style) and a linter (to flag bugs
and bad patterns). The Python version uses Ruff for both (Python ADR-0011).

Go ships its own formatter, `gofmt`, and nearly all Go code is formatted with
it. There are no style options to choose, so there is nothing to configure.

Go also ships `go vet`, which reports likely mistakes such as format strings
that do not match their arguments. It covers a small set of checks.

Options considered for wider linting:

- **`go vet` alone**: built in, but narrow.
- **staticcheck** on its own: a widely used set of checks for bugs and
  simplifications.
- **golangci-lint**: runs many linters, and can include `go vet` and
  staticcheck among them, in one command with one configuration file.

## Decision

- **Formatting:** `gofmt`, run on the whole project with `go fmt ./...`.
- **Linting:** golangci-lint v2, run with `golangci-lint run`.
  - Its version is pinned in `mise.toml` (ADR-0004).
  - Configuration lives in `.golangci.yml`, which starts with
    `version: "2"`.
  - The default linter set (`standard`) is used to start. Linters are added
    to it with `linters.enable` when a need comes up, and each addition is
    noted in this ADR's successor or in the commit that adds it.
- Both run as mise tasks (`mise run format`, `mise run lint`).

## Consequences

- Formatting is never debated: the output of `gofmt` is the style.
- One tool and one configuration file for linting, like Ruff in the Python
  version.
- golangci-lint is installed as a prebuilt binary through mise, as its own
  documentation recommends, rather than built with `go install`.
- Type checking is covered separately in ADR-0012.
