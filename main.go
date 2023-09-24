package main

import (
	"log"
	"task-manager-app/config"
	"task-manager-app/middleware"
	"task-manager-app/models"
	"task-manager-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the application configuration
	appConfig, err := config.Load()
	if err != nil {
		log.Fatalf("ERROR: Failed to load configuration: %v", err)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Static("/static", "./static")

	// Auto Migrate the User and Task models
	appConfig.Database.AutoMigrate(&models.User{}, &models.Task{})

	// Use the InjectDB middleware to inject the database into the context
	r.Use(middleware.InjectDB(appConfig.Database))

	// Define routes and pass appConfig as a closure to the route handlers
	routes.SetupRoutes(r, appConfig)

	// Serve static files (like CSS, JS, and images) if needed
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("ERROR: Failed to start application: %v", err)
	}
}
