package controllers

import (
	"log"
	"net/http"
	"strings"
	"task-manager-app/db"
	"task-manager-app/helpers"
	"task-manager-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Register handles user registration.
func Register(c *gin.Context) {
	// Get user details from HTML form
	firstName := c.PostForm("firstname")
	lastName := c.PostForm("lastname")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Bcrypt the password received from the HTML form
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		log.Printf("ERROR: Failed to hash password - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
	if err := db.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email already exists. Please use a different email address."})
		return
	} else if err != gorm.ErrRecordNotFound {
		log.Printf("ERROR: Failed to check existing email - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Internal server error"})
		return
	}

	// Store the user into the database
	if err := db.DB.Create(&user).Error; err != nil {
		log.Printf("ERROR: Failed to create user - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Registration failed. Please try again later."})
		return
	}

	// Registration successful
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Registration successful!"})
}

// GetCurrentUser retrieves the current user based on the JWT token.
func GetCurrentUser(c *gin.Context) {
	// Extract the JWT token from the cookie
	token, err := c.Cookie("jwt-token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing token cookie"})
		return
	}

	// Proceed with the existing code to validate and retrieve the user
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is missing"})
		return
	}

	// Extract the token from the "Bearer <token>" format
	token = strings.Replace(authorizationHeader, "Bearer ", "", 1)

	// Validate and parse the JWT token to get the user information
	user, err := helpers.GetUserFromToken(token) // Implement GetUserFromToken based on your needs
	if err != nil {
		log.Printf("ERROR: Failed to validate token - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Fetch user data from the database
	if err := db.DB.Where("id = ?", user.ID).Preload("Tasks").First(&user).Error; err != nil {
		log.Printf("ERROR: Failed to retrieve user - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// Return the user information in the response
	c.JSON(http.StatusOK, user)
}
