package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtSecret []byte

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the request header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("error : JWT token is missing")
			return
		}

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil // Use the same JWT secret key defined earlier
		})

		if err != nil || !token.Valid {
			log.Println("error : Invalid JWT token")
			return
		}

		// If the token is valid, proceed to the next handler
		c.Next()
	}
}
