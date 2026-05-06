# Section 06 - systemd service for loki-ui

## Date

2026-05-06

## Goal

Create a systemd service for `loki-ui` so the Go-based internal Loki UI runs as a managed service on the server instead of being started manually.

## Server

    Hostname: 381239
    User: deploy
    Project path: /home/deploy/apps/loki-ui
    Service file: /etc/systemd/system/loki-ui.service

## Actions performed

- Created a systemd service file for `loki-ui`.
- Configured the service to run as the `deploy` user.
- Configured the service working directory as `/home/deploy/apps/loki-ui`.
- Configured the service to load runtime variables from `/home/deploy/apps/loki-ui/.env`.
- Configured the service to execute `/home/deploy/apps/loki-ui/bin/loki-ui`.
- Enabled the service to start on boot.
- Started the service with systemd.
- Verified that the service is active.
- Verified that the service is enabled.
- Verified that `/logs` returns HTTP 200.
- Verified that `/api/logs?range=24h&limit=5` returns Laravel log data from Loki.
- Verified that `loki-ui` listens only on `127.0.0.1:18090`.
- Verified that no UFW rule exposes port `18090`.

## Service definition

The service file was created at:

    /etc/systemd/system/loki-ui.service

Service content:

    [Unit]
    Description=Internal Loki UI for Laravel logs
    After=network.target loki.service
    Wants=loki.service

    [Service]
    Type=simple
    User=deploy
    Group=deploy
    WorkingDirectory=/home/deploy/apps/loki-ui
    EnvironmentFile=/home/deploy/apps/loki-ui/.env
    ExecStart=/home/deploy/apps/loki-ui/bin/loki-ui
    Restart=on-failure
    RestartSec=3

    NoNewPrivileges=true
    PrivateTmp=true
    ProtectSystem=strict
    ProtectHome=false
    ReadWritePaths=/home/deploy/apps/loki-ui

    [Install]
    WantedBy=multi-user.target

## systemd verification

The service was enabled and started:

    sudo systemctl daemon-reload
    sudo systemctl enable --now loki-ui

Service state:

    active

Service enabled state:

    enabled

Systemd status showed:

    Active: active (running)
    Main PID: loki-ui
    ExecStart: /home/deploy/apps/loki-ui/bin/loki-ui

Journal output showed:

    loki-ui listening on 127.0.0.1:18090

## HTTP verification

The logs page was tested locally on the server:

    curl -I http://127.0.0.1:18090/logs

Result:

    HTTP/1.1 200 OK
    Content-Type: text/html; charset=utf-8

The API endpoint was tested:

    curl -sS "http://127.0.0.1:18090/api/logs?range=24h&limit=5" | python3 -m json.tool

Result:

    The API returned a JSON response containing Laravel log data from Loki.

This confirmed that the systemd-managed `loki-ui` process can query the server-local Loki instance.

## Listening address verification

The listening socket was checked with:

    sudo ss -lntp | grep ':18090' || true

Result:

    127.0.0.1:18090

Security conclusion:

    loki-ui listens only on loopback.
    loki-ui does not bind to 0.0.0.0.
    loki-ui is not directly exposed to the public internet.

## Firewall verification

UFW status showed:

    Status: active
    Default: deny incoming, allow outgoing, deny routed

Allowed inbound rules:

    22/tcp
    80/tcp from the allowed client IP

There was no rule for:

    18090/tcp

Security conclusion:

    No public firewall port was opened for loki-ui.
    Access must remain through SSH tunneling unless proper authentication, authorization, TLS, rate limiting, and audit logging are implemented.

## Runtime data verification

The API response returned a Laravel log entry with:

    log_type: http_request
    request_id: present
    service: Laravel Log2 Loki
    level: INFO
    labels.job: laravel
    labels.service_name: laravel-log2-loki
    labels.environment: production
    labels.host: 381239

This confirmed the working path:

    Laravel -> log file -> Alloy -> Loki -> loki-ui systemd service

## Important issues found

The API response returned data, but some top-level fields were not correctly populated:

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

Conclusion:

    The systemd deployment is working, but the Go DTO extraction logic needs to be updated to map nested Laravel log fields into the top-level API response fields.

The journal also showed the startup message twice:

    loki-ui listening on 127.0.0.1:18090
    loki-ui listening on 127.0.0.1:18090

Conclusion:

    The duplicate startup log should be reviewed later in the Go code or service behavior.

## Security notes

- The service runs as `deploy`, not root.
- The service reads configuration from `.env`.
- The `.env` file is not tracked by Git.
- The service is loopback-only.
- No UFW port was opened for 18090.
- The UI can expose sensitive log data and must remain internal-only.
- Cookie values and session hashes must not be copied into public documentation.

## Result

The `loki-ui` application now runs as a managed systemd service on the server.

The service is enabled, active, loopback-only, and able to query Laravel logs from Loki.

## Resume relevance

This section demonstrates operational deployment of a Go observability tool on Ubuntu using systemd.

The work demonstrates:

- Linux service management with systemd
- secure service user selection
- environment-based runtime configuration
- loopback-only service exposure
- firewall verification
- integration with Loki-backed Laravel logs
- operational validation using curl, journalctl, ss, and ufw
- identification of a real DTO mapping bug for future improvement

## Remaining work

- Commit and push this documentation from the local development machine.
- Add SSH tunnel usage documentation.
- Test browser access through SSH local port forwarding.
- Fix top-level API field extraction from nested context.http.
- Add HTTP server timeouts.
- Add graceful shutdown.
- Add authentication before any public or shared access.
