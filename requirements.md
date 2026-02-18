## Task Manager CRUD API – Requirements

### Functional Requirements

- **Core entity: Task**
  - `id`: unique identifier (number or UUID)
  - `title`: short text, required
  - `description`: longer text, optional
  - `completed`: boolean flag
  - `createdAt`: timestamp when the task was created

- **Core API endpoints (MVP target)**
  - `POST /tasks`: create a new task
  - `GET /tasks`: list all tasks
  - `GET /tasks/:id`: get a single task by id
  - `PUT /tasks/:id`: update an existing task
  - `DELETE /tasks/:id`: delete a task

### Day 2 – [x] Health check endpoint

- **Objective**: Have a running Fiber server with a simple health check endpoint.
- **Endpoint**: `GET /health`
- **Behavior**:
  - Returns HTTP `200 OK`
  - Returns JSON body similar to: `{"status": "ok"}`
- **Done when**:
  - `go run .` starts the Fiber server without errors.
  - Visiting `http://localhost:3000/health` (or your chosen port) returns the JSON response.

### Day 3 – [x] GET /tasks (in-memory)

- **Objective**: Implement `GET /tasks` that returns a hard-coded list of tasks (no database yet).
- **Endpoint**: `GET /tasks`
- **Behavior**:
  - Returns HTTP `200 OK`
  - Returns JSON array of task objects, e.g.:
    - `[{ "id": 1, "title": "Learn Go", "completed": false }, { "id": 2, "title": "Build Task API", "completed": true }]`
  - Data is stored in memory (a Go slice), not in a database.
- **Done when**:
  - `go run .` starts the server without errors.
  - `curl http://localhost:3000/tasks` returns the JSON array of tasks.

### Day 4 – [x] GET /tasks/:id

- **Objective**: Implement `GET /tasks/:id` to fetch a single task by its `id` from the same in-memory slice.
- **Endpoint**: `GET /tasks/:id`
- **Behavior**:
  - Returns HTTP `200 OK` with the matching task as JSON when `id` exists.
  - Returns HTTP `404 Not Found` with a simple JSON error, e.g. `{"error": "task not found"}`, when `id` does not exist.
- **Done when**:
  - `curl http://localhost:3000/tasks/1` returns the JSON for task with `id = "1"`.
  - `curl http://localhost:3000/tasks/999` returns `404` with the JSON error.

### Day 5 – [x] POST /tasks

- **Objective**: Implement `POST /tasks` to create a new task in memory.
- **Endpoint**: `POST /tasks`
- **Behavior**:
  - Accepts JSON body with at least `title` (string) and optional `description`.
  - Appends a new `Task` to the in-memory `tasks` slice with:
    - A new `id` (e.g. incrementing string `"3"`, `"4"`, ...).
    - `completed` defaulting to `false`.
    - `createdAt` set to the current time.
  - Returns HTTP `201 Created` with the created task as JSON.
- **Done when**:
  - `curl -i -X POST http://localhost:3000/tasks -H "Content-Type: application/json" -d '{"title":"New Task","description":"Test"}'`
    returns `201` and the created task JSON.

### Day 6 – [x] PUT /tasks/:id

- **Objective**: Implement `PUT /tasks/:id` to update an existing task in memory.
- **Endpoint**: `PUT /tasks/:id`
- **Behavior**:
  - Accepts JSON body with updatable fields (e.g. `title`, `description`, `completed`).
  - Finds the task by `id` in the `tasks` slice.
  - Updates the task’s fields while keeping `id` and `createdAt` unchanged.
  - Returns HTTP `200 OK` with the updated task as JSON when `id` exists.
  - Returns HTTP `404 Not Found` with `{"error": "task not found"}` when `id` does not exist.
- **Done when**:
  - `curl -i -X PUT http://localhost:3000/tasks/1 -H "Content-Type: application/json" -d '{"title":"Updated title","completed":true}'`
    returns `200` and the updated task JSON.
  - `curl -i -X PUT http://localhost:3000/tasks/999 -H "Content-Type: application/json" -d '{"title":"X"}'`
    returns `404` with the error JSON.

### Day 7 – [x] DELETE /tasks/:id

- **Objective**: Implement `DELETE /tasks/:id` to delete a task from memory.
- **Endpoint**: `DELETE /tasks/:id`
- **Behavior**:
  - Finds the task by `id` in the `tasks` slice.
  - Removes the task from the slice.
  - Returns HTTP `200 OK` with a confirmation message, e.g. `{"message": "task deleted"}`.
  - Returns HTTP `404 Not Found` with `{"error": "task not found"}` when `id` does not exist.
- **Done when**:
  - `curl -i -X DELETE http://localhost:3000/tasks/1` returns `200` with the confirmation JSON.
  - `curl http://localhost:3000/tasks` no longer includes the deleted task.
  - `curl -i -X DELETE http://localhost:3000/tasks/999` returns `404` with the error JSON.

### Day 8 – [x] Query filtering

- **Objective**: Add query parameter filtering to `GET /tasks` so clients can filter by `completed` status.
- **Endpoint**: `GET /tasks?completed=true` or `GET /tasks?completed=false`
- **Behavior**:
  - If no query parameter is provided, return all tasks (existing behavior).
  - If `?completed=true` is provided, return only tasks where `completed` is `true`.
  - If `?completed=false` is provided, return only tasks where `completed` is `false`.
  - Returns HTTP `200 OK` with a JSON array of matching tasks (may be empty `[]`).
- **Done when**:
  - `curl http://localhost:3000/tasks` returns all tasks.
  - `curl http://localhost:3000/tasks?completed=true` returns only completed tasks.
  - `curl http://localhost:3000/tasks?completed=false` returns only incomplete tasks.

### Day 9 – [x] Refactor project structure

- **Objective**: Refactor the project by splitting `main.go` into separate files with clear responsibilities.
- **Target structure**:
  ```
  helloworld/
  ├── main.go          # App setup and server start only
  ├── models.go        # Task struct and validator struct
  ├── handlers.go      # All route handler functions
  ├── routes.go        # Route registration (maps paths to handlers)
  ```
- **Rules**:
  - All files stay in `package main` (no sub-packages yet).
  - The API behavior must remain exactly the same — no new features.
  - Each file should have a single, clear purpose.
- **Done when**:
  - `go run .` starts the server without errors.
  - All existing curl commands still work identically.
  - `main.go` is short and only contains app config and `app.Listen`.

### Day 10 – [x] SQLite with GORM

- **Objective**: Replace the in-memory slice with a SQLite database using GORM.
- **What you'll learn**: Go ORM basics, database migrations, and persistent storage.
- **Steps**:
  - Install GORM and the SQLite driver (`gorm.io/gorm`, `gorm.io/driver/sqlite`).
  - Configure GORM to auto-migrate the `Task` model on startup.
  - Replace all slice operations in handlers with GORM queries (`db.Create`, `db.Find`, `db.First`, `db.Save`, `db.Delete`).
- **Done when**:
  - `go run .` starts the server and creates a `tasks.db` file.
  - All existing curl commands work identically.
  - Data persists across server restarts.

### Day 11 – [x] Unit tests

- **Objective**: Write unit tests for your API endpoints using Go's `testing` package and `net/http/httptest`.
- **What you'll learn**: Go testing conventions, table-driven tests, and HTTP testing.
- **Steps**:
  - Create `handlers_test.go` with tests for each endpoint.
  - Use table-driven tests (a slice of test cases) for different scenarios (success, not found, bad input).
  - Test at least: `GET /health`, `POST /tasks` (valid + invalid), `GET /tasks/:id` (found + not found), `DELETE /tasks/:id`.
- **Done when**:
  - `go test ./...` passes with all tests green.
  - Each endpoint has at least one success and one failure test case.

### Day 12 – [x] Logging & CORS middleware

- **Objective**: Add middleware for request logging and CORS.
- **What you'll learn**: How Fiber middleware works, request lifecycle, and cross-origin resource sharing.
- **Steps**:
  - Add Fiber's built-in `logger` middleware to log every request (method, path, status, duration).
  - Add Fiber's `cors` middleware so the API can be called from a browser frontend.
  - Create a custom middleware that adds a `X-Request-ID` header to every response (use `uuid` package).
- **Done when**:
  - Every request prints a log line in the server terminal.
  - Responses include `Access-Control-Allow-Origin` and `X-Request-ID` headers.

### Day 13 – [x] Environment config

- **Objective**: Add environment-based configuration using a `.env` file.
- **What you'll learn**: The `godotenv` package, `os.Getenv`, and separating config from code.
- **Steps**:
  - Install `github.com/joho/godotenv`.
  - Move hardcoded values (port, database path) into a `.env` file.
  - Create a `config.go` that loads and exposes configuration values.
  - Add `.env` to `.gitignore` and create a `.env.example` with placeholder values.
- **Done when**:
  - Changing `PORT=4000` in `.env` makes the server start on port 4000.
  - The app works with default values if `.env` is missing.

### Day 14 – [ ] Pagination

- **Objective**: Add pagination to `GET /tasks`.
- **What you'll learn**: Query parameters, GORM `Offset`/`Limit`, and pagination response patterns.
- **Steps**:
  - Accept `?page=1&limit=10` query parameters (default: page 1, limit 10).
  - Use GORM's `.Offset().Limit()` to fetch only the requested page.
  - Return a JSON response with `tasks`, `page`, `limit`, and `total` count.
- **Done when**:
  - `curl http://localhost:3000/tasks?page=1&limit=1` returns 1 task and correct pagination metadata.
  - `curl http://localhost:3000/tasks` still works with default pagination.

### Day 15 – [ ] JWT authentication

- **Objective**: Add JWT authentication to protect write endpoints.
- **What you'll learn**: JSON Web Tokens, auth middleware, and protecting routes.
- **Steps**:
  - Add a `POST /login` endpoint that accepts a hardcoded username/password and returns a JWT token.
  - Create an auth middleware that validates the JWT from the `Authorization: Bearer <token>` header.
  - Protect `POST`, `PUT`, and `DELETE /tasks` routes — `GET` routes remain public.
- **Done when**:
  - `POST /tasks` without a token returns `401 Unauthorized`.
  - `POST /login` returns a JWT token.
  - `POST /tasks` with a valid `Authorization` header succeeds.

### Day 16 – [ ] Docker

- **Objective**: Dockerize the application.
- **What you'll learn**: Dockerfile, multi-stage builds, and container basics.
- **Steps**:
  - Create a `Dockerfile` with a multi-stage build (build stage with `golang` image, run stage with `alpine`).
  - Create a `docker-compose.yml` for easy local development.
  - Ensure the SQLite database file is persisted using a Docker volume.
- **Done when**:
  - `docker compose up` starts the API and all curl commands work against it.
  - Restarting the container preserves task data.

### Day 17 – [ ] Swagger docs

- **Objective**: Add Swagger/OpenAPI documentation to the API.
- **What you'll learn**: API documentation with `swaggo/swag`, annotation comments, and auto-generated docs.
- **Steps**:
  - Install `swaggo/swag` and `swaggo/fiber-swagger`.
  - Add Swagger annotation comments to each handler function.
  - Serve the Swagger UI at `/swagger/*`.
- **Done when**:
  - Visiting `http://localhost:3000/swagger/index.html` shows interactive API docs.
  - All endpoints are documented with request/response examples.

### Day 18 – [ ] Deploy to cloud

- **Objective**: Deploy the API to a free cloud platform.
- **What you'll learn**: Cloud deployment, environment variables in production, and basic DevOps.
- **Steps**:
  - Choose a platform (Railway, Render, or Fly.io — all have free tiers).
  - Configure production environment variables.
  - Deploy and verify all endpoints work with the public URL.
- **Done when**:
  - The API is accessible at a public URL.
  - All curl commands work against the deployed version.
