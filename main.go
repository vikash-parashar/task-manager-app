package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"task-manager-app/database"
	"task-manager-app/models"
	"task-manager-app/render"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	jwtSecret []byte
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Get the JWT secret key from the environment variables
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

}

func main() {
	// Setup the database
	DB, err := database.SetupDatabase()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	db = DB
	// Auto Migrate the Task model
	DB.AutoMigrate(&models.Task{}, &models.User{})

	r := gin.Default()
	r.Static("/static", "./static")

	// Routes
	r.GET("/", HomePage)
	r.GET("/register", RegisterPage)
	r.GET("/login", LoginPage)
	r.POST("/register", Register)
	r.POST("/login", Login)

	r.POST("/tasks", GetAllTask)
	r.POST("/task/create", CreateTask)
	r.POST("/task/update/{id}", UpdateTask)
	r.POST("/task/delete/{id}", DeleteTask)
	// Serve static files (like CSS, JS, and images) if needed
	if err := r.Run(":8080"); err != nil {
		log.Println("failed to start application")
		os.Exit(0)
	}
}

func HomePage(c *gin.Context) {
	// Check if the user is authenticated
	userID, exists := c.Get("user_id")
	if !exists {
		// User is not authenticated, render the "home" template
		render.RenderTemplate(c, "home", nil)
		return
	}

	// User is authenticated, proceed with authenticated user actions
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch tasks for the user
	tasks, err := GetTasksByUserID(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Render the "todo" template with tasks
	render.RenderTemplate(c, "todo", tasks)
}

func LoginPage(c *gin.Context) {
	// Render the HTML page using the template cache
	render.RenderTemplate(c, "login", nil)
}

func RegisterPage(c *gin.Context) {
	// Render the HTML page using the template cache
	render.RenderTemplate(c, "register", nil)
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by email
	if result := db.Where("email = ?", user.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		// You can add more claims here as needed
	})

	// Sign the token with your secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	// Send the JWT token as a response
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Register(c *gin.Context) {
	// Get user details from HTML form
	firstName := c.PostForm("firstname")
	lastName := c.PostForm("lastname")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Bcrypt the password received from the HTML form
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": err.Error()})
		return
	}

	// Create a User struct with the form values
	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Email:     email,
		Password:  hashedPassword,
	}
	fmt.Println(user)
	// Check if the email already exists in the database
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email already exists. Please use a different email address."})
		return
	} else if err != gorm.ErrRecordNotFound {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Internal server error"})
		return
	}

	// Store the user into the database
	user.Password = hashedPassword // Store the hashed password
	if err := db.Create(&user).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Registration failed. Please try again later."})
		return
	}

	// Registration successful
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Registration successful!"})
}

func hashPassword(password string) ([]byte, error) {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the task
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	var task models.Task

	// Find the task by ID
	if result := db.Where("id = ?", taskID).First(&task); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the task
	db.Model(&task).Updates(updatedTask)

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	var task models.Task

	// Find the task by ID
	if result := db.Where("id = ?", taskID).First(&task); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Delete the task
	db.Delete(&task)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func GetAllTask(c *gin.Context) {
	var tasks []models.Task
	// TODO: Fetch tasks from your PostgreSQL database using Gorm
	// For example:
	// db.Find(&tasks)

	// Return the tasks as a JSON response
	c.JSON(http.StatusFound, gin.H{"tasks": tasks})
}

func GetTasksByUserID(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
