# apicactl — PR Review Context
Repo-specific context for reviewing PRs in ApicaSystem/apicactl. Read alongside the org-wide review rules.

## What this is
- A Go CLI for the Apica Ascent observability platform: stream/query/search logs, events and event rules, dashboards, alerts, licenses, resource listing, log pattern-signature extraction, and Grafana dashboard import.
- Go module `github.com/ApicaSystem/apicactl`. go.mod declares `go 1.12`, but that directive is stale — deps (grpc 1.64, protobuf 1.34) need modern Go; Dockerfile builds on golang:1.13.7.
- CLI framework: spf13/cobra (commands) + spf13/viper (config) + manifoldco/promptui (interactive prompts).
- Config/auth: stored in `~/.apicactl/apicactl.toml` (viper/TOML). Auth is either an API token (`config set-token`, viper key `uitoken`, sent as `Authorization: Key <token>` or `Bearer`) or UI user+password (`config set-credential`, stored base64-encoded, used to log in for a session cookie). Cluster endpoint, port, and default namespace also live in the config file. The gRPC path (grpc-gateway) forwards auth as a cookie in outgoing metadata.

## Layout
- `cmd/` — cobra command definitions: root, config, get, create, delete, logs, tail, license.
- `services/` — per-resource API call logic (applications, events, eventrules, query, namespace, processes, federation, license, realtime log stream).
- `api/` — protobuf: `proto/` sources, `v1/` generated Go stubs, `swagger/` generated JSON.
- `utils/` — HTTP client, config keys, protocol/URL resource routing, output formatting, net-trace.
- `grpc_utils/` — gRPC context and auth-cookie setup.
- `types/`, `defines/` — data models (dashboard, alert, query, grafana, datasource) and resource/endpoint enums.
- `ui/`, `converter/` — dashboard/visualization/widget builders and Grafana import conversion.
- `test/` — tests (jarcoal/httpmock + stretchr/testify), run via `run_test.sh`.
- `docs/` — generated cobra markdown command docs.

## Build / conventions
- Go modules. `make build` = `go build`. Entry point is `logiqctl.go` (package main) — NOT the `apicactl.go` the README's build steps reference.
- Generated code: `make generate` / `generate_grpc.sh` runs protoc + grpc-gateway + swagger to regenerate `api/v1/` and `api/swagger/` from `api/proto/`. Do not hand-edit generated output.
- `make revendor` runs `go mod vendor`, but `vendor/` is gitignored (not committed) — it will not appear in diffs.
- Cross-compilation via `goxbuild.sh` (gox); tests via `run_test.sh` (`go test ./test/...`).
- Historical naming: this was LogIQ's `logiqctl` before Apica. Expect leftover `logiq`/`flash` references and mixed copyright headers (Logiq.ai 2020 vs apica.io 2024) — cosmetic, not bugs.

## Review focus (priority order)
1. Go correctness — nil derefs; error handling (many paths `os.Exit(-1)` or print-and-return and swallow errors, so verify errors aren't silently dropped); unclosed `http.Response` bodies (`MakeApiCall` returns the response and `GetResponseString` reads but does not `Close` it); goroutine leaks.
2. CLI arg/flag correctness — manual arg-count checks (`len(args)`), flag parsing, and `--output` (table|json|yaml) / `--time-format` (relative|epoch|RFC3339) value handling.
3. Security — auth token and credential handling/storage (token stored plaintext, creds only base64-encoded); secrets leaking into output/logs (`config view` prints the token; `--nettrace` dumps request headers); TLS verification changes.
4. API/contract compat — command and subcommand/flag names consumed by users and scripts (config, get, create, delete, logs, tail, license). Renaming a command/flag or changing output-format shape is a breaking change.

## Do NOT flag (known-intentional)
- Generated code: `api/v1/**` and `api/swagger/**` (regenerated from `api/proto/` by `generate_grpc.sh`).
- `go.sum` (committed, large); `vendor/` (gitignored, absent from diffs).
- Pervasive `tls.Config{InsecureSkipVerify: true}` in `utils/client.go`, `utils/url.go`, `grpc_utils/grpc.go` — established baseline. Flag only a change that alters TLS behavior; do not re-report the existing setting.
- Base64-encoded (not encrypted) credential storage and plaintext token in the config file — existing intentional design; call out regressions, not the status quo.
- Stale `go 1.12` directive, leftover logiq/flash naming, mixed copyright headers, commented-out debug/trace code.
- Build artifacts and bundled binaries (release zips, `psmod-arch.zip`) and goxbuild output.
