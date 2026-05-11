# Section 17 — Responsive request log view

## Date

```text
2026-05-09
```

## Goal

This section improved the `/logs` page layout for smaller screens without changing Go logic, Loki queries, API contracts, DTOs, or the security model.

The goal was to keep the log viewer usable on tablets and mobile-width browsers while preserving the existing internal-only access model.

## Problem

The `/logs` page was mainly designed for desktop. On smaller widths, the filters became cramped, the header did not wrap well, the search input used a fixed width, the summary bar did not adapt cleanly, and the table could overflow in a way that made the page hard to use.

## Root cause

The CSS only had a limited responsive rule for hiding the sidebar. The main layout, filters, summary bar, and log rows were not fully adapted for tablet or mobile views.

## Files changed

The work was limited to the frontend template layer.

The main changed file was:

```text
templates/logs.tmpl
```

## Changes made

The page layout was adjusted so that the log view is easier to read on smaller screens.

The changes improved:

```text
filter wrapping
header layout
summary bar wrapping
mobile log row readability
table/container overflow behavior
small-screen spacing
```

## Scope boundaries

This section did not change:

```text
Go backend logic
Loki API client
LogQL query generation
Laravel logging
Alloy configuration
Loki configuration
systemd service
firewall rules
SSH tunnel model
```

## Validation

The page was checked visually after the CSS and template changes.

The goal was to confirm that the page remained usable on desktop while improving smaller-width layouts.

## Security notes

This was a UI-only improvement. It did not add authentication, authorization, TLS, rate limiting, audit logging, or public access.

The service still had to remain internal-only and reachable through the existing SSH tunnel model.

## Final result

The log view became more usable on smaller screens while preserving the previous behavior and security assumptions.

## Resume-safe statement

```text
Improved the responsive layout of an internal Go-based Loki log viewer, making request logs and filters more usable on smaller screens without changing backend logic, Loki queries, or the internal-only access model.
```

## Limitations

This section did not add new observability features. It only improved the frontend layout and readability.
