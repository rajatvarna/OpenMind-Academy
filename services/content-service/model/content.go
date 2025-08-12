package model

import "time"

// Course represents a collection of lessons.
type Course struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AuthorID    int64     `json:"author_id"` // Foreign key to the User service's users table
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Lesson represents a single educational unit within a course.
type Lesson struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	TextContent string    `json:"text_content"`
	VideoURL    string    `json:"video_url"` // URL to the video in GCS
	CourseID    int64     `json:"course_id"` // Foreign key to the courses table
	Position    int       `json:"position"`  // For ordering lessons within a course
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// --- API Request/Response Structs ---

// CreateCourseRequest defines the payload for creating a new course.
type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required,min=5"`
	Description string `json:"description" binding:"required,min=10"`
}

// CreateLessonRequest defines the payload for creating a new lesson.
type CreateLessonRequest struct {
	Title       string `json:"title" binding:"required,min=5"`
	TextContent string `json:"text_content" binding:"required"`
	CourseID    int64  `json:"course_id" binding:"required"`
	Position    int    `json:"position"` // Optional, can be auto-managed
}
