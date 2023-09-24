package controllers

import (
	"log"
	"net/http"
	"task-manager-app/config"
	"task-manager-app/helpers"
	"task-manager-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTask creates a new task.
func CreateTask(c *gin.Context, appConfig *config.AppConfig) {
	// Declare a variable to hold the task data
	var task models.Task

	// Check if the JSON data in the request body can be bound to the task struct
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Printf("ERROR: Failed to bind JSON request data - %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the JWT token from the cookie
	token, err := c.Cookie("jwt-token")
	if err != nil {
		log.Printf("ERROR: Failed to extract JWT token from cookie - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token cookie"})
		return
	}

	// Get the user from the token (You may need to implement this logic)
	user, err := helpers.GetUserFromToken(token)
	if err != nil {
		log.Printf("ERROR: Invalid token - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Assuming you have extracted the user's ID from the token correctly
	log.Println("User from token:", user.Email)

	// Set the user ID in the task
	task.UserID = user.ID

	// Generate a UUID for the task
	task.ID = uuid.New()
	// by default set task status to pending
	task.Status = models.Pending

	// Create the task in the database
	if err := appConfig.Database.Create(&task).Error; err != nil {
		log.Printf("ERROR: Failed to create task - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// Return the created task with a 201 status code
	c.JSON(http.StatusCreated, task)
}

// GetAllTasks retrieves all tasks for the authenticated user.
func GetAllTasks(c *gin.Context, appConfig *config.AppConfig) {
	// Extract the JWT token from the cookie
	token, err := c.Cookie("jwt-token")
	if err != nil {
		log.Printf("ERROR: Failed to extract JWT token from cookie - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token cookie"})
		return
	}

	// Get the user ID from the token using your GetUserFromToken function
	user, err := helpers.GetUserFromToken(token)
	if err != nil {
		log.Printf("ERROR: Invalid token - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Retrieve tasks associated with the user
	var tasks []models.Task
	if err := appConfig.Database.Where("user_id = ?", user.ID).Find(&tasks).Error; err != nil {
		log.Printf("ERROR: Failed to fetch tasks - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// UpdateTask updates the status and priority of a task.
func UpdateTask(c *gin.Context, appConfig *config.AppConfig) {
	// Extract the JWT token from the cookie
	token, err := c.Cookie("jwt-token")
	if err != nil {
		log.Printf("ERROR: Failed to extract JWT token from cookie - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token cookie"})
		return
	}

	// Get the user ID from the token using your GetUserFromToken function
	user, err := helpers.GetUserFromToken(token)
	if err != nil {
		log.Printf("ERROR: Invalid token - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract the task ID from the URL route parameters
	taskID := c.Param("id")

	// Check if the task exists and belongs to the user
	var task models.Task
	if err := appConfig.Database.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task).Error; err != nil {
		log.Printf("ERROR: Task not found or does not belong to the user - %s\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or does not belong to the user"})
		return
	}

	// Parse the request body to get the updated status and priority
	var updateData struct {
		Status   string `json:"status"`
		Priority string `json:"priority"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		log.Printf("ERROR: Failed to bind JSON request data - %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if updateData.Status != "" {
		// Update the task status
		task.Status = updateData.Status
	}

	if updateData.Priority != "" {
		// Update the task priority
		task.Priority = updateData.Priority
	}

	if err := appConfig.Database.Save(&task).Error; err != nil {
		log.Printf("ERROR: Failed to update task - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

// DeleteTask deletes a task.
func DeleteTask(c *gin.Context, appConfig *config.AppConfig) {
	// Extract the JWT token from the cookie
	token, err := c.Cookie("jwt-token")
	if err != nil {
		log.Printf("ERROR: Failed to extract JWT token from cookie - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token cookie"})
		return
	}

	// Get the user ID from the token using your GetUserFromToken function
	user, err := helpers.GetUserFromToken(token)
	if err != nil {
		log.Printf("ERROR: Invalid token - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract the task ID from the URL route parameters
	taskID := c.Param("id")

	// Check if the task exists and belongs to the user
	var task models.Task
	if err := appConfig.Database.Where("id = ? AND user_id = ?", taskID, user.ID).First(&task).Error; err != nil {
		log.Printf("ERROR: Task not found or does not belong to the user - %s\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or does not belong to the user"})
		return
	}

	// Delete the task
	if err := appConfig.Database.Delete(&task).Error; err != nil {
		log.Printf("ERROR: Failed to delete task - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
