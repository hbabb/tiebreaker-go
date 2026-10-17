# 1. Record architecture decisions

Date: 2026-10-16

## Status

Accepted

## Context

Tiebreaker is built from a Frontend Mentor product challenge, which provides a
product spec but no design or technical stack. This repository is a Go build
of the same app as [hbabb/tiebreaker](https://github.com/hbabb/tiebreaker),
which is written in Python with Litestar.

Every technical choice is made by this project, and the reasoning behind those
choices needs to stay available to anyone reading the code later, including
future contributors and reviewers of this portfolio project. Where a decision
differs from the Python version, the reason for the difference also needs to
be recorded.

## Decision

We will record architecture decisions as Architecture Decision Records
(ADRs), using the format described by Michael Nygard: Title, Status, Context,
Decision, and Consequences.

ADRs live in `docs/adr/`, numbered sequentially with four digits
(`0001-…`, `0002-…`). An accepted ADR is not edited to change its decision.
When a decision changes, a new ADR is written, and the old one's status is
updated to "Superseded by ADR-NNNN".

ADR numbers in this repository follow the same topics as the Python
repository where possible, so ADR-0013 covers the database in both. The
numbering is otherwise independent.

Day-to-day rules that follow from these decisions (setup steps, commit
format, branching) are written in `CONTRIBUTING.md`.

## Consequences

- The reasoning behind each decision, including the options rejected, is
  kept in the repository next to the code.
- A reader can compare a Go decision with its Python counterpart by number.
- Each significant decision takes a few minutes to write up.