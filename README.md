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

| Part      | Choice                                                 | Decision                                                         |
| --------- | ------------------------------------------------------ | ---------------------------------------------------------------- |
| Language  | Go 1.27                                                | [ADR-0004](docs/adr/0004-use-mise-and-go-modules-for-tooling.md) |
| Routing   | chi                                                    | [ADR-0002](docs/adr/0002-use-chi-for-routing.md)                 |
| Templates | templ                                                  | [ADR-0019](docs/adr/0019-use-templ-for-templates.md)             |
| Browser   | Alpine.js and Alpine AJAX                              | [ADR-0021](docs/adr/0021-use-alpine-and-alpine-ajax.md)          |
| Database  | SQLite, with sqlc for queries and goose for migrations | [ADR-0013](docs/adr/0013-use-sqlite-with-sqlc-and-goose.md)      |
| Settings  | Environment variables, secrets from Doppler            | [ADR-0017](docs/adr/0017-load-settings-from-environment.md)      |
| Tooling   | mise                                                   | [ADR-0004](docs/adr/0004-use-mise-and-go-modules-for-tooling.md) |

## Getting started

Requirements: [mise](https://mise.jdx.dev/), which installs Go and every
other tool the project uses.

Install the pinned tools:

```sh
mise install
```

Run the tests:

```sh
mise run test
```

Full setup, including Doppler, is in [CONTRIBUTING.md](CONTRIBUTING.md).

## Commands

| Command              | What it does                                           |
| -------------------- | ------------------------------------------------------ |
| `mise run test`      | Runs the test suite                                    |
| `mise run generate`  | Generates Go code from `.templ` files                  |
| `mise run format`    | Formats `.templ` and Go files                          |
| `mise run lint`      | Checks code and formatting                             |
| `mise run check`     | Checks dependencies for known vulnerabilities          |
| `mise run precommit` | Runs format, generate, lint, test, and check, in order |

## Documentation

- [CONTRIBUTING.md](CONTRIBUTING.md): setup, branching, commit messages, and
  code conventions.
- [`docs/adr/`](docs/adr/): every architecture decision and the reasoning
  behind it.
