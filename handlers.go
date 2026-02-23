package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func loginHandler(c fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Username == "admin" && req.Password == "secret" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
		}
		return c.JSON(fiber.Map{"token": tokenString})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
}

func healthHandler(c fiber.Ctx) error {
	return c.JSON(map[string]string{"status": "ok"})
}

func getTasksHandler(c fiber.Ctx) error {
	var tasks []Task
	var total int64

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	offset := (page - 1) * limit

	if c.Query("completed") != "" {
		completed, err := strconv.ParseBool(c.Query("completed"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameter 'completed'"})
		}
		db.Model(&Task{}).Where("completed = ?", completed).Count(&total)
		db.Where("completed = ?", completed).Offset(offset).Limit(limit).Find(&tasks)
	} else {
		db.Model(&Task{}).Count(&total)
		db.Offset(offset).Limit(limit).Find(&tasks)
	}

	return c.JSON(fiber.Map{
		"tasks": tasks,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func getTaskByIdHandler(c fiber.Ctx) error {
	id := c.Params("id")

	var task Task
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}
	return c.JSON(task)
}

func createTaskHandler(c fiber.Ctx) error {
	task := new(Task)

	if err := c.Bind().Body(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := db.Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Task created successfully", "task": task})
}

func updateTaskHandler(c fiber.Ctx) error {
	id := c.Params("id")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID format"})
	}

	var task Task
	if err := db.First(&task, uint(idUint)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	if err := c.Bind().Body(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	if err := db.Save(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}

	return c.Status(fiber.StatusOK).JSON(task)
}

func deleteTaskHandler(c fiber.Ctx) error {
	id := c.Params("id")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID format"})
	}

	var task Task

	if err := db.First(&task, uint(idUint)).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	if err := db.Delete(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	return c.JSON(fiber.Map{"message": "Task deleted successfully"})
}
