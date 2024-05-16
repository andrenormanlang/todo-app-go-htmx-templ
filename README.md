
# Golang, HTMX & Templ To-Do List Application

This is a simple web-based To-Do List application implemented in Go. The application supports adding, updating, and deleting tasks. It uses an HTML template for rendering the user interface and stores tasks in a JSON file.

## Features

- View the list of tasks
- Add a new task
- Update an existing task
- Delete a task
- Tasks are persisted in a JSON file

## Prerequisites

- Go 1.16 or later

## Getting Started

### 1. Clone the Repository

```sh
git clone <repository-url>
cd <repository-directory>
```

### 2. Install Dependencies

No additional dependencies need to be installed.

### 3. Directory Structure

```
.
├── main.go
├── tasks.json
├── templates
│   ├── index.templ
│   └── todo.templ
├── static
│   └── script.js
└── README.md
```

### 4. JSON File

Ensure a `tasks.json` file exists in the root directory. This file is used to store the tasks. If it doesn't exist, it will be created upon starting the application.

Example content of `tasks.json`:
```json
[
    {
        "id": 1,
        "text": "Sample Task 1"
    },
    {
        "id": 2,
        "text": "Sample Task 2"
    }
]
```

### 5. Run the Application

```sh
go run main.go
```

The application will start a web server on `localhost:8080`.

### 6. Access the Application

Open a web browser and navigate to `http://localhost:8080`.

## Handlers

### 1. `indexHandler`

Renders the main page showing the list of tasks.

### 2. `addTaskHandler`

Handles adding a new task. Expects a POST request with a `text` field.

### 3. `deleteTaskHandler`

Handles deleting a task. Expects a POST request with an `id` field.

### 4. `updateTaskHandler`

Handles updating a task. Expects a POST request with `id` and `text` fields.

## HTML Templates

### `templates/index.templ`

Main template rendering the list of tasks.

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>To-Do List</title>
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.7.0"></script>
    <script src="/static/script.js"></script>
</head>
<body class="bg-gray-100 flex flex-col items-center">
    <h1 class="text-2xl font-bold text-gray-800 mt-10">To-Do List</h1>
    <form id="addTaskForm" hx-post="/add" hx-target="#tasks" hx-swap="beforeend" class="my-6">
        <input type="text" name="text" placeholder="New task" required class="p-2 border border-gray-300 rounded">
        <button type="submit" class="p-2 bg-green-500 text-white rounded hover:bg-green-600">Add Task</button>
    </form>
    <ul id="tasks" class="w-96">
        {{range .}}
            {{template "task" .}}
        {{end}}
    </ul>
</body>
</html>
```
