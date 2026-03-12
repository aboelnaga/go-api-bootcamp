# Day 11 – Unit Tests

## Concepts

### Go's built-in `testing` package
No third-party framework needed. Test files end in `_test.go` and test functions are named `TestXxx`:
```go
func TestGetHealth(t *testing.T) { ... }
```
Run all tests: `go test ./...`
Run a specific test: `go test -run TestGetHealth ./...`

### Table-driven tests
Instead of one test per scenario, define a slice of cases and loop over them:
```go
tests := []struct {
    name       string
    method     string
    url        string
    wantStatus int
}{
    {"health ok", "GET", "/health", 200},
    {"not found", "GET", "/tasks/999", 404},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // ...
    })
}
```
This pattern scales cleanly — adding a new case is one line.

### Testing Fiber without a real server
`app.Test(req)` runs a request through the Fiber app in-process. No port is opened, no network involved:
```go
req := httptest.NewRequest("GET", "/health", nil)
resp, err := app.Test(req)
```

### In-memory SQLite for tests
Using `:memory:` as the DSN gives each test run a fresh, isolated database — no leftover data between tests:
```go
db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
db.AutoMigrate(&Task{})
```

### `setupTestApp` helper
Centralizing app creation in a helper avoids repeating setup in every test:
```go
func setupTestApp() *fiber.App {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&Task{})
    app := fiber.New()
    // register routes with this db...
    return app
}
```

### `getAuthToken` helper
For protected routes, a helper that calls `POST /login` and returns the JWT keeps tests clean:
```go
func getAuthToken(app *fiber.App) string {
    body := `{"username":"admin","password":"secret"}`
    req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    // parse token from response...
}
```

## Key takeaways
- Table-driven tests are idiomatic Go — one loop, many scenarios
- `app.Test(req)` tests Fiber handlers without starting a real HTTP server
- `:memory:` SQLite gives each test a clean database with no teardown needed
- Helper functions (`setupTestApp`, `getAuthToken`) eliminate boilerplate across test cases

## Q&A

**Q:** Why use `httptest.NewRequest` instead of `http.NewRequest`?
**A:** `httptest.NewRequest` is designed for testing — it panics on invalid input (making test bugs obvious) and sets sensible defaults.

**Q:** Does `app.Test` run middleware too?
**A:** Yes — it goes through the full Fiber middleware stack, just like a real request.

**Q:** Why `:memory:` instead of a test `.db` file?
**A:** A file-based DB persists between runs and can cause tests to interfere with each other. `:memory:` is always fresh, fast, and automatically cleaned up when the process exits.
