package main

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func healthHandler(c fiber.Ctx) error {
	return c.JSON(map[string]string{"status": "ok"})
}

func getTasksHandler(c fiber.Ctx) error {
	var tasks []Task
	query := c.Query("completed")

	if query != "" {
		completed, err := strconv.ParseBool(query)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameter 'completed'"})
		}
		db.Where("completed = ?", completed).Find(&tasks)
	} else {
		db.Find(&tasks)
	}

	return c.JSON(tasks)
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
