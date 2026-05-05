# Section 03 - Local build and runtime verification

## Date

2026-05-05

## Goal

Verify that the `loki-ui` project can be built from source and executed locally with environment-based configuration.

## Actions performed

- Installed a working Go toolchain in WSL.
- Verified that the original local Go tarballs were corrupted or incomplete.
- Downloaded a clean Go tarball from the official Go download endpoint.
- Verified the downloaded Go archive with `gzip -t`.
- Verified the SHA256 checksum.
- Installed Go under `/usr/local/go`.
- Added Go to the shell `PATH`.
- Built the `loki-ui` binary from source.
- Verified that the build artifact is ignored by Git.
- Created a local `.env` file from `.env.example`.
- Verified that `.env` is ignored by Git.
- Started the `loki-ui` binary locally.
- Tested the `/logs` route.
- Tested the `/api/logs` route.
- Verified that the app listens only on `127.0.0.1:18090`.
- Confirmed that local Loki was not running on `127.0.0.1:3100`.

## Go installation evidence

The existing local Go archives were not reliable.

One archive failed with:

    gzip: unexpected end of file

Another archive failed with:

    gzip: invalid compressed data--format violated

A clean Go archive was downloaded and verified.

Downloaded archive:

    go1.26.2.linux-amd64.tar.gz

Downloaded size:

    64M

SHA256:

    990e6b4bbba816dc3ee129eaeaf4b42f17c2800b88a2166c265ac1a200262282

Installed Go version:

    go version go1.26.2 linux/amd64

## Build evidence

Build command:

    go build -o bin/loki-ui ./cmd/loki-ui

Build artifact:

    bin/loki-ui

Artifact details:

    -rwxrwxrwx 1 masoud masoud 13M May  5 14:13 bin/loki-ui

File type:

    ELF 64-bit LSB executable, x86-64, dynamically linked, with debug_info, not stripped

## Git hygiene verification

The build artifact was ignored by Git:

    .gitignore:11:/bin/     bin
    .gitignore:11:/bin/     bin/loki-ui

The working tree remained clean after the build.

## Runtime configuration

The local `.env` file was created from `.env.example`.

Expected runtime values:

    LOKI_URL=http://127.0.0.1:3100
    LISTEN_ADDR=127.0.0.1:18090
    UI_TIMEZONE=Asia/Tehran

The `.env` file is intentionally ignored by Git.

## Local Loki availability test

Command:

    curl -sS http://127.0.0.1:3100/ready

Result:

    curl: (28) Failed to connect to 127.0.0.1 port 3100 after 134132 ms: Couldn't connect to server

Conclusion:

    Loki was not running locally in WSL. Therefore full log-query verification must be performed later on the server where Loki is already installed and receiving Laravel logs.

## Runtime route test

Command:

    curl -I http://127.0.0.1:18090/logs

Result:

    HTTP/1.1 502 Bad Gateway

This was expected because `loki-ui` tried to query local Loki at `127.0.0.1:3100`, but Loki was not running locally.

## API route test

Command:

    curl -sS "http://127.0.0.1:18090/api/logs?range=1h&limit=5"

Result:

    loki error: do request: Get "http://127.0.0.1:3100/loki/api/v1/query_range?direction=BACKWARD&end=1777979453958252194&limit=5&query=%7Bjob%3D%22laravel%22%7D+%7C%3D+%22log_type%5C%22%3A%5C%22http_request%5C%22%22&start=1777975853958252194": context deadline exceeded (Client.Timeout exceeded while awaiting headers)

Conclusion:

    The Go application started correctly and attempted to query Loki, but the configured Loki endpoint was unavailable in the local WSL environment.

## Listening address verification

Command:

    ss -lntp | grep ':18090' || true

Result:

    LISTEN 0      4096        127.0.0.1:18090      0.0.0.0:*    users:(("loki-ui",pid=4806,fd=3))

Security conclusion:

    The application listened only on 127.0.0.1:18090.
    It did not bind to 0.0.0.0.
    This matches the intended internal-only security model.

## Result

Local build verification passed.

Runtime startup verification passed.

Full Loki query verification could not be completed locally because Loki was not running on `127.0.0.1:3100`.

The next real integration test must be done on the server where:

    Loki listens on 127.0.0.1:3100
    Alloy ships Laravel logs into Loki
    Laravel produces structured JSON request logs

## Resume relevance

This step proves that the Go UI can be built from source, runs with environment-based configuration, and enforces localhost-only runtime exposure during local testing.

It also documents a realistic dependency failure: the application was operational, but its external dependency, Loki, was unavailable in the local WSL environment.

## Remaining work

- Commit this verification document.
- Push the repository to GitHub.
- Clone or pull the repository on the server.
- Build the app on the server.
- Test the app against server-local Loki.
- Create a systemd service for `loki-ui`.
- Verify SSH tunnel access from the local machine.
