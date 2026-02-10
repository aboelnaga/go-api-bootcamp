package main

import "github.com/gofiber/fiber/v3"

// import "log"

func main() {
	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(map[string]string{"status": "ok"})
	})

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
