package main

import "github.com/gofiber/fiber/v3"

func startServer() error {
	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	return app.Listen(":3000")
}
