# 10. Use JavaScript with JSDoc types

Date: 2026-10-16

## Status

Accepted

## Context

ADR-0003 limits client-side script to browser-only features HTMX does not
cover, such as copying the poll link, the native share sheet, and the winner
reveal. That script can be written as plain JavaScript, JavaScript with JSDoc
types, or TypeScript.

Options considered:

- **Plain JavaScript**: no build step, but no type checking.
- **TypeScript**: type checked, but browsers cannot run it directly. It needs
  a compile step, which brings Node, npm, and a `package.json` into the
  project.
- **JavaScript with JSDoc types**: type annotations are written in `/** */`
  comments. Browsers run the file as-is, and the TypeScript language server
  in the editor reads the comments and reports type errors.

This is the same decision as Python ADR-0010.

## Decision

Client-side script is written as JavaScript with JSDoc type annotations.

- Files live in `static/js/` and are served exactly as written.
- There is no build step, no `package.json`, and no `node_modules/` in the
  project.

## Consequences

- No compile step between editing a file and seeing it in the browser.
- Because `static/` is embedded into the Go binary (ADR-0005), a changed
  script appears after the app is rebuilt and restarted.
- Type errors are reported in the editor. The editor's TypeScript language
  server needs Node on the developer's machine; nothing Node-related is added
  to the repository.
- JavaScript types are not checked in CI. Doing so would require TypeScript
  on the CI runner.
- Whether type checking is turned on per file (`// @ts-check`) or for the
  whole project (`jsconfig.json` with `checkJs`) is not yet decided.
