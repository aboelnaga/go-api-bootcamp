# Day 14 – Pagination

## What we built

| File | Change |
|------|--------|
| `handlers.go` | `getTasksHandler` updated with `page`, `limit`, `offset`, `total` |
| `handlers_test.go` | Added `paginatedResponse` struct, updated `TestGetTasks`, added pagination test cases |

---

## How pagination works

Client sends:
```
GET /tasks?page=2&limit=10
```

Server responds:
```json
{
  "tasks": [...],
  "page": 2,
  "limit": 10,
  "total": 47
}
```

Frontend uses `total` to calculate total pages:
```js
totalPages = Math.ceil(total / limit)  // ceil(47/10) = 5
```

Without `total` the client has no idea how many pages exist.

---

## The offset formula

```
offset = (page - 1) * limit
```

| page | limit | offset | rows fetched |
|------|-------|--------|--------------|
| 1    | 10    | 0      | rows 1–10    |
| 2    | 10    | 10     | rows 11–20   |
| 3    | 10    | 20     | rows 21–30   |

---

## GORM query — order matters

`.Offset()` and `.Limit()` must come **before** `.Find()`:

```go
// correct
db.Offset(offset).Limit(limit).Find(&tasks)

// wrong — Find executes immediately, Offset/Limit are ignored
db.Find(&tasks).Offset(offset).Limit(limit)
```

---

## Counting with a filter

When `completed` filter is active, `Count` must include the same `Where`:

```go
// correct — counts only filtered tasks
db.Model(&Task{}).Where("completed = ?", completed).Count(&total)

// wrong — counts all tasks, ignores filter
db.Model(&Task{}).Count(&total)
```

`Model(&Task{})` is required so GORM knows which table to count from.

---

## Reading query params with defaults

`c.Query("page", "1")` returns `"1"` if `?page` is absent or empty.
This prevents the empty-string parse failure but does NOT prevent bad input like `?page=abc`.
Error handling is still required:

```go
page, err := strconv.Atoi(c.Query("page", "1"))
if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
}

limit, err := strconv.Atoi(c.Query("limit", "10"))  // err reused — limit is new
if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
}
```

**Reusing `err`**: `:=` allows reuse when at least one variable on the left is new.

---

## Testing paginated responses

Response shape changed from `[]Task` to an object, so tests need a matching struct:

```go
type paginatedResponse struct {
    Tasks []Task `json:"tasks"`
    Page  int    `json:"page"`
    Limit int    `json:"limit"`
    Total int64  `json:"total"`
}
```

Use it with the existing generic `readBody` helper:

```go
var result paginatedResponse
readBody(t, resp, &result)
if len(result.Tasks) != tt.expectedCount { ... }
```

No new helper needed — `readBody` accepts `any`.

---

## Go tips from this session

### Testing a single Go expression quickly

| Option | How |
|--------|-----|
| Go Playground | [play.golang.org](https://go.dev/play/) — paste and run |
| Scratch file | `go run /tmp/scratch.go` |
| Quick test | `go test -v -run TestScratch` in project |

### git: untrack a previously committed file

```bash
git rm --cached path/to/file
```

`--cached` stops git tracking it but leaves the file on disk.
After this, `.gitignore` takes effect for that file.
