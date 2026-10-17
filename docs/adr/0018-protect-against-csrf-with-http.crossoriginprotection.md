# 18. Protect against CSRF with http.CrossOriginProtection

Date: 2026-10-16

## Status

Proposed

## Context

Cross-Site Request Forgery (CSRF) is when another website causes a visitor's
browser to send a request to this app, such as casting a vote or closing a
poll, using the visitor's existing session.

The Python version uses Litestar's token-based protection: the server signs a
token with a secret, and every form and HTMX request must send it back
(Python ADR-0003).

Go 1.25 added `http.CrossOriginProtection` to the standard library. It does
not use tokens. Instead:

- It reads the `Sec-Fetch-Site` header, which modern browsers add to every
  request to say whether it came from the same site.
- If that header is missing, it compares the `Origin` header with the `Host`
  header.
- If neither header is present, the request is treated as not coming from a
  browser and is allowed.
- Only unsafe methods (such as `POST`) are checked. `GET`, `HEAD`, and
  `OPTIONS` always pass.
- A cross-origin unsafe request is rejected with `403 Forbidden`.

Limitations:

- Browsers send `Sec-Fetch-Site` only to HTTPS sites and `localhost`, so
  production must use HTTPS.
- Older browsers may send neither header.
- It requires the same origin (scheme, host, and port), which is stricter
  than same-site.

Options considered:

- **A token library** such as `gorilla/csrf`: matches the Python approach,
  but needs a signing secret and a token in every form and HTMX request.
- **`http.CrossOriginProtection`**: standard library, no secret, no token in
  the markup.

## Decision

- The app's handler is wrapped with `http.CrossOriginProtection`.
- Session cookies are set with `SameSite=Lax`, as a second layer for browsers
  that send neither header.
- Production is served over HTTPS only.

## Consequences

- No CSRF token in forms or in `hx-headers` on `<body>`, and no CSRF secret
  in settings (ADR-0017).
- Protection depends on browser headers, so very old browsers rely on the
  `SameSite` cookie setting alone.
- Request-level tests must set `Sec-Fetch-Site` or `Origin` on test requests
  to check that cross-origin `POST` requests are refused, the same way the
  Python version tests that a `POST` without a token is refused.
