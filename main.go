package main

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Completed   bool
}

var tasks = make(map[int]Task)
var taskID int = 1

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	task.ID = taskID
	taskID++

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
		return
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	taskList := make([]Task, 0, len(tasks))

	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(taskList)
	if err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("POST /tasks", createTask)
	http.HandleFunc("GET /tasks", getTasks)
	http.ListenAndServe(":8080", nil)
}
