package main
import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)
const storageFile = "tasks.json"

func loadTasks() ([]Task, error){
	data, err :=os.ReadFile(storageFile)
	if errors.Is(err, os.ErrNotExist){
		data=[]byte("[]")
	
	if err:=os.WriteFile(storageFile,data,0644); err != nil{
		return nil, fmt.Errorf("cannot create tasks file: %w",err)
		}
	} else if err != nil{
		return nil, fmt.Errorf("cannot read task file: %w",err)
	}

	var tasks []Task

	if err:=json.Unmarshal(data, &tasks); err!=nil{
		return nil, fmt.Errorf("invalid tasks JSON: %w",err)
	}
	if tasks == nil{
		return nil, fmt.Errorf("tasks JSON must contain an array")
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error{
	if tasks == nil {
		tasks = []Task{}
	}

	data, err := json.MarshalIndent(tasks,""," ")
	if err != nil {
		return fmt.Errorf("cannot encode tasks: %w",err)
	}

	tempFile, err := os.CreateTemp(".", "tasks-*.tmp")
	if err != nil {
		return fmt.Errorf("cannot create temporary file: %w",err)
	}

	tempName := tempFile.Name()
	defer os.Remove(tempName)

	if err := tempFile.Chmod(0644); err!= nil{
		tempFile.Close()
		return fmt.Errorf("cannot set file permissions: %w",err)
	}

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return fmt.Errorf("cannot write temporary file: %w",err)
	}

	if err := tempFile.Close(); err != nil {
		tempFile.Close()
		return fmt.Errorf("cannot close temporary file: %w",err)
	}

	if err := os.Rename(tempName, storageFile); err != nil {
        return fmt.Errorf("cannot replace tasks file: %w", err)
    }

	return nil

}