# Day 04 – GET /tasks/:id

## Concepts

### Route parameters
The `:id` syntax declares a named URL segment. Fiber captures its value at runtime:
```go
app.Get("/tasks/:id", getTaskHandler)

func getTaskHandler(c fiber.Ctx) error {
    id := c.Params("id")
    // id is a string, e.g. "1"
}
```

### Searching a slice
Go has no built-in "find by field" — you loop manually:
```go
func getTaskById(id string) (Task, bool) {
    for _, task := range tasks {
        if task.ID == id {
            return task, true
        }
    }
    return Task{}, false
}
```
Returning a `(value, bool)` pair is the Go idiom for "found or not found".

### 404 Not Found
When an ID doesn't exist, return 404 with a JSON error body:
```go
task, found := getTaskById(id)
if !found {
    return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
}
return c.JSON(task)
```

### Multiple return values
Go functions can return multiple values. The `ok` / `found` pattern is idiomatic for lookups:
```go
value, ok := someMap[key]
task, found := getTaskById(id)
```

## Key takeaways
- `c.Params("id")` retrieves a named URL segment as a string
- The `(value, bool)` return pattern is idiomatic for lookups that may fail
- `c.Status(code).JSON(...)` sets both the status code and the response body
- Route params are always strings — parse them if you need a number

## Q&A

**Q:** What happens if I call `c.Params("id")` on a route without `:id`?
**A:** It returns an empty string `""`.

**Q:** Why return `Task{}` (zero value) when not found instead of a pointer?
**A:** Returning a zero value is idiomatic when paired with a `bool`. The caller checks the bool before using the value.
