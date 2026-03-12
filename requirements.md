## Task Manager CRUD API – Requirements

### Functional Requirements

- **Core entity: Task**
  - `id`: unique identifier (number or UUID)
  - `title`: short text, required
  - `description`: longer text, optional
  - `completed`: boolean flag
  - `createdAt`: timestamp when the task was created

- **Core API endpoints (MVP target)**
  - `POST /tasks`: create a new task
  - `GET /tasks`: list all tasks
  - `GET /tasks/:id`: get a single task by id
  - `PUT /tasks/:id`: update an existing task
  - `DELETE /tasks/:id`: delete a task

---

### Bash quick reference

These are the commands you'll use most. Each day reinforces them in context.

| Command | What it does |
|---------|-------------|
| `ls -la` | List all files including hidden ones, with permissions |
| `cat <file>` | Print a file's contents to the terminal |
| `touch <file>` | Create an empty file |
| `mkdir -p <path>` | Create a directory (and any missing parents) |
| `rm <file>` | Delete a file |
| `rm -rf <dir>` | Delete a directory and everything in it (irreversible — be careful) |
| `cp <src> <dst>` | Copy a file |
| `mv <src> <dst>` | Move or rename a file |
| `pwd` | Print the current working directory |
| `cd <path>` | Change directory |
| `code <file>` | Open a file in VS Code from the terminal |
| `nano <file>` | Open a file in nano (simple terminal editor; `Ctrl+X` to exit) |
| `echo $?` | Print the exit code of the last command (0 = success) |
| `export KEY=val` | Set an environment variable for the current session |
| `env \| grep KEY` | List env vars filtered by keyword |
| `which <cmd>` | Show the full path to a command binary |
| `<cmd> &` | Run a command in the background |
| `kill <pid>` | Stop a process by its PID |
| `lsof -i :<port>` | Find what process is using a port |

---

### Day 2 – [x] Health check endpoint

- **Objective**: Have a running Fiber server with a simple health check endpoint.
- **Endpoint**: `GET /health`
- **Behavior**:
  - Returns HTTP `200 OK`
  - Returns JSON body similar to: `{"status": "ok"}`
- **Done when**:
  - `go run .` starts the Fiber server without errors.
  - Visiting `http://localhost:3000/health` (or your chosen port) returns the JSON response.

### Day 3 – [x] GET /tasks (in-memory)

- **Objective**: Implement `GET /tasks` that returns a hard-coded list of tasks (no database yet).
- **Endpoint**: `GET /tasks`
- **Behavior**:
  - Returns HTTP `200 OK`
  - Returns JSON array of task objects, e.g.:
    - `[{ "id": 1, "title": "Learn Go", "completed": false }, { "id": 2, "title": "Build Task API", "completed": true }]`
  - Data is stored in memory (a Go slice), not in a database.
- **Done when**:
  - `go run .` starts the server without errors.
  - `curl http://localhost:3000/tasks` returns the JSON array of tasks.

### Day 4 – [x] GET /tasks/:id

- **Objective**: Implement `GET /tasks/:id` to fetch a single task by its `id` from the same in-memory slice.
- **Endpoint**: `GET /tasks/:id`
- **Behavior**:
  - Returns HTTP `200 OK` with the matching task as JSON when `id` exists.
  - Returns HTTP `404 Not Found` with a simple JSON error, e.g. `{"error": "task not found"}`, when `id` does not exist.
- **Done when**:
  - `curl http://localhost:3000/tasks/1` returns the JSON for task with `id = "1"`.
  - `curl http://localhost:3000/tasks/999` returns `404` with the JSON error.

### Day 5 – [x] POST /tasks

- **Objective**: Implement `POST /tasks` to create a new task in memory.
- **Endpoint**: `POST /tasks`
- **Behavior**:
  - Accepts JSON body with at least `title` (string) and optional `description`.
  - Appends a new `Task` to the in-memory `tasks` slice with:
    - A new `id` (e.g. incrementing string `"3"`, `"4"`, ...).
    - `completed` defaulting to `false`.
    - `createdAt` set to the current time.
  - Returns HTTP `201 Created` with the created task as JSON.
- **Done when**:
  - `curl -i -X POST http://localhost:3000/tasks -H "Content-Type: application/json" -d '{"title":"New Task","description":"Test"}'`
    returns `201` and the created task JSON.

### Day 6 – [x] PUT /tasks/:id

- **Objective**: Implement `PUT /tasks/:id` to update an existing task in memory.
- **Endpoint**: `PUT /tasks/:id`
- **Behavior**:
  - Accepts JSON body with updatable fields (e.g. `title`, `description`, `completed`).
  - Finds the task by `id` in the `tasks` slice.
  - Updates the task's fields while keeping `id` and `createdAt` unchanged.
  - Returns HTTP `200 OK` with the updated task as JSON when `id` exists.
  - Returns HTTP `404 Not Found` with `{"error": "task not found"}` when `id` does not exist.
- **Done when**:
  - `curl -i -X PUT http://localhost:3000/tasks/1 -H "Content-Type: application/json" -d '{"title":"Updated title","completed":true}'`
    returns `200` and the updated task JSON.
  - `curl -i -X PUT http://localhost:3000/tasks/999 -H "Content-Type: application/json" -d '{"title":"X"}'`
    returns `404` with the error JSON.

### Day 7 – [x] DELETE /tasks/:id

- **Objective**: Implement `DELETE /tasks/:id` to delete a task from memory.
- **Endpoint**: `DELETE /tasks/:id`
- **Behavior**:
  - Finds the task by `id` in the `tasks` slice.
  - Removes the task from the slice.
  - Returns HTTP `200 OK` with a confirmation message, e.g. `{"message": "task deleted"}`.
  - Returns HTTP `404 Not Found` with `{"error": "task not found"}` when `id` does not exist.
- **Done when**:
  - `curl -i -X DELETE http://localhost:3000/tasks/1` returns `200` with the confirmation JSON.
  - `curl http://localhost:3000/tasks` no longer includes the deleted task.
  - `curl -i -X DELETE http://localhost:3000/tasks/999` returns `404` with the error JSON.

### Day 8 – [x] Query filtering

- **Objective**: Add query parameter filtering to `GET /tasks` so clients can filter by `completed` status.
- **Endpoint**: `GET /tasks?completed=true` or `GET /tasks?completed=false`
- **Behavior**:
  - If no query parameter is provided, return all tasks (existing behavior).
  - If `?completed=true` is provided, return only tasks where `completed` is `true`.
  - If `?completed=false` is provided, return only tasks where `completed` is `false`.
  - Returns HTTP `200 OK` with a JSON array of matching tasks (may be empty `[]`).
- **Done when**:
  - `curl http://localhost:3000/tasks` returns all tasks.
  - `curl http://localhost:3000/tasks?completed=true` returns only completed tasks.
  - `curl http://localhost:3000/tasks?completed=false` returns only incomplete tasks.

### Day 9 – [x] Refactor project structure

- **Objective**: Refactor the project by splitting `main.go` into separate files with clear responsibilities.
- **Target structure**:
  ```
  helloworld/
  ├── main.go          # App setup and server start only
  ├── models.go        # Task struct and validator struct
  ├── handlers.go      # All route handler functions
  ├── routes.go        # Route registration (maps paths to handlers)
  ```
- **Rules**:
  - All files stay in `package main` (no sub-packages yet).
  - The API behavior must remain exactly the same — no new features.
  - Each file should have a single, clear purpose.
- **Done when**:
  - `go run .` starts the server without errors.
  - All existing curl commands still work identically.
  - `main.go` is short and only contains app config and `app.Listen`.

### Day 10 – [x] SQLite with GORM

- **Objective**: Replace the in-memory slice with a SQLite database using GORM.
- **What you'll learn**: Go ORM basics, database migrations, and persistent storage.
- **Steps**:
  - Install GORM and the SQLite driver (`gorm.io/gorm`, `gorm.io/driver/sqlite`).
  - Configure GORM to auto-migrate the `Task` model on startup.
  - Replace all slice operations in handlers with GORM queries (`db.Create`, `db.Find`, `db.First`, `db.Save`, `db.Delete`).
- **Done when**:
  - `go run .` starts the server and creates a `tasks.db` file.
  - All existing curl commands work identically.
  - Data persists across server restarts.

### Day 11 – [x] Unit tests

- **Objective**: Write unit tests for your API endpoints using Go's `testing` package and `net/http/httptest`.
- **What you'll learn**: Go testing conventions, table-driven tests, and HTTP testing.
- **Steps**:
  - Create `handlers_test.go` with tests for each endpoint.
  - Use table-driven tests (a slice of test cases) for different scenarios (success, not found, bad input).
  - Test at least: `GET /health`, `POST /tasks` (valid + invalid), `GET /tasks/:id` (found + not found), `DELETE /tasks/:id`.
- **Done when**:
  - `go test ./...` passes with all tests green.
  - Each endpoint has at least one success and one failure test case.

### Day 12 – [x] Logging & CORS middleware

- **Objective**: Add middleware for request logging and CORS.
- **What you'll learn**: How Fiber middleware works, request lifecycle, and cross-origin resource sharing.
- **Steps**:
  - Add Fiber's built-in `logger` middleware to log every request (method, path, status, duration).
  - Add Fiber's `cors` middleware so the API can be called from a browser frontend.
  - Create a custom middleware that adds a `X-Request-ID` header to every response (use `uuid` package).
- **Done when**:
  - Every request prints a log line in the server terminal.
  - Responses include `Access-Control-Allow-Origin` and `X-Request-ID` headers.

### Day 13 – [x] Environment config

- **Objective**: Add environment-based configuration using a `.env` file.
- **What you'll learn**: The `godotenv` package, `os.Getenv`, and separating config from code.
- **Steps**:
  - Install `github.com/joho/godotenv`.
  - Move hardcoded values (port, database path) into a `.env` file.
  - Create a `config.go` that loads and exposes configuration values.
  - Add `.env` to `.gitignore` and create a `.env.example` with placeholder values.
- **Done when**:
  - Changing `PORT=4000` in `.env` makes the server start on port 4000.
  - The app works with default values if `.env` is missing.

### Day 14 – [x] Pagination

- **Objective**: Add pagination to `GET /tasks`.
- **What you'll learn**: Query parameters, GORM `Offset`/`Limit`, and pagination response patterns.
- **Steps**:
  - Accept `?page=1&limit=10` query parameters (default: page 1, limit 10).
  - Use GORM's `.Offset().Limit()` to fetch only the requested page.
  - Return a JSON response with `tasks`, `page`, `limit`, and `total` count.
- **Done when**:
  - `curl http://localhost:3000/tasks?page=1&limit=1` returns 1 task and correct pagination metadata.
  - `curl http://localhost:3000/tasks` still works with default pagination.

### Day 15 – [x] JWT authentication

- **Objective**: Add JWT authentication to protect write endpoints.
- **What you'll learn**: JSON Web Tokens, auth middleware, and protecting routes.
- **Steps**:
  - Add a `POST /login` endpoint that accepts a hardcoded username/password and returns a JWT token.
  - Create an auth middleware that validates the JWT from the `Authorization: Bearer <token>` header.
  - Protect `POST`, `PUT`, and `DELETE /tasks` routes — `GET` routes remain public.
- **Done when**:
  - `POST /tasks` without a token returns `401 Unauthorized`.
  - `POST /login` returns a JWT token.
  - `POST /tasks` with a valid `Authorization` header succeeds.

### Day 16 – [x] CC: CLAUDE.md and Memory

> **Claude Code refresher** — no Go code today. You'll learn how Claude Code loads instructions and remembers context across sessions.
> *Why now?* You've built a solid API foundation. Before adding DevOps complexity, it's worth understanding how to configure Claude as a project-aware assistant.

- **Objective**: Understand how Claude Code loads instructions and remembers things across sessions.
- **What you'll learn**: CLAUDE.md hierarchy, auto memory, rules/, local vs project vs user scope.
- **Key concepts**:
  - `CLAUDE.md` files are Claude's persistent instructions — loaded from the project root, `~/.claude/`, and walked up the directory tree.
  - `~/.claude/CLAUDE.md` applies to every project; `CLAUDE.md` (or `.claude/CLAUDE.md`) applies to the current repo (check into git for team-sharing).
  - `CLAUDE.local.md` is your personal project-level override (gitignored automatically).
  - Auto memory: Claude writes notes to `~/.claude/projects/<project>/memory/MEMORY.md`; first 200 lines load at every session start.
  - `.claude/rules/` lets you split instructions into modular files with optional `paths:` front-matter to scope rules to specific directories.
- **Exercises**:
  - Open `~/.claude/CLAUDE.md` (create it if missing) and add 2–3 personal global preferences (e.g. "always use 2-space indentation").
  - In this repo, improve one section of the existing `CLAUDE.md`.
  - Run `/memory` inside a session and explore what Claude has saved.
  - Run `/init` in a fresh directory and see what it generates.
- **Bash practice**:
  - `ls -la ~/.claude/` — list all files (including hidden ones) in your Claude config directory
  - `cat ~/.claude/CLAUDE.md` — print a file to the terminal without opening an editor
  - `mkdir -p ~/.claude/rules` — create a directory and any missing parents in one step
  - `code ~/.claude/CLAUDE.md` — open a file in VS Code from the terminal
- **Done when**:
  - You can explain the precedence order: managed policy → project → user → local.
  - You've verified your global `~/.claude/CLAUDE.md` is being loaded (ask Claude what its instructions are).
  - You know how to view, edit, and disable auto memory.

### Day 17 – [x] Docker

- **Objective**: Dockerize the application.
- **What you'll learn**: Dockerfile, multi-stage builds, and container basics.
- **Steps**:
  - Create a `Dockerfile` with a multi-stage build (build stage with `golang` image, run stage with `alpine`).
  - Create a `docker-compose.yml` for easy local development.
  - Ensure the SQLite database file is persisted using a Docker volume.
- **Bash practice**:
  - `touch Dockerfile` — create an empty file (then open it with `code Dockerfile`)
  - `docker ps` — list running containers and their IDs
  - `docker logs <container_id>` — stream stdout/stderr from a container
  - `docker exec -it <container_id> sh` — open an interactive shell inside a running container
- **Done when**:
  - `docker compose up` starts the API and all curl commands work against it.
  - Restarting the container preserves task data.

### Day 18 – [ ] Swagger docs

- **Objective**: Add Swagger/OpenAPI documentation to the API.
- **What you'll learn**: API documentation with `swaggo/swag`, annotation comments, and auto-generated docs.
- **Steps**:
  - Install `swaggo/swag` and `swaggo/fiber-swagger`.
  - Add Swagger annotation comments to each handler function.
  - Serve the Swagger UI at `/swagger/*`.
- **Bash practice**:
  - `go install github.com/swaggo/swag/cmd/swag@latest` — install a Go CLI tool to your `$GOPATH/bin`
  - `which swag` — confirm the binary is on your `$PATH` (should print its full path)
  - `swag init` — generate the `docs/` folder from annotation comments
- **Done when**:
  - Visiting `http://localhost:3000/swagger/index.html` shows interactive API docs.
  - All endpoints are documented with request/response examples.

### Day 19 – [ ] Deploy to cloud

- **Objective**: Deploy the API to a free cloud platform.
- **What you'll learn**: Cloud deployment, environment variables in production, and basic DevOps.
- **Steps**:
  - Choose a platform (Railway, Render, or Fly.io — all have free tiers).
  - Configure production environment variables.
  - Deploy and verify all endpoints work with the public URL.
- **Bash practice**:
  - `export PORT=4000` — set an env var for the current shell session only
  - `env | grep PORT` — list all env vars and filter by keyword using a pipe
  - `curl -I https://your-app.example.com/health` — send a HEAD request (headers only, no body) to check if a deployment is live
- **Done when**:
  - The API is accessible at a public URL.
  - All curl commands work against the deployed version.

### Day 20 – [ ] CC: Skills and Cross-Project Commands

> **Claude Code refresher** — no Go code today. You'll build custom slash commands that travel with you across every project.
> *Why now?* You've just deployed a real API and have real workflows (commit, review, deploy). Time to automate them.

- **Objective**: Create custom slash commands that work in any project.
- **What you'll learn**: Skill frontmatter, `$ARGUMENTS`, user-level vs project-level skills, dynamic context injection with `!`.
- **Key concepts**:
  - Skills live in `~/.claude/skills/<name>/SKILL.md` (global) or `.claude/skills/<name>/SKILL.md` (project).
  - Invoke with `/<name>` in any Claude Code session.
  - Frontmatter controls model, tools, context, and invocation permissions.
  - `$ARGUMENTS` passes arguments; `` !`command` `` runs shell commands before Claude sees the prompt.
  - `disable-model-invocation: true` means only you can trigger it (Claude won't auto-invoke it).
- **Exercises**:
  1. Create `~/.claude/skills/commit/SKILL.md` — a `/commit` skill that reads `git diff --cached` and writes a conventional commit message.
  2. Create `~/.claude/skills/review/SKILL.md` — a `/review` skill that checks code for security issues and code quality.
  3. Create a project-level `.claude/skills/day/SKILL.md` that wraps the existing day skill logic.
  4. Test all three skills in a session.
- **Bash practice**:
  - `mkdir -p ~/.claude/skills/commit` — create nested directories in one command
  - `touch ~/.claude/skills/commit/SKILL.md` — create the file, then `code` it to edit
  - `cat ~/.claude/skills/commit/SKILL.md` — quickly read a small file without leaving the terminal
- **Skill frontmatter reference**:
  ```yaml
  ---
  name: commit
  description: Write a conventional commit after reviewing staged diff
  disable-model-invocation: true
  allowed-tools: Bash(git *)
  ---
  1. Run: !`git diff --cached`
  2. Write a conventional commit message (feat:, fix:, docs:, refactor:, test:)
  3. Commit with: git commit -m "<message>"
  ```
- **Done when**:
  - `/commit` works from any project directory.
  - You understand the difference between user-level and project-level skills.
  - You can explain when to use `disable-model-invocation: true`.

---

## Phase 2 — Go deeper on Go

### Day 21 – [ ] Graceful shutdown

- **Objective**: Handle OS signals so the server finishes in-flight requests before stopping.
- **What you'll learn**: `os/signal`, `context`, and production-safe shutdown patterns.
- **Bash practice**:
  - `lsof -i :3000` — find which process (and its PID) is listening on port 3000
  - `kill -SIGTERM <pid>` — send a graceful shutdown signal (ask the process to stop cleanly)
  - `ps aux | grep helloworld` — list running processes and filter by name to find the PID
- **Done when**:
  - `Ctrl+C` waits for active requests to finish before exiting.
  - A log line confirms graceful shutdown.

### Day 22 – [ ] Background jobs with goroutines

- **Objective**: Run a task asynchronously after an HTTP request returns.
- **What you'll learn**: Goroutines, channels, and Go's concurrency model in a practical context.
- **Bash practice**:
  - `go run . &` — start the server in the background so you keep your terminal free
  - `jobs` — list all background jobs running in the current shell session
  - `fg %1` — bring background job #1 back to the foreground
- **Done when**:
  - Creating a task triggers a background goroutine (e.g. logs creation asynchronously).
  - The HTTP response is not delayed by the background work.

### Day 23 – [ ] CC: Settings, Permissions, and Hooks

> **Claude Code refresher** — no Go code today. You'll configure Claude Code's automation behavior using settings and hooks.
> *Why now?* You're working with goroutines and concurrency — hooks like auto-run tests after edits become genuinely useful here.

- **Objective**: Configure Claude Code's behavior at a fine-grained level using settings files and hooks.
- **What you'll learn**: Settings hierarchy, permission rules, hook lifecycle events, hook types (command/http/prompt/agent).
- **Key concepts**:
  - Settings precedence: managed → CLI flags → local (`.claude/settings.local.json`) → project (`.claude/settings.json`) → user (`~/.claude/settings.json`).
  - Permission rules use glob syntax: `"allow": ["Bash(npm run *)", "Read(~/.zshrc)"]`.
  - Hooks fire on lifecycle events: `PreToolUse`, `PostToolUse`, `UserPromptSubmit`, `PreCompact`, `Stop`, etc.
  - Hook I/O: JSON on stdin, exit code 0 = proceed, exit code 2 = block (write reason to stderr).
- **Hook lifecycle events**:
  | Event | Use for |
  |-------|---------|
  | `PostToolUse` (matcher: `Edit\|Write`) | Auto-format files after edits |
  | `PreToolUse` (matcher: `Bash`) | Block dangerous commands |
  | `PreCompact` | Save session notes before context compaction |
  | `Stop` | Run tests after Claude finishes a task |
  | `Notification` | Desktop alerts when Claude needs attention |
- **Exercises**:
  1. Open `~/.claude/settings.json` and add a permission rule that asks before `git push`.
  2. Use `/hooks` to add a `PostToolUse` hook that runs `go fmt ./...` after any Go file is edited.
  3. Add a `PreCompact` hook that writes a "compacting…" message to a log file.
  4. Test by editing a `.go` file and confirm formatting runs automatically.
- **Bash practice**:
  - `chmod +x myscript.sh` — make a shell script executable (required before you can run it with `./`)
  - `./myscript.sh` — run a script in the current directory (the `./` tells the shell to look here, not in `$PATH`)
  - `echo $?` — print the exit code of the last command; hooks use this (0 = success, 2 = block)
- **Done when**:
  - You can read `.claude/settings.json` and explain what each section does.
  - At least one hook is working end-to-end.
  - You can use `/status` to see which settings are active and where they came from.

---

## Phase 3 — System design foundations

### Day 24 – [ ] Load balancer with Nginx

- **Objective**: Run two instances of the API behind an Nginx reverse proxy.
- **What you'll learn**: Load balancing concepts, upstream configuration, and horizontal scaling basics.
- **Bash practice**:
  - `nginx -t` — test your Nginx config for syntax errors before reloading
  - `tail -f /var/log/nginx/access.log` — stream a log file in real time (`Ctrl+C` to stop)
  - `curl -w "%{http_code}" -o /dev/null -s http://localhost` — print only the HTTP status code of a request
- **Done when**:
  - Two API instances run on different ports.
  - Nginx distributes requests between them in round-robin.

### Day 25 – [ ] Proper package structure

- **Objective**: Refactor into `internal/` sub-packages with interfaces and dependency injection.
- **What you'll learn**: Go package conventions, interfaces, and testable architecture.
- **Bash practice**:
  - `find . -name "*.go" -not -path "*/vendor/*"` — list all Go files in the project recursively
  - `tree -I vendor` — show the directory tree (install with `brew install tree` on macOS)
  - `wc -l handlers.go` — count the number of lines in a file
- **Done when**:
  - Handlers depend on interfaces, not concrete GORM types.
  - Each package has a single clear responsibility.

### Day 26 – [ ] CC: Plan → Execute → Iterate Workflow

> **Claude Code refresher** — no Go code today. You'll master structured, safe development using Claude Code's plan mode and iteration tools.
> *Why now?* The upcoming PostgreSQL migration and distributed systems work are non-trivial refactors. Plan mode is exactly the right tool for those.

- **Objective**: Master structured, safe development using Claude Code's plan mode and iteration tools.
- **What you'll learn**: Plan mode, checkpoints, rewinding, context compaction, extended thinking.
- **Key concepts**:
  - **Plan mode**: Claude can only read files and ask questions — it cannot make changes. Toggle with `Shift+Tab` or `--permission-mode plan`.
  - **Normal mode → Auto-Accept mode → Plan mode**: cycle with `Shift+Tab`.
  - **Checkpoints**: use `/checkpoint` or `Esc+Esc` to save state; use `/rewind` to restore.
  - **Extended thinking**: `Option+T` (macOS) toggles deeper reasoning; set effort with `CLAUDE_CODE_EFFORT_LEVEL=high`.
  - **Context compaction**: `/compact <instructions>` summarizes the session; `PreCompact` hook saves notes before it happens.
  - **SpecKit equivalent**: Claude Code does not have a dedicated spec runner, but you achieve the same effect by: (1) writing a spec in a `.md` file, (2) having Claude read it in Plan mode, (3) iterating in Normal mode with checkpoints.
- **Recommended workflow for non-trivial features**:
  ```
  1. Write a spec file (e.g. docs/spec-feature.md) describing requirements
  2. Enter Plan mode (Shift+Tab twice)
  3. Ask Claude to read the spec and create an implementation plan
  4. Refine the plan (ask follow-up questions)
  5. Exit Plan mode → Claude implements
  6. Run tests; if failing, /rewind and try a different approach
  7. /compact to free context, then continue iteration
  ```
- **Exercises**:
  1. Enter Plan mode, ask Claude to analyze `handlers.go` and propose a refactor — verify it makes no changes.
  2. Write a `docs/spec-postgres-migration.md` with requirements for Day 27's PostgreSQL migration.
  3. Use Plan mode with the spec to generate an implementation plan; export it to a file.
  4. Implement at least the first step of the plan using `/checkpoint` before each change.
  5. Use `/rewind` to undo one step and try a different approach.
- **Bash practice**:
  - `git stash` — temporarily shelve uncommitted changes so you can safely switch context
  - `git stash pop` — restore the most recently stashed changes
  - `git diff HEAD~1` — diff the current state against the commit before the last one
- **Done when**:
  - You can switch between Normal/Auto-Accept/Plan modes confidently.
  - You've completed at least one full Plan → Execute → Verify cycle.
  - You know how to use checkpoints and rewind to recover from mistakes.

### Day 27 – [ ] PostgreSQL

- **Objective**: Swap SQLite for PostgreSQL.
- **What you'll learn**: Connection pooling, indexes, query explain plans, and production databases.
- **Bash practice**:
  - `psql -U postgres -d taskdb` — connect to a PostgreSQL database interactively in the terminal
  - `\dt` — list all tables (run this inside the `psql` shell)
  - `pg_dump -U postgres taskdb > backup.sql` — export your entire database to a SQL file
- **Done when**:
  - All endpoints work with Postgres.
  - An index is added and verified with `EXPLAIN ANALYZE`.

---

## Phase 4 — Distributed systems

### Day 28 – [ ] Message queue (Kafka or RabbitMQ)

- **Objective**: Publish an event when a task is created; consume it in a separate service.
- **What you'll learn**: Producer/consumer pattern, async communication, event-driven architecture.
- **Bash practice**:
  - `docker compose logs -f kafka` — follow live logs for a specific Docker Compose service
  - `docker compose exec kafka kafka-topics.sh --list --bootstrap-server localhost:9092` — run a command inside a running service container
  - `curl -s http://localhost:3000/tasks | jq .` — pipe JSON output through `jq` to pretty-print it (`brew install jq`)
- **Done when**:
  - Task creation publishes a message to a queue.
  - A consumer service reads and logs the message.

### Day 29 – [ ] Caching with Redis

- **Objective**: Cache the `GET /tasks` response in Redis.
- **What you'll learn**: Cache-aside pattern, TTL, and cache invalidation.
- **Bash practice**:
  - `redis-cli ping` — confirm Redis is running (should return `PONG`)
  - `redis-cli keys "*"` — list all keys currently in Redis
  - `redis-cli ttl tasks:all` — check the remaining TTL (in seconds) of a specific cached key
- **Done when**:
  - Repeated `GET /tasks` calls are served from Redis.
  - Cache is invalidated when a task is created or updated.

### Day 30 – [ ] CC: MCP Servers

> **Claude Code refresher** — no Go code today. You'll extend Claude Code with external tools via MCP servers.
> *Why now?* You're using Redis, Kafka, and Postgres. MCP servers let Claude query these same services directly during development.

- **Objective**: Extend Claude Code with external tools via MCP (Model Context Protocol) servers.
- **What you'll learn**: MCP concepts, installing servers, scopes (local/project/user), OAuth, `@` resource mentions.
- **Key concepts**:
  - MCP servers expose tools, resources, and prompts to Claude Code via a standardized protocol.
  - Three transports: `http` (recommended), `sse` (deprecated), `stdio` (local process).
  - Scopes: `--scope local` (this project, personal), `--scope project` (`.mcp.json`, team-shared), `--scope user` (all projects).
  - Resources can be referenced with `@server:resource-uri` in prompts.
- **Installation examples**:
  ```bash
  # HTTP server
  claude mcp add --transport http github https://api.githubcopilot.com/mcp/

  # Local stdio server
  claude mcp add --transport stdio postgres -- npx -y @modelcontextprotocol/server-postgres $DB_URL

  # Project-scoped (committed to git)
  claude mcp add --scope project --transport http notion https://mcp.notion.com/mcp
  ```
- **Exercises**:
  1. Run `claude mcp list` to see currently installed servers.
  2. Install the GitHub MCP server (requires a GitHub token): `claude mcp add --transport http github https://api.githubcopilot.com/mcp/`.
  3. In a session, use the GitHub server to fetch an issue or PR from a real repo.
  4. Create a `.mcp.json` in this project for the Postgres MCP server pointing at your local DB.
  5. Use `/mcp` in a session to authenticate and explore available tools.
- **Bash practice**:
  - `claude mcp list` — list installed MCP servers directly from the terminal
  - `jq '.mcpServers' .mcp.json` — extract a specific field from a JSON config file
  - `curl -H "Authorization: Bearer $TOKEN" https://api.example.com` — use an env var in a curl command (the shell substitutes `$TOKEN` before sending)
- **Done when**:
  - At least one MCP server is installed and you've successfully used it from a session.
  - You understand the difference between local, project, and user scopes.
  - You can explain when to commit `.mcp.json` to git (team tools) vs keep in `~/.claude.json` (personal).

### Day 31 – [ ] Full-text search with Elasticsearch

- **Objective**: Add full-text search to tasks via Elasticsearch.
- **What you'll learn**: Indexing, search queries, and when to use a search engine vs a DB.
- **Bash practice**:
  - `curl -X GET "localhost:9200/_cat/indices?v"` — list all Elasticsearch indices with stats
  - `curl -s -X POST "localhost:9200/tasks/_search" -H "Content-Type: application/json" -d '{"query":{"match_all":{}}}' | jq .hits.total` — run a search and extract a nested field from the result
  - `| jq '.hits.hits[].\_source'` — extract the `_source` field from every search result
- **Done when**:
  - `GET /tasks?q=keyword` returns tasks matching the keyword via Elasticsearch.

### Day 32 – [ ] Columnar DB with ClickHouse

- **Objective**: Stream task events into ClickHouse and build a simple analytics query.
- **What you'll learn**: Columnar storage, OLAP vs OLTP, and analytics query patterns.
- **Bash practice**:
  - `clickhouse-client --query "SHOW TABLES"` — run a ClickHouse query directly from the terminal
  - `clickhouse-client --query "SELECT count() FROM task_events"` — count rows without opening an interactive session
  - `time curl -s http://localhost:3000/tasks > /dev/null` — measure how long an HTTP request takes (the `time` prefix works with any command)
- **Done when**:
  - Task creation events are inserted into ClickHouse.
  - A query returns task creation counts grouped by day.

---

## Phase 5 — Microservices

### Day 33 – [ ] Split into services

- **Objective**: Extract auth and notifications into separate services.
- **What you'll learn**: Service boundaries, inter-service HTTP/gRPC communication, and the real cost of microservices.
- **Bash practice**:
  - `go run ./cmd/auth &` — run a specific `main` package (sub-directory) in the background
  - `lsof -i -P -n | grep LISTEN` — list every port currently being listened on, with process names
  - `curl -s http://localhost:8081/health | jq .` — test a service and pretty-print the JSON response in one line
- **Done when**:
  - Auth service issues JWTs independently.
  - Task service validates tokens by calling the auth service.

### Day 34 – [ ] CC: Global Developer Toolkit

> **Claude Code refresher** — no Go code today. You'll build a personal library of reusable skills, settings, and hooks that work across all your projects.
> *Why now?* You're managing multiple services. A global toolkit that works everywhere becomes essential, not optional.

- **Objective**: Build a personal `~/.claude/` toolkit that travels to every project you work in.
- **What you'll learn**: Structuring a personal toolkit, plugin distribution, common patterns the community uses.
- **Your global toolkit structure**:
  ```
  ~/.claude/
  ├── CLAUDE.md                    # Global personal instructions
  ├── settings.json                # Global permissions and hooks
  ├── skills/
  │   ├── commit/SKILL.md          # Conventional commits
  │   ├── review/SKILL.md          # Code review
  │   ├── fix-issue/SKILL.md       # Fix GitHub issue by number
  │   ├── pr/SKILL.md              # Create PR with description
  │   ├── test-debug/SKILL.md      # Debug failing tests
  │   └── doc/SKILL.md             # Generate/update docs
  ├── rules/
  │   ├── code-style.md            # Language-agnostic style rules
  │   ├── git-workflow.md          # Commit and PR conventions
  │   └── security.md             # Security checklist
  └── agents/
      └── reviewer/AGENT.md       # Dedicated code review agent
  ```
- **Common community skills**:
  | Skill | What it does |
  |-------|-------------|
  | `/commit` | Stage, write conventional commit, commit |
  | `/pr` | Push branch, open PR with summary |
  | `/fix-issue <n>` | Fetch GitHub issue, implement fix, create PR |
  | `/review` | Security + quality review of changed files |
  | `/test-debug` | Analyze failing tests and suggest fixes |
  | `/doc` | Add/update docstrings and README sections |
  | `/simplify` | (built-in) Parallel code quality review |
  | `/batch <instruction>` | (built-in) Large-scale parallel changes |
  | `/debug` | (built-in) Troubleshoot Claude Code session |
- **Exercises**:
  1. Build out at least 3 skills in `~/.claude/skills/` that you'll actually use daily.
  2. Add global hooks: `PostToolUse` for auto-formatting (language-aware), `Stop` to run tests.
  3. Write `~/.claude/rules/git-workflow.md` with your commit and branching conventions.
  4. Test each skill from this Go project and at least one other project directory.
- **Bash practice**:
  - `ln -s ~/.claude/settings.json ~/dotfiles/claude-settings.json` — create a symlink so a file lives in two places at once without duplicating it
  - `git init && git remote add origin <url>` — turn your `~/.claude/` into a versioned dotfiles repo
  - `source ~/.zshrc` — reload your shell config without closing and reopening the terminal
- **Done when**:
  - Your `~/.claude/` directory is version-controlled (push it to a private dotfiles repo).
  - You have at least 3 working global skills.
  - You can onboard a new project in under 5 minutes using your global toolkit.

### Day 35 – [ ] API gateway

- **Objective**: Add an API gateway in front of all services.
- **What you'll learn**: Routing, auth delegation, rate limiting at the gateway level.
- **Bash practice**:
  - `curl -v http://localhost:8080/tasks` — verbose mode: shows the full request and response headers, useful for debugging routing
  - `watch -n1 "curl -s http://localhost:8080/health"` — re-run a command every second to monitor a live service
  - `ab -n 100 -c 10 http://localhost:8080/tasks` — basic load test: 100 requests, 10 at a time (Apache Bench, pre-installed on macOS)
- **Done when**:
  - All client requests go through the gateway.
  - The gateway forwards to the correct service based on the path.

### Day 36 – [ ] Service discovery

- **Objective**: Services find each other dynamically instead of hardcoded URLs.
- **What you'll learn**: Consul or Kubernetes-style service discovery basics.
- **Bash practice**:
  - `consul members` — list all nodes registered with Consul
  - `dig @127.0.0.1 -p 8600 task-service.service.consul` — resolve a service name through Consul's built-in DNS server
  - `ss -tlnp` — list all open TCP listening ports with process names (modern replacement for `netstat`)
- **Done when**:
  - Services register themselves on startup.
  - The gateway resolves service addresses dynamically.

### Day 37 – [ ] CC: Multi-Agent Workflows

> **Claude Code refresher** — no Go code today. You'll learn how Claude Code spawns and coordinates subagents to parallelize work and protect context.
> *Why now?* You've built a distributed system across multiple services. Multi-agent workflows mirror that same "divide and conquer" pattern — but for AI-assisted development tasks.

- **Objective**: Understand how Claude Code uses subagents and how to design multi-agent workflows.
- **What you'll learn**: The Task tool, subagent types, parallel vs sequential agents, context isolation, the Claude Agent SDK.
- **Key concepts**:
  - **Subagents** are separate Claude instances spawned by the main session via the Task tool. Each gets its own context window.
  - **Why subagents?** Two reasons: (1) parallelism — run multiple tasks at once; (2) context isolation — a subagent's output doesn't bloat the main conversation.
  - **Subagent types** built into Claude Code: `Bash`, `Explore`, `Plan`, `general-purpose`, `claude-code-guide`, and more. Each has a fixed toolset.
  - **Session storage**: subagent conversations are saved as `subagents/agent-<id>.jsonl` nested under the parent session's folder in `~/.claude/projects/`.
  - **Claude Agent SDK**: the programmatic API for building multi-agent pipelines outside of Claude Code — useful for automating complex workflows in CI or tooling.
  - **When to use parallel agents**: independent tasks (e.g. "search for X" and "search for Y" simultaneously). Never parallelize tasks that depend on each other's output.
- **Exercises**:
  1. In a session, ask Claude to use two parallel subagents: one to explore `handlers.go` and one to explore `routes.go` — observe how results are merged back.
  2. Open `~/.claude/projects/<this-project>/` and find the `subagents/` folder from a past session. Read one JSONL file and identify the message types (user, assistant, tool use).
  3. Read the [Claude Agent SDK docs](https://docs.anthropic.com/en/docs/claude-code/sdk) and identify one workflow from your own projects you could automate with it.
  4. Use the built-in `/batch` command to make a parallel change across multiple files (e.g. add a comment to every handler function).
- **Bash practice**:
  - `ls ~/.claude/projects/<project>/` — list session files for a project
  - `ls ~/.claude/projects/<project>/*/subagents/` — list all subagent JSONL files across sessions
  - `wc -l ~/.claude/projects/<project>/<uuid>.jsonl` — count how many messages are in a session
- **Done when**:
  - You can explain the difference between the main agent and a subagent.
  - You understand when to run agents in parallel vs sequentially.
  - You've located and inspected a real subagent JSONL file from a past session.
