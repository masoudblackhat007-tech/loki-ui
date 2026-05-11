# Section 11 - Remove duplicate startup logging

## Date

2026-05-06

## Goal

Remove duplicate startup log messages from `loki-ui`.

After adding the explicit Go `http.Server` in Section 09, the service logged the startup message twice:

    loki-ui listening on 127.0.0.1:18090
    loki-ui listening on 127.0.0.1:18090

Only one process and one listening socket existed, so this was not a port-binding issue.

## Root cause

The startup message was logged in two places:

    cmd/loki-ui/main.go
    internal/httpserver/server.go

The log in `internal/httpserver/server.go` is the better location because that package creates and starts the HTTP server.

## File changed

    cmd/loki-ui/main.go

## Fix

Removed the duplicate startup log line from `cmd/loki-ui/main.go`.

The `start` function now delegates to `httpserver.Start(addr)` without logging the same startup message twice.

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
    journalctl -u loki-ui -n 20 --no-pager

## Expected runtime result

After restart, the journal should show only one startup line for the current service start:

    loki-ui listening on 127.0.0.1:18090

The service must remain:

    active/running
    bound only to 127.0.0.1:18090
    inaccessible through public UFW port 18090

## Security notes

This change is not a security control by itself.

It improves operational clarity. Duplicate startup logs can confuse incident analysis by making it look like two startup paths, two listeners, or two process starts occurred.

Clean logs matter because observability tooling should not create misleading signals.

## Result

Duplicate startup logging was removed.

Startup logging now has a single source of truth in the HTTP server startup path.

## Resume relevance

This section demonstrates operational cleanup by tracing and removing misleading duplicate service logs after a server startup refactor.
