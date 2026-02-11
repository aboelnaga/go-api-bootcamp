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

### Day 2 – First Goal

- **Objective**: Have a running Fiber server with a simple health check endpoint.
- **Endpoint**: `GET /health`
- **Behavior**:
  - Returns HTTP `200 OK`
  - Returns JSON body similar to: `{"status": "ok"}`
- **Done when**:
  - `go run .` starts the Fiber server without errors.
  - Visiting `http://localhost:3000/health` (or your chosen port) returns the JSON response.

### Day 3 – Next Goal

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

