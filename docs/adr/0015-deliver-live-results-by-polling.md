# 15. Deliver live results by polling

Date: 2026-10-16

## Status

Superseded by [ADR-0021](0021-use-alpine-and-alpine-ajax.md).

## Context

The challenge requires results that update live while a poll is open.

Options considered:

- **Server-Sent Events (SSE)**: the browser keeps one connection open and
  the server pushes new results when a vote arrives. This needs the separate
  `htmx-ext-sse` extension. It is limited to 6 connections per site in the
  browser on HTTP/1.1, and if the app ever ran as more than one process, a
  vote received by one process would have to reach connections held by the
  others.
- **Polling**: the browser asks for the latest results on a fixed interval
  using HTMX's `hx-trigger="every …"`.

This is the same decision as Python ADR-0015.

## Decision

Live results are delivered by HTMX polling, for these reasons:

- Every update is an ordinary request, so it works on any host.
- There is no need to pass messages between processes.
- It avoids the browser's limit of 6 SSE connections per site on HTTP/1.1.
- The server stops polling when a poll closes by replying with status `286`,
  which HTMX treats as "stop polling". `net/http` has no named constant for
  286, so the app defines one.

## Consequences

- Results reach the browser up to one polling interval after a vote is cast.
- A request is sent every interval, whether or not the results have changed.
  These are reads, which SQLite serves alongside writes in WAL mode
  (ADR-0013).
- The results fragment is the `Results` templ component, rendered on its own
  for each polling request and inside the full vote page (ADR-0003,
  ADR-0019).
- The polling interval is not yet decided.
