package db

import (
	"errors"
	"log"
	"task-manager-app/models"

	"gorm.io/gorm"
)

// getUserByEmail retrieves a user by their email address.
func GetUserByEmail(email string, db *gorm.DB) (*models.User, error) {
	// Create a new User instance to store the result
	user := &models.User{}

	// Query the database to find the user by email
	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("INFO: User not found for email: %s\n", email)
			// Return a custom error if the user is not found
			return nil, errors.New("user not found")
		}
		log.Printf("ERROR: Failed to retrieve user by email - %s\n", err)
		// Return the error for any other database issues
		return nil, err
	}

	// Do not log sensitive user details
	log.Printf("INFO: Retrieved user by email: %s\n", email)

	// Return the user object and nil error if successful
	return user, nil
}
