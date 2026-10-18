# 14. Use test-driven development

Date: 2026-10-16

## Status

Accepted

The test library is amended by
[ADR-0020](0020-use-testify-for-test-assertions.md).

## Context

The challenge defines exact rules that can be checked automatically. For
example: a poll has between 2 and 10 options, votes are refused after the
closing time, a declined suggestion never appears on the ballot, and a
closed poll can be reopened.

In test-driven development (TDD), a test is written first and fails because
the feature does not exist yet; then just enough code is written to make it
pass; then the code is cleaned up with the test still passing.

This is the same decision as Python ADR-0014.

## Decision

Features are built test first.

- **Test runner:** the standard library's `testing` package, run with
  `go test ./...`. No test library is added to start.
- **Request-level tests:** the standard library's `net/http/httptest`
  package builds a request, sends it through the app's `http.Handler`, and
  records the response, with no real server running. Tests check the status
  code, headers, and returned HTML.
- **Rule tests:** functions in the `rules` package have no web or database
  code (ADR-0005), so they are tested by calling them directly.
- **Table-driven tests:** where one function is checked against many inputs,
  such as the 2-to-10 option rule, the cases are listed in a slice and run in
  a loop with `t.Run`, which is the common Go pattern.
- Tests live next to the code they test, in `_test.go` files.

## Consequences

- Every feature arrives with tests that describe its rules.
- Go has no `assert`. A check is a plain `if` that calls `t.Errorf` (report
  and continue) or `t.Fatalf` (report and stop this test).
- A new test that calls code that does not exist yet fails to compile. That
  is the expected first failure (ADR-0012).
- Tests check the HTML that comes back, not how the page looks. Layout and
  visual design are checked by eye.
- Tests will run in CI on every pull request once CI is set up (ADR-0008).