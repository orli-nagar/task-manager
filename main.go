package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type createTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
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

func createTask(c *gin.Context) {
	var req createTaskRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mutex.Lock()
	task := NewTask(taskID, req.Title, req.Description)
	tasks[task.ID] = task
	taskID++
	mutex.Unlock()

	c.JSON(http.StatusCreated, task)
}

func updateTask(c *gin.Context) {
	var req updateTaskRequest

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	mutex.Lock()

	task, ok := tasks[id]
	if !ok {
		mutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
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

	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	mutex.Lock()
	_, ok := tasks[idInt]
	if !ok {
		mutex.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	delete(tasks, idInt)
	mutex.Unlock()
	c.Status(http.StatusNoContent)

}

func getTasks(c *gin.Context) {

	mutex.RLock()
	taskList := make([]Task, 0, len(tasks))

	for _, task := range tasks {
		taskList = append(taskList, task)
	}
	mutex.RUnlock()

	c.JSON(http.StatusOK, taskList)
}

func main() {
	gin.EnableJsonDecoderDisallowUnknownFields()
	r := gin.Default()
	r.POST("/tasks", createTask)
	r.GET("/tasks", getTasks)
	r.PATCH("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	r.Run(":8080")
}
