# Section 12 - Graceful shutdown for loki-ui

## Date

2026-05-06

## Goal

Add graceful shutdown support to `loki-ui`.

Before this section, the HTTP server started and stopped directly through `ListenAndServe`.

That meant systemd restarts or stops could terminate the process without giving the server an explicit shutdown path.

## Problem

The service is managed by systemd and may be restarted during deploys or maintenance.

Without graceful shutdown:

- active requests may be interrupted abruptly
- the listener is closed only by process termination
- there is no explicit shutdown timeout
- service lifecycle behavior is less predictable

For an internal observability UI, this is not catastrophic, but it is still weak operational behavior.

## Files changed

    cmd/loki-ui/main.go
    internal/httpserver/server.go

## Fix

`cmd/loki-ui/main.go` now creates a signal-aware context using:

    signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

The context is passed to:

    httpserver.Start(ctx, addr)

`internal/httpserver/server.go` now:

- starts `server.ListenAndServe()` in a goroutine
- waits for either server error or context cancellation
- handles `http.ErrServerClosed` as expected shutdown behavior
- calls `server.Shutdown()` with a bounded timeout
- logs shutdown request and completion

The shutdown timeout is:

    10 seconds

## Why this design

systemd sends SIGTERM when stopping or restarting a service.

The application now handles SIGTERM explicitly and asks the HTTP server to shut down cleanly instead of relying only on process termination.

This makes deploy and restart behavior more predictable.

## Local validation

Before validation, the local WSL Go environment had a broken `GOROOT` value injected from a Windows toolchain path:

    /mnt/c/Users/1SKY.IR/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.4.windows-386

This caused errors such as:

    go: no such tool "vet"
    go: no such tool "compile"

The local shell was corrected with:

    unset GOROOT
    export GOTOOLCHAIN=local
    hash -r

The corrected Go toolchain values were:

    GOROOT=/usr/local/go
    GOTOOLDIR=/usr/local/go/pkg/tool/linux_amd64
    GOTOOLCHAIN=local

Validation commands run on local WSL:

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

## Runtime verification

After deploy, verify startup:

    journalctl -u loki-ui -n 30 --no-pager

Expected startup log:

    loki-ui listening on 127.0.0.1:18090

Verify stop/restart behavior:

    sudo systemctl restart loki-ui
    journalctl -u loki-ui -n 30 --no-pager

Expected graceful shutdown logs around restart:

    loki-ui shutdown requested
    loki-ui shutdown completed

Network and firewall checks:

    ss -ltnp | grep 18090 || true
    sudo ufw status numbered

Expected security result:

    loki-ui must still listen only on 127.0.0.1:18090
    UFW must not expose 18090/tcp publicly

Functional verification through SSH tunnel:

    curl -I http://127.0.0.1:18090/logs

    curl -s 'http://127.0.0.1:18090/api/logs?range=24h&limit=1' \
      | jq '.logs[0] | {method, route, status, duration_ms, log_type}'

## Security notes

Graceful shutdown is not an access-control feature.

It does not make `loki-ui` safe for public exposure.

The existing security rule still applies:

- keep `loki-ui` bound to loopback
- do not open port 18090 in UFW
- use SSH tunnel access only
- add authentication, authorization, TLS, rate limiting, and audit logging before any broader exposure

## Result

`loki-ui` now has explicit graceful shutdown handling for SIGINT and SIGTERM.

The HTTP server shuts down with a bounded timeout instead of relying only on abrupt process termination.

## Resume relevance

This section demonstrates production-style service lifecycle handling for a Go-based internal observability tool by adding signal-aware graceful shutdown and validating build behavior.
