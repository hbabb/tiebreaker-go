# 12. Rely on the Go compiler for type checking

Date: 2026-10-16

## Status

Accepted

## Context

The Python version needs a separate type checker, Pyright, because Python
does not check type hints when code runs (Python ADR-0012).

Go is statically typed: every value has a type known when the code is
compiled, and the compiler refuses to build code with a type error. Editors
with Go support report the same errors while code is being written.

## Decision

There is no separate type checker. Type checking is done by the Go compiler
on every build and every test run (`go build ./...`, `go test ./...`).

## Consequences

- Type errors cannot reach a running app, because code with a type error
  does not compile.
- In test-driven development (ADR-0014), a new test that calls a function
  that does not exist yet fails as a compile error. That is the expected
  "red" step.
- Code that the compiler cannot check, such as field names used inside
  templates (ADR-0019) and SQL text, is checked by tests and by sqlc
