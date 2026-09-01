package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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

type updateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

var tasks = make(map[int]Task)
var taskID int = 1
var mutex = &sync.RWMutex{}

func NewTask(id int, title string, description string) Task {
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
	task := NewTask(taskID, req.Title, req.Description)
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

func updateTask(w http.ResponseWriter, r *http.Request) {
	var req updateTaskRequest
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	task, ok := tasks[id]
	if !ok {
		mutex.Unlock()
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
	}

	tasks[id] = task
	mutex.Unlock()

	response, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to marshal task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	_, ok := tasks[id]
	if !ok {
		mutex.Unlock()
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	delete(tasks, id)
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Task deleted successfully"}`))

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

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", createTask)
	mux.HandleFunc("GET /tasks", getTasks)
	mux.HandleFunc("PATCH /tasks/{id}", updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", deleteTask)
	return mux
}

func main() {
	http.ListenAndServe(":8080", newRouter())
}
