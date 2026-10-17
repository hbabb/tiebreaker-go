# 17. Load settings from the environment, with secrets in Doppler

Date: 2026-10-16

## Status

Accepted

## Context

The app needs settings that differ between machines and must not be
committed, starting with the database location (ADR-0013). The Python version
reads settings with pydantic-settings and keeps secrets in Doppler (Python
ADR-0017 and ADR-0018).

Python ADR-0018 recorded two problems with reading settings by hand:

- **Secrets printed with the settings.** Printing or logging the settings
  exposed the secret value.
- **Bad values not caught clearly.** A missing secret failed with an
  unclear error, and an empty secret was accepted.

Both problems apply in Go. Go's `fmt` package and its structured logger,
`log/slog`, print every field of a struct.

Options considered for reading settings:

- **The standard library**: `os.LookupEnv` reports whether each variable is
  set, and `errors.Join` combines several problems into one error. Every
  setting needs its own hand-written checks.
- **`github.com/caarlos0/env/v11`**: fills a struct from environment
  variables using struct tags (annotations on each field, such as
  `env:"DATABASE_PATH"`). It applies a prefix to every name, supports
  defaults (`envDefault`), rejects missing values (`required`) and empty
  values (`notEmpty`), and returns all problems together. It has no way to
  hide a value when printed.

Options considered for supplying secrets on a developer machine:

- **A `.env` file**: a local, git-ignored file of secret values.
- **Doppler**: `doppler run -- <command>` fetches the project's secrets and
  passes them to the command as environment variables, so no secret is
  stored in the repository folder. The Python version already uses it.

## Decision

- Settings are read in `internal/config` with `caarlos0/env`, using the
  `TIEBREAKER_` prefix, into an exported `Settings` struct.
  - Required settings have no default and are tagged `required,notEmpty`.
  - Non-secret settings, such as the development database path, have
    defaults.
- Secrets use a `Secret` type, which hides its value when printed:
  - `String()` returns `**********`, which `fmt` uses for `%v`, `%+v`,
    and `%s`.
  - `GoString()` returns the same mask, which `fmt` uses for `%#v`.
  - `MarshalText()` returns the same mask, which JSON output uses, including
    `slog`'s JSON handler.
  - Code reads the real value only where it is needed, by converting it to a
    `string`.
- `Settings` fields are exported (capitalized). `caarlos0/env` cannot set
  unexported fields, and `fmt` cannot call `String()` on a value held in an
  unexported field, so an unexported secret field would print in plain text.
- `Load` accepts the environment as a map of names to values. The app passes
  the real environment; tests pass their own map, so tests run without
  Doppler or network access.
- Secrets are stored in Doppler. The repository folder is linked to the
  Doppler project with `doppler setup`. The development server is started
  through Doppler, wrapped in the mise `dev` task. There is no `.env` file.

## Consequences

- A missing or empty required setting stops the app at startup with one
  error that names every problem.
- Secrets stay hidden when settings are printed or logged, and a test checks
  this for `fmt` and `slog`.
- Settings are not read-only. Go has no equivalent of a frozen dataclass, so
  any code holding the `Settings` value could change a field. The app loads
  settings once in `main` and passes them down.
- With CSRF handled without a signing secret (ADR-0018), no setting needs a
  secret yet. The `Secret` type and Doppler are in place for the first one.
- New dependency: `github.com/caarlos0/env/v11`.
- Running the development server requires the Doppler CLI (pinned in
  `mise.toml`, ADR-0004) and a logged-in account.
