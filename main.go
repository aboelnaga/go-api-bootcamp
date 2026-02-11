package main

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

var tasks = []Task{
	{ID: "1", Title: "Learn Go", Description: "Learn Go", Completed: false, CreatedAt: time.Now()},
	{ID: "2", Title: "Build Task API", Description: "Build Task API", Completed: true, CreatedAt: time.Now()},
}

func main() {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	getTaskById := func(id string) (Task, error) {
		for _, task := range tasks {
			if task.ID == id {
				return task, nil
			}
		}
		return Task{}, fmt.Errorf("task not found")
	}

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(map[string]string{"status": "ok"})
	})

	app.Get("/tasks", func(c fiber.Ctx) error {
		return c.JSON(tasks)
	})

	app.Get("/tasks/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		task, err := getTaskById(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		return c.JSON(task)
	})

	app.Post("/tasks", func(c fiber.Ctx) error {
		task := new(Task)

		if err := c.Bind().Body(task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		task.ID = fmt.Sprintf("%d", len(tasks)+1)
		task.CreatedAt = time.Now()
		tasks = append(tasks, *task)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully", "task": task})
	})

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
