# 16. Use Go doc comments

Date: 2026-10-16

## Status

Accepted

## Context

The Python version documents code with Google-style docstrings, chosen from
several competing formats (Python ADR-0016).

Go has one documented format for doc comments, described at
<https://go.dev/doc/comment>, used by `go doc`, pkg.go.dev, and editors. It
has no `Args:` or `Returns:` sections. Parameters and results are described
in plain sentences, and types are already in the function signature.

## Decision

Code is documented with Go doc comments.

- Every exported (capitalized) name has a doc comment directly above it,
  written as full sentences that start with the name:
  `// Load reads settings from the environment.`
- Each package has one package comment starting with `Package`:
  `// Package config reads application settings from the environment.`
- Functions that return a boolean are described with "reports whether".
- Comments inside function bodies are used for lines whose purpose is not
  obvious from the code.

## Consequences

- There is no format to choose and no checker configuration to write; `gofmt`
  normalizes doc comment layout.
- Whether a linter is turned on to require doc comments on exported names is
  left to ADR-0011's linter configuration.
