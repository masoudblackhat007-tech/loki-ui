# Section 01 - Git initialization for loki-ui

## Date
2026-05-05

## Goal
Initialize Git for the Go-based `loki-ui` project in the correct project path and create a clean source-only baseline commit.

## Project path

```text
/mnt/d/project-learn/loki-ui
```

## Context

The project had two local copies:

```text
/mnt/d/project-learn/loki-ui
/home/masoud/projects/loki-ui
```

The correct source-of-truth path for this project was selected as:

```text
/mnt/d/project-learn/loki-ui
```

The copy under `/home/masoud/projects/loki-ui` was not used as the main working copy.

## Actions performed

- Confirmed the correct project path.
- Checked that `/mnt/d/project-learn/loki-ui` was not already a Git repository.
- Created `.gitignore`.
- Created `.env.example`.
- Initialized Git.
- Renamed the default branch to `main`.
- Verified ignored files with `git check-ignore`.
- Created the first baseline commit.
- Detected that the initial `.gitignore` pattern `loki-ui` incorrectly ignored `cmd/loki-ui/main.go`.
- Fixed the ignore rule from `loki-ui` to `/loki-ui`.
- Amended the initial commit so the repository starts from a clean baseline.

## Files tracked in Git

```text
.env.example
.gitignore
cmd/loki-ui/main.go
go.mod
internal/httpserver/handler.go
internal/httpserver/server.go
internal/loki/client.go
internal/loki/types.go
templates/layout.tmpl
templates/log_detail.tmpl
templates/logs.tmpl
```

## Files intentionally ignored

```text
.env
.idea/
bin/
loki-ui
loki-ui.new
```

## Security decisions

- The real `.env` file is not tracked in Git.
- Build artifacts are not tracked in Git.
- IDE metadata is not tracked in Git.
- Runtime logs are not tracked in Git.
- The repository stores source code and safe configuration examples only.
- Runtime configuration is documented through `.env.example`, not through the real `.env`.

## Important mistake found and fixed

The original ignore pattern was:

```gitignore
loki-ui
```

This pattern also ignored the Go entrypoint path:

```text
cmd/loki-ui/main.go
```

That would have created a broken Go repository because the main executable entrypoint would be missing from Git.

The fixed pattern is:

```gitignore
/loki-ui
```

This only ignores the root build artifact:

```text
/loki-ui
```

and does not ignore:

```text
cmd/loki-ui/main.go
```

## Verification commands

```bash
git check-ignore -v cmd/loki-ui/main.go || echo "cmd/loki-ui/main.go is NOT ignored"

git check-ignore -v .env .idea bin loki-ui loki-ui.new

git ls-files

git log --oneline -1
```

## Verification results

### Entrypoint is not ignored

```text
cmd/loki-ui/main.go is NOT ignored
```

### Sensitive and generated files are ignored

```text
.gitignore:2:.env       .env
.gitignore:7:.idea/     .idea
.gitignore:11:/bin/     bin
.gitignore:12:/loki-ui  loki-ui
.gitignore:13:/loki-ui.new      loki-ui.new
```

### Files tracked by Git

```text
.env.example
.gitignore
cmd/loki-ui/main.go
go.mod
internal/httpserver/handler.go
internal/httpserver/server.go
internal/loki/client.go
internal/loki/types.go
templates/layout.tmpl
templates/log_detail.tmpl
templates/logs.tmpl
```

### Initial commit

```text
4ac811d Initial loki-ui project
```

## Result

Git was initialized successfully in the correct project path.

The initial repository baseline now includes:

- Go entrypoint
- internal Go packages
- HTML templates
- `go.mod`
- `.gitignore`
- `.env.example`

The repository intentionally excludes:

- real environment files
- IDE metadata
- build artifacts
- runtime logs

## Resume relevance

This step created a clean and secure source-control baseline for the internal Go-based Loki UI project.

The work demonstrates:

- Git repository initialization from an existing local project
- secure `.gitignore` design
- prevention of accidental environment file leakage
- prevention of binary artifacts entering source control
- validation of ignore rules before continuing development
- correction of a real `.gitignore` mistake before pushing to a remote repository

## Remaining work

- Add a project README.
- Create or connect a GitHub repository.
- Push the `main` branch to GitHub.
- Add local build instructions.
- Add server deployment instructions.
- Add systemd service documentation.
- Test the Go UI against local Loki.