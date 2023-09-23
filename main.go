package main

import (
	"log"
	"task-manager-app/config"
	"task-manager-app/db"
	"task-manager-app/middleware"
	"task-manager-app/models"
	"task-manager-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ERROR: Failed to load configuration: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	r.Use(gin.Recovery())
	r.Static("/static", "./static")

	// Setup the database using the configuration variables
	dbInstance, err := db.Setup()
	if err != nil {
		log.Fatalf("ERROR: Failed to connect to database: %v", err)
	}

	// Auto Migrate the User and Task models
	dbInstance.AutoMigrate(&models.User{}, &models.Task{})
	cfg.Database = dbInstance
	// Define routes
	routes.SetupRoutes(r, dbInstance)

	// Serve static files (like CSS, JS, and images) if needed
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("ERROR: Failed to start application: %v", err)
	}
}
