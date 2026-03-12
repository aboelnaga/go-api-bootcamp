# Day 17 – Docker

## Concepts

### What is Docker?
Docker packages your app and everything it needs to run into a **container** — an isolated, reproducible environment that behaves the same everywhere. "Works on my machine" stops being a problem.

Key terms:
- **Image**: a blueprint, built from a `Dockerfile`
- **Container**: a running instance of an image
- **Layer**: each instruction in a `Dockerfile` creates a cached layer — unchanged layers are reused on rebuild

### Multi-stage builds
A Go app only needs the Go compiler to **build**. At **runtime** you just need the compiled binary.

Multi-stage builds use two `FROM` statements in one `Dockerfile`:
1. A full **build stage** (large Go image with gcc, tools, etc.) that compiles the binary
2. A tiny **run stage** (e.g. `alpine`) that only receives the compiled binary

Result: ~20MB final image instead of ~800MB.

```dockerfile
FROM golang:1.23-alpine AS builder   # stage 1 — compile
...
FROM alpine:latest                   # stage 2 — run
COPY --from=builder /app/server .    # copy only the binary
```

### CGO and SQLite
`gorm.io/driver/sqlite` uses `mattn/go-sqlite3`, which is written in C and requires **CGO** (Go's C foreign-function interface). This has consequences for Docker:

- You **cannot** use `CGO_ENABLED=0` (a fully static Go build) with this driver
- The build stage needs a C compiler (`gcc`) and the C standard library headers (`musl-dev`)
- The final image needs a compatible C runtime

**Solution**: build on `golang:alpine` (add `gcc` + `musl-dev`), run on `alpine:latest`.
Both stages use **musl libc** (Alpine's default), so the compiled binary runs correctly in the final image without any extra packages.

### Docker volumes for SQLite persistence
Containers are ephemeral — their filesystem is destroyed when the container is removed. To persist the SQLite database file across restarts, mount a **volume**:

```yaml
volumes:
  - ./data:/app/data   # host path : container path
```

Set `DB_PATH=/app/data/tasks.db` so GORM writes the file inside the mounted directory. The file then lives on the host and survives container restarts.

### Layer caching trick
Copy `go.mod` and `go.sum` **before** copying the rest of the source, then run `go mod download` as a separate step:

```dockerfile
COPY go.mod go.sum ./
RUN go mod download   # this layer is cached until go.mod/go.sum change
COPY . .
RUN go build ...
```

If only your source files change (not dependencies), Docker reuses the cached `go mod download` layer and rebuilds only the compile step. This significantly speeds up iterative builds.

### .dockerignore
Works like `.gitignore` but for Docker's build context. Prevents unnecessary files from being sent to the Docker daemon:

```
.git
*.db
data/
.env
bin/
```

---

## Key takeaways
- Multi-stage builds = small images, no build tools in production
- SQLite via GORM requires CGO → must install `gcc` and `musl-dev` in the build stage
- SQLite data must live in a Docker volume or it disappears on container restart
- Copy `go.mod`/`go.sum` first to get fast layer caching on rebuilds
- `.dockerignore` keeps the build context lean

---

## Session Q&A

**Q:** What does `-rf` mean in `rm -rf`?
**A:** `-r` = recursive (delete directory and all contents), `-f` = force (no confirmation, ignore missing files). Dangerous — no undo.

**Q:** What does `-p` mean in `mkdir -p`?
**A:** Parents — create missing parent directories in the path, and don't error if the directory already exists. Safe to use habitually.

**Q:** What does `CMD` stand for in a Dockerfile?
**A:** Just "command" — the command to run when the container starts.

**Q:** What is the difference between `CMD` and `ENTRYPOINT`?
**A:** `CMD` can be overridden at runtime (`docker run <image> ./other`). `ENTRYPOINT` is fixed unless you pass `--entrypoint`. For simple apps, `CMD` is the right choice.

**Q:** What does `grep -o` do?
**A:** `-o` = print only the matching part, not the whole line. Used with a pipe to extract a specific field from JSON output.

**Q:** If I change code, does the container auto-update?
**A:** No. Docker bakes your code into the image at build time. You must run `docker compose up --build` to pick up changes. Use `air` or `go run .` during development; Docker is for packaging and deployment.

**Q:** Does `docker compose restart` keep the container running after?
**A:** Yes — it stops then starts the container again. It stays running until you explicitly run `docker compose down` (stop + remove) or `docker compose stop`.

**Q:** Does `docker compose down` stop all containers on the machine?
**A:** No — only the containers defined in the `docker-compose.yml` in the current directory. Use `docker ps` to see all running containers across all projects, and `docker ps -a` to include stopped ones.

**Q:** What other Docker concepts are there to learn?
**A:** The basics covered here are sufficient for now. More concepts come up naturally in later days: Docker networks (Day 24 — Nginx), running a database container (Day 27 — PostgreSQL), and multi-service compose files (Day 28 — Kafka/RabbitMQ).

## Gotchas

- **Go version mismatch**: The `tool` directive in `go.mod` was introduced in Go 1.24. Using `golang:1.23-alpine` in the Dockerfile caused `go mod download` to fail with `unknown directive: tool`. Fix: match the Go image version to your `go.mod` version (`golang:1.25-alpine`).
- **Volume permission error**: Docker can't create the `./data` host directory automatically on macOS. Create it manually with `mkdir -p data` before running `docker compose up`.

## Original Q&A

**Q:** What does `.gitkeep` do?
**A:** It's a convention, not a Git feature. Git doesn't track empty directories — it only tracks files. Placing an empty `.gitkeep` file inside a folder forces Git to include the folder in the repository. The name `.gitkeep` is just a widely-understood convention; the file could be named anything. Once real files exist in the folder, `.gitkeep` can be deleted.

**Q:** How do I add instructions to CLAUDE.md to save notes during `/day` sessions before context is auto-compacted?
**A:** Two parts:
1. Add a `## Session notes` section to `CLAUDE.md` describing the file naming convention (`notes/dayXX-<slug>.md`), what to capture (concepts, Q&A, gotchas), and the note format.
2. Create `.claude/settings.json` with a `PreCompact` hook — a prompt that fires automatically before context compaction, telling Claude to flush the session content to the notes file.

---

## Files to create

| File | Purpose |
|------|---------|
| `Dockerfile` | Multi-stage build (builder + alpine run stage) |
| `docker-compose.yml` | Service definition with port mapping and volume |
| `.dockerignore` | Exclude `.git`, `*.db`, `.env`, `bin/` from build context |

## Verify when done
```bash
docker compose up --build
curl http://localhost:3000/health
curl http://localhost:3000/tasks
# restart the container and verify task data persists
docker compose restart
curl http://localhost:3000/tasks
```
