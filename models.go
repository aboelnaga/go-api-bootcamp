package main

import (
	"time"

	"github.com/go-playground/validator/v10"
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
