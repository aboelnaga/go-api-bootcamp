package main

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

// import "log"

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
}

func main() {
	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(map[string]string{"status": "ok"})
	})

	app.Get("/tasks", func(c fiber.Ctx) error {
		tasks := []Task{
			{ID: "1", Title: "Learn Go", Description: "Learn Go", Completed: false, CreatedAt: time.Now()},
			{ID: "2", Title: "Build Task API", Description: "Build Task API", Completed: true, CreatedAt: time.Now()},
		}
		return c.JSON(tasks)
	})

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
