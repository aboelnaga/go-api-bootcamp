# Day 13 – Environment Config Notes

## What we did
- Installed `github.com/joho/godotenv`
- Created `.env` and `.env.example` with `PORT` and `DB_PATH`
- Created `config.go` with a `Config` struct and `LoadConfig()` function
- Updated `main.go` to use `cfg.DBPath` and `cfg.Port`
- Updated `.gitignore`

## Key concepts

### godotenv.Load()
- Reads `.env` file and sets OS environment variables
- Returns an error if file is missing — but ignore it, that's intentional
- Does NOT override env vars already set in the OS (OS always wins)

### os.Getenv
- Reads a single environment variable by name
- Returns empty string `""` if not set — use that to apply defaults

### Default values pattern
```go
port := os.Getenv("PORT")
if port == "" {
    port = "3000"
}
```

### Override from terminal (without touching .env)
```bash
PORT=4000 go run .        # inline, one-off
export PORT=4000          # for the whole terminal session
unset PORT                # clear it
```

### Multiple environments
```go
godotenv.Load(".env.local", ".env")  // .env.local wins if both define the same key
```

## .gitignore rules
- `.env` → gitignored (may contain secrets)
- `.env.example` → committed (template for teammates)
- `*.db`, `*.db-shm`, `*.db-wal` → gitignored (SQLite files)
- `bin/` → gitignored (compiled binaries)
- `tmp/` → gitignored (air live-reload output)

## go.mod concepts
- **First require block** — direct dependencies (packages your code imports)
- **Second require block `// indirect`** — transitive dependencies (deps of your deps)
- `go get` adds packages but doesn't clean up organization
- `go mod tidy` separates direct vs indirect and removes unused entries
- `tool` directive (Go 1.24+) — tracks CLI tools like `air` tied to the project
- Both `go.mod` and `go.sum` should be committed; never manually edit `go.sum`

## Build-time variables (bonus)
```bash
go build -ldflags "-X main.Version=1.2.3" -o bin/app .
```
Bakes a value into the binary at compile time — useful for version strings, not secrets.
