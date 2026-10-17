# 2. Use chi for routing

Date: 2026-10-16

## Status

Accepted

## Context

The Tiebreak challenge calls for a deployed product with a database, creator
authentication, a public vote page that needs no account, and live-updating
results. The Python version uses Litestar, which bundles routing, sessions,
CSRF protection, templating, and rate limiting (Python ADR-0002).

Go has no single framework that is the default choice the way Litestar or
Django is in Python. Most Go web apps are built on the standard library's
`net/http` package, with small libraries added for specific jobs.

Since Go 1.22, `http.ServeMux` (the standard library router) matches the HTTP
method and path wildcards in one pattern, such as `GET /polls/{id}`. It does
not have route groups (a set of routes sharing a path prefix and middleware)
or ready-made middleware, such as request logging and panic recovery. Those
would be written by hand.

Options considered:

- **net/http alone**: no dependency. Route groups and middleware such as
  request logging and panic recovery are written by hand.
- **chi** (`github.com/go-chi/chi/v5`): a router built on `net/http`. Its
  handlers and middleware use the standard `net/http` types, so code written
  for one works with the other. It adds route groups (`Group`, `Route`,
  `Mount`) and a `middleware` package with about 27 ready-made middlewares,
  such as request logging, panic recovery, request IDs, and timeouts.
- **A full framework** such as Gin or Echo: these use their own handler types,
  so code is tied to the framework.

## Decision

We will route requests with chi v5.

- Routes are registered on a chi router in `internal/server` (ADR-0005),
  grouped by URL prefix and shared middleware with `Route` and `Group`.
- Middleware comes from chi's `middleware` package where it has what the app
  needs, added as each need comes up. Middleware the app writes itself uses
  the standard `func(http.Handler) http.Handler` form.
- Features chi does not cover, such as sessions and rate limiting, are added
  with libraries that also use the standard `net/http` types, as each feature
  needs them.

## Consequences

- Request logging, panic recovery, and similar middleware are used, not
  written. Writing them by hand can come later, as a learning exercise.
- Routes that share middleware, such as the creator-only pages, are grouped
  in one place.
- Handlers are plain `net/http` handlers, so they do not depend on chi. Moving
  to `http.ServeMux` later would change only how routes are registered. That
  would be recorded in a new ADR.
- Route wildcards are read with `chi.URLParam(r, "id")`.
- New dependency: `github.com/go-chi/chi/v5`.

