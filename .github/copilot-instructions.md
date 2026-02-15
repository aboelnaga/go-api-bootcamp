# Copilot instructions for this repository

Purpose: help AI coding agents be productive quickly when changing or extending this Task API.

Summary
- Framework: Fiber (github.com/gofiber/fiber/v3). Server listens on port `:3000` by default.
- Key files: `main.go`, `routes.go`, `handlers.go`, `models.go`, `requirements.md`, `go.mod`.

Quick start (developer)
- Run locally: `go run .` from repo root (starts Fiber on :3000). See `main.go` for app config.
- Build: `go build -o bin/helloworld .`
- Tests: `go test ./...` (no tests currently; add `handlers_test.go` when adding tests).
- Debug: use Delve: `dlv debug --headless --listen=:2345 --api-version=2`.

Architecture notes (what to know before editing)
- Routing: `setupRoutes(app *fiber.App)` in `routes.go` registers endpoints to handlers (e.g. `/tasks`, `/health`).
- Handlers: `handlers.go` contains all route logic. Handlers accept `fiber.Ctx` and return JSON via `c.JSON(...)` or status+JSON.
- Models & validation: `models.go` defines `Task` and a `structValidator` that wraps `go-playground/validator`. `main.go` wires the validator into `fiber.Config{StructValidator: ...}`.
- Persistence: currently the code uses an in-memory `tasks` slice (see `models.go`) but also opens a GORM SQLite DB (`gorm.Open(sqlite.Open("tasks.db"))`). Be careful: both exist simultaneously — switching to DB-backed operations requires replacing slice logic in `handlers.go`.

Project-specific conventions
- Single-package app: all files are `package main` (no sub-packages). Keep this consistent unless doing an intentional refactor.
- Route registration pattern: add handler functions in `handlers.go`, then register them in `routes.go` via `setupRoutes`.
- Error responses: follow existing style — return `c.Status(<code>).JSON(fiber.Map{"error": "..."})`.
- Body binding: handlers use `c.Bind().Body(target)` pattern to parse JSON into structs — match existing usage when adding new endpoints.
- ID values: current IDs are strings (incrementing numbers as strings). Keep ID type consistent unless migrating to UUIDs (update requirements and handlers together).

Dependencies & environment
- Go version declared in `go.mod`: `go 1.25.7` — some newer stdlib packages (e.g., `slices`) are used.
- DB driver: GORM + SQLite are present in `go.mod`. `models.go` panics at init if DB open fails — handle carefully when changing startup.

Examples (how to make common edits)
- Add a new GET route: add handler in `handlers.go` with signature `func myHandler(c fiber.Ctx) error { ... }` then register in `routes.go` with `app.Get("/my", myHandler)`.
- Convert a handler to use GORM: replace in-memory operations (iterate or `slices.Delete`) with `db.Find`, `db.First`, `db.Create`, etc. Remove or phase out the global `tasks` slice.
- Add validation: add tags to `Task` fields in `models.go` and rely on the existing `structValidator` (it is already wired in `main.go`).

Gotchas discovered
- `models.go` both defines `tasks` slice and opens `tasks.db` with GORM; changing persistence requires updating handlers and removing slice usage to avoid data duplication.
- `go.mod` uses `go 1.25.7` so CI or contributors must have a compatible Go toolchain.
- The code imports the standard `slices` package (Go 1.21+). If contributors use older Go versions tests/builds will fail.

What to document or add next (suggestions for maintainers)
- Add a short `README.md` with the quick-start commands and port information.
- Add basic unit tests under `handlers_test.go` and a simple CI workflow that runs `go test ./...` and `go vet`.
- Add `.env.example` and a `config.go` if moving away from hardcoded port/DB paths.

Questions for the maintainer (for clarifying agent actions)
- Should we preserve the in-memory `tasks` slice when adding GORM, or migrate fully to DB-backed operations?
- Preferred ID format: keep string numeric IDs, or migrate to UUIDs (`github.com/google/uuid`)?

If anything here is unclear or you'd like a different level of detail, tell me which section to expand or revise.
