# Section 10 - systemd hardening for loki-ui

## Date

2026-05-06

## Goal

Harden the `loki-ui` systemd service without exposing the service publicly and without breaking the SSH-tunneled operational workflow.

The service must remain internal-only:

    127.0.0.1:18090

Port `18090` must not be opened in UFW.

## Starting point

Before this section, the service already had basic hardening:

    User=deploy
    Group=deploy
    Restart=on-failure
    RestartSec=3
    NoNewPrivileges=true
    PrivateTmp=true
    ProtectSystem=strict
    ProtectHome=false
    ReadWritePaths=/home/deploy/apps/loki-ui

Runtime validation before additional hardening showed:

    loki-ui active/running
    loki-ui listening only on 127.0.0.1:18090
    UFW had no rule for 18090/tcp
    SSH-tunneled /logs returned HTTP 200
    SSH-tunneled /api/logs returned Laravel log DTOs correctly

## Problem

`systemd-analyze security loki-ui` showed that the service still had unnecessary access to several host resources.

Important gaps included:

    ProtectHome=no
    PrivateDevices not enabled
    ProtectKernelTunables not enabled
    ProtectKernelModules not enabled
    ProtectKernelLogs not enabled
    ProtectControlGroups not enabled
    CapabilityBoundingSet not restricted
    AmbientCapabilities not restricted
    RestrictAddressFamilies not restricted
    IPAddressDeny / IPAddressAllow not configured

For an internal observability UI, this is still too permissive.

Internal-only does not mean safe-by-default.

## File changed on server

The systemd service file was edited on the server:

    /etc/systemd/system/loki-ui.service

A backup was created before editing:

    /etc/systemd/system/loki-ui.service.bak-section10

Backup command:

    sudo cp /etc/systemd/system/loki-ui.service /etc/systemd/system/loki-ui.service.bak-section10

## Final service hardening configuration

The service was updated to include the following hardening directives:

    NoNewPrivileges=true
    PrivateTmp=true
    PrivateDevices=true

    ProtectSystem=strict
    ProtectHome=read-only
    ProtectControlGroups=true
    ProtectKernelTunables=true
    ProtectKernelModules=true
    ProtectKernelLogs=true

    CapabilityBoundingSet=
    AmbientCapabilities=

    RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX
    IPAddressDeny=any
    IPAddressAllow=localhost

    RestrictNamespaces=true
    RestrictSUIDSGID=true
    RestrictRealtime=true
    LockPersonality=true
    SystemCallArchitectures=native
    RemoveIPC=true
    UMask=0077

## Validation commands

Syntax validation:

    sudo systemd-analyze verify /etc/systemd/system/loki-ui.service

Note:

    systemd-analyze verify printed a warning about snapd.service RestartMode.
    The warning was not from loki-ui.service.

Reload and restart:

    sudo systemctl daemon-reload
    sudo systemctl restart loki-ui
    sudo systemctl status loki-ui --no-pager

Hardening inspection:

    systemctl show loki-ui \
      -p User \
      -p Group \
      -p NoNewPrivileges \
      -p PrivateTmp \
      -p PrivateDevices \
      -p ProtectSystem \
      -p ProtectHome \
      -p ProtectControlGroups \
      -p ProtectKernelTunables \
      -p ProtectKernelModules \
      -p ProtectKernelLogs \
      -p CapabilityBoundingSet \
      -p AmbientCapabilities \
      -p RestrictAddressFamilies \
      -p IPAddressDeny \
      -p IPAddressAllow \
      -p RestrictNamespaces \
      -p RestrictSUIDSGID \
      -p RestrictRealtime \
      -p LockPersonality \
      -p SystemCallArchitectures \
      -p RemoveIPC \
      -p UMask \
      -p Restart \
      -p RestartUSec

Security analysis:

    systemd-analyze security loki-ui | head -n 70

Network and firewall validation:

    ss -ltnp | grep 18090 || true
    sudo ufw status numbered

Functional validation through SSH tunnel:

    curl -I http://127.0.0.1:18090/logs

    curl -s 'http://127.0.0.1:18090/api/logs?range=24h&limit=1' \
      | jq '.logs[0] | {method, route, status, duration_ms, log_type}'

## Verified systemd properties after hardening

The final service properties included:

    Restart=on-failure
    RestartUSec=3s
    IPAddressAllow=127.0.0.0/8 ::1/128
    IPAddressDeny=0.0.0.0/0 ::/0
    UMask=0077
    CapabilityBoundingSet=
    AmbientCapabilities=
    User=deploy
    Group=deploy
    RemoveIPC=yes
    PrivateTmp=yes
    PrivateDevices=yes
    ProtectKernelTunables=yes
    ProtectKernelModules=yes
    ProtectKernelLogs=yes
    ProtectControlGroups=yes
    ProtectHome=read-only
    ProtectSystem=strict
    NoNewPrivileges=yes
    SystemCallArchitectures=native
    LockPersonality=yes
    RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX
    RestrictRealtime=yes
    RestrictSUIDSGID=yes
    RestrictNamespaces=yes

## Runtime result

After hardening, the service remained active:

    Active: active (running)

The service continued listening only on loopback:

    127.0.0.1:18090

UFW remained closed for port 18090:

    no 18090/tcp rule

The SSH-tunneled UI test passed:

    HTTP/1.1 200 OK

The SSH-tunneled API test passed and returned populated DTO fields:

    method: HEAD
    route: generated::3a9ow6Wh45d46nQj
    status: 200
    duration_ms: 3
    log_type: http_request

## Security conclusion

The service is more constrained after this section.

The hardening reduces unnecessary host access by restricting:

- Linux capabilities
- ambient capabilities
- device access
- writable home access
- kernel tunable/module/log access
- cgroup filesystem access
- namespace creation
- SUID/SGID behavior
- realtime scheduling
- non-native syscall architectures
- non-localhost network access

The service still is not safe for public exposure.

Missing controls still include:

- authentication
- authorization
- TLS
- application-level rate limiting
- audit logging
- CSRF protection if state-changing actions are ever added

## Result

`loki-ui` now runs under a stricter systemd sandbox while continuing to serve the internal UI and API through SSH tunnel access.

The service remains loopback-only and UFW still does not expose port 18090.

## Resume relevance

This section demonstrates practical Linux service hardening for a Go-based internal observability tool using systemd sandboxing, capability removal, network egress restrictions, and runtime validation.