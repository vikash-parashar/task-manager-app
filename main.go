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

// Define a struct for the login request data
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CustomClaims represents custom claims in the JWT token.
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	// Add other claims as needed
	jwt.StandardClaims
}

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
	r.GET("/todo", TodoPage)

	r.GET("/tasks", GetAllTasks)
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
	tokenString := c.Query("token")
	if tokenString == "" {
		// Token is not present, render the "home" template
		render.RenderTemplate(c, "home", nil)
		return
	}

	// Render the "todo" template with tasks
	render.RenderTemplate(c, "todo", nil)
}

func LoginPage(c *gin.Context) {
	// Render the HTML page using the template cache
	render.RenderTemplate(c, "login", nil)
}

func TodoPage(c *gin.Context) {
	// Fetch the token from the URL query parameters
	token := c.DefaultQuery("token", "")

	if token == "" {
		// Token is not available, redirect to /home
		c.Redirect(http.StatusSeeOther, "/")
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
	// Declare a variable to hold the task data
	var task models.Task

	// Check if the JSON data in the request body can be bound to the task struct
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the token from the URL query parameters
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is missing"})
		return
	}

	// Get the user from the token (You may need to implement this logic)
	user, err := GetUserFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	// Assuming you have extracted the user's ID from the token correctly
	log.Println("User from token:", user.Email)

	// Set the user ID in the task
	task.UserID = user.ID

	// Generate a UUID for the task
	task.ID = uuid.New()

	// Create the task in the database
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// Return the created task with a 201 status code
	c.JSON(http.StatusCreated, task)
}

func GetAllTasks(c *gin.Context) {
	// Extract the token from the URL query parameters
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is missing"})
		return
	}

	// Get the user ID from the token using your GetUserFromToken function
	user, err := GetUserFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Retrieve tasks associated with the user
	var tasks []models.Task
	if err := db.Where("user_id = ?", user.ID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func getUserByEmail(email string) (*models.User, error) {
	// Create a new User instance to store the result
	user := &models.User{}

	// Query the database to find the user by email
	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return a custom error if the user is not found
			return nil, errors.New("user not found")
		}
		// Return the error for any other database issues
		return nil, err
	}

	// Return the user object and nil error if successful
	return user, nil
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

func Login(c *gin.Context) {
	var loginRequest LoginRequest

	// Parse the request data based on content type
	if err := parseLoginRequest(c, &loginRequest); err != nil {
		return
	}

	// Check the credentials and get the user
	user, err := checkCredentials(loginRequest.Email, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	// Generate a JWT token
	token, err := generateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "Failed to generate JWT token."})
		return
	}

	log.Println("success: JWT token generated")
	log.Println("token: ", token)

	// Redirect based on success
	// if token != "" {
	// 	c.Redirect(http.StatusSeeOther, "/todo?token="+token) // Redirect to /todo with the token
	// } else {
	// 	c.Redirect(http.StatusSeeOther, "/home") // Redirect to /home
	// }
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// parseLoginRequest parses the login request data based on content type
func parseLoginRequest(c *gin.Context, loginRequest *LoginRequest) error {
	contentType := c.Request.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		// Parse JSON request body
		if err := c.ShouldBindJSON(loginRequest); err != nil {
			log.Println("error: failed to bind JSON request data")
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Invalid request"})
			return err
		}
	case "application/x-www-form-urlencoded":
		// Parse form-encoded data
		if err := c.ShouldBind(loginRequest); err != nil {
			log.Println("error: failed to bind form-encoded data")
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Invalid request"})
			return err
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Invalid content type"})
		return nil
	}

	return nil
}

// checkCredentials checks the user's credentials and returns the user if valid
func checkCredentials(email, password string) (*models.User, error) {
	// Replace this with your database lookup logic to retrieve the user by email
	user, err := getUserByEmail(email)
	if err != nil {
		return &models.User{}, err
	}
	// Compare the password from the request with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println("error : password is incorrect")
		log.Println("error : failed to compare hash password with provided password")
		return &models.User{}, err
	}
	return user, nil
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

// GetUserFromToken extracts user information from a JWT token.
func GetUserFromToken(tokenString string) (*models.User, error) {
	// Parse the JWT token
	claims, err := parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Assuming you have a User struct in your application
	user := &models.User{
		ID:    claims.UserID,
		Email: claims.Email,
		// Add other user information as needed
	}

	return user, nil
}

// CheckTokenValidity checks if a JWT token is valid.
func CheckTokenValidity(tokenString string) error {
	// Parse the JWT token
	_, err := parseToken(tokenString)
	return err
}

func parseToken(tokenString string) (*CustomClaims, error) {
	// Define the JWT secret key (you should securely store this)
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract custom claims
	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	}

	return nil, errors.New("failed to extract claims from token")
}
