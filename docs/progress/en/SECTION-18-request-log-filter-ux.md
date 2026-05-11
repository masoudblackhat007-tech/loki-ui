# Section 18 — Request log filter UX improvements

## Goal

This section improved the user experience of the `Requests` page in `loki-ui` so HTTP request logs from multiple Laravel projects became easier to inspect, filter, and explain.

Before this section, the page fetched logs from Loki, but once multiple Laravel projects were producing logs at the same time, the flat list made it harder to identify which project each request belonged to and to quickly narrow down relevant rows.

The goal was to improve the frontend experience only, without changing the backend, Loki, Alloy, Laravel logging, or the security model.

## Main file changed

```text
templates/logs.tmpl
```

## Scope

The work was limited to the `Requests` page frontend.

The following were not changed:

```text
Go backend APIs
LogQL generation
Loki configuration
Alloy configuration
Laravel logging code
systemd service
firewall rules
SSH tunnel access model
```

## Starting point

Before this section, the `Requests` page displayed the fetched HTTP request logs as a mostly flat list.

That was acceptable for one Laravel project, but it became less readable after a second Laravel project was added to the same observability pipeline.

## Project-based grouping

The request logs were grouped by project.

Each project now has a separate card with its own summary and local filtering controls.

Project detection uses the service field and supports Loki labels such as:

```text
labels.service_name
```

## Per-project summaries

Each project card displays separate counts for:

```text
total
2xx
4xx
5xx
```

This makes it possible to quickly distinguish successful requests, client-side errors, and server-side errors for each Laravel project.

## Local project filters

Each project card has local filters that operate only on rows already fetched for that project.

The local filters include:

```text
Search this project
Method
Status
```

These filters do not send new queries to Loki. They only filter the already fetched rows in the browser.

## Server-side Loki filters

The top filters remained server-side filters.

They change the query sent to Loki and include fields such as:

```text
Project label
Lookback range
Max rows fetched
Request ID
Loki text contains
```

The UI text was clarified to make the difference between server-side Loki filters and local project filters explicit.

## Click-to-filter behavior

Clicking request fields now helps fill relevant filters.

The behavior includes:

```text
Path fills the local project search
Time fills the local project search
Verb fills the local method filter
Status fills the local status filter
Service fills the server-side Project label filter
```

The Service click behavior was intentionally changed to use the low-cardinality Loki label `service_name`, which is more appropriate for server-side filtering than copying display text into local search.

## Clearable inputs

Input fields gained clear buttons.

When an input has a value, an `×` button appears. Clicking it clears the input, keeps focus on the input, and triggers the correct UI update events.

This applies to both global server-side filters and local project filters.

## Focus preservation

A focus-loss issue was fixed for project-level filters.

Before the fix, typing a character could trigger a re-render and remove focus from the input.

The fix stores the active input and cursor position before rendering and restores them after rendering.

A short debounce was added to reduce unnecessary renders during typing.

## Global search cleanup

The old global search field was removed because it no longer had a clear role after per-project local filters were added.

This simplified the UI and avoided duplicated or confusing filtering behavior.

## Limit semantics

The UI clarifies that `Max rows fetched` is not the total number of logs in Loki.

It is only the maximum number of rows returned by the current Loki query.

If the value is too low, some logs may never reach the UI and therefore cannot be filtered locally.

The backend still limits invalid, zero, negative, or overly large values to safe defaults and maximums.

## Commits in this section

The section included several frontend-focused commits:

```text
2676526 Improve project-based request log view
701fcec Add click-to-filter request log fields
74e49dc Add clear buttons to request log filters
50a37d0 Improve request log filter interactions
e442aa5 Fix request log service filter behavior
659b883 Remove unused global request log filters
955c482 Clarify request log server filters
17f726f Document request log filter UX improvements
```

## Validation

After template changes, `go test ./...` completed successfully locally.

After UI commits, the changes were pushed to GitHub, pulled on the server, built, and the `loki-ui` systemd service was restarted.

For the documentation-only commit, only a server pull was needed; no build or restart was required.

## Security notes

This section did not change the security model.

The UI still had to remain internal-only behind SSH tunneling.

Port `18090` still must not be opened publicly.

No authentication, authorization, TLS, rate limiting, audit logging, or public reverse proxy was added in this section.

## Final result

The `Requests` page became project-aware, easier to filter, and easier to explain.

Existing functionality such as auto-refresh, responsive layout, click-to-filter, clearable inputs, and server-side Loki filtering was preserved or improved.

## Resume-safe statement

```text
Improved a Go-based internal Loki request log viewer with project-grouped request cards, per-project local filters, click-to-filter interactions, clearable inputs, and clearer server-side Loki filter semantics without changing backend APIs, Loki, Alloy, Laravel logging, or the internal SSH tunnel access model.
```

## Limitations

This section did not add charts, AI analysis, anomaly detection, authentication, authorization, TLS, rate limiting, or audit logging.
