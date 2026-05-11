# Section 08 - Fix nested HTTP context DTO extraction

## Date

2026-05-06

## Goal

Fix incorrect top-level DTO fields returned by `loki-ui` for Laravel HTTP request logs.

The API already returned Laravel logs from Loki, but these top-level DTO fields were wrong:

    route: empty
    method: empty
    status: 0
    duration_ms: 0

The correct values existed inside the nested Laravel log context:

    context.http.method
    context.http.path
    context.http.route
    context.http.status_code
    context.http.duration_ms

## Root cause

The Go handler extracted HTTP request fields from the root of `context`.

The code expected values like:

    context.method
    context.route
    context.status_code
    context.duration_ms

But the Laravel JSON log structure stores these fields under:

    context.http.method
    context.http.route
    context.http.path
    context.http.status_code
    context.http.duration_ms

Because of this mismatch, Go type assertions failed and the DTO kept zero values.

## File changed

    internal/httpserver/handler.go

## Fix

The handler now reads HTTP request fields from the nested `context.http` map while keeping fallback support for older flat context fields.

The fix applies to:

- `/logs` server-rendered view
- `/api/logs` JSON DTO output
- `/logs/detail` detail page DTO extraction
- request/response JSON shown on the detail page

Helper functions were added for safe extraction:

- `firstString`
- `firstNumber`
- `firstValue`
- `numberValue`

## Local validation

Commands run on local WSL:

    go test ./...
    go build -o bin/loki-ui ./cmd/loki-ui

Result:

    go test ./... passed
    go build completed successfully

Observed test output:

    ?       loki-ui/cmd/loki-ui             [no test files]
    ?       loki-ui/internal/httpserver     [no test files]
    ?       loki-ui/internal/loki           [no test files]

## Server deploy commands

These commands must be run on the server after pushing the commit:

    cd /home/deploy/apps/loki-ui
    git pull --ff-only
    go test ./...
    go build -o bin/loki-ui ./cmd/loki-ui
    sudo systemctl restart loki-ui
    sudo systemctl status loki-ui --no-pager

## Security checks

These commands must be run on the server:

    ss -ltnp | grep 18090 || true
    sudo ufw status numbered

Expected security result:

    loki-ui must still listen only on 127.0.0.1:18090
    UFW must not expose 18090/tcp publicly

## API verification through SSH tunnel

This command must be run from the local machine while the SSH tunnel is active:

    curl -s 'http://127.0.0.1:18090/api/logs?range=24h&limit=1' \
      | jq '.logs[0] | {method, route, status, duration_ms, request_id, log_type}'

Expected result:

    method is not empty
    route is not empty
    status is not 0
    duration_ms is not 0
    log_type is http_request

Example expected shape:

    method: HEAD
    route: generated::3a9ow6Wh45d46nQj
    status: 200
    duration_ms: 3
    log_type: http_request

## Security notes

- `loki-ui` remains internal-only.
- Port `18090` must not be opened in UFW.
- Access must continue through SSH tunnel only.
- Raw logs must not be copied into public documentation.
- Request IDs, session hashes, IPs, user agents, headers, payloads, and raw context fields can be correlation-capable and should be sanitized before documentation.

## Result

The DTO extraction bug was fixed by aligning Go-side mapping with the actual nested Laravel JSON log structure.

The end-to-end observability path remains:

    Laravel
    -> JSON log file
    -> Alloy
    -> Loki
    -> loki-ui
    -> SSH tunnel
    -> local browser

## Resume relevance

This section demonstrates debugging and fixing a real observability data-mapping bug across Laravel structured logs, Loki storage, and a Go-based internal log UI.
