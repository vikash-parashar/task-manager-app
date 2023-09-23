package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// CustomClaims represents custom claims in the JWT token.
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	// Add other claims as needed
	jwt.StandardClaims
}
