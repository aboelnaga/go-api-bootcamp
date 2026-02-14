package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	setupRoutes(app)

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
