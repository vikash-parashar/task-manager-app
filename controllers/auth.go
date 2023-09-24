package controllers

import (
	"log"
	"net/http"
	"task-manager-app/config"
	"task-manager-app/helpers"
	"task-manager-app/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user login and JWT token generation.
func Login(c *gin.Context, appConfig *config.AppConfig) {
	var loginRequest models.LoginRequest

	// Parse the request data based on content type
	if err := helpers.ParseLoginRequest(c, &loginRequest); err != nil {
		log.Printf("ERROR: Failed to parse login request - %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	// Retrieve the user data from the database based on the user's email
	dbUser, err := helpers.GetUserByEmail(loginRequest.Email, appConfig)
	if err != nil {
		log.Printf("ERROR: Failed to get user from the database - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": "Invalid email or password."})
		return
	}

	// Compare the provided password with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword(dbUser.Password, []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": "Invalid email or password."})
		return
	}

	// At this point, the login is successful

	// Generate a JWT token using the user's ID
	token, err := helpers.GenerateJWTToken(dbUser.ID)
	if err != nil {
		log.Printf("ERROR: Failed to generate JWT token - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "Failed to generate JWT token."})
		return
	}

	// Respond with the generated token in the response body
	response := gin.H{
		"status":  "success",
		"message": "Authentication successful",
		"data": gin.H{
			"user_id":      dbUser.ID,
			"firstname":    dbUser.FirstName,
			"email":        dbUser.Email,
			"access_token": token,
			"expires_in":   1800, // Token expiration time in seconds (30 minutes)
		},
	}

	c.JSON(http.StatusOK, response)
}

// Logout function to clear the JWT token from client-side local storage
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "logout successfully"})
}

// client side code for log out

// <script>

// Function to clear the JWT token from local storage

// function logout() {
//     localStorage.removeItem('jwt-token');
//     // Redirect the user to the logout page or home page
//     window.location.href = '/logout'; // You can specify the desired logout page
// }

// </script>
