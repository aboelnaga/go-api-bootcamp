# Day 08 – Query Filtering on GET /tasks

## Concepts

### Reading query parameters in Fiber
```go
c.Query("completed")          // returns "" if not provided
c.Query("page", "1")          // returns "1" as default if not provided
```

### Handling optional query parameters
Always check if the parameter is present before parsing:
```go
query := c.Query("completed")
if query != "" {
    completed, err := strconv.ParseBool(query)
    ...
}
```

### `strconv.ParseBool`
Parses `"true"`, `"false"`, `"1"`, `"0"` and more. Returns an error for anything else (e.g. `"yes"`, `"invalid"`, `""`).
Better than manual if/else string comparison.

## Key takeaways
- Check for empty string before parsing — `ParseBool("")` returns an error
- Go naming convention: local variables use `camelCase`, not `PascalCase` (uppercase = exported/public)
- Remove debug `fmt.Println` logs once a feature is confirmed working
- Always check errors before using the result — don't log or use values when error may be non-nil

## Q&A

**Q:** Will logs from `fmt.Println` appear in curl output?
**A:** No. Logs go to the server terminal (stdout). Curl only receives what you send via `c.JSON()` or similar response methods.

**Q:** Where do server logs appear?
**A:** In the same terminal window where you ran `go run .` or `air`.
