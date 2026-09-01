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

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

var tasks = make(map[int]Task)
var taskID int = 1
var mutex = &sync.RWMutex{}

func newTask(id int, title string, description string) Task {
	return Task{
		ID:          id,
		Title:       title,
		Description: description,
		Completed:   false,
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	task := newTask(taskID, req.Title, req.Description)
	tasks[task.ID] = task
	taskID++
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
