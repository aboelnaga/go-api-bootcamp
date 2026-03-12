# Day 02 – Health Check Endpoint

## Concepts

### What is Fiber?
Fiber is a Go web framework inspired by Express.js. It wraps Go's `fasthttp` package and provides a simple, Express-like API for building HTTP servers.

```go
app := fiber.New()
app.Get("/health", healthHandler)
app.Listen(":3000")
```

### Defining a route handler
Every handler in Fiber takes a `fiber.Ctx` and returns an `error`:
```go
func healthHandler(c fiber.Ctx) error {
    return c.JSON(fiber.Map{"status": "ok"})
}
```

### `fiber.Map`
`fiber.Map` is just `map[string]interface{}` — a shorthand for building arbitrary JSON responses without defining a struct.

### Go modules
`go.mod` declares your module name and dependencies. `go mod tidy` downloads missing packages and removes unused ones.

### Running the server
`go run .` compiles and runs all `.go` files in the current directory. The `.` means "this package".

## Key takeaways
- Fiber's minimal API: `fiber.New()` → register routes → `app.Listen()`
- Handler signature is always `func(c fiber.Ctx) error`
- `fiber.Map` is the quick way to return JSON without a dedicated struct
- `go run .` is how you run a multi-file package during development

## Q&A

**Q:** Why does `app.Listen` block?
**A:** It starts an HTTP server loop that waits for connections. Everything after it won't run until the server stops.

**Q:** What port should I use?
**A:** Any unused port above 1024. 3000 is conventional for local development.
