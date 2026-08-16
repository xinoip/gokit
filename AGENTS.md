# Repository Guidelines

## Project Structure & Module Organization

This repository is a collection of intentionally isolated Go packages under the module `github.com/xinoip/gokit`. Core packages live at the repository root:

- `api/` defines Huma-based API registration and endpoint helpers.
- `httprpc/` provides RPC request/response conventions.
- `mux/` contains Chi router construction and middleware configuration.
- `server/` provides HTTP server lifecycle helpers.
- `migrate/` wraps database migration behavior.
- `testutil/` contains Postgres, Redis, API, and context test helpers.

Each package keeps documentation in a local `README.md`. The `example/` directory is a separate module demonstrating a notes API with migrations, sqlc persistence, and Redis caching. Avoid dependencies between root packages unless required; packages are designed to be copied independently.

## Build, Test, and Development Commands

Run commands from the repository root unless noted otherwise:

- `make all` runs dependency cleanup, linting, and the full test suite.
- `make tidy` updates `go.mod` and `go.sum` with `go mod tidy`.
- `make lint` runs the configured `golangci-lint` checks and formatters.
- `make test` runs `go test -coverprofile=coverage.out -race ./...`.
- `make clean` clears the Go test cache.
- `make doc` opens local package documentation with `pkgsite`.

For the sample app, use `cd example && make up-dependencies`, then `make migrate` and `make run`. `make gen` refreshes generated sqlc code.

## Coding Style & Naming Conventions

Use tabs as produced by `gofmt`, short lowercase package names, PascalCase for exported identifiers, and camelCase for unexported identifiers. Exported packages and symbols require doc comments. Keep lines within 120 characters. Before submitting, run `make lint`; it applies `gofmt`, `goimports`, `golines`, and a broad linter set.

## Testing Guidelines

Place tests beside implementation files as `*_test.go`. Prefer table-driven unit tests for behavior and Go `Example...` tests for public API documentation. Use `testify` where assertions improve clarity. No fixed coverage threshold is configured, but new behavior and regressions should be covered. Always run race-enabled `make test`; integration helpers may require local Postgres or Redis.

## Commit & Pull Request Guidelines

Recent commits use concise imperative subjects, commonly `feat: ...` or `fix: ...`; a package prefix may precede the type (for example, `rpc:feat: ...`). Keep each commit focused. Pull requests should explain the motivation and behavior change, mention affected packages, link relevant issues, and report `make all` results. Include generated files when inputs change and call out API or migration compatibility concerns; screenshots are only needed for documentation or UI-visible changes.
