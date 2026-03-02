# Day 06 – PUT /tasks/:id

## Concepts

### Route parameter syntax
Fiber requires a `/` before the parameter name:
```go
app.Put("/tasks/:id", ...)  // correct
app.Put("/tasks:id", ...)   // wrong — Fiber won't match the route
```

### Value vs pointer in Go structs
When a function returns a struct by value, you get a **copy**. Modifying the copy does not affect the original.
```go
task, _ := getTaskById(id)   // task is a COPY
c.Bind().Body(&task)         // updates the copy only — original slice unchanged
```

### Updating a slice element in-place
To update the actual element, operate on the slice index directly:
```go
c.Bind().Body(&tasks[index]) // writes into the slice element directly
```

### The & operator
`&` is the address-of operator. It gives the callee a pointer to the variable so it can modify the original:
```go
c.Bind().Body(&task)    // Body receives a *Task and can write into it
c.Bind().Body(task)     // Body receives a copy — changes are lost
```

### Protecting fields on update
`c.Bind().Body` overwrites every matching JSON field. Save and restore fields that must not change:
```go
originalID := tasks[index].ID
originalCreatedAt := tasks[index].CreatedAt
c.Bind().Body(&tasks[index])
tasks[index].ID = originalID
tasks[index].CreatedAt = originalCreatedAt
```

## Key takeaways
- Returning a struct from a function gives a value copy — changes don't propagate back
- Always use `&` when passing a target to a binding/unmarshalling function
- Protect immutable fields (id, createdAt) by saving and restoring them after binding
- A new helper `getTaskIndexById` was added instead of changing `getTaskById` (which was used elsewhere returning a copy — acceptable for reads)

## Q&A

**Q:** Why not return a pointer or index from `getTaskById`?
**A:** Returning a `*Task` from a `[]Task` slice is fragile — if the slice grows and Go reallocates it, the pointer may be invalid. Returning an index is the safest approach with a value slice.

**Q:** Why create a new helper instead of changing `getTaskById`?
**A:** `getTaskById` was already used in `GET /tasks/:id` where a copy is fine. No need to change what already works.

**Q:** What happens if & is removed from `c.Bind().Body`?
**A:** Go gives a compile error — `Body()` expects a pointer to write into.
