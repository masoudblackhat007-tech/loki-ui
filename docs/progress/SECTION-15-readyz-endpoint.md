# Section 15 - Add readiness endpoint for Loki dependency

## Date

2026-05-07

## Goal

Add a `/readyz` endpoint to `loki-ui` to check whether the service can reach its main dependency: Loki.

This is separate from `/healthz`.

## Endpoint roles

`/healthz` answers:

    Is the loki-ui HTTP process alive?

`/readyz` answers:

    Can loki-ui reach Loki?

These checks are intentionally separate.

A process can be alive while its dependency is unavailable.

## Security requirements

The readiness endpoint must not expose:

- `LOKI_URL`
- raw Loki errors
- raw Loki responses
- environment variables
- internal paths
- stack traces
- log data

The response must be intentionally small:

    ready

or:

    not ready

## Files changed

    internal/loki/client.go
    internal/httpserver/handler.go
    internal/httpserver/server.go

## Implementation

A lightweight Loki readiness method was added to the Loki client:

    Client.Ready(ctx) error

It calls Loki's lightweight readiness endpoint:

    /ready

The HTTP handler was added as:

    Handler.Readyz

Supported methods:

    GET
    HEAD

Unsupported methods return:

    405 Method Not Allowed
    Allow: GET, HEAD

If Loki is ready, `/readyz` returns:

    HTTP/1.1 200 OK
    Content-Type: text/plain; charset=utf-8

    ready

If Loki is not ready or unreachable, `/readyz` returns:

    HTTP/1.1 503 Service Unavailable
    Content-Type: text/plain; charset=utf-8

    not ready

The handler uses a short bounded timeout:

    2 seconds

## Local validation

Commands run on local WSL:

    gofmt -w internal/loki/client.go internal/httpserver/handler.go internal/httpserver/server.go
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

Because Loki was not running locally in WSL, `/readyz` correctly returned not ready.

Health endpoint:

    curl -i http://127.0.0.1:18091/healthz

Result:

    HTTP/1.1 200 OK

    ok

Readiness endpoint:

    curl -i http://127.0.0.1:18091/readyz

Result:

    HTTP/1.1 503 Service Unavailable

    not ready

HEAD readiness check:

    curl -I http://127.0.0.1:18091/readyz

Result:

    HTTP/1.1 503 Service Unavailable

Unsupported method:

    curl -i -X POST http://127.0.0.1:18091/readyz

Result:

    HTTP/1.1 405 Method Not Allowed
    Allow: GET, HEAD

## Server deploy commands

These commands must be run on the server after pushing the commit:

    cd /home/deploy/apps/loki-ui
    git pull --ff-only
    go test ./...
    go build -o bin/loki-ui ./cmd/loki-ui
    sudo systemctl restart loki-ui
    sudo systemctl status loki-ui --no-pager

## Server validation commands

Because Loki is running on the server at `127.0.0.1:3100`, `/readyz` should return ready on the server:

    curl -i http://127.0.0.1:18090/readyz
    curl -I http://127.0.0.1:18090/readyz
    curl -i -X POST http://127.0.0.1:18090/readyz

Network and firewall checks:

    ss -ltnp | grep 18090 || true
    sudo ufw status numbered

Expected security result:

    loki-ui must still listen only on 127.0.0.1:18090
    UFW must not expose 18090/tcp publicly

## Security notes

`/readyz` is a shallow dependency readiness check.

It does not return raw error details.

It does not expose the Loki URL.

It does not query Laravel logs.

It does not return Loki response bodies.

This endpoint is useful for operational readiness checks, but it does not make `loki-ui` safe for public exposure.

The service must remain loopback-only and accessed through SSH tunnel until authentication, authorization, TLS, rate limiting, and audit logging are added.

## Result

`loki-ui` now has a separate readiness endpoint for Loki dependency availability.

The project now has separate liveness and readiness checks:

    /healthz
    /readyz

## Resume relevance

This section demonstrates adding dependency-aware readiness checking to a Go internal observability service while avoiding metadata leakage and keeping liveness separate from readiness.
