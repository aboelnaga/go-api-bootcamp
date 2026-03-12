# Day 03 – GET /tasks (In-Memory)

## Concepts

### Go structs
A struct defines the shape of a data record. JSON field names are controlled with struct tags:
```go
type Task struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"createdAt"`
}
```

### Package-level variables
Declaring the tasks slice at the package level makes it accessible to all handlers in the same package:
```go
var tasks = []Task{
    {ID: "1", Title: "Learn Go", Completed: false},
    {ID: "2", Title: "Build Task API", Completed: true},
}
```

### Returning a slice as JSON
`c.JSON()` serializes any Go value to JSON and sets `Content-Type: application/json`:
```go
func getTasksHandler(c fiber.Ctx) error {
    return c.JSON(tasks)
}
```
An empty slice `[]Task{}` serializes to `[]`, not `null` — important for API consumers.

### In-memory storage
Data lives in a slice in the process's memory. It's fast and simple, but resets every time the server restarts. This is fine for learning; later replaced with a database.

## Key takeaways
- Struct tags (`json:"..."`) control how Go fields serialize to/from JSON
- Package-level slice acts as a simple in-memory store
- `c.JSON(value)` handles serialization and the correct Content-Type header
- In-memory state is lost on restart — acceptable for now, replaced with SQLite on Day 10

## Q&A

**Q:** Why use `[]Task{}` instead of `nil` as the initial value?
**A:** A nil slice serializes to `null` in JSON. An empty slice `[]Task{}` serializes to `[]`, which is what API clients expect.

**Q:** Can multiple handlers access the same `tasks` slice?
**A:** Yes — package-level variables are accessible to all functions in the same package.
