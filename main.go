package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"task-manager-app/database"
	"task-manager-app/models"
	"task-manager-app/render"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
	r.POST("/todo", TodoPage)

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
	// Check if a token is present in the request
	tokenString, exists := c.Get("token")
	if !exists {
		// Token is not present, render the "home" template
		render.RenderTemplate(c, "home", nil)
		return
	}

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString.(string), func(token *jwt.Token) (interface{}, error) {
		// Replace with your JWT secret key
		return []byte("your-secret-key"), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Token is valid, proceed with authenticated user actions

	// Fetch tasks for the user (you will need to extract user ID from the token)
	userID := getUserIdFromToken(token)
	tasks, err := GetTasksByUserID(userID)
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

func TodoPage(c *gin.Context) {
	// Check if a token is available in the request payload
	token := c.PostForm("token") // Assuming the token is sent as a POST parameter

	if token == "" {
		// Token is not available, redirect to /home
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}

	// Token is available, render the HTML page using the template cache
	render.RenderTemplate(c, "todo", nil)
}

func RegisterPage(c *gin.Context) {
	// Render the HTML page using the template cache
	render.RenderTemplate(c, "register", nil)
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

	// Generate a UUID for the user
	userID := uuid.New()

	// Create a User struct with the form values and generated UUID
	user := models.User{
		ID:        userID,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Email:     email,
		Password:  hashedPassword,
	}

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

func getUserByEmail(email string) (*models.User, error) {
	// Create a new User instance to store the result
	user := &models.User{}

	// Query the database to find the user by email
	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return a custom error if the user is not found
			return nil, errors.New("User not found")
		}
		// Return the error for any other database issues
		return nil, err
	}

	// Return the user object and nil error if successful
	return user, nil
}

func Login(c *gin.Context) {
	// Define a struct to parse the JSON request body
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Parse the JSON request body into a LoginRequest struct
	var loginRequest LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Retrieve user from the database by email
	user, err := getUserByEmail(loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user exists with this email"})
		return
	}

	// Compare the password from the request with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect"})
		return
	}

	// If everything is okay, generate a JWT token based on UserID (UUID)
	token, err := generateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	// Send the JWT token to the frontend in the response
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func generateJWTToken(userID uuid.UUID) (string, error) {
	// Create a new JWT token with a custom claim (e.g., user's ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),                       // Convert UUID to string
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (adjust as needed)
	})

	// Sign the token with your secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUserIdFromToken(token *jwt.Token) uint {
	claims := token.Claims.(jwt.MapClaims)
	// Extract the user ID from the token claims
	userID, _ := claims["user_id"].(uint)
	return userID
}
