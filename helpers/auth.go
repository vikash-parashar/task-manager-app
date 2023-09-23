package helpers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"task-manager-app/config"
	"task-manager-app/db"
	"task-manager-app/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) ([]byte, error) {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR: Failed to hash password - %s\n", err)
		return nil, err
	}
	return hashedPassword, nil
}

// GetUserFromToken extracts user information from a JWT token.
func GetUserFromToken(tokenString string) (*models.User, error) {
	// Parse the JWT token
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Assuming you have a User struct in your application
	user := &models.User{
		ID:    claims.UserID,
		Email: claims.Email,
		// Add other user information as needed
	}

	return user, nil
}

// CheckTokenValidity checks if a JWT token is valid.
func CheckTokenValidity(tokenString string) error {
	// Parse the JWT token
	_, err := ParseToken(tokenString)
	return err
}

func ParseToken(tokenString string) (*models.CustomClaims, error) {
	// Define the JWT secret key (you should securely store this)
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract custom claims
	if claims, ok := token.Claims.(*models.CustomClaims); ok {
		return claims, nil
	}

	return nil, errors.New("failed to extract claims from token")
}

// parseLoginRequest parses the login request data based on content type
func ParseLoginRequest(c *gin.Context, loginRequest *models.LoginRequest) error {
	contentType := c.Request.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		// Parse JSON request body
		if err := c.ShouldBindJSON(loginRequest); err != nil {
			log.Printf("ERROR: Failed to bind JSON request data - %s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Invalid request"})
			return err
		}
	case "application/x-www-form-urlencoded":
		// Parse form-encoded data
		if err := c.ShouldBind(loginRequest); err != nil {
			log.Printf("ERROR: Failed to bind form-encoded data - %s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Invalid request"})
			return err
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "Invalid content type"})
		return nil
	}

	return nil
}

// checkCredentials checks the user's credentials and returns the user if valid
func CheckCredentials(email, password string, database *gorm.DB) (*models.User, error) {
	// Replace this with your database lookup logic to retrieve the user by email
	user, err := db.GetUserByEmail(email, database)
	if err != nil {
		log.Printf("ERROR: Failed to fetch user by email - %s\n", err)
		return &models.User{}, err
	}

	// Compare the password from the request with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println("ERROR: Password is incorrect")
		log.Printf("ERROR: Failed to compare hash password with provided password - %s\n", err)
		return &models.User{}, err
	}

	return user, nil
}

func GenerateJWTToken(userID uuid.UUID) (string, error) {
	// Create a new JWT token with a custom claim (e.g., user's ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),                         // Convert UUID to string
		"exp":     time.Now().Add(time.Minute * 30).Unix(), // Token expiration time (adjust as needed)
	})
	mySecret, _ := config.Load()
	// Sign the token with your secret key
	tokenString, err := token.SignedString(mySecret.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
