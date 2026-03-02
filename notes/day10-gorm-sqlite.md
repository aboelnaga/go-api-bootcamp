# Day 10 – SQLite with GORM

## Concepts

### What GORM is
An ORM (Object-Relational Mapper) — lets you use Go structs to interact with a database instead of raw SQL.

### Setup
```go
db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
db.AutoMigrate(&Task{})   // creates the table from your struct
```
`tasks.db` is a file created in the project root. Delete it and GORM recreates it on next startup.

### GORM is database-agnostic
Only the driver line changes. Handler code stays the same for SQLite, PostgreSQL, MySQL, etc:
```go
gorm.Open(sqlite.Open("tasks.db"), ...)        // SQLite
gorm.Open(postgres.Open("host=..."), ...)      // PostgreSQL
gorm.Open(mysql.Open("user:pass@..."), ...)    // MySQL
```

### Common GORM operations
```go
db.Find(&tasks)                    // SELECT all
db.Where("completed = ?", true).Find(&tasks)  // SELECT with condition
db.First(&task, id)                // SELECT by primary key
db.Create(&task)                   // INSERT
db.Save(&task)                     // UPDATE all fields
db.Delete(&Task{}, id)             // DELETE by primary key
result.Error                       // check for errors
result.RowsAffected                // check how many rows were affected
```

### Error handling
GORM returns `*gorm.DB` for chaining. Extract the error at the end:
```go
if err := db.First(&task, id).Error; err != nil { ... }
```

### SQL injection prevention
Use `?` placeholders — never concatenate user input into queries:
```go
db.Where("id = ?", id)    // safe — parameterized
db.Where("id = " + id)    // DANGEROUS — SQL injection
```

### Global `db` variable
Declared in `models.go` at package level:
```go
var db *gorm.DB
```
All handlers can access it directly since they share `package main`.

## Key takeaways
- `tasks.db` is just a file — data persists across restarts
- `db.AutoMigrate` reads your struct and creates/updates the table automatically
- Replace all slice operations with GORM equivalents when migrating
- `strconv.ParseUint` should validate ID params before passing to GORM
- `gorm.ErrRecordNotFound` is returned by `First` when no row matches

## Bugs encountered and fixed

### Bug 1: Variable shadowing with `:=`
```go
// main.go — WRONG: creates local db, global db stays nil
db, err := gorm.Open(...)

// CORRECT: assigns to the global db
var err error
db, err = gorm.Open(...)
```
Using `:=` when `db` is already declared at package level creates a new local variable that shadows the global. The global stays `nil`, causing a nil pointer panic in handlers.

### Bug 2: Double pointer with `new(Task)`
```go
task := new(Task)       // task is already *Task
db.Create(&task)        // &task is **Task — GORM can't update ID field

// Fix: use var instead
var task Task
db.Create(&task)        // &task is *Task — correct
```

### Bug 3: Old helpers still referencing `tasks` slice
After migrating to GORM, the old `getTaskById` and `getTaskIndexById` helpers still referenced the removed `tasks` global. They were deleted and replaced with inline GORM queries.

### Bug 4: Missing nil check for `db` in handlers
Handlers crashed because `db` was nil (caused by Bug 1 above). After fixing the assignment, all handlers work correctly.

## Q&A

**Q:** Where is the SQLite database file created?
**A:** In the project root — the same directory where you run `go run .`.

**Q:** Is GORM coupled to SQLite?
**A:** No. Swap the driver import and `gorm.Open(...)` line. All handler code stays the same.

**Q:** Should I handle errors from `db.Create`?
**A:** Yes. Any GORM operation can fail. Always check `.Error` and return a proper HTTP error response.

**Q:** For DELETE, should I fetch first then delete, or delete directly?
**A:** Direct delete + check `RowsAffected` is more efficient (1 query vs 2). If `RowsAffected == 0`, return 404.

**Q:** When should I use `db.First` vs `db.Where`?
**A:** Use `db.First(&task, id)` for primary key lookups (shortest). Use `db.Where("col = ?", val).First(&task)` for any other condition.
