# Task Tracker CLI

A simple command-line task management application written in Go.

The application allows users to create, update, delete, and track tasks. All tasks are stored locally in a JSON file.

## Features

- Add new tasks
- Update task descriptions
- Delete tasks
- Mark tasks as in progress or done
- List all tasks
- Filter tasks by status
- Store tasks in a local JSON file
- Handle invalid commands and input errors
- Run automated tests

## Technologies

- Go (Golang)
- JSON
- Go standard library
- Git

No external dependencies are required.

## Installation

Clone the repository:

```bash
git clone https://github.com/an4nasik007/task-tracker-cli.git
cd task-tracker-cli
```

Build the application:

**macOS / Linux**

```bash
go build -o task-cli .
```

**Windows**

```bash
go build -o task-cli.exe .
```

## Usage

### Add a task

```bash
./task-cli add "Learn Go"
```

Output:

```text
Task added successfully (ID: 1)
```

### Update a task

```bash
./task-cli update 1 "Learn Go basics"
```

### Delete a task

```bash
./task-cli delete 1
```

### Mark a task as in progress

```bash
./task-cli mark-in-progress 1
```

### Mark a task as done

```bash
./task-cli mark-done 1
```

### List all tasks

```bash
./task-cli list
```

### Filter tasks by status

```bash
./task-cli list todo
./task-cli list in-progress
./task-cli list done
```

On Windows, use `.\task-cli.exe` instead of `./task-cli`.

## Task Structure

Each task contains the following fields:

```json
{
  "id": 1,
  "description": "Learn Go",
  "status": "todo",
  "createdAt": "2026-11-22T12:00:00Z",
  "updatedAt": "2026-11-22T12:00:00Z"
}
```

Supported statuses:

- `todo`
- `in-progress`
- `done`

Tasks are stored in `tasks.json` in the current working directory. The file is created automatically if it does not exist.

## Testing

Run all automated tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Check test coverage:

```bash
go test -cover ./...
```

Run static analysis:

```bash
go vet ./...
```

## Project Structure

```text
task-tracker-cli/
├── main.go
├── task.go
├── storage.go
├── main_test.go
├── storage_test.go
├── go.mod
├── .gitignore
└── README.md
```

## License

No license has been specified yet.