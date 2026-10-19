# 21. Use Alpine.js and Alpine AJAX

Date: 2026-10-18

## Status

Accepted

Supersedes [ADR-0003](0003-server-rendered-html-with-htmx.md) and
[ADR-0015](0015-deliver-live-results-by-polling.md).

## Context

ADR-0003 rendered HTML on the server and used HTMX to fetch fragments and
swap them into the page. ADR-0015 delivered live results with HTMX polling
(`hx-trigger="every …"`), stopped by a `286` status.

The app also needs behavior that runs in the browser with no request to the
server, mostly on phones and on the creator side:

- adding and removing option fields on the create form (2 to 10 options)
- the mobile menu
- copying the poll link and opening the native share sheet
- the winner reveal

Alpine.js covers that behavior. Using it alongside HTMX means two libraries
that both change the page. When HTMX replaces part of the page, any Alpine
state held in the replaced elements, such as whether a menu is open, is lost.

Options considered:

- **Alpine.js and HTMX**: HTMX keeps its built-in polling and `286` stop
  signal. The two libraries have to be kept from conflicting when HTMX swaps
  content.
- **Alpine.js and Alpine AJAX**: Alpine AJAX is an Alpine.js plugin that
  lets forms and links send requests and swap the returned HTML into the
  page. One library and one set of conventions. It has no built-in polling,
  so polling is written in the app's own components.

The Python version stays on HTMX (Python ADR-0003 and ADR-0015).

## Decision

We will render HTML on the server and use Alpine.js with the Alpine AJAX
plugin for browser behavior and partial updates.

Rendering:

- Templates are written with templ (ADR-0019).
- Full pages are wrapped in a shared layout component, which places each
  page's content with `{ children... }`.
- Fragments that Alpine AJAX swaps in are templ components rendered on their
  own (for example `Results`). The full page uses the same component, so the
  fragment and the page always render the same markup.
- Routes return HTML, not JSON. There is no public API, so routes are not
  versioned. If a JSON API is added later, it gets its own versioned
  location (`/api/v1/`).
- CSRF protection is covered in ADR-0018.

Requests:

- Forms and links name the element to update with `x-target`, using that
  element's `id`. Requests started from code use the `$ajax` helper.
- Validation errors are returned with status `422` and sent to their own
  element with `x-target.422`.
- Alpine AJAX sends `X-Alpine-Request: true` with each request, and
  `X-Alpine-Target` with the ids being updated. Handlers read
  `X-Alpine-Request` to return only the fragment instead of the full page.

Live results:

- The `Results` component schedules the next request itself. While the poll
  is open, it renders with an `x-init` timer that calls `$ajax` once after
  the polling interval, targeting `results`. The response replaces the
  element, and the new element schedules the next request.
- When the poll is closed, the server renders `Results` without the timer.
  Nothing schedules another request, so polling stops. No special status
  code is used.
- The URL to poll is passed in a `data-url` attribute and read with
  `$el.dataset.url`.

Browser-only behavior:

- The create form, mobile menu, copy link, share sheet, and winner reveal
  are written with Alpine.js.

## Consequences

- One codebase and one deployment for both the creator side and the vote
  page.
- Pages are usable as soon as the HTML arrives.
- Two script files replace `htmx.min.js` in `static/js/`, both embedded in
  the compiled binary: Alpine AJAX (`v0.12.7` when added) and Alpine.js
  (`v3.17.4` when added). Both are loaded with `defer`, and Alpine AJAX's
  documentation requires its script to be included before Alpine's.
- Alpine AJAX is pre-1.0, so an upgrade can include breaking changes.
  Upgrades are done on purpose, one version at a time, with the full test
  suite run after each.
- There is no Alpine library on the server. The request headers the app
  needs are read directly with `net/http`.
- Alpine.js evaluates attribute values such as `x-init` as JavaScript.
  Values from the database are never written into those attributes. They are
  passed in `data-` attributes, which templ escapes.
- Alpine AJAX sets `aria-busy="true"` on the elements being updated during a
  request and disables a form's submit button while it sends, so loading
  states need no extra code.
- Results reach the browser up to one polling interval after a vote is cast.
  A request is sent every interval, whether or not the results have changed.
  These are reads, which SQLite serves alongside writes in WAL mode
  (ADR-0013).
- If a polling request fails, for example when a phone loses signal, no new
  `Results` element arrives, so polling stops for that visitor until the
  page is reloaded. Whether to show a message or retry on `ajax:error` is
  not yet decided.
- The polling interval is not yet decided.
- Confirmation dialogs, such as before closing a poll, have no built-in
  feature. They use the `ajax:before` event with the browser's `confirm()`.
- ADR-0010 still applies to separate script files in `static/js/`.
  Expressions written in Alpine attributes are not covered by JSDoc type
  checking.
- The `_` prefix for fragments from the Python version does not carry over:
  in templ, a fragment is an ordinary component.
