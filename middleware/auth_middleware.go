package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// AuthMiddleware checks the validity of the JWT token from the request header and authorizes the user
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the request header
		tokenString := c.GetHeader("Authorization")

		// Check if the token is missing or doesn't start with "Bearer "
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
			c.Abort()
			return
		}

		// Extract the token without the "Bearer " prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method and provide the key used for signing
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid signing method")
			}
			// Replace "your-secret-key" with your actual secret key
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Token is valid, you can access claims from token.Claims
		claims, _ := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Store the user ID or other claims in the context for later use
		c.Set("user_id", userID)

		// Continue with the request
		c.Next()
	}
}

func InjectDB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", db) // Set the "DB" key in the context to the database instance
		c.Next()        // Continue processing the request
	}
}
