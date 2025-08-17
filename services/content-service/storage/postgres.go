package storage

import (
	"context"

	"github.com/free-education/content-service/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

/*
Expected Database Schema:

CREATE TABLE IF NOT EXISTS courses (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    author_id BIGINT NOT NULL, -- This would be a foreign key in a real setup, but for microservices, we just store the ID.
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS lessons (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    text_content TEXT,
    video_url VARCHAR(255),
    transcript_url VARCHAR(255),
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS course_reviews (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL, -- Would be a FK to users table
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_id, user_id) -- A user can only review a course once
);

CREATE TABLE IF NOT EXISTS learning_paths (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS learning_path_courses (
    path_id BIGINT NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    step INTEGER NOT NULL,
    PRIMARY KEY (path_id, course_id)
);

CREATE TABLE IF NOT EXISTS quizzes (
    id BIGSERIAL PRIMARY KEY,
    lesson_id BIGINT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    questions JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (lesson_id) -- A lesson can only have one quiz
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
		RETURNING id, title, description, author_id, is_featured, created_at, updated_at
	`
	var newCourse model.Course
	err := s.db.QueryRow(ctx, query, course.Title, course.Description, authorID).Scan(
		&newCourse.ID,
		&newCourse.Title,
		&newCourse.Description,
		&newCourse.AuthorID,
		&newCourse.IsFeatured,
		&newCourse.CreatedAt,
		&newCourse.UpdatedAt,
	)
	return &newCourse, err
}

// GetCourse retrieves a single course by its ID.
func (s *ContentStore) GetCourse(ctx context.Context, courseID int64) (*model.Course, error) {
	query := `SELECT id, title, description, author_id, is_featured, created_at, updated_at FROM courses WHERE id = $1`
	var course model.Course
	err := s.db.QueryRow(ctx, query, courseID).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.AuthorID,
		&course.IsFeatured,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	return &course, err
}

// GetLesson retrieves a single lesson by its ID.
func (s *ContentStore) GetLesson(ctx context.Context, lessonID int64) (*model.Lesson, error) {
	query := `
		SELECT id, title, text_content, video_url, course_id, position, created_at, updated_at
		FROM lessons
		WHERE id = $1
	`
	var lesson model.Lesson
	err := s.db.QueryRow(ctx, query, lessonID).Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.TextContent,
		&lesson.VideoURL,
		&lesson.CourseID,
		&lesson.Position,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)
	return &lesson, err
}

// GetFeaturedCourses retrieves a list of all featured courses.
func (s *ContentStore) GetFeaturedCourses(ctx context.Context) ([]model.Course, error) {
	rows, err := s.db.Query(ctx, "SELECT id, title, description, author_id, is_featured, created_at, updated_at FROM courses WHERE is_featured = TRUE ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []model.Course
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.AuthorID, &c.IsFeatured, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

// --- Quiz Storage Functions ---

// CreateQuiz creates a new quiz in the database.
func (s *ContentStore) CreateQuiz(ctx context.Context, quiz *model.Quiz) (*model.Quiz, error) {
	query := `
		INSERT INTO quizzes (lesson_id, title, questions)
		VALUES ($1, $2, $3)
		RETURNING id, lesson_id, title, questions, created_at, updated_at
	`
	var newQuiz model.Quiz
	err := s.db.QueryRow(ctx, query, quiz.LessonID, quiz.Title, quiz.Questions).Scan(
		&newQuiz.ID,
		&newQuiz.LessonID,
		&newQuiz.Title,
		&newQuiz.Questions,
		&newQuiz.CreatedAt,
		&newQuiz.UpdatedAt,
	)
	return &newQuiz, err
}

// GetQuizByLessonID retrieves a quiz for a given lesson.
func (s *ContentStore) GetQuizByLessonID(ctx context.Context, lessonID int64) (*model.Quiz, error) {
	query := `
		SELECT id, lesson_id, title, questions, created_at, updated_at
		FROM quizzes
		WHERE lesson_id = $1
	`
	var quiz model.Quiz
	err := s.db.QueryRow(ctx, query, lessonID).Scan(
		&quiz.ID,
		&quiz.LessonID,
		&quiz.Title,
		&quiz.Questions,
		&quiz.CreatedAt,
		&quiz.UpdatedAt,
	)
	return &quiz, err
}

// DeleteCourse deletes a course and all its associated content (lessons, reviews) via cascading deletes.
func (s *ContentStore) DeleteCourse(ctx context.Context, courseID int64) error {
	query := `DELETE FROM courses WHERE id = $1`
	_, err := s.db.Exec(ctx, query, courseID)
	return err
}

// --- Learning Path Storage Functions ---

// CreateLearningPath creates a new learning path and associates courses with it in a transaction.
func (s *ContentStore) CreateLearningPath(ctx context.Context, req *model.CreateLearningPathRequest) (*model.LearningPath, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // Rollback on error

	// 1. Create the learning_paths entry
	var newPath model.LearningPath
	pathQuery := `INSERT INTO learning_paths (title, description) VALUES ($1, $2) RETURNING id, title, description, created_at`
	err = tx.QueryRow(ctx, pathQuery, req.Title, req.Description).Scan(
		&newPath.ID, &newPath.Title, &newPath.Description, &newPath.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// 2. Create the learning_path_courses entries
	for i, courseID := range req.CourseIDs {
		step := i + 1
		assocQuery := `INSERT INTO learning_path_courses (path_id, course_id, step) VALUES ($1, $2, $3)`
		_, err := tx.Exec(ctx, assocQuery, newPath.ID, courseID, step)
		if err != nil {
			return nil, err
		}
	}

	// 3. Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &newPath, nil
}

// GetLearningPathByID retrieves a learning path and its associated courses.
func (s *ContentStore) GetLearningPathByID(ctx context.Context, pathID int64) (*model.LearningPath, error) {
	// This would be a more complex query with a JOIN to get all data at once.
	// For simplicity, we'll do it in two steps.

	// 1. Get path details
	var path model.LearningPath
	pathQuery := `SELECT id, title, description, created_at FROM learning_paths WHERE id = $1`
	err := s.db.QueryRow(ctx, pathQuery, pathID).Scan(&path.ID, &path.Title, &path.Description, &path.CreatedAt)
	if err != nil {
		return nil, err
	}

	// 2. Get associated courses
	courseQuery := `
		SELECT c.id, c.title, c.description, c.author_id, c.is_featured, c.created_at, c.updated_at
		FROM courses c
		JOIN learning_path_courses lpc ON c.id = lpc.course_id
		WHERE lpc.path_id = $1
		ORDER BY lpc.step ASC
	`
	rows, err := s.db.Query(ctx, courseQuery, pathID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []model.Course
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.AuthorID, &c.IsFeatured, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	path.Courses = courses

	return &path, nil
}

// GetAllCourses retrieves a paginated list of all courses.
func (s *ContentStore) GetAllCourses(ctx context.Context, cursor int64, limit int) ([]model.Course, error) {
	query := `
		SELECT id, title, description, author_id, is_featured, created_at, updated_at
		FROM courses
		WHERE id > $1
		ORDER BY id ASC
		LIMIT $2
	`
	rows, err := s.db.Query(ctx, query, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []model.Course
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.AuthorID, &c.IsFeatured, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
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

// CreateReview adds a new course review to the database.
func (s *ContentStore) CreateReview(ctx context.Context, req *model.CreateReviewRequest) (*model.Review, error) {
	query := `
		INSERT INTO course_reviews (course_id, user_id, rating, review)
		VALUES ($1, $2, $3, $4)
		RETURNING id, course_id, user_id, rating, review, created_at
	`
	var newReview model.Review
	err := s.db.QueryRow(ctx, query, req.CourseID, req.UserID, req.Rating, req.Review).Scan(
		&newReview.ID,
		&newReview.CourseID,
		&newReview.UserID,
		&newReview.Rating,
		&newReview.Review,
		&newReview.CreatedAt,
	)
	return &newReview, err
}

// GetReviewsForCourse retrieves a paginated list of reviews for a given course.
func (s *ContentStore) GetReviewsForCourse(ctx context.Context, courseID int64, cursor int64, limit int) ([]model.Review, error) {
	query := `
		SELECT id, course_id, user_id, rating, review, created_at
		FROM course_reviews
		WHERE course_id = $1 AND id > $2
		ORDER BY id ASC
		LIMIT $3
	`
	rows, err := s.db.Query(ctx, query, courseID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var review model.Review
		if err := rows.Scan(&review.ID, &review.CourseID, &review.UserID, &review.Rating, &review.Review, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// UpdateLessonTranscript updates the transcript_url for a specific lesson.
func (s *ContentStore) UpdateLessonTranscript(ctx context.Context, lessonID int64, transcriptURL string) error {
	query := `UPDATE lessons SET transcript_url = $1, updated_at = NOW() WHERE id = $2`
	_, err := s.db.Exec(ctx, query, transcriptURL, lessonID)
	return err
}

// GetCoursesForUser retrieves all courses created by a specific user.
func (s *ContentStore) GetCoursesForUser(ctx context.Context, userID int64) ([]model.Course, error) {
	rows, err := s.db.Query(ctx, "SELECT id, title, description, author_id, is_featured, created_at, updated_at FROM courses WHERE author_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []model.Course
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.AuthorID, &c.IsFeatured, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}
