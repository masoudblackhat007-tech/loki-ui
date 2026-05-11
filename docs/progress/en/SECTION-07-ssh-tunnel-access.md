# Section 07 - SSH tunnel access to loki-ui

## Date

2026-05-06

## Goal

Verify secure local access to the server-side `loki-ui` service through SSH local port forwarding without exposing port `18090` publicly.

## Access model

The `loki-ui` service runs on the server and listens only on:

    127.0.0.1:18090

The local machine accesses it through an SSH tunnel:

    local machine -> SSH tunnel -> server 127.0.0.1:18090 -> loki-ui -> Loki -> Laravel logs

No public firewall port is opened for `loki-ui`.

## Commands used

SSH tunnel command:

    ssh -L 18090:127.0.0.1:18090 laravel-server

Alternative explicit SSH command:

    ssh -i ~/.ssh/laravel_log2_loki_ed25519 -o IdentitiesOnly=yes -L 18090:127.0.0.1:18090 deploy@91.107.169.146

Local HTTP test:

    curl -I http://127.0.0.1:18090/logs

Local API test:

    curl -sS "http://127.0.0.1:18090/api/logs?range=24h&limit=5"

Server listening socket verification:

    sudo ss -lntp | grep ':18090' || true

Server firewall verification:

    sudo ufw status verbose

## Local access verification

The local machine successfully accessed the logs page through the SSH tunnel.

Result:

    HTTP/1.1 200 OK
    Content-Type: text/html; charset=utf-8

The local machine also accessed the API endpoint through the SSH tunnel.

Result:

    The API returned a JSON response containing Laravel log data from Loki.

The returned log entry included:

    log_type: http_request
    service: Laravel Log2 Loki
    level: INFO
    request_id: present
    labels.job: laravel
    labels.service_name: laravel-log2-loki
    labels.environment: production
    labels.host: 381239

Sensitive and correlation-capable values from the raw API response were intentionally not copied into this document.

## Server socket verification

The server showed `loki-ui` listening only on:

    127.0.0.1:18090

Security conclusion:

    loki-ui is bound to loopback only.
    loki-ui is not listening on 0.0.0.0.
    loki-ui is not directly reachable from the public internet.

## Firewall verification

UFW was active.

Default policy:

    deny incoming
    allow outgoing
    deny routed

Allowed inbound rules:

    22/tcp
    80/tcp from the allowed client IP

There was no firewall rule for:

    18090/tcp

Security conclusion:

    No public firewall access was added for loki-ui.
    Access is restricted to SSH local port forwarding.

## Verified path

This section verified the full access path:

    local machine
    -> SSH tunnel
    -> server-local loki-ui
    -> server-local Loki
    -> Laravel structured logs

## Important issue still present

The API response returned the log entry successfully, but some top-level DTO fields were still not populated correctly:

    route: empty
    method: empty
    status: 0
    duration_ms: 0

The correct values exist inside the nested Laravel log context:

    context.http.method
    context.http.path
    context.http.route
    context.http.status_code
    context.http.duration_ms

This confirms that the remaining problem is not deployment or tunneling. It is a Go-side extraction/mapping issue.

## Security notes

- The UI exposes log data and must remain internal-only.
- The service should not be exposed through a public port.
- The service should not be put behind a public reverse proxy until authentication, authorization, TLS, rate limiting, and audit logging are implemented.
- Raw API responses may contain sensitive or correlation-capable data such as session hashes, request metadata, headers, payloads, and raw log context.
- Public documentation and screenshots should sanitize sensitive fields.

## Result

SSH tunnel access to `loki-ui` works.

The local machine can access the UI and API without opening port `18090` publicly.

The service remains loopback-only on the server.

## Resume relevance

This section demonstrates secure internal access to an observability UI using SSH local port forwarding.

The work demonstrates:

- secure access to an internal-only service
- SSH tunneling for operational tooling
- verification of loopback-only binding
- firewall validation
- end-to-end access from local machine to Laravel logs through Loki
- disciplined avoidance of public exposure for sensitive log tooling

## Remaining work

- Commit and push this documentation from the local development machine.
- Fix Go DTO extraction for nested Laravel log context fields.
- Add systemd hardening improvements if needed.
- Add HTTP server timeouts.
- Add graceful shutdown.
- Add authentication before any shared or public access.