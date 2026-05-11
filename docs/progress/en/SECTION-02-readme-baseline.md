# Section 02 - README and project baseline documentation

## Date

2026-05-05

## Goal

Add a baseline README for the `loki-ui` project so the project purpose, architecture, runtime model, security assumptions, limitations, and verification commands are documented from the beginning.

## Actions performed

- Added `README.md`.
- Documented the purpose of `loki-ui`.
- Documented its relationship with the `laravel-log2-loki` Laravel project.
- Documented the expected runtime environment variables.
- Documented the intended localhost-only security model.
- Documented current routes.
- Documented current features.
- Documented current limitations.
- Added local build and run instructions.
- Added verification commands.
- Added Git hygiene notes.
- Added a project roadmap.

## Security decisions documented

- `loki-ui` is internal-only at this stage.
- `loki-ui` must listen on `127.0.0.1:18090`.
- Loki must remain on `127.0.0.1:3100`.
- Grafana must remain on `127.0.0.1:3000`.
- Access should happen through SSH local port forwarding.
- No public firewall port should be opened for `loki-ui`.
- The project is not safe to expose publicly until authentication, authorization, TLS, rate limiting, and audit logging are added.

## Important limitations documented

The README explicitly states that the current project does not yet include:

    authentication
    authorization
    TLS
    CSRF protection
    rate limiting
    audit logging
    user/session management
    public network hardening
    graceful shutdown
    explicit HTTP server read/write timeouts

## Files added

    README.md
    docs/progress/SECTION-02-readme-baseline.md

## Verification commands

    git status --short
    git log --oneline -4

## Result

The project now has a baseline README that explains what the project is, how it fits into the Laravel/Loki stack, how it should be run, and what security boundaries must not be violated.

## Resume relevance

This step documents the project as an internal observability tool instead of leaving it as unexplained source code.

The work demonstrates:

- technical documentation
- security-aware project description
- architecture explanation
- explicit limitation tracking
- preparation for repeatable deployment
- documentation of operational assumptions

## Remaining work

- Push repository to GitHub.
- Add server deployment documentation.
- Add systemd service.
- Test the app against the real local Loki instance.
- Start hardening the Go HTTP server.
