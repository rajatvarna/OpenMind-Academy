package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/free-education/content-service/model"
	"github.com/gin-gonic/gin"
)

// MockContentStore is a mock implementation of the ContentStore for testing.
type MockContentStore struct {
	CreateCourseFunc           func(ctx context.Context, course *model.CreateCourseRequest, authorID int64) (*model.Course, error)
	GetCourseFunc              func(ctx context.Context, courseID int64) (*model.Course, error)
	DeleteCourseFunc           func(ctx context.Context, courseID int64) error
	CreateLessonFunc           func(ctx context.Context, lesson *model.CreateLessonRequest) (*model.Lesson, error)
	GetLessonsByCourseFunc     func(ctx context.Context, courseID int64) ([]model.Lesson, error)
	CreateReviewFunc           func(ctx context.Context, req *model.CreateReviewRequest, userID int64) (*model.Review, error)
	GetReviewsForCourseFunc    func(ctx context.Context, courseID int64, cursor int64, limit int) ([]model.Review, error)
	GetFeaturedCoursesFunc     func(ctx context.Context) ([]model.Course, error)
	CreateLearningPathFunc     func(ctx context.Context, req *model.CreateLearningPathRequest) (*model.LearningPath, error)
	GetLearningPathByIDFunc    func(ctx context.Context, pathID int64) (*model.LearningPath, error)
	GetAllCoursesFunc          func(ctx context.Context, cursor int64, limit int) ([]model.Course, error)
	UpdateLessonTranscriptFunc func(ctx context.Context, lessonID int64, transcriptURL string) error
	GetCoursesForUserFunc      func(ctx context.Context, userID int64) ([]model.Course, error)
	GetLessonFunc              func(ctx context.Context, lessonID int64) (*model.Lesson, error)
	CreateQuizFunc             func(ctx context.Context, quiz *model.Quiz) (*model.Quiz, error)
	GetQuizByLessonIDFunc      func(ctx context.Context, lessonID int64) (*model.Quiz, error)
}


func (m *MockContentStore) CreateCourse(ctx context.Context, course *model.CreateCourseRequest, authorID int64) (*model.Course, error) {
	return m.CreateCourseFunc(ctx, course, authorID)
}

func (m *MockContentStore) GetCourse(ctx context.Context, courseID int64) (*model.Course, error) {
	return m.GetCourseFunc(ctx, courseID)
}

func (m *MockContentStore) DeleteCourse(ctx context.Context, courseID int64) error {
	return m.DeleteCourseFunc(ctx, courseID)
}

func (m *MockContentStore) CreateLesson(ctx context.Context, lesson *model.CreateLessonRequest) (*model.Lesson, error) {
	return m.CreateLessonFunc(ctx, lesson)
}

func (m *MockContentStore) GetLessonsByCourse(ctx context.Context, courseID int64) ([]model.Lesson, error) {
	return m.GetLessonsByCourseFunc(ctx, courseID)
}

func (m *MockContentStore) CreateReview(ctx context.Context, req *model.CreateReviewRequest, userID int64) (*model.Review, error) {
	return m.CreateReviewFunc(ctx, req, userID)
}

func (m *MockContentStore) GetReviewsForCourse(ctx context.Context, courseID int64, cursor int64, limit int) ([]model.Review, error) {
	return m.GetReviewsForCourseFunc(ctx, courseID, cursor, limit)
}

func (m *MockContentStore) GetFeaturedCourses(ctx context.Context) ([]model.Course, error) {
	return m.GetFeaturedCoursesFunc(ctx)
}

func (m *MockContentStore) CreateLearningPath(ctx context.Context, req *model.CreateLearningPathRequest) (*model.LearningPath, error) {
	return m.CreateLearningPathFunc(ctx, req)
}

func (m *MockContentStore) GetLearningPathByID(ctx context.Context, pathID int64) (*model.LearningPath, error) {
	return m.GetLearningPathByIDFunc(ctx, pathID)
}

func (m *MockContentStore) GetAllCourses(ctx context.Context, cursor int64, limit int) ([]model.Course, error) {
	return m.GetAllCoursesFunc(ctx, cursor, limit)
}

func (m *MockContentStore) UpdateLessonTranscript(ctx context.Context, lessonID int64, transcriptURL string) error {
	return m.UpdateLessonTranscriptFunc(ctx, lessonID, transcriptURL)
}

func (m *MockContentStore) GetCoursesForUser(ctx context.Context, userID int64) ([]model.Course, error) {
	return m.GetCoursesForUserFunc(ctx, userID)
}

func (m *MockContentStore) GetLesson(ctx context.Context, lessonID int64) (*model.Lesson, error) {
	return m.GetLessonFunc(ctx, lessonID)
}

func (m *MockContentStore) CreateQuiz(ctx context.Context, quiz *model.Quiz) (*model.Quiz, error) {
	return m.CreateQuizFunc(ctx, quiz)
}

func (m *MockContentStore) GetQuizByLessonID(ctx context.Context, lessonID int64) (*model.Quiz, error) {
	return m.GetQuizByLessonIDFunc(ctx, lessonID)
}

func TestCreateCourseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful course creation", func(t *testing.T) {
		mockStore := &MockContentStore{
			CreateCourseFunc: func(ctx context.Context, course *model.CreateCourseRequest, authorID int64) (*model.Course, error) {
				return &model.Course{
					ID:          1,
					Title:       course.Title,
					Description: course.Description,
					AuthorID:    authorID,
				}, nil
			},
		}

		apiHandler := NewAPI(mockStore, "") // QnAServiceURL not needed for this test

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the authenticated user ID in the context
		c.Set("userID", int64(123))

		body := map[string]string{
			"title":       "Test Course",
			"description": "A great course.",
		}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/courses", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.CreateCourseHandler(c)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d; got %d", http.StatusCreated, w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["title"] != "Test Course" {
			t.Errorf("expected title %s; got %s", "Test Course", response["title"])
		}
		if response["author_id"] != float64(123) {
			t.Errorf("expected author_id %d; got %v", 123, response["author_id"])
		}
	})
}

func TestCreateLessonHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful lesson creation by author", func(t *testing.T) {
		mockStore := &MockContentStore{
			GetCourseFunc: func(ctx context.Context, courseID int64) (*model.Course, error) {
				return &model.Course{ID: 1, AuthorID: 123}, nil
			},
			CreateLessonFunc: func(ctx context.Context, lesson *model.CreateLessonRequest) (*model.Lesson, error) {
				return &model.Lesson{ID: 1, Title: lesson.Title, CourseID: lesson.CourseID}, nil
			},
		}
		apiHandler := NewAPI(mockStore, "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", int64(123))

		body := map[string]interface{}{
			"title":      "Test Lesson",
			"course_id":  1,
			"text_content": "Lesson content.",
		}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/lessons", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.CreateLessonHandler(c)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d; got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("Forbidden lesson creation by non-author", func(t *testing.T) {
		mockStore := &MockContentStore{
			GetCourseFunc: func(ctx context.Context, courseID int64) (*model.Course, error) {
				return &model.Course{ID: 1, AuthorID: 999}, nil // Different author
			},
		}
		apiHandler := NewAPI(mockStore, "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", int64(123))

		body := map[string]interface{}{
			"title":      "Test Lesson",
			"course_id":  1,
			"text_content": "Lesson content.",
		}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/lessons", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.CreateLessonHandler(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status %d; got %d", http.StatusForbidden, w.Code)
		}
	})
}

func TestDeleteCourseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful course deletion by author", func(t *testing.T) {
		mockStore := &MockContentStore{
			GetCourseFunc: func(ctx context.Context, courseID int64) (*model.Course, error) {
				return &model.Course{ID: 1, AuthorID: 123}, nil
			},
			DeleteCourseFunc: func(ctx context.Context, courseID int64) error {
				return nil
			},
		}
		apiHandler := NewAPI(mockStore, "")

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("userID", int64(123))
			c.Next()
		})
		router.DELETE("/api/v1/courses/:courseId", apiHandler.DeleteCourseHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/courses/1", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d; got %d", http.StatusNoContent, w.Code)
		}
	})

	t.Run("Forbidden course deletion by non-author", func(t *testing.T) {
		mockStore := &MockContentStore{
			GetCourseFunc: func(ctx context.Context, courseID int64) (*model.Course, error) {
				return &model.Course{ID: 1, AuthorID: 999}, nil // Different author
			},
		}
		apiHandler := NewAPI(mockStore, "")

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("userID", int64(123))
			c.Next()
		})
		router.DELETE("/api/v1/courses/:courseId", apiHandler.DeleteCourseHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/courses/1", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status %d; got %d", http.StatusForbidden, w.Code)
		}
	})
}
