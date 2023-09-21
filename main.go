package main

import (
	"net/http"
	"os"
	"task-manager-app/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
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

var (
	db        *gorm.DB
	jwtSecret []byte
)

func main() {
	// Setup the database
	db, err := database.SetupDatabase()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Auto Migrate the Task model
	db.AutoMigrate(&Task{}, &User{})

	r := gin.Default()

	// Routes
	r.GET("/", Home)
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.POST("/task/create", CreateTask)
	r.POST("/task/update/{id}", UpdateTask)
	r.POST("/task/delete/{id}", DeleteTask)

	r.Run(":8080")
}
func Home(c *gin.Context) {
	c.HTML(http.StatusAccepted, "home.html", nil)
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email already exists
	var existingUser User
	if result := db.Where("email = ?", user.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash the password (you should use a proper password hashing library)
	// For simplicity, we'll just store the plain text password here
	// In a production environment, you should hash it securely.

	// Create the user
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func CreateTask(c *gin.Context) {
	var task Task
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
	var task Task

	// Find the task by ID
	if result := db.Where("id = ?", taskID).First(&task); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var updatedTask Task
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
	var task Task

	// Find the task by ID
	if result := db.Where("id = ?", taskID).First(&task); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Delete the task
	db.Delete(&task)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by email
	if result := db.Where("email = ?", user.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// In a real application, you should compare the hashed password here.
	// For simplicity, we're just comparing plain text passwords.
	if user.Password != user.Password {
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
