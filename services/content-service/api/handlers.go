package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/free-education/content-service/model"
	"github.com/gin-gonic/gin"
)

// API holds the dependencies for the API handlers.
type API struct {
	ContentStore ContentStore
	QnAServiceURL string
}

// NewAPI creates a new API struct.
func NewAPI(store ContentStore, qnaServiceURL string) *API {
	return &API{ContentStore: store, QnAServiceURL: qnaServiceURL}
}

// CreateCourseHandler handles the creation of a new course.
// The author ID is taken from the authenticated user's context.
func (a *API) CreateCourseHandler(c *gin.Context) {
	var req model.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	authorID := c.MustGet("userID").(int64)

	course, err := a.ContentStore.CreateCourse(c.Request.Context(), &req, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// GetCourseHandler handles retrieving a single course and its associated lessons.
// This is a public endpoint.
func (a *API) GetCourseHandler(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := a.ContentStore.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		// This could be sql.ErrNoRows, which we treat as a 404.
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
// It includes an authorization check to ensure the user is the course author.
func (a *API) CreateLessonHandler(c *gin.Context) {
	var req model.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	userID := c.MustGet("userID").(int64)

	// Authorization check: Ensure the user is the author of the course.
	course, err := a.ContentStore.GetCourse(c.Request.Context(), req.CourseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	if course.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add a lesson to this course"})
		return
	}

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

	userID := c.MustGet("userID").(int64)

	review, err := a.ContentStore.CreateReview(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// getPaginationParams is a helper function to parse cursor and limit from query params.
// It sets default values and enforces a maximum limit to prevent abuse.
func getPaginationParams(c *gin.Context, defaultLimit int) (int64, int) {
	cursor, _ := strconv.ParseInt(c.Query("cursor"), 10, 64)
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}
	// Enforce a maximum limit to prevent clients from requesting too much data.
	if limit > 100 {
		limit = 100
	}
	return cursor, limit
}

// GetReviewsHandler handles fetching a paginated list of reviews for a specific course.
func (a *API) GetReviewsHandler(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	cursor, limit := getPaginationParams(c, 5) // Default limit of 5 for reviews

	reviews, err := a.ContentStore.GetReviewsForCourse(c.Request.Context(), courseID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews for course"})
		return
	}

	var nextCursor int64 = 0
	if len(reviews) > 0 {
		nextCursor = reviews[len(reviews)-1].ID
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       reviews,
		"next_cursor": nextCursor,
	})
}

// GetFeaturedCoursesHandler handles fetching all featured courses.
// This is a public endpoint.
func (a *API) GetFeaturedCoursesHandler(c *gin.Context) {
	courses, err := a.ContentStore.GetFeaturedCourses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get featured courses"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// --- Learning Path Handlers ---

// CreateLearningPathHandler handles the creation of a new learning path.
// This is an authenticated endpoint.
func (a *API) CreateLearningPathHandler(c *gin.Context) {
	var req model.CreateLearningPathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Note: In a real app, you might want to add an author_id to learning paths
	// and perform an authorization check here.
	path, err := a.ContentStore.CreateLearningPath(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create learning path"})
		return
	}
	c.JSON(http.StatusCreated, path)
}

// GetLearningPathHandler retrieves a learning path and its courses.
// This is a public endpoint.
func (a *API) GetLearningPathHandler(c *gin.Context) {
	pathID, err := strconv.ParseInt(c.Param("pathId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path ID"})
		return
	}
	path, err := a.ContentStore.GetLearningPathByID(c.Request.Context(), pathID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Learning path not found"})
		return
	}
	c.JSON(http.StatusOK, path)
}

// DeleteCourseHandler handles deleting a course.
// It includes an authorization check to ensure only the course author can delete it.
func (a *API) DeleteCourseHandler(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	userID := c.MustGet("userID").(int64)

	// Authorization check: Get the course and verify the author.
	course, err := a.ContentStore.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	if course.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this course"})
		return
	}

	err = a.ContentStore.DeleteCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllCoursesHandler handles fetching a paginated list of all courses.
func (a *API) GetAllCoursesHandler(c *gin.Context) {
	cursor, limit := getPaginationParams(c, 10) // Default limit of 10 for courses

	courses, err := a.ContentStore.GetAllCourses(c.Request.Context(), cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	var nextCursor int64 = 0
	if len(courses) > 0 {
		nextCursor = courses[len(courses)-1].ID
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       courses,
		"next_cursor": nextCursor,
	})
}

// UpdateTranscriptHandler handles updating the transcript URL for a lesson.
// It includes an authorization check to ensure the user is the course author.
func (a *API) UpdateTranscriptHandler(c *gin.Context) {
	lessonID, err := strconv.ParseInt(c.Param("lessonId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	userID := c.MustGet("userID").(int64)

	// Authorization check: Fetch the lesson, then the course, then check the author.
	lesson, err := a.ContentStore.GetLesson(c.Request.Context(), lessonID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}
	course, err := a.ContentStore.GetCourse(c.Request.Context(), lesson.CourseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	if course.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this lesson"})
		return
	}

	var req struct {
		TranscriptURL string `json:"transcript_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = a.ContentStore.UpdateLessonTranscript(c.Request.Context(), lessonID, req.TranscriptURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transcript URL"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCoursesForUserHandler handles fetching all courses created by a specific user.
// This is a public endpoint.
func (a *API) GetCoursesForUserHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	courses, err := a.ContentStore.GetCoursesForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses for user"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// --- Quiz Handlers ---

// CreateQuizHandler generates and saves a new quiz for a lesson.
func (a *API) CreateQuizHandler(c *gin.Context) {
	var req model.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 1. Get the lesson content to generate the quiz from
	lesson, err := a.ContentStore.GetLesson(c.Request.Context(), req.LessonID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	// 2. Call the Q&A service to generate the quiz
	qnaReqBody := map[string]string{"text_content": lesson.TextContent}
	reqBytes, _ := json.Marshal(qnaReqBody)

	resp, err := http.Post(a.QnAServiceURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Q&A service"})
		return
	}
	defer resp.Body.Close()

	var generatedQuiz struct {
		Questions []model.QuizQuestion `json:"questions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&generatedQuiz); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode quiz from Q&A service"})
		return
	}

	// 3. Save the generated quiz to our database
	newQuiz := &model.Quiz{
		LessonID:  req.LessonID,
		Title:     req.Title,
		Questions: generatedQuiz.Questions,
	}

	savedQuiz, err := a.ContentStore.CreateQuiz(c.Request.Context(), newQuiz)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save generated quiz"})
		return
	}

	c.JSON(http.StatusCreated, savedQuiz)
}

// GetQuizByLessonIDHandler retrieves the quiz for a specific lesson.
func (a *API) GetQuizByLessonIDHandler(c *gin.Context) {
	lessonID, err := strconv.ParseInt(c.Param("lessonId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	quiz, err := a.ContentStore.GetQuizByLessonID(c.Request.Context(), lessonID)
	if err != nil {
		// Could be sql.ErrNoRows
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found for this lesson"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}
