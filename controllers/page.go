package controllers

import (
	"log"
	"net/http"
	"task-manager-app/render"

	"github.com/gin-gonic/gin"
)

// HomePage renders the home page.
func HomePage(c *gin.Context) {
	// Extract the JWT token from the cookie
	token, err := c.Cookie("jwt-token")
	if err != nil || token == "" {
		// Log the error and return an unauthorized response
		log.Printf("ERROR: Failed to extract JWT token from cookie - %s\n", err)
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token cookie"})
		// Token is not present, render the "home" template
		render.RenderTemplate(c, "home", nil)
		return
	}

	// Render the "todo" template with tasks
	render.RenderTemplate(c, "todo", nil)
}

// LoginPage renders the login page.
func LoginPage(c *gin.Context) {
	// Render the HTML page using the template cache
	render.RenderTemplate(c, "login", nil)
}

// TodoPage renders the todo page.
func TodoPage(c *gin.Context) {
	// Extract the JWT token from the cookie
	token, err := c.Cookie("jwt-token")
	if err != nil {
		// Log the error and return an unauthorized response
		log.Printf("ERROR: Failed to extract JWT token from cookie - %s\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token cookie"})
		return
	}

	if token == "" {
		// Token is not available, redirect to /home
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	// Token is available, render the HTML page using the template cache
	render.RenderTemplate(c, "todo", nil)
}

// RegisterPage renders the register page.
func RegisterPage(c *gin.Context) {
	// Render the HTML page using the template cache
	render.RenderTemplate(c, "register", nil)
}
