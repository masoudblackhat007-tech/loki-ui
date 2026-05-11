# Section 19 — Per-project HTTP status trend charts

## Goal

This section updated the `Requests` page in `loki-ui` so each Laravel project card displays its own HTTP status trend chart.

Before this section, the `Requests` page already grouped logs by project and showed per-project counts for `2xx`, `4xx`, and `5xx`. However, it did not show when those statuses happened over time.

The goal was to let the user quickly see, per project, when successful requests, client-side errors, and server-side errors occurred before reading the detailed request table.

## Main file changed

```text
templates/logs.tmpl
```

## Scope

This section was frontend-only.

The following were not changed:

```text
internal/httpserver/handler.go
internal/httpserver/server.go
internal/loki/client.go
internal/loki/types.go
Alloy configuration
Loki configuration
Laravel logging code
systemd service
server firewall configuration
SSH tunnel model
```

There were no backend, API, LogQL, Alloy, Loki, Laravel logging, or security model changes.

The access model remained:

```text
browser -> SSH tunnel -> 127.0.0.1:18090 on server -> loki-ui -> Loki
```

Port `18090` still must not be opened publicly.

## Starting point

At the end of Section 18, the `Requests` page already had:

```text
project-based grouping
per-project local filters
click-to-filter
clearable inputs
focus-preserving re-render
server-side Loki filters
auto-refresh
responsive layout
```

Each project card had counters for:

```text
2xx
4xx
5xx
```

But only final counts were visible. The UI did not show whether those requests were concentrated in a specific minute or spread across time.

## UX decision

A separate chart was added inside each project card.

The final order inside each project card became:

```text
project header
HTTP status trend chart
project local filters
request log table
```

The chart is shown above the local filters and table so the user can understand the status distribution before drilling into request rows.

## Chart type

The selected chart type was:

```text
stacked bar chart
```

Each bar represents one one-minute time bucket.

Each bar is split into three status groups:

```text
2xx
4xx
5xx
```

Meaning of colors:

```text
Green  -> 2xx
Orange -> 4xx
Red    -> 5xx
```

The total height of a bar represents the total number of visible requests in that minute.

The stacked segments show how many of those requests were `2xx`, `4xx`, or `5xx`.

## Data source

The chart is computed per project from the same rows that are visible inside that project card.

This means the chart responds to local project filters.

If the user applies `Search this project`, `Method`, or `Status`, the chart is recalculated from the currently visible rows.

The chart does not send a new Loki query.

It only uses data already fetched from:

```text
/api/logs
```

## Time bucket logic

Each log timestamp is read and grouped into its one-minute bucket.

For each bucket, these counts are calculated:

```text
2xx count
4xx count
5xx count
total count
```

The chart displays at most the last 12 buckets.

That means it shows up to 12 recent minute-level bars for the current visible dataset.

## Frontend functions added

The chart logic added frontend JavaScript functions inside `templates/logs.tmpl`.

The bucket-building function:

```text
buildStatusChartBuckets
```

Its job:

```text
read log timestamp
group logs by minute
count 2xx, 4xx, and 5xx per minute
return the last 12 buckets
```

The rendering function:

```text
renderStatusChart
```

Its job:

```text
build the SVG chart
draw grid lines
draw stacked bars
draw axis labels
draw legend
handle empty chart state
```

## UI output

After deployment, each project card shows a chart titled:

```text
HTTP status trend
```

The chart description says:

```text
Stacked per-minute request counts for the currently visible rows.
```

The legend displays:

```text
2xx
4xx
5xx
```

## Example interpretation

If a project card shows:

```text
Showing 84 of 84 logs
2xx: 42
4xx: 42
5xx: 0
```

It means that among the currently fetched and visible rows for that project:

```text
42 requests were successful
42 requests had client-side errors
0 requests had server-side errors
```

If there is no red segment in the chart, then the current visible rows contain no `5xx` responses.

If orange segments are high, the `4xx` rate is high. If those requests were intentional test paths, the result is expected. If they were real traffic, they may indicate missing routes, broken links, wrong client requests, bot scans, probing, or invalid endpoints.

## Important limitation

The chart does not show the real total count of all logs stored in Loki.

It only shows the logs returned by the current query and visible in the UI.

It is affected by server-side filters:

```text
Project label
Lookback range
Max rows fetched
Request ID
Loki text contains
```

It is also affected by local project filters:

```text
Search this project
Method
Status
```

If `Max rows fetched` is too low, some logs may never reach the UI and therefore will not appear in the chart.

## Failed first attempt

The first implementation attempt rewrote `templates/logs.tmpl` incorrectly.

The real project expected a template named:

```text
logs
```

The broken version removed the correct `{{ define "logs" }}` structure and used incompatible template definitions.

The runtime error was:

```text
render error: html/template: "logs" is undefined
```

This was not a Loki, Alloy, backend, or Laravel logging issue.

The issue was limited to the `Requests` page template structure.

The failed commit was:

```text
b9acec1 Add per-project request status charts
```

This commit must not be documented as a successful chart implementation.

## Recovery and final fix

The broken template was recovered from the previous healthy version.

Then `templates/logs.tmpl` was rewritten again while preserving the real template structure:

```gotemplate
{{ define "logs" }}
...
{{ end }}
```

The final implementation used pure SVG and did not add external dependencies, assets, packages, or build pipeline changes.

## Local validation notes

Before the final commit, Git showed only this file as changed:

```text
templates/logs.tmpl
```

A local run without `LOKI_URL` failed with an environment error:

```text
panic: LOKI_URL is required
```

That was not a template error.

A later local run could not bind the default port because it was already in use:

```text
listen tcp 127.0.0.1:18090: bind: address already in use
```

Because of local runtime limitations, the final UI validation was performed on the server.

## Successful commit

The final successful commit was:

```text
4f17fb0 Add per-project request status charts
```

This commit was pushed to GitHub.

## Server deployment

Deployment was performed under:

```text
/home/deploy/apps/loki-ui
```

The server moved from the failed commit to the successful commit:

```text
b9acec1 -> 4f17fb0
```

The server build completed successfully, the service was restarted, and systemd showed:

```text
Active: active (running)
```

## Final result

The `Requests` page now displays a separate HTTP status trend chart for each Laravel project.

The chart shows:

```text
2xx
4xx
5xx
```

Each bar represents one minute.

The chart is recalculated from the currently visible rows in that project card.

Existing features were preserved:

```text
project grouping
project summary
local project filters
click-to-filter
clearable inputs
focus-preserving render
server-side Loki filters
auto-refresh
responsive layout
```

## Security notes

No security model change was made in this section.

The following remain true:

```text
loki-ui must stay internal-only
port 18090 must not be opened publicly
access must stay behind SSH tunnel
observability UI must not be exposed to the internet
```

This section must not be presented as authentication, authorization, TLS, rate limiting, audit logging, or hardening work.

It is a UX improvement for faster HTTP log interpretation.

## Limitations

The chart only works on fetched UI data.

It does not show the total real count of all Loki logs.

It does not perform anomaly detection.

It does not include AI-based analysis.

It groups statuses only into:

```text
2xx
4xx
5xx
```

## Resume-safe statement

```text
Added per-project HTTP status trend charts to an internal Loki-based Laravel request log viewer, using frontend-only SVG stacked bar charts for 2xx, 4xx, and 5xx request counts while preserving the existing internal-only SSH tunnel access model and avoiding backend, Loki, Alloy, or Laravel logging changes.
```

This claim is limited to this section.

It must not be presented as a backend change, security model change, authentication feature, or AI analysis feature.
