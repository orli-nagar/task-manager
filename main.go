package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var tasks = make(map[int]Task)
var taskID int = 1
var mutex = &sync.RWMutex{}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task.Completed = false

	if task.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	task.ID = taskID
	taskID++
	tasks[task.ID] = task
	mutex.Unlock()

	response, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to marshal task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func getTasks(w http.ResponseWriter, r *http.Request) {

	mutex.RLock()
	taskList := make([]Task, 0, len(tasks))

	for _, task := range tasks {
		taskList = append(taskList, task)
	}
	mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(taskList)
	if err != nil {
		http.Error(w, "Failed to marshal tasks", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	http.HandleFunc("POST /tasks", createTask)
	http.HandleFunc("GET /tasks", getTasks)
	http.ListenAndServe(":8080", nil)
}
