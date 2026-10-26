# Tiebreaker

A group polling app: create a poll, share the link, and let everyone vote
from their phone with no account.

This is the Go version of Tiebreaker, built from a
[Frontend Mentor](https://www.frontendmentor.io) product challenge. The
original Python version lives at
[hbabb/tiebreaker](https://github.com/hbabb/tiebreaker).

## Status

Early development. The app serves a home page and its static files. Polls,
voting, and accounts are not built yet.

## Stack

| Part      | Choice                                                 | Decision                                                            |
| --------- | ------------------------------------------------------ | ------------------------------------------------------------------- |
| Language  | Go 1.27                                                | [ADR-0023](docs/adr/0023-use-just-and-go-tool-directive-for-tooling.md) |
| Routing   | chi                                                    | [ADR-0002](docs/adr/0002-use-chi-for-routing.md)                    |
| Templates | templ                                                  | [ADR-0019](docs/adr/0019-use-templ-for-templates.md)                |
| Browser   | Alpine.js and Alpine AJAX                              | [ADR-0021](docs/adr/0021-use-alpine-and-alpine-ajax.md)             |
| Database  | SQLite, with sqlc for queries and goose for migrations | [ADR-0013](docs/adr/0013-use-sqlite-with-sqlc-and-goose.md)         |
| Settings  | Environment variables, secrets from Doppler            | [ADR-0017](docs/adr/0017-load-settings-from-environment.md)         |
| Tooling   | Just, and Go's `tool` directive in `go.mod`            | [ADR-0023](docs/adr/0023-use-just-and-go-tool-directive-for-tooling.md) |

## Getting started

Requirements:

- [Go](https://go.dev/) 1.27 or newer
- [Just](https://just.systems/)
- [golangci-lint](https://golangci-lint.run/)
- The [Doppler CLI](https://docs.doppler.com/docs/cli)

templ, goose, sqlc, and govulncheck are pinned in `go.mod` and run with
`go tool`. Go downloads and builds them on first use.

Run the tests:

```sh
just test
```

Full setup, including Doppler, is in [CONTRIBUTING.md](CONTRIBUTING.md).

## Commands

| Command           | What it does                                                        |
| ----------------- | ------------------------------------------------------------------- |
| `just`            | Lists every recipe                                                  |
| `just test`       | Runs the test suite                                                 |
| `just templ-fmt`  | Formats `.templ` files                                              |
| `just generate`   | Generates Go code from `.templ` files                               |
| `just fmt`        | Formats Go files                                                    |
| `just lint`       | Checks code and formatting                                          |
| `just vuln`       | Checks dependencies for known vulnerabilities                       |
| `just tidy`       | Cleans up `go.mod` and `go.sum`                                     |
| `just pre-commit` | Runs tidy, templ-fmt, generate, fmt, lint, vuln, and test, in order |

## Documentation

- [CONTRIBUTING.md](CONTRIBUTING.md): setup, branching, commit messages, and
  code conventions.
- [`docs/adr/`](docs/adr/): every architecture decision and the reasoning
  behind it.
