# Section 09 - HTTP server timeouts

## Date

2026-05-06

## Goal

Add explicit HTTP server timeouts to `loki-ui`.

Before this section, the server used the default Go HTTP server behavior through `http.ListenAndServe`.

That default is not strict enough for an operational tool, even when the service is loopback-only.

## Problem

The previous server startup used:

    http.ListenAndServe(addr, mux)

This does not configure explicit server-level timeouts.

Without explicit timeouts, slow or broken clients can keep connections open longer than necessary.

For an internal observability UI, this is still a bad default because:

- the UI exposes operational log data
- the service is small and should fail predictably
- SSH tunnel clients or local clients can still behave badly
- internal-only does not mean safe-by-default

## File changed

    internal/httpserver/server.go

## Fix

Replaced raw `http.ListenAndServe` usage with an explicit `http.Server`.

Configured timeouts:

    ReadHeaderTimeout: 5 seconds
    ReadTimeout:       15 seconds
    WriteTimeout:      30 seconds
    IdleTimeout:       60 seconds

## Why these values

`ReadHeaderTimeout` limits how long a client can take to send request headers.

`ReadTimeout` limits total request read time.

`WriteTimeout` limits response write time.

`IdleTimeout` limits how long idle keep-alive connections can remain open.

These values are conservative enough for a local/internal log UI and stricter than the default behavior.

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

## Runtime verification through SSH tunnel

This command must be run from the local machine while the SSH tunnel is active:

    curl -I http://127.0.0.1:18090/logs

Expected result:

    HTTP/1.1 200 OK

## Security notes

- This change does not add authentication.
- This change does not make `loki-ui` safe for public exposure.
- Port `18090` must remain closed in UFW.
- Access must continue through SSH local port forwarding.
- Timeouts reduce one class of resource-exhaustion risk but do not replace authentication, authorization, TLS, rate limiting, or audit logging.

## Result

`loki-ui` now uses an explicit Go `http.Server` with configured read, write, header, and idle timeouts.

The service remains internal-only.

## Resume relevance

This section demonstrates hardening a Go internal observability service by replacing default HTTP server behavior with explicit timeout controls and validating the change through build, test, and deployment checks.
