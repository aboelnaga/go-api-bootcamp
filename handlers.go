package main
import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

func getTaskById(id string) (Task, error) {
	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return Task{}, fmt.Errorf("task not found")
}

func getTaskIndexById(id string) (int, error) {
	for i, task := range tasks {
		if task.ID == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("task not found")
}

func queryTasksByCompleted(isCompleted string) ([]Task, error) {
	queryResult := []Task{}
	queryValue, err := strconv.ParseBool(isCompleted)

	if err != nil {
		return nil, fmt.Errorf("not valid query")
	}

	for _, task := range tasks {
		if task.Completed == queryValue {
			queryResult = append(queryResult, task)
		}
	}

	return queryResult, nil
}

func healthHandler (c fiber.Ctx) error {
	return c.JSON(map[string]string{"status": "ok"})
}

func getTasksHandler (c fiber.Ctx) error {
	completed := c.Query("completed")

	if completed != "" {
		queryResult, err := queryTasksByCompleted(completed)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(queryResult)
	}

	return c.JSON(tasks)
}

func getTaskByIdHandler (c fiber.Ctx) error {
	id := c.Params("id")

	task, err := getTaskById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}
	return c.JSON(task)
}

func createTaskHandler (c fiber.Ctx) error {
	task := new(Task)

	if err := c.Bind().Body(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	task.ID = fmt.Sprintf("%d", len(tasks)+1)
	task.CreatedAt = time.Now()
	tasks = append(tasks, *task)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Task created successfully", "task": task})
}

func updateTaskHandler (c fiber.Ctx) error {
	id := c.Params("id")

	index, err := getTaskIndexById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	originalID := tasks[index].ID
	originalCreatedAt := tasks[index].CreatedAt

	if err := c.Bind().Body(&tasks[index]); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	tasks[index].ID = originalID
	tasks[index].CreatedAt = originalCreatedAt

	return c.Status(fiber.StatusOK).JSON(tasks[index])
}

func deleteTaskHandler (c fiber.Ctx) error {
	id := c.Params("id")

	index, err := getTaskIndexById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	tasks = slices.Delete(tasks, index, index+1)

	return c.JSON(fiber.Map{"message": "Task deleted successfully"})
}