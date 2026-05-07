# Section 14 - Add health endpoint for loki-ui

## Date

2026-05-07

## Goal

Add a lightweight health endpoint to `loki-ui`.

Before this section, basic service checks used:

    /logs

That endpoint is a real UI page and may depend on template rendering or future application logic.

A dedicated health endpoint gives a smaller and clearer smoke test target.

## Endpoint

The new endpoint is:

    /healthz

## Design

The endpoint is intentionally minimal.

It does not query Loki.

It does not expose:

- version information
- hostname
- environment variables
- file paths
- Loki URL
- application internals

It only confirms that the `loki-ui` HTTP process is alive and able to serve requests.

## File changed

    internal/httpserver/server.go

## Implementation

A new route was registered:

    mux.HandleFunc("/healthz", healthz)

The handler supports:

    GET
    HEAD

The handler rejects unsupported methods with:

    405 Method Not Allowed
    Allow: GET, HEAD

Successful GET response:

    HTTP/1.1 200 OK
    Content-Type: text/plain; charset=utf-8

    ok

Successful HEAD response:

    HTTP/1.1 200 OK
    Content-Type: text/plain; charset=utf-8

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

## Local runtime validation

The binary was run locally on a temporary loopback port:

    LISTEN_ADDR=127.0.0.1:18091 LOKI_URL=http://127.0.0.1:3100 ./bin/loki-ui

GET health check:

    curl -i http://127.0.0.1:18091/healthz

Result:

    HTTP/1.1 200 OK
    Content-Type: text/plain; charset=utf-8

    ok

HEAD health check:

    curl -I http://127.0.0.1:18091/healthz

Result:

    HTTP/1.1 200 OK
    Content-Type: text/plain; charset=utf-8

Unsupported method check:

    curl -i -X POST http://127.0.0.1:18091/healthz

Result:

    HTTP/1.1 405 Method Not Allowed
    Allow: GET, HEAD

Graceful shutdown was also observed after stopping the local process with Ctrl+C:

    loki-ui shutdown requested
    loki-ui shutdown completed

## Server deploy commands

These commands must be run on the server after pushing the commit:

    cd /home/deploy/apps/loki-ui
    git pull --ff-only
    go test ./...
    go build -o bin/loki-ui ./cmd/loki-ui
    sudo systemctl restart loki-ui
    sudo systemctl status loki-ui --no-pager

## Server validation commands

Health check on the server:

    curl -i http://127.0.0.1:18090/healthz
    curl -I http://127.0.0.1:18090/healthz
    curl -i -X POST http://127.0.0.1:18090/healthz

Network and firewall checks:

    ss -ltnp | grep 18090 || true
    sudo ufw status numbered

Expected security result:

    loki-ui must still listen only on 127.0.0.1:18090
    UFW must not expose 18090/tcp publicly

## Security notes

The health endpoint is intentionally shallow.

It verifies local HTTP process health only.

It does not verify Loki health, Laravel health, Alloy ingestion, or end-to-end log availability.

That is deliberate.

A deep dependency health check should be a separate endpoint or separate operational test, not mixed with a basic process liveness endpoint.

The endpoint must not expose sensitive runtime metadata.

## Result

`loki-ui` now has a lightweight `/healthz` endpoint suitable for smoke tests and future monitoring integration.

## Resume relevance

This section demonstrates adding a minimal operational health endpoint to a Go internal observability service while avoiding unnecessary dependency coupling or metadata exposure.
