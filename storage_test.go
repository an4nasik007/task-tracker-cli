
package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"
)


func TestLoadTasksCreatesFile(t *testing.T) {
	setupTest(t)

	tasks, err := loadTasks()
	if err != nil {
		t.Fatalf("loadTasks failed: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}

	data, err := os.ReadFile(storageFile)
	if err != nil {
		t.Fatalf("tasks.json was not created: %v", err)
	}

	if string(data) != "[]" {
		t.Errorf("expected empty JSON array, got %q", string(data))
	}
}


func TestLoadTasksExistingFile(t *testing.T) {
	setupTest(t)

	now := time.Now().UTC().Truncate(time.Second)

	expected := []Task{
		{
			ID:          1,
			Description: "Learn Go",
			Status:      StatusTodo,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Description: "Write tests",
			Status:      StatusInProgress,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	data, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(storageFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	actual, err := loadTasks()
	if err != nil {
		t.Fatalf("loadTasks failed: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("loaded tasks do not match expected tasks")
	}
}


func TestSaveTasks(t *testing.T) {
	setupTest(t)

	now := time.Now().UTC().Truncate(time.Second)

	expected := []Task{
		{
			ID:          1,
			Description: "Learn Golang",
			Status:      StatusTodo,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	if err := saveTasks(expected); err != nil {
		t.Fatalf("saveTasks failed: %v", err)
	}

	data, err := os.ReadFile(storageFile)
	if err != nil {
		t.Fatalf("cannot read saved file: %v", err)
	}

	var actual []Task

	if err := json.Unmarshal(data, &actual); err != nil {
		t.Fatalf("saved file contains invalid JSON: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("saved tasks do not match expected tasks")
	}
}


func TestSaveEmptyTasks(t *testing.T) {
	setupTest(t)

	if err := saveTasks([]Task{}); err != nil {
		t.Fatalf("saveTasks failed: %v", err)
	}

	data, err := os.ReadFile(storageFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "[]" {
		t.Errorf("expected [], got %q", string(data))
	}
}


func TestSaveNilTasks(t *testing.T) {
	setupTest(t)

	if err := saveTasks(nil); err != nil {
		t.Fatalf("saveTasks failed: %v", err)
	}

	data, err := os.ReadFile(storageFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "[]" {
		t.Errorf("expected [], got %q", string(data))
	}
}

func TestSaveTasksOverwrites(t *testing.T) {
	setupTest(t)

	now := time.Now().UTC()

	first := []Task{
		{
			ID:          1,
			Description: "Old task",
			Status:      StatusTodo,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	if err := saveTasks(first); err != nil {
		t.Fatal(err)
	}

	second := []Task{
		{
			ID:          2,
			Description: "New task",
			Status:      StatusDone,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	if err := saveTasks(second); err != nil {
		t.Fatal(err)
	}

	actual, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if len(actual) != 1 {
		t.Fatalf("expected 1 task, got %d", len(actual))
	}

	if actual[0].ID != 2 {
		t.Errorf("expected ID 2, got %d", actual[0].ID)
	}

	if actual[0].Description != "New task" {
		t.Error("old task was not replaced")
	}
}

func TestInvalidJSON(t *testing.T) {
	setupTest(t)

	invalidData := []byte("{invalid json}")

	if err := os.WriteFile(storageFile, invalidData, 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadTasks(); err == nil {
		t.Error("expected error for invalid JSON")
	}

	data, err := os.ReadFile(storageFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != string(invalidData) {
		t.Error("invalid JSON file was unexpectedly modified")
	}
}

func TestLoadTasksRejectsNull(t *testing.T) {
	setupTest(t)

	if err := os.WriteFile(storageFile, []byte("null"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadTasks(); err == nil {
		t.Error("expected error for JSON null")
	}
}


func TestLoadTasksRejectsObject(t *testing.T) {
	setupTest(t)

	data := []byte(`{"id":1,"description":"Test"}`)

	if err := os.WriteFile(storageFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadTasks(); err == nil {
		t.Error("expected error for JSON object instead of array")
	}
}


func TestLoadTasksEmptyFile(t *testing.T) {
	setupTest(t)

	if err := os.WriteFile(storageFile, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := loadTasks(); err == nil {
		t.Error("expected error for empty JSON file")
	}
}


func TestLoadTasksUnreadablePath(t *testing.T) {
	setupTest(t)

	if err := os.Mkdir(storageFile, 0755); err != nil {
		t.Fatal(err)
	}

	if _, err := loadTasks(); err == nil {
		t.Error("expected error when storage path is a directory")
	}
}


func TestSaveTasksInvalidPath(t *testing.T) {
	setupTest(t)

	if err := os.Mkdir(storageFile, 0755); err != nil {
		t.Fatal(err)
	}

	err := saveTasks([]Task{})
	if err == nil {
		t.Error("expected error when storage path is a directory")
	}
}


func TestSaveLoadRoundTrip(t *testing.T) {
	setupTest(t)

	now := time.Now().UTC().Truncate(time.Second)

	expected := []Task{
		{
			ID:          1,
			Description: "Buy groceries",
			Status:      StatusTodo,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Description: "Write Go tests",
			Status:      StatusInProgress,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          3,
			Description: "Finish project",
			Status:      StatusDone,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	if err := saveTasks(expected); err != nil {
		t.Fatal(err)
	}

	actual, err := loadTasks()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("tasks changed after save/load round trip")
	}
}
