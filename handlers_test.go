package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type paginatedResponse struct {
	Tasks []Task `json:"tasks"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Total int64  `json:"total"`
}

func setupTestApp(t *testing.T) *fiber.App {
	t.Helper() // makes error messages point to the test, not this function

	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	db.AutoMigrate(&Task{})

	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})
	setupRoutes(app)

	return app
}

// readBody is a helper that reads the response body and unmarshals JSON into target.
// For example: var tasks []Task; readBody(t, resp, &tasks)
func readBody(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v\nbody was: %s", err, string(body))
	}
}

func getAuthToken(t *testing.T, app *fiber.App) string {
	t.Helper()
	loginBody := `{"username":"admin","password":"secret"}`
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from /login, got %d", resp.StatusCode)
	}

	var result map[string]any
	readBody(t, resp, &result)
	token, ok := result["token"].(string)
	if !ok {
		t.Fatal("expected token in login response")
	}
	return token
}

func TestLoginHandler(t *testing.T) {
	app := setupTestApp(t)

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "valid credentials return token",
			body:         `{"username":"admin","password":"secret"}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid credentials return 401",
			body:         `{"username":"admin","password":"wrong"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "missing fields return 401",
			body:         `{"username":"admin"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "missing fields return 400",
			body:         ``,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			if tt.expectedCode == http.StatusOK {
				var result map[string]any
				readBody(t, resp, &result)
				if _, ok := result["token"].(string); !ok {
					t.Error("expected token in response")
				}
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	app := setupTestApp(t)
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestCreateTaskHandler uses table-driven tests to check success and failure cases.
// Notice the pattern: each test case defines a name, input body, and expected status code.
func TestCreateTaskHandler(t *testing.T) {
	tests := []struct {
		name         string
		body         string // raw JSON string to send
		expectedCode int
	}{
		{
			name:         "valid task returns 201",
			body:         `{"title":"Learn Go","description":"A test task"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "missing title returns 400",
			body:         `{"description":"no title here"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "empty body returns 400",
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid JSON returns 400",
			body:         `not json at all`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp(t)        // fresh DB each subtest
			token := getAuthToken(t, app) // get a valid token for auth

			req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}
		})
	}
}

// TestCreateTaskIncrementalIDs verifies that task IDs are auto-incremented by the database.
func TestCreateTaskIncrementalIDs(t *testing.T) {
	app := setupTestApp(t)

	// Create 3 tasks and collect their IDs
	titles := []string{"First", "Second", "Third"}
	var ids []float64
	token := getAuthToken(t, app) // get a valid token for auth

	for _, title := range titles {
		body := `{"title":"` + title + `"}`
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201 for task '%s', got %d", title, resp.StatusCode)
		}

		var result map[string]any
		readBody(t, resp, &result)
		task := result["task"].(map[string]any)
		ids = append(ids, task["id"].(float64))
	}

	// Verify IDs are 1, 2, 3
	for i, id := range ids {
		expected := float64(i + 1)
		if id != expected {
			t.Errorf("task %d: expected id %.0f, got %.0f", i+1, expected, id)
		}
	}
}

// TestGetTasks uses a seed function per test case to set up different data scenarios.
// This is a common pattern: seed controls what's in the DB, then you check the response.
func TestGetTasks(t *testing.T) {
	tests := []struct {
		name          string
		seed          func() // populates the DB before the request
		query         string // appended to /tasks, e.g. "?completed=true"
		expectedCode  int
		expectedCount int
	}{
		{
			name:          "empty database returns empty array",
			seed:          func() {},
			query:         "",
			expectedCode:  http.StatusOK,
			expectedCount: 0,
		},
		{
			name: "returns all seeded tasks",
			seed: func() {
				db.Create(&Task{Title: "Task A"})
				db.Create(&Task{Title: "Task B", Completed: true})
			},
			query:         "",
			expectedCode:  http.StatusOK,
			expectedCount: 2,
		},
		{
			name: "filter completed=true",
			seed: func() {
				db.Create(&Task{Title: "Incomplete"})
				db.Create(&Task{Title: "Done", Completed: true})
			},
			query:         "?completed=true",
			expectedCode:  http.StatusOK,
			expectedCount: 1,
		},
		{
			name: "filter completed=false",
			seed: func() {
				db.Create(&Task{Title: "Incomplete"})
				db.Create(&Task{Title: "Done", Completed: true})
			},
			query:         "?completed=false",
			expectedCode:  http.StatusOK,
			expectedCount: 1,
		},
		{
			name:         "invalid completed param returns 400",
			seed:         func() {},
			query:        "?completed=notabool",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "filter page=1&limit=2",
			seed: func() {
				db.Create(&Task{Title: "Incomplete"})
				db.Create(&Task{Title: "Done", Completed: true})
				db.Create(&Task{Title: "Another Incomplete"})
			},
			query:         "?page=1&limit=2",
			expectedCode:  http.StatusOK,
			expectedCount: 2,
		},
		{
			name: "filter page=2&limit=2",
			seed: func() {
				db.Create(&Task{Title: "Incomplete"})
				db.Create(&Task{Title: "Done", Completed: true})
				db.Create(&Task{Title: "Another Incomplete"})
			},
			query:         "?page=2&limit=2",
			expectedCode:  http.StatusOK,
			expectedCount: 1, // only 1 task on the second page
		},
		{
			name:         "invalid page/limit returns 400",
			seed:         func() {},
			query:        "?page=abc&limit=def",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "filter completed=true&page=1&limit=3",
			seed: func() {
				db.Create(&Task{Title: "Incomplete"})
				db.Create(&Task{Title: "Done", Completed: true})
				db.Create(&Task{Title: "Another Incomplete"})
			},
			query:         "?completed=true&page=1&limit=3",
			expectedCode:  http.StatusOK,
			expectedCount: 1, // only 1 task on the first page
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp(t)
			tt.seed() // populate test data

			req, _ := http.NewRequest(http.MethodGet, "/tasks"+tt.query, nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			// Only check the count for successful responses
			if tt.expectedCode == http.StatusOK {
				var result paginatedResponse
				readBody(t, resp, &result)
				if len(result.Tasks) != tt.expectedCount {
					t.Errorf("expected %d tasks, got %d", tt.expectedCount, len(result.Tasks))
				}
			}
		})
	}
}

// TestGetTaskById tests fetching a single task — both found and not found.
func TestGetTaskById(t *testing.T) {
	tests := []struct {
		name         string
		seed         func()
		id           string
		expectedCode int
	}{
		{
			name: "existing task returns 200",
			seed: func() {
				db.Create(&Task{Title: "Test Task"})
			},
			id:           "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "non-existent task returns 404",
			seed:         func() {},
			id:           "999",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid id returns 404",
			seed:         func() {},
			id:           "abc",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp(t)
			tt.seed()

			req, _ := http.NewRequest(http.MethodGet, "/tasks/"+tt.id, nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			// For the success case, verify we got a task back with a title
			if tt.expectedCode == http.StatusOK {
				var task Task
				readBody(t, resp, &task)
				if task.Title == "" {
					t.Error("expected task to have a title")
				}
			}
		})
	}
}

// TestUpdateTask tests PUT /tasks/:id — success, not found, and invalid id.
func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name         string
		seed         func()
		id           string
		body         string
		expectedCode int
	}{
		{
			name: "update existing task returns 200",
			seed: func() {
				db.Create(&Task{Title: "Old Title"})
			},
			id:           "1",
			body:         `{"title":"Updated Title","completed":true}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "non-existent task returns 404",
			seed:         func() {},
			id:           "999",
			body:         `{"title":"Doesn't matter"}`,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid id returns 400",
			seed:         func() {},
			id:           "abc",
			body:         `{"title":"X"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp(t)
			tt.seed()
			token := getAuthToken(t, app)

			req, _ := http.NewRequest(http.MethodPut, "/tasks/"+tt.id, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			// For the success case, verify the task was actually updated
			if tt.expectedCode == http.StatusOK {
				var task Task
				readBody(t, resp, &task)
				if task.Title != "Updated Title" {
					t.Errorf("expected title 'Updated Title', got '%s'", task.Title)
				}
				if !task.Completed {
					t.Error("expected completed to be true")
				}
			}
		})
	}
}

// TestDeleteTask tests DELETE /tasks/:id — success, not found, and invalid id.
func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name         string
		seed         func()
		id           string
		expectedCode int
	}{
		{
			name: "delete existing task returns 200",
			seed: func() {
				db.Create(&Task{Title: "To Delete"})
			},
			id:           "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "non-existent task returns 404",
			seed:         func() {},
			id:           "999",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid id returns 400",
			seed:         func() {},
			id:           "abc",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp(t)
			tt.seed()
			token := getAuthToken(t, app) // get a valid token for auth

			req, _ := http.NewRequest(http.MethodDelete, "/tasks/"+tt.id, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			// Verify the task is actually gone from the database
			if tt.expectedCode == http.StatusOK {
				var count int64
				db.Model(&Task{}).Count(&count)
				if count != 0 {
					t.Errorf("expected 0 tasks after delete, got %d", count)
				}
			}
		})
	}
}
