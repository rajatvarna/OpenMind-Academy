package storage

import (
	"context"
	"log"

	"github.com/free-education/user-service/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

/*
Expected Database Schema:

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_lesson_progress (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id BIGINT NOT NULL, -- In a monolith, this would be a foreign key to lessons.
    completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, lesson_id)
);

*/

// UserStore handles database operations for users.
type UserStore struct {
	db *pgxpool.Pool
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

// CreateUser creates a new user in the database after hashing their password.
func (s *UserStore) CreateUser(ctx context.Context, userReq *model.RegistrationRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, first_name, last_name, created_at, updated_at
	`

	var newUser model.User
	err = s.db.QueryRow(ctx, query, userReq.Email, string(hashedPassword), userReq.FirstName, userReq.LastName).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return &newUser, nil
}

// GetUserByEmail retrieves a user by their email address.
func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, created_at, updated_at
		FROM users WHERE email = $1
	`
	var user model.User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		// It's common for "no rows" error to occur, which isn't a server error.
		// The handler should decide what to do with this.
		return nil, err
	}

	return &user, nil
}

// CheckPassword compares a plaintext password with the stored hash.
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// MarkLessonAsComplete marks a lesson as completed for a user.
// It uses an 'ON CONFLICT DO NOTHING' clause to handle cases where the entry already exists.
func (s *UserStore) MarkLessonAsComplete(ctx context.Context, userID int64, lessonID int64) error {
	query := `
		INSERT INTO user_lesson_progress (user_id, lesson_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, lesson_id) DO NOTHING
	`
	_, err := s.db.Exec(ctx, query, userID, lessonID)
	return err
}

// GetCompletedLessonsForUser retrieves a list of completed lesson IDs for a user.
func (s *UserStore) GetCompletedLessonsForUser(ctx context.Context, userID int64) ([]int64, error) {
	query := `SELECT lesson_id FROM user_lesson_progress WHERE user_id = $1`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var completedLessonIDs []int64
	for rows.Next() {
		var lessonID int64
		if err := rows.Scan(&lessonID); err != nil {
			return nil, err
		}
		completedLessonIDs = append(completedLessonIDs, lessonID)
	}

	return completedLessonIDs, nil
}
