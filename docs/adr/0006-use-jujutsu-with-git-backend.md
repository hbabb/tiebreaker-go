# 6. Use Jujutsu with a Git backend

Date: 2026-10-16

## Status

Accepted

## Context

The project is hosted on GitHub, which requires Git. Jujutsu (jj) is a
version control tool that stores its history in Git, so it works with GitHub
while offering a different local workflow.

This is the same decision as Python ADR-0006.

## Decision

We will use jj for local version control in a colocated repository
(created with `jj git init`), with GitHub as the remote (`origin`).

## Consequences

- Anyone can clone and work with plain Git; jj is a local choice.
- jj does not run Git hooks. `pre-commit` and `commit-msg` hooks installed in
  `.git/hooks/` never fire for commits created by jj. Checks that would
  normally run as hooks (such as commit message format) are enforced in CI
  instead.
- jj snapshots every file in the working copy automatically, so `.gitignore`
  must be in place before generated files appear, such as the compiled
  binary, coverage reports, and the SQLite database in `data/`.
- Commit signing is configured in jj's own config, separate from Git's
  (see ADR-0009).