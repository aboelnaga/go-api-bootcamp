# Day 16 – Docker

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

## Q&A

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
