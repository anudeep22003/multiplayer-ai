# AGENTS.md

## Overview

`multiplayer-ai` (`github.com/anudeep22003/multiplayer-ai`) is a production Go backend server built with Go `1.27.2` and a standard-library-first philosophy.

As the codebase grows, organize new packages under `internal/` around clear responsibilities or domains rather than introducing premature layers or generic `util` packages.

## Commands

Run all commands from `backend/` so `.env.local` and module paths resolve properly; see `Makefile` for common targets.

## Engineering Guidelines

### 1. Dependencies & Standard Library
- Prefer Go standard library packages (`net/http` with method+pattern ServeMux routing, `context`, `errors`, `log/slog`, `net/http/httptest`) before reaching for external frameworks.
- Ask or justify before adding new third-party dependencies to `go.mod`.

### 2. Configuration & Secrets
- All runtime settings belong in `internal.Config` using `env` / `envDefault` tags.
- Exported environment variables take precedence over `.env.local` (`godotenv.Load` does not overwrite set env vars).
- Whenever adding, renaming, or removing a `Config` field, update `.env.sample` in the same change.
- Never commit, hardcode, or log secrets or connection strings (e.g., `DATABASE_URL`); avoid reading `.env.local` unless explicitly requested.

### 3. HTTP, Concurrency & Production Readiness
- **Timeouts & Shutdown**: Every server, external HTTP client, and DB pool must have explicit timeouts. Long-running work and shutdown paths must respect `context.Context` deadlines.
- **Context Propagation**: Always pass `r.Context()` down into service, database, and outbound network calls.
- **Error Handling**: Wrap internal errors with context (`fmt.Errorf("...: %w", err)`), inspect with `errors.Is` / `errors.As`, and never expose internal error details or stack traces in client-facing HTTP responses.
- **Route & Wiring Structure**: Domain handlers live in `internal/handlers/<domain>` with a `NewMux(...) *http.ServeMux` constructor and narrow consumer-defined interfaces. `internal/handlers` holds generic assembly and middleware (`Assemble`, `RouteGroup`, `Chain`); `internal/app` is the single composition root that wires concrete dependencies and middleware onto route groups.
- **Entrypoint Discipline**: Keep `cmd/server/main.go` (HTTP server), `cmd/api/main.go` (AWS Lambda), and `cmd/migrate/main.go` (goose CLI) focused strictly on configuration loading, resource lifecycle, and transport/command startup; keep shared route/dependency wiring in `internal/app`.

### 4. Database, Migrations & Queries
- **Migrations**: Add sequential SQL files in `migrations/*.sql` with `-- +goose Up` / `-- +goose Down` markers; apply via `cmd/migrate` using `DIRECT_DATABASE_URL` (keep `DIRECT_DATABASE_URL` out of `internal.Config`).
- **Queries (`sqlc`)**: Put static queries in `queries/<domain>.sql` and run `make sqlc` (`sqlc generate`); never edit `internal/db/*.go` manually, and always commit generated files.
- **Consumer Interfaces**: Handlers depend on narrow domain interfaces satisfied by `*db.Queries` (wired in `internal/app`), rather than importing `pgx` directly.
- **Transactions**: For multi-query atomicity, begin a `pgx.Tx` from `*pgxpool.Pool` in a small service/repository and run `queries.WithTx(tx)`.
- **Raw SQL Escape Hatch**: Drop down to `*pgxpool.Pool` + `pgx.CollectRows` only for dynamic queries (variable `WHERE` clauses, dynamic sorting) that `sqlc` cannot express cleanly.

### 5. Code Style & Testing
- Accept small interfaces at the consumer side; return concrete structs from constructors.
- Write deterministic, table-driven tests (`*_test.go`) alongside implementation code, using `net/http/httptest` for HTTP handlers.
- Always verify changes with `go vet ./...`, `go test -race ./...`, and `go build ./...` before considering work complete.
