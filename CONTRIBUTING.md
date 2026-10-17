# Contributing to Tiebreaker

This guide covers how to set up the project, how work moves from a change to
`main`, and how commit messages are written. The reasoning behind each rule
is recorded in [`docs/adr/`](docs/adr/).

## Setup

Requirements: [mise](https://mise.jdx.dev/). mise installs Go and the
project's command-line tools ([ADR-0004](docs/adr/0004-use-mise-and-go-modules-for-tooling.md)).

Clone the repository and move into it:

```sh
git clone git@github.com:hbabb/tiebreaker-go.git
```

```sh
cd tiebreaker-go
```

Install the pinned Go version and tools (set in `mise.toml`):

```sh
mise install
```

Link the folder to the Doppler project, which supplies secrets as
environment variables ([ADR-0017](docs/adr/0017-load-settings-from-the-environment.md)):

```sh
doppler setup
```

Start the development server:

```sh
mise run dev
```

Go downloads the libraries listed in `go.mod` the first time the app is
built. There is no separate install step.

## Branching

This project uses GitHub Flow ([ADR-0008](docs/adr/0008-use-github-flow.md)).

- `main` is always deployable.
- Work happens on a short-lived bookmark created from `main`.
- Changes reach `main` only through a pull request.

Bookmark names use the commit type and a short topic:

```
feat/poll-creation
fix/double-vote
docs/contributing
```

With jj, start a change from `main`:

```sh
jj new main
```

When the work is ready, name the bookmark:

```sh
jj bookmark create feat/poll-creation -r @
```

Push it:

```sh
jj git push --bookmark feat/poll-creation
```

Then open a pull request into `main` on GitHub.

## Commit messages

Commits follow [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/)
([ADR-0007](docs/adr/0007-use-conventional-commits.md)).

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

- **description**: short, lowercase, written as an instruction ("add", not
  "added"), with no period at the end.
- **body**: optional. Explains what changed and why, separated from the
  description by a blank line.
- **footer**: optional. For example `Refs: #12`, or `BREAKING CHANGE:`
  followed by a description.

### Types

| Type       | Use for                                               |
| ---------- | ----------------------------------------------------- |
| `feat`     | a new feature                                         |
| `fix`      | a bug fix                                             |
| `docs`     | documentation only                                    |
| `style`    | formatting that doesn't change behavior               |
| `refactor` | code changes that neither fix a bug nor add a feature |
| `perf`     | performance improvements                              |
| `test`     | adding or correcting tests                            |
| `build`    | dependencies or the build system                      |
| `ci`       | CI configuration                                      |
| `chore`    | maintenance that doesn't fit another type             |

### Scopes

Scopes match the project's packages
([ADR-0005](docs/adr/0005-layer-based-project-layout.md)):

| Scope      | Covers                                                   |
| ---------- | -------------------------------------------------------- |
| `config`   | `internal/config`: reading settings from the environment |
| `database` | `internal/database`, `migrations/`, and sqlc queries     |
| `rules`    | `internal/rules`: the app's rules as plain functions     |
| `server`   | `internal/server`: routes, middleware, and handlers      |
| `views`    | `internal/views`: templ components                       |
| `static`   | `static/`: CSS, JavaScript, and images                   |

A scope is optional. Leave it off when a change spans several packages or
the whole project.

### Examples

```
feat(views): add poll results fragment
fix(rules): reject polls with more than 10 options
docs: add contributing guide
build: add caarlos0/env
```

### One type per commit

If a change fits more than one type, split it into separate commits.

### Enforcement

jj does not run Git hooks, so a local `commit-msg` hook cannot check the
format ([ADR-0006](docs/adr/0006-use-jujutsu-with-git-backend.md)). The format
will be checked in CI on pull requests.

## Test-driven development

Features are built test first
([ADR-0014](docs/adr/0014-use-test-driven-development.md)), and the commit
history records that order:

1. Write a test for the new behavior and run `go test ./...` to confirm it
   fails. A test that calls code that does not exist yet fails to compile;
   that counts as the failing step.
2. Commit the failing test on its own, with the `test` type:
   `test: add failing home page test`
3. Write the code that makes the test pass, and commit it with the `feat` or
   `fix` type.
4. If the code needs cleaning up, do it with the tests still passing, and
   commit it with the `refactor` type.

A commit that contains only a failing test is expected. The pull request as
a whole must pass its tests.

## Formatting and linting

Code is formatted with `gofmt` and linted with golangci-lint
([ADR-0011](docs/adr/0011-use-gofmt-and-golangci-lint.md)):

```sh
mise run format
```

```sh
mise run lint
```

## Signed commits

Commits pushed to this repository should be signed
([ADR-0009](docs/adr/0009-sign-commits-with-ssh.md)). SSH signing is the
project's approach; any signing method that GitHub shows as **Verified** is
accepted.

