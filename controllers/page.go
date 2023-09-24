package controllers

import (
	"log"
	"net/http"
	"strings"
	"task-manager-app/helpers"
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
		render.RenderTemplate(c, "home", nil)
		// Token is not present, render the "home" template

		return
	} else {
		http.RedirectHandler("/todo", http.StatusTemporaryRedirect)
		return
	}

}

// LoginPage renders the login page.
func LoginPage(c *gin.Context) {
	// Render the HTML page using the template cache
	render.RenderTemplate(c, "login", nil)

}

// TodoPage renders the todo page.
func TodoPage(c *gin.Context) {
	// Get the JWT token from the request header
	tokenString := c.GetHeader("Authorization")

	// Check if the token is missing or doesn't start with "Bearer "
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		c.Redirect(http.StatusSeeOther, "/") // Redirect if token is missing or invalid
		return
	}

	// Extract the token without the "Bearer " prefix
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if !helpers.IsValidToken(tokenString) {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	// Token is valid, render the HTML page using the template cache
	render.RenderTemplate(c, "todo", nil)
}

// RegisterPage renders the register page.
func RegisterPage(c *gin.Context) {

	// Render the HTML page using the template cache
	render.RenderTemplate(c, "register", nil)

}
