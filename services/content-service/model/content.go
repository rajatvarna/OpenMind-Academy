package model

import "time"

// Course represents a collection of lessons, forming an educational module.
type Course struct {
	// The unique identifier for the course.
	ID int64 `json:"id"`
	// The title of the course.
	Title string `json:"title"`
	// A detailed description of the course content and objectives.
	Description string `json:"description"`
	// The ID of the user who authored the course.
	AuthorID int64 `json:"author_id"`
	// A flag to indicate if the course should be featured on the homepage.
	IsFeatured bool `json:"is_featured,omitempty"`
	// The timestamp when the course was created.
	CreatedAt time.Time `json:"created_at"`
	// The timestamp when the course was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// Lesson represents a single educational unit within a course.
type Lesson struct {
	// The unique identifier for the lesson.
	ID int64 `json:"id"`
	// The title of the lesson.
	Title string `json:"title"`
	// The main text content of the lesson.
	TextContent string `json:"text_content"`
	// The URL to an optional video for the lesson, likely stored in a cloud bucket.
	VideoURL string `json:"video_url"`
	// The URL to an optional transcript file for the video.
	TranscriptURL string `json:"transcript_url,omitempty"`
	// The ID of the course this lesson belongs to.
	CourseID int64 `json:"course_id"`
	// The numerical position of the lesson within the course for ordering.
	Position int `json:"position"`
	// The timestamp when the lesson was created.
	CreatedAt time.Time `json:"created_at"`
	// The timestamp when the lesson was last updated.
	UpdatedAt time.Time `json:"updated_at"`
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
	// The unique identifier for the review.
	ID int64 `json:"id"`
	// The ID of the course being reviewed.
	CourseID int64 `json:"course_id"`
	// The ID of the user who wrote the review.
	UserID int64 `json:"user_id"`
	// The rating given by the user, from 1 to 5.
	Rating int `json:"rating"`
	// The text content of the review.
	Review string `json:"review"`
	// The timestamp when the review was created.
	CreatedAt time.Time `json:"created_at"`
}

// CreateReviewRequest defines the payload for creating a new review.
type CreateReviewRequest struct {
	CourseID int64  `json:"course_id" binding:"required"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Review   string `json:"review"`
}

// --- Learning Path Structs ---

// LearningPath is a curated sequence of courses designed to guide a user through a topic.
type LearningPath struct {
	// The unique identifier for the learning path.
	ID int64 `json:"id"`
	// The title of the learning path.
	Title string `json:"title"`
	// A description of what the learning path covers.
	Description string `json:"description"`
	// The list of courses in the path, ordered by step. This is populated on retrieval.
	Courses []Course `json:"courses,omitempty"`
	// The timestamp when the learning path was created.
	CreatedAt time.Time `json:"created_at"`
}

// CreateLearningPathRequest defines the payload for creating a new learning path.
type CreateLearningPathRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	// An ordered list of course IDs that make up the path.
	CourseIDs []int64 `json:"course_ids" binding:"required"`
}

// --- Quiz Structs ---

// Quiz represents a set of questions associated with a lesson to test understanding.
type Quiz struct {
	// The unique identifier for the quiz.
	ID int64 `json:"id"`
	// The ID of the lesson this quiz is for.
	LessonID int64 `json:"lesson_id"`
	// The title of the quiz.
	Title string `json:"title"`
	// The list of questions that make up the quiz.
	Questions []QuizQuestion `json:"questions"`
	// The timestamp when the quiz was created.
	CreatedAt time.Time `json:"created_at"`
	// The timestamp when the quiz was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// QuizQuestion represents a single question within a quiz.
type QuizQuestion struct {
	// The type of question (e.g., 'multiple-choice', 'true-false').
	Type string `json:"type"`
	// The text of the question.
	Question string `json:"question"`
	// A list of possible answers for multiple-choice questions.
	Options []string `json:"options,omitempty"`
	// The correct answer to the question.
	CorrectAnswer string `json:"correct_answer"`
}

// CreateQuizRequest defines the payload for creating a new quiz.
type CreateQuizRequest struct {
	LessonID int64  `json:"lesson_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
}
