package model

import "time"

// Thread represents a discussion thread associated with a course.
type Thread struct {
	ID        int64     `json:"id"`
	CourseID  int64     `json:"course_id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// Post represents a single post or reply within a thread.
type Post struct {
	ID        int64     `json:"id"`
	ThreadID  int64     `json:"thread_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// --- API Request Structs ---

// CreateThreadRequest defines the payload for creating a new thread.
type CreateThreadRequest struct {
	CourseID int64  `json:"course_id" binding:"required"`
	UserID   int64  `json:"user_id" binding:"required"` // From JWT
	Title    string `json:"title" binding:"required,min=10"`
	// The initial post content will be sent in a separate request for simplicity
}

// CreatePostRequest defines the payload for creating a new post (reply).
type CreatePostRequest struct {
	ThreadID int64  `json:"thread_id" binding:"required"`
	UserID   int64  `json:"user_id" binding:"required"` // From JWT
	Content  string `json:"content" binding:"required,min=1"`
}
