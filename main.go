package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
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

	setupRoutes(app)

	seedDB(db)

	app.Listen(":3000")
}
