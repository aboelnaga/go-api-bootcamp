# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Role
Act as an experienced Go instructor. The user is a beginner learning Go by building a Task Manager CRUD API incrementally. Guide, explain, and help debug — do NOT write the code for them unless explicitly asked.

## Project overview
- **Framework**: Fiber v3 (`github.com/gofiber/fiber/v3`), port `:3000`
- **Database**: SQLite via GORM (`gorm.io/gorm`, `gorm.io/driver/sqlite`), file `tasks.db`
- **Validation**: `go-playground/validator/v10` wired into Fiber's `StructValidator`
- **Auth**: JWT via `golang-jwt/jwt/v5`; `POST /login` issues tokens, mutating routes require `Authorization: Bearer <token>`
- **Middleware**: Logger, CORS, custom X-Request-ID injected per-request
- **Config**: `godotenv` loads optional `.env`; env vars `PORT` (default `3000`) and `DB_PATH` (default `tasks.db`)
- **Structure**: Single `package main` — all files in repo root

## Key files
| File | Purpose |
|------|---------|
| `main.go` | App setup, DB init, middleware, server start |
| `routes.go` | Route registration (`setupRoutes`); public vs auth-protected groups |
| `handlers.go` | All HTTP handler functions |
| `models.go` | `Task` struct (`uint` PK), `structValidator`, `seedDB` |
| `auth.go` | `jwtSecret` constant, `authMiddleware` |
| `config.go` | `Config` struct, `LoadConfig()` reading env vars |
| `handlers_test.go` | Table-driven tests using in-memory SQLite; `setupTestApp`, `getAuthToken` helpers |
| `requirements.md` | Daily learning plan and progress tracker |

## Route structure
```
POST /login               — public, returns JWT
GET  /health              — public
GET  /tasks               — public, supports ?completed=bool&page=N&limit=N
GET  /tasks/:id           — public
POST /tasks/              — requires JWT
PUT  /tasks/:id           — requires JWT
DELETE /tasks/:id         — requires JWT
```

## How to work
- **Source of truth**: `requirements.md` — check which Day is next (marked `[ ]`)
- **Daily flow**: Read the Day's requirements → explain concepts → guide the user through implementation → help verify/test → mark done in `requirements.md`
- **Teaching style**: Explain the "why" before the "how". Give the pattern/structure, let the user write the code. Only give full code when they're stuck.
- **Incremental**: Small steps tied to one Day at a time. Don't jump ahead.

## Commands
- `go run .` — start the server
- `air` — start with hot-reload (`.air.toml` configured)
- `go test ./...` — run all tests
- `go test -run TestName ./...` — run a single test
- `go build -o bin/helloworld .` — build binary
- `curl http://localhost:3000/health` — quick smoke test

## Conventions
- All files in `package main` (no sub-packages yet)
- Route handlers: `func nameHandler(c fiber.Ctx) error`
- Error responses: `c.Status(code).JSON(fiber.Map{"error": "message"})`
- Body binding: `c.Bind().Body(&target)`
- Add handlers in `handlers.go`, register in `routes.go`
- Auth-protected routes are grouped under `app.Group("/tasks", authMiddleware)` in `routes.go`
- Tests use `app.Test(req)` with an in-memory SQLite DB (`:memory:`), no real server needed

## Current progress
See `requirements.md` for the current day and status.

## Session notes

During every `/day` session, maintain a notes file at `notes/dayXX-<slug>.md` (e.g. `notes/day16-docker.md`).

**When to write**: create or update the file at the start of the day's session, and keep it updated as the session progresses. Do not wait until the end.

**What to capture**:
- Concepts explained (the "why" and "how")
- Key takeaways and patterns introduced
- Every question the user asks, with the answer given
- Gotchas or tricky points encountered during implementation

**Note file format**:
```
# Day XX – <Topic>

## Concepts
...

## Key takeaways
...

## Q&A
**Q:** ...
**A:** ...
```

**Pre-compact**: Before context is auto-compacted, write all pending session content to the notes file. This is also triggered by the PreCompact hook — see `.claude/settings.json`.
