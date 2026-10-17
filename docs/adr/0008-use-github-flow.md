# 8. Use GitHub Flow for branching

Date: 2026-10-16

## Status

Accepted

## Context

Work needs to reach `main` through pull requests. Options considered:

- **Git Flow**, with a long-lived `develop` branch: features merge into
  `develop`, which later merges into `main`. Its author, Vincent Driessen,
  added a note in 2020 recommending a simpler workflow such as GitHub Flow
  for continuously delivered web apps, and keeping Git Flow for software
  that is explicitly versioned or has multiple versions supported at once.
- **GitHub Flow**: short-lived feature branches created from `main` and
  merged back into `main` by pull request.

Tiebreaker is a single web app with one version in production. This is the
same decision as Python ADR-0008.

## Decision

We will use GitHub Flow.

- `main` is always deployable.
- Each change is made on a short-lived bookmark created from `main`, named
  by type and topic (for example `feat/poll-creation`).
- Changes reach `main` only through pull requests.
- There is no `develop` bookmark.

## Consequences

- No second long-lived branch to keep in sync with `main`.
- Keeping `main` deployable depends on checks running on every pull request,
  so CI is required before this rule can be enforced.