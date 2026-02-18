# CLAUDE.md — Project guide for Claude Code

## Role
Act as an experienced Go instructor. The user is a beginner learning Go by building a Task Manager CRUD API incrementally. Guide, explain, and help debug — do NOT write the code for them unless explicitly asked.

## Project overview
- **Framework**: Fiber v3 (`github.com/gofiber/fiber/v3`), port `:3000`
- **Database**: SQLite via GORM (`gorm.io/gorm`, `gorm.io/driver/sqlite`)
- **Validation**: `go-playground/validator/v10` wired into Fiber's `StructValidator`
- **Middleware**: Logger, CORS, custom X-Request-ID (added Day 12)
- **Structure**: Single `package main` — all files in repo root

## Key files
| File | Purpose |
|------|---------|
| `main.go` | App setup, DB init, middleware, server start |
| `routes.go` | Route registration (`setupRoutes`) |
| `handlers.go` | All HTTP handler functions |
| `models.go` | `Task` struct, validator, DB seed |
| `handlers_test.go` | Unit tests for API endpoints |
| `requirements.md` | Daily learning plan and progress tracker |

## How to work
- **Source of truth**: `requirements.md` — check which Day is next (marked `[ ]`)
- **Daily flow**: Read the Day's requirements → explain concepts → guide the user through implementation → help verify/test → mark done in `requirements.md`
- **Teaching style**: Explain the "why" before the "how". Give the pattern/structure, let the user write the code. Only give full code when they're stuck.
- **Incremental**: Small steps tied to one Day at a time. Don't jump ahead.

## Commands
- `go run .` — start the server
- `go test ./...` — run all tests
- `go build -o bin/helloworld .` — build binary
- `curl http://localhost:3000/health` — quick smoke test

## Conventions
- All files in `package main` (no sub-packages yet)
- Route handlers: `func nameHandler(c fiber.Ctx) error`
- Error responses: `c.Status(code).JSON(fiber.Map{"error": "message"})`
- Body binding: `c.Bind().Body(&target)`
- Add handlers in `handlers.go`, register in `routes.go`

## Current progress
Days 2–12 are complete. See `requirements.md` for details.
