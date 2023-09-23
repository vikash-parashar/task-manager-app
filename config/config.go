package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Config represents application configuration.
type AppConfig struct {
	JWTSecret []byte
	Database  *gorm.DB
}

// Load loads configuration from environment variables and .env file.
func Load() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: Error loading .env file")
	}

	return &AppConfig{
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),
	}, nil
}
