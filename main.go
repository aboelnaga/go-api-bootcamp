package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})

	if err != nil {
		panic("couldn't connect to DB")
	}

	db.AutoMigrate(&Task{})

	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(func(c fiber.Ctx) error {
		requestId := uuid.New().String()
		c.Set("X-Request-ID", requestId)
		return c.Next()
	})

	setupRoutes(app)

	seedDB(db)

	app.Listen(":3000")
}
