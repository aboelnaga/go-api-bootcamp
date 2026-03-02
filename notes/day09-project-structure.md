# Day 09 – Refactor Project Structure

## Concepts

### Multiple files in the same package
All `.go` files in the same directory with `package main` form one package. They can see each other's types, functions, and variables without importing each other.
```
go run .     ← builds ALL .go files in the directory together
go build .   ← same
```

### File responsibilities
| File | Contents |
|------|---------|
| `main.go` | App config, middleware, `app.Listen` only |
| `models.go` | Structs, global variables, seed data |
| `handlers.go` | Named handler functions + helpers |
| `routes.go` | `setupRoutes(app)` — maps paths to handlers |

### Package-level functions vs anonymous functions
Inside `main()`, you can use `:=` for closures:
```go
getTaskById := func(id string) (Task, error) { ... }  // ok inside main()
```
At package level, use `func` declarations:
```go
func getTaskById(id string) (Task, error) { ... }     // correct at package level
var getTaskById := func(...) { ... }                   // INVALID at package level
```

### Imports per file
Each file declares its own imports for the packages it directly uses. There is no "import file into another file" in Go — packages are the unit of sharing.

### Naming: exported vs unexported
- Uppercase first letter = **exported** (visible outside the package) → used for public API
- Lowercase first letter = **unexported** (package-private) → used for internal helpers

All names in this project start lowercase since everything is `package main`.

## Key takeaways
- Same package = same namespace, no cross-file imports needed
- `var x := func()` is invalid at package level — use `func x()` instead
- Each file needs its own imports for packages it uses directly
- Named functions in `handlers.go` can be referenced from `routes.go` because they share `package main`

## Q&A

**Q:** Do I need to import other files in the project?
**A:** No. Files in the same package automatically see each other. Only import external packages (stdlib or third-party).

**Q:** Do I need to change function names when moving from inline to named?
**A:** You need to give them names (they were anonymous before). Convention: `getTasksHandler`, `createTaskHandler`, etc.
