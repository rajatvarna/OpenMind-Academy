package model

import "time"

// Course represents a collection of lessons.
type Course struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AuthorID    int64     `json:"author_id"` // Foreign key to the User service's users table
	IsFeatured  bool      `json:"is_featured,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Lesson represents a single educational unit within a course.
type Lesson struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	TextContent    string    `json:"text_content"`
	VideoURL       string    `json:"video_url"` // URL to the video in GCS
	TranscriptURL  string    `json:"transcript_url,omitempty"` // URL to the transcript file
	CourseID       int64     `json:"course_id"` // Foreign key to the courses table
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

// Review represents a user's review and rating for a course.
type Review struct {
	ID        int64     `json:"id"`
	CourseID  int64     `json:"course_id"`
	UserID    int64     `json:"user_id"`
	Rating    int       `json:"rating"`
	Review    string    `json:"review"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateReviewRequest defines the payload for creating a new review.
type CreateReviewRequest struct {
	CourseID int64  `json:"course_id" binding:"required"`
	UserID   int64  `json:"user_id" binding:"required"` // In a real app, this would come from the JWT
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Review   string `json:"review"`
}

// --- Learning Path Structs ---

// LearningPath defines the structure for a learning path.
type LearningPath struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Courses     []Course  `json:"courses,omitempty"` // Populated on retrieval
	CreatedAt   time.Time `json:"created_at"`
}

// CreateLearningPathRequest defines the payload for creating a new learning path.
type CreateLearningPathRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	CourseIDs   []int64 `json:"course_ids" binding:"required"` // Ordered list of course IDs
}
