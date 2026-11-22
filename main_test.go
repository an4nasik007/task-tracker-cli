
package main

import (
	"os"
	"testing"
)


func setupTest(t *testing.T) {
	t.Helper()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	tempDir := t.TempDir()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})
}


func TestAddTask(t *testing.T) {
	setupTest(t)

	err := addTask("Learn Go")
	if err != nil {
		t.Fatalf("addTask failed: %v", err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatalf("loadTasks failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}

	if task.Description != "Learn Go" {
		t.Errorf("expected description 'Learn Go', got %q", task.Description)
	}

	if task.Status != StatusTodo {
		t.Errorf("expected status %q, got %q", StatusTodo, task.Status)
	}

	if task.CreatedAt.IsZero() {
		t.Error("createdAt must not be empty")
	}

	if task.UpdatedAt.IsZero() {
		t.Error("updatedAt must not be empty")
	}

	if !task.CreatedAt.Equal(task.UpdatedAt) {
		t.Error("createdAt and updatedAt should be equal for a new task")
	}
}

func TestAddMultipleTasks(t *testing.T) {
	setupTest(t)

	if err := addTask("First task"); err != nil {
		t.Fatal(err)
	}

	if err := addTask("Second task"); err != nil {
		t.Fatal(err)
	}

	if err := addTask("Third task"); err != nil {
		t.Fatal(err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	for i, task := range tasks {
		expectedID := i + 1

		if task.ID != expectedID {
			t.Errorf(
				"task %d: expected ID %d, got %d",
				i, expectedID, task.ID,
			)
		}
	}
}


func TestAddEmptyTask(t *testing.T) {
	setupTest(t)

	descriptions := []string{
		"",
		" ",
		"   ",
		"\t\n",
	}

	for _, description := range descriptions {
		if err := addTask(description); err == nil {
			t.Errorf(
				"expected error for empty description %q",
				description,
			)
		}
	}
}

func TestAddTaskTrimsSpaces(t *testing.T) {
	setupTest(t)

	if err := addTask("   Learn Go   "); err != nil {
		t.Fatal(err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Description != "Learn Go" {
		t.Errorf(
			"expected trimmed description, got %q",
			tasks[0].Description,
		)
	}
}

func TestUpdateTask(t *testing.T) {
	setupTest(t)

	if err := addTask("Old description"); err != nil {
		t.Fatal(err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	originalCreatedAt := tasks[0].CreatedAt
	originalUpdatedAt := tasks[0].UpdatedAt

	if err := updateTask(1, "New description"); err != nil {
		t.Fatal(err)
	}

	tasks, err = loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if tasks[0].Description != "New description" {
		t.Errorf(
			"expected updated description, got %q",
			tasks[0].Description,
		)
	}

	if !tasks[0].CreatedAt.Equal(originalCreatedAt) {
		t.Error("createdAt must not change after update")
	}

	if tasks[0].UpdatedAt.Before(originalUpdatedAt) {
		t.Error("updatedAt cannot be earlier than previous updatedAt")
	}
}


func TestUpdateNonexistentTask(t *testing.T) {
	setupTest(t)

	err := updateTask(999, "New description")

	if err == nil {
		t.Error("expected error when updating nonexistent task")
	}
}


func TestUpdateEmptyDescription(t *testing.T) {
	setupTest(t)

	if err := addTask("Original description"); err != nil {
		t.Fatal(err)
	}

	err := updateTask(1, "   ")
	if err == nil {
		t.Error("expected error for empty description")
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if tasks[0].Description != "Original description" {
		t.Error("original description should remain unchanged")
	}
}

func TestDeleteTask(t *testing.T) {
	setupTest(t)

	if err := addTask("Task to delete"); err != nil {
		t.Fatal(err)
	}

	if err := deleteTask(1); err != nil {
		t.Fatal(err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestDeleteMiddleTask(t *testing.T) {
	setupTest(t)

	for _, description := range []string{"First", "Second", "Third"} {
		if err := addTask(description); err != nil {
			t.Fatal(err)
		}
	}

	if err := deleteTask(2); err != nil {
		t.Fatal(err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].ID != 1 || tasks[1].ID != 3 {
		t.Errorf(
			"expected IDs [1, 3], got [%d, %d]",
			tasks[0].ID,
			tasks[1].ID,
		)
	}
}

func TestDeleteNonexistentTask(t *testing.T) {
	setupTest(t)

	if err := deleteTask(999); err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestMarkTaskInProgress(t *testing.T) {
	setupTest(t)

	if err := addTask("Learn testing"); err != nil {
		t.Fatal(err)
	}

	if err := markTask(1, StatusInProgress); err != nil {
		t.Fatal(err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if tasks[0].Status != StatusInProgress {
		t.Errorf(
			"expected %q, got %q",
			StatusInProgress,
			tasks[0].Status,
		)
	}
}

func TestMarkTaskDone(t *testing.T) {
	setupTest(t)

	if err := addTask("Finish project"); err != nil {
		t.Fatal(err)
	}

	if err := markTask(1, StatusDone); err != nil {
		t.Fatal(err)
	}

	tasks, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if tasks[0].Status != StatusDone {
		t.Errorf(
			"expected %q, got %q",
			StatusDone,
			tasks[0].Status,
		)
	}
}

func TestMarkNonexistentTask(t *testing.T) {
	setupTest(t)

	if err := markTask(999, StatusDone); err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestInvalidID(t *testing.T) {
	invalidIDs := []string{
		"",
		"abc",
		"0",
		"-1",
		"1.5",
		"1abc",
	}

	for _, value := range invalidIDs {
		if _, err := parseID(value); err == nil {
			t.Errorf("expected error for ID %q", value)
		}
	}
}

func TestValidID(t *testing.T) {
	validIDs := []struct {
		input    string
		expected int
	}{
		{"1", 1},
		{"10", 10},
		{"999", 999},
	}

	for _, test := range validIDs {
		id, err := parseID(test.input)

		if err != nil {
			t.Errorf("unexpected error for %q: %v", test.input, err)
			continue
		}

		if id != test.expected {
			t.Errorf(
				"input %q: expected %d, got %d",
				test.input,
				test.expected,
				id,
			)
		}
	}
}

func TestListInvalidStatus(t *testing.T) {
	setupTest(t)

	invalidStatuses := []string{
		"invalid",
		"completed",
		"pending",
		"TODO",
	}

	for _, status := range invalidStatuses {
		if err := listTasks(status); err == nil {
			t.Errorf(
				"expected error for invalid status %q",
				status,
			)
		}
	}
}


func TestListValidStatuses(t *testing.T) {
	setupTest(t)

	if err := addTask("Test task"); err != nil {
		t.Fatal(err)
	}

	validStatuses := []string{
		"",
		StatusTodo,
		StatusInProgress,
		StatusDone,
	}

	for _, status := range validStatuses {
		if err := listTasks(status); err != nil {
			t.Errorf(
				"unexpected error for status %q: %v",
				status,
				err,
			)
		}
	}
}
