# Section 05 - Server clone, build, and Loki integration test

## Date

2026-05-06

## Goal

Clone the `loki-ui` repository on the production server, build it from source, configure it to use the local Loki instance, and verify that it can display Laravel logs shipped through Alloy.

## Server

```
Hostname: 381239
User: deploy
Project path: /home/deploy/apps/loki-ui
Repository: git@github.com:masoudblackhat007-tech/loki-ui.git
```

## Actions performed

* Connected to the server as the `deploy` user.
* Verified that Loki, Grafana, and Alloy were active and enabled.
* Verified Loki readiness on `127.0.0.1:3100`.
* Verified Alloy readiness on `127.0.0.1:12345`.
* Verified Grafana health on `127.0.0.1:3000`.
* Installed Go 1.26.2 on the server from the official Go tarball.
* Verified the Go tarball checksum before installation.
* Cloned the `loki-ui` repository from GitHub.
* Created a server-side `.env` file from `.env.example`.
* Built the Go application on the server.
* Verified that `.env` and `bin/loki-ui` were ignored by Git.
* Ran `loki-ui` manually on the server.
* Verified that it listened only on `127.0.0.1:18090`.
* Generated a Laravel HTTP request to create a fresh structured log entry.
* Verified that Laravel wrote a JSON log file for the current date.
* Verified that Alloy could read the Laravel log file.
* Verified that Loki had the expected labels.
* Restarted Alloy so it picked up the current Laravel log file.
* Queried Loki directly with a 24-hour range.
* Verified that `loki-ui` returned Laravel logs through `/api/logs`.

## Observability service health

The following services were active and enabled:

```
loki
grafana-server
alloy
```

Loki readiness:

```
ready
```

Alloy readiness:

```
Alloy is ready.
```

Grafana health:

```
database: ok
version: 13.0.1
```

## Go installation

The server did not initially have Go installed.

The official Go archive was downloaded:

```
go1.26.2.linux-amd64.tar.gz
```

The archive size was approximately:

```
64M
```

The SHA256 checksum was verified:

```
990e6b4bbba816dc3ee129eaeaf4b42f17c2800b88a2166c265ac1a200262282
```

Installed Go version:

```
go version go1.26.2 linux/amd64
```

Go was installed under:

```
/usr/local/go
```

## Repository clone

The repository was cloned into:

```
/home/deploy/apps/loki-ui
```

Git status after clone:

```
On branch main
Your branch is up to date with 'origin/main'.

nothing to commit, working tree clean
```

Latest remote commit at the time of clone:

```
d109b12 Document GitHub remote push
```

## Server environment file

The server `.env` file was created from `.env.example`.

Runtime values:

```
LOKI_URL=http://127.0.0.1:3100
LISTEN_ADDR=127.0.0.1:18090
UI_TIMEZONE=Asia/Tehran
```

The `.env` file permission was set to:

```
600
```

The `.env` file is ignored by Git.

## Server build

Build command:

```
go build -o bin/loki-ui ./cmd/loki-ui
```

Build artifact:

```
bin/loki-ui
```

Artifact size:

```
13M
```

File type:

```
ELF 64-bit LSB executable, x86-64, statically linked, with debug_info, not stripped
```

Git ignore verification:

```
.gitignore:2:.env       .env
.gitignore:11:/bin/     bin
.gitignore:11:/bin/     bin/loki-ui
```

The working tree remained clean after creating `.env` and building the binary.

## Laravel log generation

A request was sent to the Laravel application:

```
curl -I http://91.107.169.146
```

Laravel returned:

```
HTTP/1.1 200 OK
X-Request-Id: present
```

The response also included cookies. Cookie values were intentionally not documented because they are sensitive runtime data.

Laravel wrote a JSON log entry to:

```
/var/www/laravel-log2-loki/storage/logs/laravel-2026-05-06.log
```

The log entry contained:

```
message: http_request
log_type: http_request
request_id: present
service: Laravel Log2 Loki
environment: production
http.method: HEAD
http.path: /
http.status_code: 200
http.duration_ms: 3
level_name: INFO
```

Sensitive or correlation-capable values such as cookies and session hashes are intentionally not copied into this documentation.

## Alloy verification

The Alloy user can read the current Laravel log file:

```
alloy can read today's laravel log
```

Alloy groups included:

```
www-data
```

Alloy config path:

```
/etc/alloy/config.alloy
```

Relevant Alloy source path:

```
/var/www/laravel-log2-loki/storage/logs/laravel-*.log
```

Configured Loki labels:

```
job=laravel
service_name=laravel-log2-loki
environment=production
host=381239
```

Loki write endpoint:

```
http://127.0.0.1:3100/loki/api/v1/push
```

## Loki label verification

Loki labels included:

```
environment
filename
host
job
service_name
```

Job values included:

```
laravel
```

Service name values included:

```
laravel-log2-loki
```

## Loki query verification

A direct Loki query for `{job="laravel"}` initially returned an empty result for the default short range.

After generating a Laravel request and restarting Alloy, a 24-hour `query_range` returned a Laravel log stream.

Returned stream labels included:

```
environment=production
filename=/var/www/laravel-log2-loki/storage/logs/laravel-2026-05-06.log
host=381239
job=laravel
service_name=laravel-log2-loki
```

This proved that the pipeline was working:

```
Laravel log file -> Alloy -> Loki
```

## loki-ui runtime test

The application was started manually:

```
set -a
source .env
set +a
./bin/loki-ui
```

Runtime log:

```
loki-ui listening on 127.0.0.1:18090
```

The `/logs` route returned:

```
HTTP/1.1 200 OK
```

The `/api/logs?range=1h&limit=5` route returned:

```
logs: []
```

The empty result was caused by the query range being too short for the available ingested log data.

The `/api/logs?range=24h&limit=5` route returned one Laravel log entry.

This proved that `loki-ui` can query server-local Loki and return Laravel logs.

## Listening address verification

The server showed:

```
127.0.0.1:18090
```

Security conclusion:

```
loki-ui listened only on loopback.
loki-ui did not bind to 0.0.0.0.
No public firewall port was opened.
```

## Important issue found

The API response returned a log entry, but some top-level fields were not populated correctly:

```
route: empty
method: empty
status: 0
duration_ms: 0
```

The same values existed inside the nested context:

```
context.http.method
context.http.path
context.http.route
context.http.status_code
context.http.duration_ms
```

Conclusion:

```
The server integration works, but the Go DTO/view extraction logic needs to be updated to read nested Laravel log fields from context.http.
```

This should be handled in a later code improvement section.

## Security notes

* `loki-ui` was not exposed publicly.
* No UFW rule was added for port 18090.
* Access should remain through SSH tunnel until authentication, authorization, TLS, rate limiting, and audit logging are implemented.
* Cookie values and session hashes must not be copied into public documentation.
* The UI can expose sensitive log context and must remain internal-only.

## Result

Server clone and build succeeded.

The Laravel-to-Loki-to-loki-ui integration was verified successfully using a 24-hour query range.

The current state proves:

```
GitHub -> server clone -> Go build -> local Loki query -> Laravel log display
```

## Resume relevance

This section demonstrates a practical internal observability integration:

* deployed a Go UI beside a Laravel application
* built the Go binary from source on Ubuntu
* configured runtime through `.env`
* validated service dependencies
* verified Alloy file ingestion
* queried Loki directly
* verified that the Go UI returned Laravel logs
* preserved a localhost-only security boundary
* identified a real schema-mapping bug for future improvement

## Remaining work

* Commit and push this documentation from the local development machine.
* Stop the manually running `loki-ui` process.
* Create a systemd service for `loki-ui`.
* Add SSH tunnel usage documentation.
* Fix Go DTO extraction for nested `context.http` fields.
* Add HTTP server timeouts and graceful shutdown.
