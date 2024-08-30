package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Exclude password from JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request payload for creating a new user
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserResponse represents the response payload after creating a new user
type CreateUserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserLoginRequest represents the request payload for logging in a user
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthTokenResponse represents the response payload containing the JWT token
type AuthTokenResponse struct {
	Token string `json:"token"`
}

// JWTClaims represents the claims for the JWT token
type JWTClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}
