# Day 01 – Hello, World & Project Setup

## Concepts

### Go modules
A Go module is a collection of packages managed together. `go mod init` creates a `go.mod` file that declares the module name and Go version:
```bash
go mod init example.com/helloworld
```
The module name is used for imports within the project. For a single `package main` project it doesn't matter much, but it becomes important once you add sub-packages.

### `package main` and `func main()`
Every executable Go program must have exactly one `package main` with a `func main()`. That's the entry point:
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### `go run` vs `go build`
- `go run main.go` — compiles and runs in one step, no binary left behind. Good for development.
- `go build -o bin/app .` — compiles to a binary you can distribute and run directly.

### The `fmt` package
`fmt` is Go's standard formatting package. `fmt.Println` writes to stdout with a newline. It's the simplest way to verify your program runs.

### Writing a test
Go's `testing` package is built in — no third-party library needed:
```go
func TestSomething(t *testing.T) {
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```
Run all tests with `go test ./...`.

### `go mod tidy`
Downloads missing dependencies listed in `go.mod` and removes unused ones. Run it after adding or removing imports.

## Key takeaways
- Every Go executable starts from `package main` + `func main()`
- `go.mod` is the project's dependency manifest — always commit it
- `go run .` runs the whole package; `go build` produces a binary
- Tests live in `_test.go` files and use the built-in `testing` package

## Q&A

**Q:** Why `example.com/helloworld` as the module name?
**A:** It's a placeholder — Go module names are conventionally a URL-like path, but they don't need to resolve to a real URL unless you're publishing the module.

**Q:** Do I need to install a test framework?
**A:** No. Go's `testing` package covers unit tests, benchmarks, and examples out of the box.
