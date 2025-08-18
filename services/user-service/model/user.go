package model

import "time"

// User represents a user in the system.
// The `json` tags define how the struct is serialized to/from JSON.
type User struct {
	// The unique identifier for the user.
	ID int64 `json:"id"`
	// The user's email address, used for login. Must be unique.
	Email string `json:"email"`
	// The salted and hashed password. It is never exposed to the client (`json:"-"`).
	PasswordHash string `json:"-"`
	// The user's first name.
	FirstName string `json:"first_name"`
	// The user's last name.
	LastName string `json:"last_name"`
	// The URL to the user's profile picture. Can be empty.
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
	// The role of the user (e.g., 'user', 'admin'). Determines permissions.
	Role string `json:"role"`
	// A flexible JSONB field for storing user-specific settings, like theme.
	Preferences map[string]interface{} `json:"preferences"`
	// The timestamp when the user was created.
	CreatedAt time.Time `json:"created_at"`
	// The timestamp when the user was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// RegistrationRequest represents the data required to register a new user.
// The `binding` tags are used by Gin for request validation.
type RegistrationRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest represents the data required for a user to log in.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the payload sent back to the client after a successful login.
type LoginResponse struct {
	// The JWT token used for authenticating subsequent requests.
	Token string `json:"token"`
}

// --- Quiz Attempt Structs ---

// QuizAttempt represents a record of a user's attempt at a quiz.
type QuizAttempt struct {
	// The unique identifier for this quiz attempt.
	ID int64 `json:"id"`
	// The ID of the user who made the attempt.
	UserID int64 `json:"user_id"`
	// The ID of the quiz that was attempted. This links to the content service.
	QuizID int64 `json:"quiz_id"`
	// The score the user achieved on the quiz.
	Score int `json:"score"`
	// A JSON string containing the user's answers.
	Answers string `json:"answers"`
	// The timestamp when the attempt was recorded.
	CreatedAt time.Time `json:"created_at"`
}

// CreateQuizAttemptRequest defines the payload for submitting a new quiz attempt.
type CreateQuizAttemptRequest struct {
	QuizID  int64  `json:"quiz_id" binding:"required"`
	Score   int    `json:"score" binding:"required"`
	Answers string `json:"answers" binding:"required"`
}
