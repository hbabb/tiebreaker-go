# 3. Server-rendered HTML with HTMX

Date: 2026-10-16

## Status

Superseded by [ADR-0021](0021-use-alpine-and-alpine-ajax.md).

## Context

The vote page is the most-visited screen and is usually opened on a phone
from a link in a group chat. It needs to load fast and work with no account.
The app also needs partial page updates: casting a vote without a reload,
live results, and moderating suggestions.

Options considered:

- **A separate single-page app** (React or Vue) calling a JSON API: two
  codebases, an API contract to keep in sync, and more JavaScript sent to
  phones.
- **Server-rendered HTML with HTMX**: the server returns HTML, and HTMX
  attributes in the markup fetch small HTML fragments and swap them into the
  page.

This is the same decision as Python ADR-0003.

## Decision

We will render HTML on the server and use HTMX for partial updates.

- Templates are written with templ (ADR-0019).
- Full pages are wrapped in a shared layout component, which places each
  page's content with `{ children... }`.
- Fragments that HTMX swaps in are templ components rendered on their own
  (for example `Results`). The full page uses the same component, so the
  fragment and the page always render the same markup.
- Routes return HTML, not JSON. There is no public API, so routes are not
  versioned. If a JSON API is added later, it gets its own versioned
  location (`/api/v1/`).
- JavaScript is limited to browser-only features HTMX does not cover, such as
  copying the poll link, the native share sheet, and the winner reveal
  (ADR-0010).
- CSRF protection is covered in ADR-0018. Unlike the Python version, no CSRF
  token is sent with HTMX requests.

## Consequences

- One codebase and one deployment for both the creator side and the vote
  page.
- Pages are usable as soon as the HTML arrives.
- The HTMX script file (`htmx.min.js`, version 2.0.11, the same file as the
  Python version) is served from `static/js/` and embedded in the
  compiled binary.
- There is no HTMX library on the server. The few HTMX request and response
  headers the app needs, such as checking `HX-Request`, are read and written
  directly with `net/http`.
- How live results are delivered is recorded in ADR-0015.
- The Python version marks fragments with a `_` file prefix
  (`_results.html`). That convention does not carry over: in templ, a
  fragment is an ordinary component.
