package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// Import your models package (e.g., "your/package/models")
)

var DB *gorm.DB

// Setup initializes the database connection.
func Setup() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: Error loading .env file")
	}

	// Retrieve database connection information from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create a database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open a database connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("ERROR: Failed to connect to the database: %v\n", err)
		return nil, err
	}

	log.Println("INFO: Database connection established")

	return DB, nil
}

// GetDB returns the database instance.
func GetDB() *gorm.DB {
	return DB
}
