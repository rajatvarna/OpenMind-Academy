package api

import (
	"net/http"
	"strconv"

	"github.com/free-education/content-service/model"
	"github.com/free-education/content-service/storage"
	"github.com/gin-gonic/gin"
)

// API holds the dependencies for the API handlers.
type API struct {
	ContentStore *storage.ContentStore
}

// NewAPI creates a new API struct.
func NewAPI(store *storage.ContentStore) *API {
	return &API{ContentStore: store}
}

// CreateCourseHandler handles the creation of a new course.
func (a *API) CreateCourseHandler(c *gin.Context) {
	var req model.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// In a real application, authorID would come from a JWT middleware.
	// c.GetInt64("userID")
	// For now, we'll use a placeholder or expect it in the request for testing.
	authorID := int64(1) // Placeholder author ID

	course, err := a.ContentStore.CreateCourse(c.Request.Context(), &req, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// GetCourseHandler handles retrieving a single course and its lessons.
func (a *API) GetCourseHandler(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := a.ContentStore.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		// Could be sql.ErrNoRows
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	lessons, err := a.ContentStore.GetLessonsByCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve lessons for course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"course":  course,
		"lessons": lessons,
	})
}

// CreateLessonHandler handles the creation of a new lesson.
func (a *API) CreateLessonHandler(c *gin.Context) {
	var req model.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Here you would add validation to ensure the author of the lesson
	// is the same as the author of the course.

	lesson, err := a.ContentStore.CreateLesson(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lesson"})
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

// CreateReviewHandler handles submitting a new review for a course.
func (a *API) CreateReviewHandler(c *gin.Context) {
	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	// In a real app, req.UserID would be populated from the JWT claims, not the request body.

	review, err := a.ContentStore.CreateReview(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetReviewsHandler handles fetching all reviews for a specific course.
func (a *API) GetReviewsHandler(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	reviews, err := a.ContentStore.GetReviewsForCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews for course"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// GetFeaturedCoursesHandler handles fetching all featured courses.
func (a *API) GetFeaturedCoursesHandler(c *gin.Context) {
	courses, err := a.store.GetFeaturedCourses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get featured courses"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// UpdateTranscriptHandler handles updating the transcript URL for a lesson.
func (a *API) UpdateTranscriptHandler(c *gin.Context) {
	lessonID, err := strconv.ParseInt(c.Param("lessonId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	var req struct {
		TranscriptURL string `json:"transcript_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = a.store.UpdateLessonTranscript(c.Request.Context(), lessonID, req.TranscriptURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transcript URL"})
		return
	}

	c.Status(http.StatusNoContent)
}
