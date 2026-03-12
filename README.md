# Go API Bootcamp

A 37-day hands-on learning path for Go backend development. You build a **Task Manager REST API** from scratch, incrementally adding real-world features — databases, auth, middleware, Docker, distributed systems, and microservices.

Each day introduces one concept, one implementation, and one set of working curl commands to verify it.

---

## What you build

A production-ready Task Manager API with:

- Full CRUD over tasks (`POST`, `GET`, `PUT`, `DELETE`)
- SQLite persistence via GORM
- JWT authentication on mutating routes
- Query filtering, pagination, and sorting
- Request logging, CORS, and custom `X-Request-ID` middleware
- Environment-based configuration
- Table-driven unit tests with in-memory SQLite
- Docker + docker-compose with volume-persisted data

---

## Stack

| Layer | Tech |
|-------|------|
| Framework | [Fiber v3](https://github.com/gofiber/fiber) |
| Database | SQLite via [GORM](https://gorm.io) |
| Auth | JWT via [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) |
| Validation | [go-playground/validator v10](https://github.com/go-playground/validator) |
| Config | [godotenv](https://github.com/joho/godotenv) |
| Hot reload | [air](https://github.com/air-verse/air) |

---

## API routes

```
POST   /login          — get a JWT token (username: admin, password: secret)
GET    /health         — health check
GET    /tasks          — list tasks (?completed=bool&page=N&limit=N)
GET    /tasks/:id      — get a single task
POST   /tasks/         — create task (requires JWT)
PUT    /tasks/:id      — update task (requires JWT)
DELETE /tasks/:id      — delete task (requires JWT)
```

---

## Quick start

### Local

```bash
cp .env.example .env
go run .
curl http://localhost:3000/health
```

### Docker

```bash
mkdir -p data
docker compose up --build
curl http://localhost:3000/health
```

### Hot reload

```bash
air
```

---

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | HTTP port |
| `DB_PATH` | `tasks.db` | Path to the SQLite database file |

Copy `.env.example` to `.env` and adjust as needed.

---

## Running tests

```bash
go test ./...
go test -run TestName ./...   # run a single test
```

Tests use an in-memory SQLite database — no setup required.

---

## Project structure

```
.
├── main.go            # App setup, middleware, server start
├── routes.go          # Route registration
├── handlers.go        # HTTP handler functions
├── models.go          # Task struct, validator, seed data
├── auth.go            # JWT secret, auth middleware
├── config.go          # Config struct, env loading
├── handlers_test.go   # Table-driven tests
├── Dockerfile         # Multi-stage build
├── docker-compose.yml # Local dev with volume
├── .env.example       # Environment variable template
├── requirements.md    # Full 37-day curriculum
└── notes/             # Per-day session notes
```

---

## Learning path

| Phase | Days | Topics |
|-------|------|--------|
| **Foundations** | 1–9 | Go basics, Fiber, in-memory CRUD, project structure |
| **Production basics** | 10–17 | SQLite/GORM, tests, middleware, auth, config, pagination, Docker |
| **Tooling** | 16, 20, 23, 26 | Claude Code: CLAUDE.md, skills, hooks, plan mode |
| **Go deeper** | 21–22 | Graceful shutdown, goroutines |
| **System design** | 24–27 | Nginx, package structure, PostgreSQL |
| **Distributed systems** | 28–32 | Kafka, Redis, Elasticsearch, ClickHouse |
| **Microservices** | 33–37 | Service split, API gateway, service discovery, multi-agent workflows |

See [`requirements.md`](requirements.md) for the full day-by-day breakdown and progress tracker.

Per-day session notes (concepts, Q&A, gotchas) live in [`notes/`](notes/).

---

## Progress

- [x] Days 1–17 complete
- [ ] Days 18–37 in progress
