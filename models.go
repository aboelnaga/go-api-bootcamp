package main

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
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

func seedDB(db *gorm.DB) {
	var count int64

	db.Model(&Task{}).Count(&count)

	if count == 0 {
		db.Create(&Task{Title: "Learn Go", Description: "Learn Go"})
		db.Create(&Task{Title: "Build task AAPI", Description: "Build. task API", Completed: true})
	}
}
