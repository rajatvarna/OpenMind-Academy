package model

import "time"

// User represents a user in the database.
// The `json` tags control how the struct is serialized to JSON for API responses.
// Note that PasswordHash is omitted from JSON responses for security.
type User struct {
	ID           int64                  `json:"id"`
	Email        string                 `json:"email"`
	PasswordHash string                 `json:"-"` // Omit from JSON output
	FirstName    string                 `json:"first_name"`
	LastName     string                 `json:"last_name"`
	Role         string                 `json:"role"`
	Preferences  map[string]interface{} `json:"preferences"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// RegistrationRequest represents the payload for a user registration request.
type RegistrationRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest represents the payload for a user login request.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the payload sent back after a successful login.
type LoginResponse struct {
	Token string `json:"token"`
}
