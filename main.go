package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/orli-nagar/task-manager/database"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Completed   bool   `json:"completed" db:"completed"`
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

func NewTask(id int, title string, description string) Task {
	return Task{
		ID:          id,
		Title:       title,
		Description: description,
		Completed:   false,
	}
}

func insertTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task := NewTask(0, req.Title, req.Description)
	err := database.Db.QueryRow(
		"INSERT INTO tasks (title, description, completed) VALUES ($1, $2, $3) RETURNING id",
		task.Title, task.Description, task.Completed,
	).Scan(&task.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, task)
}

func getTasks(c *gin.Context) {

	var taskList []Task
	err := database.Db.Select(&taskList, "SELECT id, title, description, completed FROM tasks")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, taskList)
}

func updateTask(c *gin.Context) {
	var req updateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var task Task

	err = database.Db.Get(
		&task,
		"SELECT id, title, description, completed FROM tasks WHERE id = $1",
		id,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err != nil {
		c.Status(http.StatusInternalServerError)
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

	_, err = database.Db.Exec(
		`UPDATE tasks
		 SET title = $1, description = $2, completed = $3
		 WHERE id = $4`,
		task.Title,
		task.Description,
		task.Completed,
		task.ID,
	)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result, err := database.Db.Exec("DELETE FROM tasks WHERE id = $1", idInt)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.Status(http.StatusNoContent)

}

func main() {
	database.ConnectDatabase()
	database.CreateTasksTable()
	gin.EnableJsonDecoderDisallowUnknownFields()
	r := gin.Default()
	r.POST("/tasks", insertTask)
	r.GET("/tasks", getTasks)
	r.PATCH("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	r.Run(":8080")
}
