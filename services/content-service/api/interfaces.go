package api

import (
	"context"

	"github.com/free-education/content-service/model"
)

// ContentStore defines the interface for content storage operations.
// This allows us to mock the storage layer for testing.
type ContentStore interface {
	CreateCourse(ctx context.Context, course *model.CreateCourseRequest, authorID int64) (*model.Course, error)
	GetCourse(ctx context.Context, courseID int64) (*model.Course, error)
	DeleteCourse(ctx context.Context, courseID int64) error
	CreateLesson(ctx context.Context, lesson *model.CreateLessonRequest) (*model.Lesson, error)
	GetLessonsByCourse(ctx context.Context, courseID int64) ([]model.Lesson, error)
	CreateReview(ctx context.Context, req *model.CreateReviewRequest, userID int64) (*model.Review, error)
	GetReviewsForCourse(ctx context.Context, courseID int64, cursor int64, limit int) ([]model.Review, error)
	GetFeaturedCourses(ctx context.Context) ([]model.Course, error)
	CreateLearningPath(ctx context.Context, req *model.CreateLearningPathRequest) (*model.LearningPath, error)
	GetLearningPathByID(ctx context.Context, pathID int64) (*model.LearningPath, error)
	GetAllCourses(ctx context.Context, cursor int64, limit int) ([]model.Course, error)
	UpdateLessonTranscript(ctx context.Context, lessonID int64, transcriptURL string) error
	GetCoursesForUser(ctx context.Context, userID int64) ([]model.Course, error)
	GetLesson(ctx context.Context, lessonID int64) (*model.Lesson, error)
	CreateQuiz(ctx context.Context, quiz *model.Quiz) (*model.Quiz, error)
	GetQuizByLessonID(ctx context.Context, lessonID int64) (*model.Quiz, error)
}
