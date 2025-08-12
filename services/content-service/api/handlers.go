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
