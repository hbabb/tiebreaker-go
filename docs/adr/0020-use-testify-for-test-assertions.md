# 20. Use testify for test assertions

Date: 2026-10-17

## Status

Accepted

Amends [ADR-0014](0014-use-test-driven-development.md), which said no test
library would be added to start.

## Context

ADR-0014 uses the standard library's `testing` package on its own. Go has
no `assert`, so each check is a plain `if` with a hand-written failure
message:

```go
if recorder.Code != http.StatusOK {
	t.Errorf("status = %d, want %d", recorder.Code, http.StatusOK)
}
```

Every check needs its own message, and how clear the failure output is
depends on how well each message was written.

Options considered:

- **The standard library only**: no dependency, but every check is an `if`
  plus a message written by hand.
- **testify** (`github.com/stretchr/testify`): a widely used assertion
  library. A check is one call, such as
  `assert.Equal(t, http.StatusOK, recorder.Code)`, and testify writes the
  failure message, labeling the expected and actual values:

  ```
  Error:      Not equal:
              expected: 200
              actual  : 418
  ```

## Decision

Tests use testify for assertions.

- **`assert`** reports a failure and lets the test keep going, like
  `t.Errorf`. It is the default.
- **`require`** reports a failure and stops the test, like `t.Fatalf`. It
  is used when the rest of the test cannot run after a failure, such as when
  a value the later checks depend on is missing.
- Values are passed expected first, then actual:
  `assert.Equal(t, expected, actual)`. Swapping them still runs, but the
  failure message labels them backwards.
- Everything else in ADR-0014 is unchanged: `go test ./...` runs the tests,
  `httptest` drives request-level tests, and table-driven tests use
  `t.Run`.

## Consequences

- Checks are one line each, and failure messages are consistent across the
  test suite.
- New dependency: `github.com/stretchr/testify` (v1.12.1 when added). It
  brings in `github.com/stretchr/objx` and `go.yaml.in/yaml/v3`. It is used
  only by `_test.go` files, so it is not compiled into the app's binary.
- Reading testify's failure output is a separate concern from how test runs
  are displayed. A tool for that, such as gotestsum, would be a separate
  decision.