package main

import "github.com/gofiber/fiber/v3"

func setupRoutes(app *fiber.App) {
	app.Get("/health", healthHandler)
	app.Get("/tasks", getTasksHandler)
	app.Get("/tasks/:id", getTaskByIdHandler)
	app.Post("/tasks", createTaskHandler)
	app.Put("/tasks/:id", updateTaskHandler)
	app.Delete("/tasks/:id", deleteTaskHandler)
}
