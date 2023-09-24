package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// AppConfig represents application configuration.
type AppConfig struct {
    JWTSecret []byte
    Database  *gorm.DB
}

var appConfig *AppConfig

// Load loads configuration from environment variables and .env file.
func Load() (*AppConfig, error) {
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
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("ERROR: Failed to connect to the database: %v\n", err)
        return nil, err
    }

    log.Println("INFO: Database connection established")

    // Initialize the AppConfig struct
    appConfig = &AppConfig{
        JWTSecret: []byte(os.Getenv("JWT_SECRET")),
        Database:  db,
    }

    return appConfig, nil
}

// GetAppConfig returns the application configuration.
func GetAppConfig() *AppConfig {
    return appConfig
}
