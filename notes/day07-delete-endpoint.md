# Day 07 – DELETE /tasks/:id

## Concepts

### Removing an element from a slice
Go has no built-in remove. The classic pattern uses `append`:
```go
slice = append(slice[:i], slice[i+1:]...)
```
The `...` unpacks the second slice so `append` treats each element individually.

### `slices.Delete` (Go 1.21+)
The standard library `slices` package provides a cleaner API:
```go
tasks = slices.Delete(tasks, index, index+1)
```
Removes elements from `index` up to (not including) `index+1` — effectively one element.
The IDE auto-imports `"slices"` when you use it.

## Key takeaways
- `slices.Delete` is the modern, preferred way to remove slice elements in Go 1.21+
- `append` trick is still valid but error-prone (off-by-one mistakes)
- Re-using `getTaskIndexById` from Day 6 keeps handlers DRY

## Q&A

**Q:** Is `go get` project-level or machine-level?
**A:** Project level. It downloads packages and records them in `go.mod` and `go.sum` — similar to `npm install` in Node.js. Not installed globally.
