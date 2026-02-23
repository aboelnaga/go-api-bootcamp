package main

import "github.com/gofiber/fiber/v3"

func setupRoutes(app *fiber.App) {
	app.Post("/login", loginHandler)
	app.Get("/health", healthHandler)
	app.Get("/tasks", getTasksHandler)
	app.Get("/tasks/:id", getTaskByIdHandler)

	tasks := app.Group("/tasks", authMiddleware)
	tasks.Post("/", createTaskHandler)
	tasks.Put("/:id", updateTaskHandler)
	tasks.Delete("/:id", deleteTaskHandler)
}
