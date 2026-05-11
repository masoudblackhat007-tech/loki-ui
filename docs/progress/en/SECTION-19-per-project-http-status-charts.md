# Section 19 — Adding per-project HTTP status charts

## Date

```text
2026-05-11
```

## Goal

In this section, the `Requests` page in `loki-ui` was improved to show a separate HTTP status chart for each Laravel project.

Before this section, the `Requests` page grouped logs by project and showed per-project counters for `2xx`, `4xx`, and `5xx`, but it did not show when those statuses occurred over time.

## Problem

After the second Laravel project was added to the log pipeline, numeric counters alone were not enough.

The UI could show how many successful or failed requests each project had, but it could not quickly show whether those requests were clustered in specific minutes or appeared as short spikes.

## Change made

A separate `HTTP status trend` chart was added to each project card.

The chart was placed above the local filters and above the request table for that project.

The chart is rendered as a `stacked bar chart` and shows `2xx`, `4xx`, and `5xx` counts per project.

Each bar represents a one-minute time bucket.

## Implementation

The change was limited to this file:

```text
templates/logs.tmpl
```

The chart uses inline SVG.

No external dependency was added.

The chart is built from the currently visible rows inside the same project card. If local filters change, the chart is recalculated from the filtered visible rows.

## Chart behavior

The chart does not send a new query to Loki.

It only uses data already fetched from `/api/logs`.

The chart displays at most the latest 12 time buckets.

Green represents `2xx`.

Orange represents `4xx`.

Red represents `5xx`.

## Failed first attempt

The first implementation attempt rewrote `templates/logs.tmpl` with the wrong template structure.

The required template named `logs` was removed, so the `Requests` page failed at runtime.

The runtime error was:

```text
html/template: "logs" is undefined
```

This was not a Loki, Alloy, backend, or Laravel logging problem.

The cause was only a broken HTML template definition.

## Recovery

For recovery, `templates/logs.tmpl` was restored on the server from the previous working version.

After build and restart, the service became active again.

The failed commit was:

```text
b9acec1 Add per-project request status charts
```

This commit must not be documented as the successful chart implementation.

## Final fix

After recovery, `templates/logs.tmpl` was rewritten again while preserving the required template structure.

The required template definition was kept:

```gotemplate
{{ define "logs" }}
```

The corrected version rendered the SVG chart successfully and the `Requests` page loaded without the previous runtime error.

The successful commit was:

```text
4f17fb0 Add per-project request status charts
```

## Validation

After the fix, the server build completed successfully.

The `loki-ui` service was restarted and systemd reported it as active.

The `Requests` page was opened through the SSH tunnel and the chart was visible for the Laravel project.

## Result

At the end of this section, every Laravel project on the `Requests` page has its own HTTP status trend chart.

The chart shows the timeline of `2xx`, `4xx`, and `5xx` requests.

Local filters, click-to-filter, clearable inputs, auto-refresh, and the responsive layout were preserved.

## Security note

This section did not change the security model.

Port `18090` must still not be public.

Access to `loki-ui` must still go through the SSH tunnel.

This section did not add authentication, authorization, TLS, rate limiting, or audit logging.

## Limitation

The chart does not show the full count of all logs stored in Loki.

It only shows data that the current query fetched and the UI currently displays.

This section did not add anomaly detection or AI analysis.

## Technical value

This section improved the readability of the `Requests` page for multiple Laravel projects and added a quick time-based view of HTTP status behavior without changing the backend, LogQL, or adding a dependency.

## Resume-safe statement

```text
Added frontend-only per-project HTTP status trend charts to an internal Loki-based Laravel request log viewer using SVG stacked bars for 2xx, 4xx, and 5xx counts while preserving the existing SSH-tunnel-only access model.
```
