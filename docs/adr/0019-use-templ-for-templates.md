# 19. Use templ for templates

Date: 2026-10-17

## Status

Accepted

## Context

ADR-0003 renders HTML on the server and swaps fragments in with HTMX. The
Python version uses Jinja2 (Python ADR-0003).

Options considered:

### `html/template` (standard library)

Templates are text files parsed when the app starts. Data is passed in as
`any`, so nothing checks it against the template until the page is rendered.

- A misspelled field such as `{{.Questoin}}` parses without error. It fails
  only when the page is rendered, after part of the page has already been
  written.
- Passing the wrong kind of data can fail silently. Rendering the vote page
  with a map instead of a `Poll` produced `<h1>Vote on </h1>` and no error.
- Template names are shared across everything parsed together. Two pages
  that each define `"content"` in one set collide, and the last one parsed
  wins. Each page needs its own template set, built with `Clone()`, and code
  to keep track of them.
- Output is escaped automatically, matching where each value appears (HTML
  text, URLs, inside `<script>`).
- No extra tool, no generate step, and covered by the Go 1 compatibility
  promise.
- GoLand supports Go templates, including field completion through a
  `gotype` comment. The compiler still does not check them.

### templ (`github.com/a-h/templ`)

Each template is written in a `.templ` file and turned into a Go function
with typed parameters by `templ generate`. The Go compiler checks the
generated code.

- A misspelled field or the wrong type is a compile error.
- A layout is a component that wraps its content with `{ children... }`, so
  pages cannot collide.
- A component is a Go function, so a fragment for an HTMX request is the
  same component called on its own, such as `Results(poll)`.
- Output is escaped automatically. It is stricter than `html/template`:
  variables are not allowed inside `<script>` or `<style>`, `href` values go
  through `templ.URL`, and `style` attributes cannot use expressions.
- It is pre-1.0 (v0.3.x, with 19 releases in the v0.3 line). Its syntax
  changed once before, from v0.1 to v0.2, with a migration command.
- Its JetBrains plugin is listed under an unverified vendor account. Its
  documentation lists official support for Neovim through its language
  server.

## Decision

We will use templ, for these reasons:

- **Mistakes are caught when the code compiles, not when a page renders.**
  This includes the silent case, where the wrong data renders a blank value
  with no error.
- **Layouts work like components in React and Vue.** A layout wraps the
  page's content with `{ children... }`, and pages are separate functions
  that cannot collide.
- **The syntax is easier to read.** Templates use Go's own `if`, `for`, and
  `switch` statements and plain HTML, rather than a separate template
  language.
- **Components are functions, which fits HTMX.** The same component renders
  inside a full page or on its own as a fragment, with no separate template
  set to manage.

How it is used:

- The templ CLI is pinned in `mise.toml` from its GitHub releases
  (`github:a-h/templ`), since templ is not in the mise registry (ADR-0004).
- Generated `_templ.go` files are committed, so the app builds with only Go
  installed. CI runs `templ generate` and fails if any file changes
  (`git diff --exit-code`), which means a generated file was out of date.
- The templ CLI and the `github.com/a-h/templ` library are kept at the same
  version.
- During development, `templ generate --watch` regenerates code on save. How
  it is combined with the development server is part of the mise `dev` task.
- Data for client-side scripts (ADR-0010), such as the poll link for the
  share button, is passed through `data-` attributes or
  `templ.JSONScript`, since templ does not allow variables inside
  `<script>`.

## Consequences

- Template mistakes stop the build instead of reaching the browser.
- Every template change needs `templ generate` before it compiles. The watch
  mode does this automatically during development.
- New dependency: `github.com/a-h/templ`. The generated code imports it and
  its `runtime` package, so it is part of the compiled binary, not only a
  development tool.
- templ is pre-1.0, so an upgrade can include breaking changes. Upgrades are
  done on purpose, one version at a time, with `templ generate` and the full
  test suite run after each.
- Components are plain Go functions, so tests render them directly into a
  buffer and check the HTML (ADR-0014).
- Where `.templ` files live in the project is recorded in ADR-0005.
- The `_` prefix for fragments from the Python version (ADR-0003) does not
  carry over. A fragment is an exported component, such as `Results`.
