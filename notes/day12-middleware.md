# Day 12 – Logging & CORS Middleware

## Concepts

### What is Middleware?
Middleware is code that runs **between** receiving a request and executing your handler. It forms a pipeline:

```
Request → [Logger] → [CORS] → [Request-ID] → Your Handler → Response
```

Each middleware can inspect/modify the request, call the next handler, then inspect/modify the response. Order matters — middleware runs in the order it is registered.

### Logger Middleware
Fiber's built-in logger logs every request: method, path, status code, and duration.
- Import: `github.com/gofiber/fiber/v3/middleware/logger`
- Usage: `app.Use(logger.New())`
- `app.Use()` registers middleware that applies to **all routes**

### CORS Middleware
CORS (Cross-Origin Resource Sharing) is a browser security feature. Browsers block JavaScript from making requests to a different domain unless the server explicitly allows it via CORS headers.
- The middleware only adds `Access-Control-Allow-Origin` when the request includes an `Origin` header
- Plain `curl` doesn't send `Origin`, so you won't see CORS headers unless you add `-H "Origin: http://example.com"`
- Import: `github.com/gofiber/fiber/v3/middleware/cors`
- Usage: `app.Use(cors.New())` — default config allows all origins (`*`)

### Custom Middleware
A Fiber v3 middleware is just a function that:
1. Takes a `fiber.Ctx`
2. Does something (modify request/response)
3. Calls `c.Next()` to pass control forward
4. Returns an error

```go
app.Use(func(c fiber.Ctx) error {
    c.Set("X-Request-ID", uuid.New().String())
    return c.Next()
})
```

`c.Set()` adds a header to the response. `c.Next()` is critical — without it the request never reaches your handler.

## Key Takeaways
- Register middleware with `app.Use()` **before** `setupRoutes()` so it wraps all routes
- Middleware order matters — logger first so it captures everything
- CORS headers only appear when an `Origin` header is present in the request
- Writing a custom middleware is just writing a regular function with `c.Next()` inside

## Q&A

**Q:** I got `not enough arguments in call to logger.New` — why?
**A:** Wrong import. `gorm.io/gorm/logger` is GORM's logger, not Fiber's. The correct import is `github.com/gofiber/fiber/v3/middleware/logger`.

**Q:** What if I need both `gorm/logger` and Fiber's `logger` in the same file?
**A:** Use an **alias** to rename one:
```go
import (
    "github.com/gofiber/fiber/v3/middleware/logger"
    gormlogger "gorm.io/gorm/logger"
)
```

**Q:** There's no `Access-Control-Allow-Origin` header in the response — is CORS working?
**A:** Yes, it's working correctly. CORS headers are only added when the request includes an `Origin` header. Simulate a browser request with: `curl -v -H "Origin: http://example.com" http://localhost:3000/health`

**Q:** What do `-v`, `-i`, `-X` mean in curl?
**A:**
- `-v` — Verbose: shows full request AND response including all headers
- `-i` — Include headers: shows response headers + body (less noisy than `-v`)
- `-X` — Method: specifies HTTP method (`GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`)

**Q:** Can I combine curl flags like `-Xi`?
**A:** You can combine flags that don't take a value (e.g., `-si`, `-siv`). Flags that take a value must come last in a combined group: `-siX POST` works, but `-Xsi` doesn't because curl reads `si` as the method name. The URL should always be last.

**Q:** Other useful curl flags?
**A:**
- `-H` — Add a request header
- `-d` — Send a request body (implies POST)
- `-s` — Silent mode (hides progress bar)
- `-L` — Follow redirects
- `-o` — Save response to file
- `-w` — Write-out specific info (e.g., status code)

**Q:** Does `CLAUDE.md` need to be inside `.claude/` folder?
**A:** No — `CLAUDE.md` belongs in the **project root**. The `.claude/` folder is for other things (commands, settings). Claude Code auto-reads `CLAUDE.md` from the root at the start of every conversation.

**Q:** What is `.claude/settings.json`?
**A:** A project-level config file for Claude Code. It lets you auto-approve tools (so you're not prompted every time), deny dangerous commands, and configure MCP servers. There's also a global version at `~/.claude/settings.json`.

**Q:** Can we save the instructor prompt as a command?
**A:** Yes — Claude Code supports custom slash commands via `.claude/commands/`. A markdown file at `.claude/commands/day.md` becomes `/day`. Using `$ARGUMENTS` in the file lets you pass a day name: `/day Day 13: Environment config`.
