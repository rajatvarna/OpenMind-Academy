package storage

import (
	"context"
	"github.com/free-education/forum-service/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

/*
Expected Database Schema:

CREATE TABLE IF NOT EXISTS threads (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    thread_id BIGINT NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

*/

// ForumStore handles database operations for the forum.
type ForumStore struct {
	db *pgxpool.Pool
}

// NewForumStore creates a new ForumStore.
func NewForumStore(db *pgxpool.Pool) *ForumStore {
	return &ForumStore{db: db}
}

// CreateThread creates a new discussion thread.
func (s *ForumStore) CreateThread(ctx context.Context, req *model.CreateThreadRequest) (*model.Thread, error) {
	query := `
		INSERT INTO threads (course_id, user_id, title) VALUES ($1, $2, $3)
		RETURNING id, course_id, user_id, title, created_at
	`
	var newThread model.Thread
	err := s.db.QueryRow(ctx, query, req.CourseID, req.UserID, req.Title).Scan(
		&newThread.ID, &newThread.CourseID, &newThread.UserID, &newThread.Title, &newThread.CreatedAt,
	)
	return &newThread, err
}

// CreatePost creates a new post in a thread.
func (s *ForumStore) CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {
	query := `
		INSERT INTO posts (thread_id, user_id, content) VALUES ($1, $2, $3)
		RETURNING id, thread_id, user_id, content, created_at
	`
	var newPost model.Post
	err := s.db.QueryRow(ctx, query, req.ThreadID, req.UserID, req.Content).Scan(
		&newPost.ID, &newPost.ThreadID, &newPost.UserID, &newPost.Content, &newPost.CreatedAt,
	)
	return &newPost, err
}

// GetThreadsForCourse retrieves all threads for a given course.
func (s *ForumStore) GetThreadsForCourse(ctx context.Context, courseID int64) ([]model.Thread, error) {
	rows, err := s.db.Query(ctx, "SELECT id, course_id, user_id, title, created_at FROM threads WHERE course_id = $1 ORDER BY created_at DESC", courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []model.Thread
	for rows.Next() {
		var t model.Thread
		if err := rows.Scan(&t.ID, &t.CourseID, &t.UserID, &t.Title, &t.CreatedAt); err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}
	return threads, nil
}

// GetPostsForThread retrieves all posts for a given thread.
func (s *ForumStore) GetPostsForThread(ctx context.Context, threadID int64) ([]model.Post, error) {
	rows, err := s.db.Query(ctx, "SELECT id, thread_id, user_id, content, created_at FROM posts WHERE thread_id = $1 ORDER BY created_at ASC", threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.ThreadID, &p.UserID, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
