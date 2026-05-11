````markdown
# Section 20 — Adding the internal bilingual documentation viewer

## Goal

In this section, an internal documentation viewer was added to the `loki-ui` project.

The goal was to make the project progress documentation readable from inside the UI instead of keeping it only as standalone Markdown files inside the repository.

This feature was added for resume-oriented documentation work, so the project sections can be read page by page in a clean, organized, and presentable interface.

## Scope of changes

The changes in this section were limited to the `loki-ui` project.

The following areas changed:

```text
internal/httpserver/handler.go
internal/httpserver/server.go
templates/docs.tmpl
docs/progress/fa/
docs/progress/en/
```

The following areas did not change:

```text
Loki query behavior
LogQL generation for request logs
Alloy configuration
Loki configuration
Laravel logging code
server firewall configuration
systemd security model
SSH tunnel access model
```

This section was not a change to the observability pipeline.

This section added an internal UI feature for reading project documentation.

## Access model

The `loki-ui` access model did not change.

The application still listens only on the server loopback address:

```text
127.0.0.1:18090
```

The access model remains:

```text
browser -> SSH tunnel -> 127.0.0.1:18090 on server -> loki-ui
```

Port `18090` must not be exposed publicly.

This section must not be presented as a public documentation portal.

It is only an internal documentation viewer inside the existing internal UI.

## State before this section

Before this section, the project documentation was stored under:

```text
docs/progress/
```

The documentation was split into separate section files from `SECTION-01` through `SECTION-19`.

The problem was that reading those files required opening them directly from the repository.

There was no in-UI reading experience for the documentation.

There was also no formal bilingual structure for Persian and English documentation.

## Design decision

A new route was added for reading documentation:

```text
/docs
```

This route was added inside the existing `loki-ui` application.

A page-by-page reading model was selected.

The user can move between documentation sections using `Back` and `Next`.

Language selection is handled through the following query parameter:

```text
lang
```

Persian route example:

```text
/docs?lang=fa&page=1
```

English route example:

```text
/docs?lang=en&page=1
```

The page number is handled through the following query parameter:

```text
page
```

## New documentation structure

To support two languages, the documentation was split into two separate directories:

```text
docs/progress/fa/
docs/progress/en/
```

The Persian version is stored under:

```text
docs/progress/fa/
```

The English version is stored under:

```text
docs/progress/en/
```

Both languages use matching file names.

Example:

```text
docs/progress/fa/SECTION-01-git-initialization.md
docs/progress/en/SECTION-01-git-initialization.md
```

This keeps pagination predictable across both languages.

## Internal UI backend changes

In `internal/httpserver/server.go`, the following route was registered:

```text
/docs
```

This route is connected to the documentation handler.

In `internal/httpserver/handler.go`, a new page data structure was added for the documentation view.

The structure provides the template with:

```text
page title
rendered documentation HTML
current page number
total page count
previous page number
next page number
current language
text direction
current file name
```

## Language selection behavior

The selected language is read from the following query parameter:

```text
lang
```

If `lang` is set to `en`, the English documentation is shown.

If `lang` is empty or invalid, the default language is Persian.

Persian is selected with:

```text
fa
```

English is selected with:

```text
en
```

## Text direction

For Persian, the text direction is set to right-to-left:

```text
rtl
```

For English, the text direction is set to left-to-right:

```text
ltr
```

This is not just a visual detail.

If Persian text is not rendered with `rtl`, mixing Persian text with technical English words, file paths, commit hashes, and commands can break readability.

Because of that, text direction was treated as a core part of this feature.

## Code block behavior

Even when the page language is Persian, code blocks and inline code must remain left-to-right.

For that reason, `templates/docs.tmpl` sets an independent direction for code elements.

This is required for readability of:

```text
file paths
commands
commit hashes
endpoints
service names
package names
LogQL
systemd unit names
```

## Markdown rendering

Documentation files are read as Markdown and rendered into simple HTML inside the UI.

No external Markdown parser was added in this section.

The current renderer supports:

```text
h1
h2
h3
paragraph
unordered list
inline code
fenced code block
```

Markdown text is escaped during rendering to avoid unsafe direct rendering.

The final output is then displayed in the template in a controlled way.

## Template changes

A new template file was created:

```text
templates/docs.tmpl
```

This template renders the documentation page.

Its main features are:

```text
centered reading card
section title
current language display
current page number
current file name
language switch buttons
back link to Requests
Back button
Next button
RTL support for Persian
LTR support for English
responsive layout for mobile
```

## Available routes

Main Persian documentation route:

```text
/docs?lang=fa&page=1
```

Main English documentation route:

```text
/docs?lang=en&page=1
```

If the user opens only:

```text
/docs
```

The default language is Persian.

## Pagination behavior

For each language, files are read from that language-specific directory.

Files are discovered using the following pattern:

```text
SECTION-*.md
```

The files are then sorted.

The page number is read from the following query parameter:

```text
page
```

If the page number is invalid, the application returns `400 Bad Request`.

If no documentation files are found for the selected language, the application returns `404 Not Found`.

## Persian fallback

To avoid breaking the Persian view during the migration period, if the new Persian directory is empty, the handler can still read the old files from the original progress directory.

This fallback was only added for Persian.

The purpose was to prevent the documentation page from completely failing while files were being moved from the old structure to the new bilingual structure.

After the bilingual structure was completed, the correct documentation paths are:

```text
docs/progress/fa/
docs/progress/en/
```

## Viewer commit

First, the internal documentation viewer was added.

The viewer commit was:

```text
dda7b46 Add internal documentation viewer
```

This commit added the `/docs` route, the documentation handler, and the initial documentation template.

## Bilingual content commit

After the viewer was added, the documentation structure was converted to bilingual content and Persian and English files were added.

The bilingual content commit was:

```text
222bcff Add bilingual documentation viewer content
```

This commit included:

```text
docs/progress/fa/
docs/progress/en/
```

It also updated the handler and template to support language selection and text direction.

## Local validation

After the changes, a local build was executed.

The build completed successfully.

This validation confirmed that the Go changes compiled and the templates parsed correctly.

## Server deployment

After commit and push, the changes were pulled on the server.

The server reached this commit:

```text
222bcff
```

A server-side build was then executed.

The `loki-ui` service was restarted.

The final service state was active.

The application still listened on the internal address:

```text
127.0.0.1:18090
```

## Final result

At the end of this section, the `loki-ui` project has an internal documentation viewer.

The viewer supports two languages:

```text
Persian
English
```

The Persian version is rendered right-to-left.

The English version is rendered left-to-right.

The user can read documentation page by page.

The user can switch between Persian and English.

Access remains internal and behind the SSH tunnel.

## Security notes

This section did not add authentication.

This section did not add authorization.

This section did not add TLS.

This section did not add rate limiting.

This section did not add new audit logging.

This section did not change the `loki-ui` security model.

Therefore, `loki-ui` must still not be exposed publicly.

Port `18090` must still not be opened in the public firewall.

This feature is acceptable only under the existing internal-only model.

## Current limitations

The current viewer uses a simple Markdown renderer.

This renderer is designed for controlled documentation stored inside the repository.

It is not designed for rendering anonymous user-supplied Markdown or public content.

The current version does not include documentation search.

The current version does not include a table of contents.

The current version does not include deep links for headings.

The current version only provides page-based navigation.

## Technical value

The value of this section was not just adding an HTML page.

The main value was converting project documentation from raw repository files into a readable in-UI experience.

It was also designed for two languages and two text directions from the beginning.

This matters because the project needs Persian documentation for learning and review, while also needing English documentation for resume-oriented presentation.

## Resume-safe statement

A defensible resume statement for this section is:

```text
Added an internal bilingual documentation viewer to a Go-based Loki UI, supporting page-by-page project documentation, Persian RTL rendering, English LTR rendering, and SSH-tunnel-only access without changing the observability pipeline or public exposure model.
```

This claim is limited to this section.

This section must not be presented as a public docs portal, authentication, authorization, TLS, rate limiting, audit logging, or a security model change.
````
