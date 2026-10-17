# 5. Layer-based project layout

Date: 2026-10-17

## Status

Accepted

## Context

Go does not enforce a folder structure. The `go` tool only requires a
`go.mod` file at the root, and treats each folder of `.go` files as one
package. A package is the unit that other code imports.

Three rules from the Go toolchain shape every layout:

- **Every folder is a package.** A new folder is a new boundary: names must
  be capitalized (exported) to be used from outside it, and other packages
  must import it.
- **Import cycles do not compile.** If package A imports B, B cannot import
  A, directly or through other packages.
- **`internal/` is enforced by the compiler.** Code outside this module
  cannot import any package under a folder named `internal`.

Tests sit next to the code they test. `go test` finds `_test.go` files in the
same folder as the package; a test in a different folder belongs to a
different package and can only reach exported names.

### Sources

- **Go's official guide**, [Organizing a Go module](https://go.dev/doc/modules/layout).
  For server projects, it recommends keeping the server's packages in
  `internal/` and its commands in `cmd/`. It takes no position on grouping by
  layer or by feature.
- **Alex Edwards, [_Let's Go_](https://lets-go.alexedwards.net/)**, a widely
  used book on Go web applications. It groups by layer: handlers and routes
  in one package, data access in `internal/models`, templates and static
  files in `ui/`.
- **Alex Edwards, [11 tips for structuring your Go projects](https://www.alexedwards.net/blog/11-tips-for-structuring-your-go-projects).**
  Start simple, create a package only when there is a real reason, and treat
  frequent import-cycle errors as a sign the structure is not working. Too
  many small packages is called one of the most common Go mistakes.
- **Ben Johnson, [Standard Package Layout](https://www.gobeyond.dev/standard-package-layout/).**
  Describes the problems with both layer and feature grouping (names that
  repeat the package, such as `polls.Poll`, and circular dependencies) and
  proposes organizing by dependency.
- **Kat Zien, [How Do You Structure Your Go Apps?](https://about.sourcegraph.com/go/gophercon-2018-how-do-you-structure-your-go-apps)**
  (GopherCon 2018). Compares flat, layer, feature, and domain-driven layouts,
  and advises the simplest structure that fits the app.
- **[go-blueprint](https://github.com/melkeydev/go-blueprint)** and Melkey's
  video [This Is The BEST Way To Structure Your GO Projects](https://www.youtube.com/watch?v=dxPakeBsgl4).
  Both group by layer: `internal/server` for routes, `internal/database` for
  data access, and one package for all templ components.
- **[golang-standards/project-layout](https://github.com/golang-standards/project-layout)**,
  a collection of community patterns. Its README states that it is not an
  official standard from the Go team.

### Options considered

- **Flat**: everything in one package. Simplest to start, but this app has
  enough distinct jobs (settings, database, rules, HTTP, templates) to
  benefit from a few packages.
- **Group by feature**: one package per part of the app, such as `polls/`
  and `voting/`, each holding its handlers, rules, and templates. This is how
  the Python version is organized. In Go it creates more packages and more
  boundaries between them, and features that share types or templates can
  import each other and form cycles.
- **Group by dependency** (Ben Johnson): domain types in the root package,
  one subpackage per outside system, wired together through interfaces.
  Suited to larger apps; more structure than this app needs.
- **Group by layer**: one package per job (database, rules, HTTP handlers,
  templates), with imports flowing in one direction. This is the layout
  shown in _Let's Go_ and go-blueprint.

## Decision

Code is grouped by layer.

### Identity

- **Product name:** Tiebreaker.
- **Repository:** `hbabb/tiebreaker-go` on GitHub.
- **Module path:** `github.com/hbabb/tiebreaker-go`, set in `go.mod`. Every
  import inside the project starts with it, for example
  `github.com/hbabb/tiebreaker-go/internal/config`.
- **Binary name:** `tiebreaker`. `go build` names a program after its
  folder, so `cmd/tiebreaker/` builds a binary called `tiebreaker`.

### Structure

Entries marked _git-ignored_ are created on a developer's machine and never
committed. The full list of ignored patterns, including editor and operating
system files, is in [`.gitignore`](../../.gitignore).

```
tiebreaker-go/
├── cmd/
│   └── tiebreaker/
│       └── main.go              ← loads settings, opens the database, builds the server, starts it
├── internal/
│   ├── config/                  ← reads settings from the environment (ADR-0017)
│   ├── database/                ← opens SQLite; sqlc-generated types and queries (ADR-0013)
│   │   └── queries/             ← hand-written .sql queries that sqlc reads
│   ├── rules/                   ← plain functions for the app's rules, such as 2 to 10 options
│   ├── server/                  ← chi router, middleware, and every HTTP handler (ADR-0002)
│   └── views/                   ← every templ component: layout, pages, fragments (ADR-0019)
├── static/
│   ├── embed.go                 ← embeds the folders below into the binary
│   ├── css/
│   ├── img/
│   └── js/                      ← htmx.min.js and JSDoc-typed scripts (ADR-0010)
├── migrations/                  ← goose .sql files, numbered with zero padding (ADR-0013)
├── docs/
│   └── adr/
├── bin/                         ← compiled binary from `go build -o bin/` (git-ignored)
├── data/                        ← SQLite database file and its WAL files (git-ignored)
├── tiebreaker                   ← compiled binary from a plain `go build` (git-ignored)
├── coverage.out                 ← test coverage report (git-ignored)
├── .gitignore
├── .golangci.yml                ← linter settings (ADR-0011)
├── CONTRIBUTING.md
├── go.mod                       ← module path, Go version, libraries
├── go.sum                       ← checksums for every library download
├── mise.toml                    ← Go and tool versions, tasks (ADR-0004)
└── sqlc.yaml                    ← sqlc settings (ADR-0013)
```

### Import direction

Imports flow one way, so no cycle can form:

```
cmd/tiebreaker → config, database, server
server         → database, rules, views, static
views          → database   (for the types it displays)
rules          → nothing in this module
database       → nothing in this module
```

Nothing imports `server` except `main`.

### Rules

- **`rules`** holds the app's rules as plain functions with no web or
  database code, such as checking that a poll has 2 to 10 options or that a
  vote arrives before the closing time. They take plain values and are
  tested by calling them directly.
- **`server`** holds routing and every handler. A handler reads the
  request, calls `rules` and `database`, and renders a component from
  `views`.
- **`views`** holds every templ component. The shared layout, the pages, and
  the fragments HTMX swaps in all live here. Components receive the data
  they display as parameters; they do not query the database.
- Files inside a package are split by part of the app when they grow, for
  example `server/polls.go` and `server/voting.go`, or `views/vote.templ`
  and `views/results.templ`. A new package is created only when there is a
  clear reason.
- Tests live next to the code they test, in `_test.go` files. There is no
  top-level `tests/` folder. Where integration tests that run the whole app
  live, if they are added, is a separate decision.
- `static/` is embedded into the compiled binary with the `embed` package.
  A `//go:embed` line can only include files in its own folder or below (it
  cannot use `..`), so the embedding is done by `static/embed.go`, a small Go
  file inside `static/` itself.
- `migrations/` is embedded the same way, so the binary can run its own
  migrations at startup (ADR-0013).
- Code generated by sqlc and templ is committed, so the app builds with only
  Go installed.
- Folders and files are created when they are needed, not up front.

## Consequences

- Six packages under `internal/`, each with one job.
- All markup is in `views/`, and all handlers are in `server/`.
- One part of the app, such as voting, is spread across `server/`,
  `views/`, `rules/`, and `database/`. Files are named for the part of the
  app they cover, so they can still be found by name.
- `server/` and `views/` grow as the app grows. Large packages are not a
  problem in Go by themselves; if one becomes hard to work with, splitting it
  is a new decision.
- `internal/` keeps every package private to this module.
- The compiled binary contains the components, static files, and
  migrations, so deploying the app means copying one file. The database
  stays outside the binary, in `data/`.
- Git does not track empty folders, so `data/` does not exist in a fresh
  clone. How it is created is not yet decided.
- The Python version groups code by feature. That layout does not carry
  over: in Go, every folder is a package boundary, and grouping by layer
  keeps the number of packages and the risk of import cycles low.
