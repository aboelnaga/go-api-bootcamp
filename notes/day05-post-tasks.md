# Day 05 – POST /tasks

## Concepts

### Binding a request body
`c.Bind().Body(&target)` deserializes the JSON request body into a Go struct. The `&` is required — Body needs a pointer to write into:
```go
var input Task
if err := c.Bind().Body(&input); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
}
```

### Input validation with struct tags
The `go-playground/validator` package reads `validate` struct tags to enforce rules:
```go
type Task struct {
    Title       string `json:"title" validate:"required"`
    Description string `json:"description"`
}
```
Common rules: `required`, `min=N`, `max=N`, `email`, `oneof=a b`.

### Generating IDs
Without a database, a simple auto-increment counter works:
```go
var nextID = 3  // start after seed data

newTask := Task{
    ID:        strconv.Itoa(nextID),
    Title:     input.Title,
    Completed: false,
    CreatedAt: time.Now(),
}
nextID++
```

### Appending to a slice
`append` returns a new slice (or the same one if there's capacity). Always reassign:
```go
tasks = append(tasks, newTask)
```

### 201 Created
A successful resource creation should return `201 Created`, not `200 OK`:
```go
return c.Status(fiber.StatusCreated).JSON(newTask)
```

## Key takeaways
- `c.Bind().Body(&input)` requires `&` — it writes into the pointed-to struct
- Always validate input at the boundary; don't trust what clients send
- Return `201 Created` (not `200`) when a new resource is created
- `append` must be reassigned: `tasks = append(tasks, item)`

## Q&A

**Q:** Why does `c.Bind().Body` need `&`?
**A:** It takes a `interface{}` (or `any`) and uses reflection to write into it. Without `&`, it gets a copy and changes are lost.

**Q:** What if the client sends extra fields not in the struct?
**A:** They're silently ignored during deserialization — only fields matching struct tags are read.

**Q:** Why not use the same ID type as the client-sent data?
**A:** The server controls ID generation — clients should never set IDs directly.
