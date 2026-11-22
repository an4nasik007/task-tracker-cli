
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	var err error

	switch command {
	case "add":
		if len(os.Args) != 3 {
			err = fmt.Errorf("usage: task-cli add \"description\"")
		} else {
			err = addTask(os.Args[2])
		}

	case "update":
		if len(os.Args) != 4 {
			err = fmt.Errorf("usage: task-cli update <id> \"description\"")
		} else {
			var id int
			id, err = parseID(os.Args[2])
			if err == nil {
				err = updateTask(id, os.Args[3])
			}
		}

	case "delete":
		if len(os.Args) != 3 {
			err = fmt.Errorf("usage: task-cli delete <id>")
		} else {
			var id int
			id, err = parseID(os.Args[2])
			if err == nil {
				err = deleteTask(id)
			}
		}

	case "mark-in-progress", "mark-done":
		if len(os.Args) != 3 {
			err = fmt.Errorf("usage: task-cli %s <id>", command)
		} else {
			var id int
			id, err = parseID(os.Args[2])

			if err == nil {
				status := StatusInProgress
				if command == "mark-done" {
					status = StatusDone
				}
				err = markTask(id, status)
			}
		}

	case "list":
		if len(os.Args) > 3 {
			err = fmt.Errorf("usage: task-cli list [todo|in-progress|done]")
		} else {
			status := ""
			if len(os.Args) == 3 {
				status = os.Args[2]
			}
			err = listTasks(status)
		}

	default:
		err = fmt.Errorf("unknown command: %s", command)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func parseID(value string) (int, error) {
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid task ID: %s", value)
	}
	return id, nil
}

func addTask(description string) error {
	description = strings.TrimSpace(description)
	if description == "" {
		return fmt.Errorf("task description cannot be empty")
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	now := time.Now().UTC()

	newTask := Task{
		ID:          maxID + 1,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)

	if err := saveTasks(tasks); err != nil {
		return err
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)
	return nil
}

func updateTask(id int, description string) error {
	description = strings.TrimSpace(description)
	if description == "" {
		return fmt.Errorf("task description cannot be empty")
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now().UTC()

			if err := saveTasks(tasks); err != nil {
				return err
			}

			fmt.Printf("Task %d updated successfully\n", id)
			return nil
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

func deleteTask(id int) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			if err := saveTasks(tasks); err != nil {
				return err
			}

			fmt.Printf("Task %d deleted successfully\n", id)
			return nil
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

func markTask(id int, status string) error {
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now().UTC()

			if err := saveTasks(tasks); err != nil {
				return err
			}

			fmt.Printf("Task %d marked as %s\n", id, status)
			return nil
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

func listTasks(status string) error {
	if status != "" &&
		status != StatusTodo &&
		status != StatusInProgress &&
		status != StatusDone {
		return fmt.Errorf("invalid status: %s", status)
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	count := 0

	for _, task := range tasks {
		if status != "" && task.Status != status {
			continue
		}

		fmt.Printf(
			"ID: %d | Description: %s | Status: %s\n",
			task.ID,
			task.Description,
			task.Status,
		)
		count++
	}

	if count == 0 {
		fmt.Println("No tasks found")
	}

	return nil
}

func printUsage() {
	fmt.Println("Task Tracker CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  task-cli add \"description\"")
	fmt.Println("  task-cli update <id> \"description\"")
	fmt.Println("  task-cli delete <id>")
	fmt.Println("  task-cli mark-in-progress <id>")
	fmt.Println("  task-cli mark-done <id>")
	fmt.Println("  task-cli list [todo|in-progress|done]")
}
