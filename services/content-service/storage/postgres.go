package storage

import (
	"context"

	"github.com/free-education/content-service/model"
	"github.comcom/jackc/pgx/v4/pgxpool"
)

/*
Expected Database Schema:

CREATE TABLE IF NOT EXISTS courses (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    author_id BIGINT NOT NULL, -- This would be a foreign key in a real setup, but for microservices, we just store the ID.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS lessons (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    text_content TEXT,
    video_url VARCHAR(255),
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

*/

// ContentStore handles database operations for content.
type ContentStore struct {
	db *pgxpool.Pool
}

// NewContentStore creates a new ContentStore.
func NewContentStore(db *pgxpool.Pool) *ContentStore {
	return &ContentStore{db: db}
}

// CreateCourse creates a new course in the database.
func (s *ContentStore) CreateCourse(ctx context.Context, course *model.CreateCourseRequest, authorID int64) (*model.Course, error) {
	query := `
		INSERT INTO courses (title, description, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, description, author_id, created_at, updated_at
	`
	var newCourse model.Course
	err := s.db.QueryRow(ctx, query, course.Title, course.Description, authorID).Scan(
		&newCourse.ID,
		&newCourse.Title,
		&newCourse.Description,
		&newCourse.AuthorID,
		&newCourse.CreatedAt,
		&newCourse.UpdatedAt,
	)
	return &newCourse, err
}

// GetCourse retrieves a single course by its ID.
func (s *ContentStore) GetCourse(ctx context.Context, courseID int64) (*model.Course, error) {
	query := `SELECT id, title, description, author_id, created_at, updated_at FROM courses WHERE id = $1`
	var course model.Course
	err := s.db.QueryRow(ctx, query, courseID).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.AuthorID,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	return &course, err
}

// CreateLesson creates a new lesson in the database.
func (s *ContentStore) CreateLesson(ctx context.Context, lesson *model.CreateLessonRequest) (*model.Lesson, error) {
	query := `
		INSERT INTO lessons (title, text_content, course_id, position)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, text_content, video_url, course_id, position, created_at, updated_at
	`
	var newLesson model.Lesson
	err := s.db.QueryRow(ctx, query, lesson.Title, lesson.TextContent, lesson.CourseID, lesson.Position).Scan(
		&newLesson.ID,
		&newLesson.Title,
		&newLesson.TextContent,
		&newLesson.VideoURL,
		&newLesson.CourseID,
		&newLesson.Position,
		&newLesson.CreatedAt,
		&newLesson.UpdatedAt,
	)
	return &newLesson, err
}

// GetLessonsByCourse retrieves all lessons for a given course, ordered by position.
func (s *ContentStore) GetLessonsByCourse(ctx context.Context, courseID int64) ([]model.Lesson, error) {
	query := `
		SELECT id, title, text_content, video_url, course_id, position, created_at, updated_at
		FROM lessons
		WHERE course_id = $1
		ORDER BY position ASC
	`
	rows, err := s.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []model.Lesson
	for rows.Next() {
		var lesson model.Lesson
		if err := rows.Scan(&lesson.ID, &lesson.Title, &lesson.TextContent, &lesson.VideoURL, &lesson.CourseID, &lesson.Position, &lesson.CreatedAt, &lesson.UpdatedAt); err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}
