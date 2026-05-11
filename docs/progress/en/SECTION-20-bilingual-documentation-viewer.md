# Section 20 — Adding the internal bilingual documentation viewer

## Date

```text
2026-05-11
```

## Goal

In this section, an internal documentation viewer was added to the `loki-ui` project.

The goal was to make the project documentation readable from inside the UI instead of keeping it only as Markdown files inside the repository.

This feature was added for resume-oriented documentation work, so the project sections can be read page by page in a clean and presentable interface.

## Problem

Before this section, the documentation existed under `docs/progress`, but reading it required opening files directly from the repository.

The project also needed both Persian and English documentation.

The Persian version had to render right-to-left, while the English version had to remain left-to-right.

Without correct text direction, Persian text mixed with commands, file paths, endpoints, and service names becomes hard to read.

## Change made

A new route was added for reading documentation:

```text
/docs
```

This route was added inside the existing `loki-ui` application.

The user can move between documentation sections with `Back` and `Next`.

The user can also switch between Persian and English.

## Documentation structure

The documentation was split into two directories:

```text
docs/progress/fa
docs/progress/en
```

The Persian version is stored under `fa`.

The English version is stored under `en`.

File names are kept consistent between both languages so pagination remains predictable.

## Implementation

In `internal/httpserver/server.go`, the `/docs` route was registered.

In `internal/httpserver/handler.go`, logic was added for reading Markdown files, selecting the language, selecting the page, and setting the text direction.

In `templates/docs.tmpl`, the documentation reading page was created.

## Language selection

Language is selected with this query parameter:

```text
lang
```

For Persian:

```text
/docs?lang=fa&page=1
```

For English:

```text
/docs?lang=en&page=1
```

If the language is missing or invalid, Persian is used as the default.

## Text direction

Persian uses `rtl`.

English uses `ltr`.

Code blocks and inline code always remain left-to-right so commands, file paths, and commit hashes stay readable.

## UI changes

The documentation page uses a centered reading card.

The header shows the current language, page number, and file name.

Language switch buttons are available in the header.

The link back to `Requests` was preserved.

The bottom of the page contains `Back` and `Next` buttons for moving between sections.

## Commits

First, the internal documentation viewer was added:

```text
dda7b46 Add internal documentation viewer
```

Then the bilingual documentation structure and content were added:

```text
222bcff Add bilingual documentation viewer content
```

Then section 20 was documented and the Persian section 19 document was cleaned up:

```text
105d2a3 Document bilingual internal docs viewer
```

## Validation

After the changes, the local build completed successfully.

The changes were pushed.

The server pulled the new commit.

The server build completed successfully.

The `loki-ui` service was restarted and remained active.

## Result

At the end of this section, `loki-ui` has an internal documentation viewer.

Documentation can be read in Persian and English.

The Persian version renders right-to-left.

The English version renders left-to-right.

The user can move page by page between documentation sections.

## Security note

This section did not change the security model.

Access must still remain behind the SSH tunnel.

Port `18090` must not be exposed publicly.

This section did not add authentication, authorization, TLS, rate limiting, or audit logging.

## Limitation

The current viewer uses a simple Markdown renderer.

The renderer is designed for controlled files stored inside the repository.

This section did not add search, a table of contents, or deep links for headings.

## Technical value

This section converted the project documentation from raw repository files into a readable in-UI experience.

It also supported two languages and two text directions from the beginning, which is important for Persian review material and English resume-oriented documentation.

## Resume-safe statement

```text
Added an internal bilingual documentation viewer to a Go-based Loki UI with page-based navigation, Persian RTL rendering, English LTR rendering, and SSH-tunnel-only access without changing the observability pipeline or public exposure model.
```
