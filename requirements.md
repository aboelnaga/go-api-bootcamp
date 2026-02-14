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

### Day 4 – Next Goal

- **Objective**: Implement `GET /tasks/:id` to fetch a single task by its `id` from the same in-memory slice.
- **Endpoint**: `GET /tasks/:id`
- **Behavior**:
  - Returns HTTP `200 OK` with the matching task as JSON when `id` exists.
  - Returns HTTP `404 Not Found` with a simple JSON error, e.g. `{"error": "task not found"}`, when `id` does not exist.
- **Done when**:
  - `curl http://localhost:3000/tasks/1` returns the JSON for task with `id = "1"`.
  - `curl http://localhost:3000/tasks/999` returns `404` with the JSON error.

### Day 5 – Goal

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

### Day 6 – Next Goal

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

### Day 7 – Next Goal

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

### Day 8 – Next Goal

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
