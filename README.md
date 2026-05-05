# loki-ui

Internal Go-based web UI for viewing Laravel structured JSON logs stored in Loki.

This project is designed to run beside the `laravel-log2-loki` application and read logs from a local Loki instance. It provides a lightweight browser interface for inspecting HTTP request logs, exceptions, DB query logs, upstream calls, request IDs, response metadata, and raw JSON log context.

## Purpose

The goal of this project is to provide a small internal observability UI for a Laravel application that already writes structured JSON logs and ships them to Loki through Grafana Alloy.

This UI is not intended to replace Grafana. It is focused on fast request-level inspection and Laravel-specific log details.

## Target architecture

    Laravel application
      -> daily JSON log files
      -> Grafana Alloy
      -> Loki
      -> loki-ui
      -> SSH tunnel
      -> local browser

## Related project

This UI is intended to work with:

    laravel-log2-loki

That Laravel project provides:

- structured JSON logging
- request logging
- request IDs
- masked client IPs
- session hash logging instead of raw session IDs
- JSON API error responses
- Alloy/Loki/Grafana integration

## Current runtime model

The application reads the following environment variables:

    LOKI_URL=http://127.0.0.1:3100
    LISTEN_ADDR=127.0.0.1:18090
    UI_TIMEZONE=Asia/Tehran

`LOKI_URL` points to the local Loki HTTP endpoint.

`LISTEN_ADDR` must stay bound to localhost unless proper authentication, authorization, TLS, and network controls are added.

`UI_TIMEZONE` controls how timestamps are displayed in the UI.

## Security model

This project is currently designed as an internal-only tool.

The expected secure deployment model is:

    loki-ui listens on 127.0.0.1:18090
    Loki listens on 127.0.0.1:3100
    Grafana listens on 127.0.0.1:3000
    Access happens through SSH local port forwarding
    No public firewall port is opened for loki-ui

Example SSH tunnel:

    ssh -L 18090:127.0.0.1:18090 laravel-server

Then open locally:

    http://127.0.0.1:18090/logs

## Current routes

    /logs
    /logs/detail
    /api/logs
    /requests
    /api/requests

`/requests` and `/api/requests` are aliases for the logs page and logs API.

## Main features

- Query Laravel logs from Loki using `query_range`
- Display recent Laravel HTTP request logs
- Filter by service, level, text, request ID, range, and limit
- Return API results as JSON
- Show request detail by `request_id`
- Display request, response, context, auth, labels, raw JSON, DB queries, and upstream calls
- Use Go `html/template` for server-rendered templates
- Keep Loki labels low-cardinality
- Keep request-specific data inside JSON fields instead of Loki labels

## Important limitations

This project currently does not provide:

- authentication
- authorization
- TLS
- CSRF protection
- rate limiting
- audit logging
- user/session management
- public network hardening
- graceful shutdown
- explicit HTTP server read/write timeouts

Because of these limitations, it must not be exposed directly to the public internet.

## Local build

    go build -o bin/loki-ui ./cmd/loki-ui

## Local run

Create a local `.env` file based on `.env.example`.

    cp .env.example .env

Load environment variables:

    set -a
    source .env
    set +a

Run:

    ./bin/loki-ui

Open:

    http://127.0.0.1:18090/logs

## Verification commands

Check that the app responds locally:

    curl -I http://127.0.0.1:18090/logs

Check API output:

    curl -sS "http://127.0.0.1:18090/api/logs?range=1h&limit=5" | python3 -m json.tool

Check Loki directly:

    curl -G -sS "http://127.0.0.1:3100/loki/api/v1/query_range" \
      --data-urlencode 'query={job="laravel"}' \
      --data-urlencode 'limit=5' \
      | python3 -m json.tool

## Git hygiene

The repository intentionally excludes:

    .env
    .idea/
    bin/
    loki-ui
    loki-ui.new
    *.log

The repository tracks source code, templates, safe examples, and documentation only.

## Roadmap

- Add deployment documentation
- Add systemd service
- Add server-side HTTP timeouts
- Add graceful shutdown
- Add basic internal authentication
- Add security headers
- Add audit logging
- Review sensitive fields shown in the detail page
- Fix template issues
- Add tests for Loki query building
