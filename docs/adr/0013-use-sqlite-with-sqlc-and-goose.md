# 13. Use SQLite with sqlc and goose

Date: 2026-10-16

## Status

Proposed

## Context

The app stores creator accounts, polls, options, suggestions, and votes. The
Python version uses SQLite through SQLAlchemy in development and left the
production database open (Python ADR-0013).

### Concurrency

The busiest moment is a group voting at once while every open vote page asks
for results on a timer (ADR-0015). That is many reads and a burst of small
writes.

SQLite in WAL (write-ahead log) mode lets readers and a writer run at the
same time: readers are not blocked by a write in progress, and each reader
sees a consistent snapshot. Only one write runs at a time; other writes wait
their turn. With `busy_timeout` set, a waiting write retries for up to that
long instead of failing at once. A vote is a single small insert, so writes
queue briefly rather than fail.

WAL mode requires every process using the database to be on the same
machine. It does not work over a network file system.

Options considered:

- **SQLite**: a file, no server. One writer at a time, unlimited concurrent
  readers in WAL mode, single machine.
- **PostgreSQL or MySQL**: a database server that handles many concurrent
  writers and many app servers on different machines. It is a second service
  to install, run, back up, and keep updated, locally and in production.

For SQLite from Go, two drivers are common: `mattn/go-sqlite3`, which needs
a C compiler (cgo), and `modernc.org/sqlite`, which is pure Go.

For database access, options considered:

- **An ORM** such as GORM: queries are built from Go method calls.
- **sqlc**: SQL is written by hand in `.sql` files, and sqlc generates
  type-safe Go functions and structs from it. A query that does not match the
  schema fails at generate time.

For migrations, sqlc does not run migrations itself, but it reads the schema
from migration files written for goose, golang-migrate, and several others,
using only the "up" part of each file.

## Decision

- **Database:** SQLite, in development and production.
- **Driver:** `modernc.org/sqlite` (pure Go, registered as `sqlite`), so no
  C compiler is needed to build.
- **Connection settings:** PRAGMAs are set in the connection string with the
  driver's `_pragma` parameter, so every pooled connection gets them:
  `journal_mode(WAL)`, `busy_timeout(5000)`, and `foreign_keys(1)` (SQLite
  does not enforce foreign keys unless this is on).
- **Queries:** sqlc, with `engine: "sqlite"` in `sqlc.yaml`. Generated code
  goes in `internal/database/`.
- **Migrations:** goose, with SQL migration files in `migrations/`. Files are
  numbered with zero padding (`00001_create_polls.sql`) because sqlc reads
  migration files in name order. The migrations are embedded into the binary
  and run by the app at startup.

## Consequences

- No database server to install or run, locally or in production. Backing up
  the data means copying one file safely.
- Votes are written one at a time. For a group voting on a poll, this queues
  writes briefly; it would only become a limit at a scale this app does not
  expect.
- The app must run as a single server. Running several app servers on
  different machines would require moving to PostgreSQL, recorded in a new
  ADR. sqlc supports PostgreSQL and MySQL, so queries could carry over with
  dialect changes.
- SQL is written by hand, so the queries are plain SQL a reader can run
  directly.
- `go build` works with `CGO_ENABLED=0`. Go's race detector (`go test
-race`) still needs cgo, which is unrelated to the database driver.
- `.gitignore` excludes `data/*.db*`, so database files are never committed.
