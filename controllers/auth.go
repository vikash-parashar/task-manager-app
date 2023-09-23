package controllers

import (
	"log"
	"net/http"
	"task-manager-app/db"
	"task-manager-app/helpers"
	"task-manager-app/models"
	"time"

	"github.com/gin-gonic/gin"
)

// Login handles user login and JWT token generation.
func Login(c *gin.Context) {
	var loginRequest models.LoginRequest

	// Parse the request data based on content type
	if err := helpers.ParseLoginRequest(c, &loginRequest); err != nil {
		log.Printf("ERROR: Failed to parse login request - %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	// Check the credentials and get the user
	user, err := helpers.CheckCredentials(loginRequest.Email, loginRequest.Password, db.DB)
	if err != nil {
		log.Printf("ERROR: Failed to check credentials - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	// Generate a JWT token
	token, err := helpers.GenerateJWTToken(user.ID)
	if err != nil {
		log.Printf("ERROR: Failed to generate JWT token - %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": "Failed to generate JWT token."})
		return
	}

	// Set the JWT token as a cookie with an expiration time of 30 minutes
	expirationTime := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{
		Name:     "jwt-token",
		Value:    token,
		Expires:  expirationTime,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}

	http.SetCookie(c.Writer, &cookie)

	// Respond with the generated token
	c.JSON(http.StatusOK, gin.H{"status": "login successfully"})
}

// Logout function to delete or expire the JWT cookie
func Logout(c *gin.Context) {
	// Delete or expire the JWT cookie by setting an expired expiration time
	expirationTime := time.Now().Add(-time.Minute * 30) // Expire the cookie by setting it to a past time
	cookie := http.Cookie{
		Name:     "jwt-token",
		Value:    "",
		Expires:  expirationTime,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}

	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, gin.H{"status": "logout successfully"})
}
