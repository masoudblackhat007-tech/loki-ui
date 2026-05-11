# Section 13 - Fix WSL Go environment contamination

## Date

2026-05-07

## Goal

Fix the local WSL Go environment so `go test` and `go build` use the Linux Go toolchain consistently.

## Problem

During Section 12, local validation failed with:

    go: no such tool "vet"
    go: no such tool "compile"

This was not caused by the project code.

The local Go environment was contaminated by a Windows Go toolchain path injected into WSL.

## Root cause

`go env` showed that `GOROOT` pointed to a Windows 32-bit toolchain path:

    /mnt/c/Users/1SKY.IR/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.4.windows-386

This caused Go running inside WSL/Linux to look for Linux tools under a Windows toolchain directory.

The broken tool directory was:

    .../pkg/tool/linux_amd64

But the expected tools such as `compile` and `vet` did not exist there.

## Evidence

Broken behavior:

    go: no such tool "vet"
    go: no such tool "compile"

The environment contained IDE-injected variables:

    _INTELLIJ_FORCE_SET_GOROOT
    _INTELLIJ_FORCE_SET_GOPATH
    _INTELLIJ_FORCE_SET_GO111MODULE

The shell startup files did not originally contain the bad `GOROOT` value.

This indicated that the bad value was injected by the IDE environment, not by the project.

## Fix

The WSL shell configuration was updated to force the Linux Go toolchain.

For bash, `~/.bashrc` was updated with:

    unset GOROOT
    export GOTOOLCHAIN=local
    export GOPATH=/mnt/d/go
    export PATH="/usr/local/go/bin:$PATH"

For fish, `~/.config/fish/config.fish` was updated with:

    set -e GOROOT
    set -gx GOTOOLCHAIN local
    set -gx GOPATH /mnt/d/go
    fish_add_path --path /usr/local/go/bin

Duplicate fish config entries were cleaned up so the Go environment block exists only once.

## Correct final Go environment

The corrected WSL Go environment is:

    GOROOT=/usr/local/go
    GOPATH=/mnt/d/go
    GOTOOLDIR=/usr/local/go/pkg/tool/linux_amd64
    GOTOOLCHAIN=local

## Validation commands

Commands run from the project directory:

    cd /mnt/d/project-learn/loki-ui
    go env GOROOT GOPATH GOTOOLDIR GOTOOLCHAIN
    go test ./...
    go build -o bin/loki-ui ./cmd/loki-ui

## Validation result

`go env` returned:

    /usr/local/go
    /mnt/d/go
    /usr/local/go/pkg/tool/linux_amd64
    local

`go test ./...` passed:

    ?       loki-ui/cmd/loki-ui             [no test files]
    ?       loki-ui/internal/httpserver     [no test files]
    ?       loki-ui/internal/loki           [no test files]

`go build -o bin/loki-ui ./cmd/loki-ui` completed successfully.

## Security and reliability notes

A contaminated local toolchain makes build results untrustworthy.

Using a Windows Go toolchain path inside WSL can cause confusing failures and may hide real project problems behind environment errors.

For this project, WSL builds must use:

    /usr/local/go

The project should not rely on IDE-injected Windows Go paths.

## Result

The local WSL Go environment now consistently uses the Linux Go toolchain.

Local build and test validation are reliable again.

## Resume relevance

This section demonstrates practical local development environment debugging by tracing failed Go builds to an IDE-injected cross-platform toolchain mismatch and fixing WSL shell configuration for reproducible builds.
