package main

import (
    "bytes"
    "encoding/json"
    "html/template"
    "log"
    "net/http"
    "os"
    "strconv"
    "sync"
)

type Task struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
}

var (
    tasks      []Task
    tasksMutex sync.Mutex
    nextID     = 1
)

func loadTasksFromFile() {
    file, err := os.Open("tasks.json")
    if err != nil {
        if os.IsNotExist(err) {
            // File does not exist, no tasks to load
            return
        }
        log.Println("Error opening JSON file:", err)
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&tasks)
    if err != nil {
        log.Println("Error decoding JSON file:", err)
    } else {
        // Update nextID based on the highest ID in the loaded tasks
        for _, task := range tasks {
            if task.ID >= nextID {
                nextID = task.ID + 1
            }
        }
    }
}

func saveTasksToFile() {
    file, err := os.Create("tasks.json")
    if err != nil {
        log.Println("Error creating JSON file:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    err = encoder.Encode(tasks)
    if err != nil {
        log.Println("Error encoding JSON to file:", err)
    }
}

func main() {
    loadTasksFromFile() // Load tasks from file at startup

    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/add", addTaskHandler)
    http.HandleFunc("/delete", deleteTaskHandler)
    http.HandleFunc("/update", updateTaskHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/index.templ", "templates/todo.templ"))
    tasksMutex.Lock()
    defer tasksMutex.Unlock()
    tmpl.Execute(w, tasks)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    text := r.FormValue("text")
    if text == "" {
        http.Error(w, "Task text cannot be empty", http.StatusBadRequest)
        return
    }

    tasksMutex.Lock()
    task := Task{ID: nextID, Text: text}
    nextID++
    tasks = append(tasks, task)
    tasksMutex.Unlock()

    // Save tasks to file
    saveTasksToFile()

    // Render the single new task
    var htmlBuffer bytes.Buffer
    tmpl := template.Must(template.New("task").ParseFiles("templates/todo.templ"))
    err := tmpl.ExecuteTemplate(&htmlBuffer, "task", task)
    if err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
        return
    }

    // Write the HTML response to the client
    w.Header().Set("Content-Type", "text/html")
    w.Write(htmlBuffer.Bytes())
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    id := r.FormValue("id")
    log.Println("Received delete request for ID:", id)

    taskID, err := strconv.Atoi(id)
    if err != nil {
        log.Println("Error converting ID to integer:", err)
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    tasksMutex.Lock()
    for i, task := range tasks {
        if task.ID == taskID {
            tasks = append(tasks[:i], tasks[i+1:]...)
            tasksMutex.Unlock()

            // Save tasks to file
            saveTasksToFile()

            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Task deleted"))
            return
        }
    }
    tasksMutex.Unlock()

    http.Error(w, "Task not found", http.StatusNotFound)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    id := r.FormValue("id")
    text := r.FormValue("text")
    if text == "" {
        http.Error(w, "Task text cannot be empty", http.StatusBadRequest)
        return
    }

    taskID, err := strconv.Atoi(id)
    if err != nil {
        log.Println("Error converting ID to integer:", err)
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    tasksMutex.Lock()
    defer tasksMutex.Unlock()
    for i, task := range tasks {
        if task.ID == taskID {
            tasks[i].Text = text

            // Save tasks to file
            saveTasksToFile()

            // Render the updated task
            var htmlBuffer bytes.Buffer
            tmpl := template.Must(template.New("task").ParseFiles("templates/todo.templ"))
            err := tmpl.ExecuteTemplate(&htmlBuffer, "task", tasks[i])
            if err != nil {
                http.Error(w, "Error rendering template", http.StatusInternalServerError)
                return
            }

            // Write the HTML response to the client
            w.Header().Set("Content-Type", "text/html")
            w.Write(htmlBuffer.Bytes())
            return
        }
    }

    http.Error(w, "Task not found", http.StatusNotFound)
}
