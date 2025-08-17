package storage

import (
	"context"
	"log"
	"time"

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
    role VARCHAR(50) NOT NULL DEFAULT 'user', -- 'user', 'moderator', 'admin'
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_lesson_progress (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id BIGINT NOT NULL, -- In a monolith, this would be a foreign key to lessons.
    completed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, lesson_id)
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    token TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS quiz_attempts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quiz_id BIGINT NOT NULL, -- Foreign key to content service's quizzes table
    score INTEGER NOT NULL,
    answers JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

*/

// PostgresUserStore handles database operations for users.
type PostgresUserStore struct {
	db *pgxpool.Pool
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *pgxpool.Pool) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

// CreateUser creates a new user in the database after hashing their password.
func (s *PostgresUserStore) CreateUser(ctx context.Context, userReq *model.RegistrationRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, preferences)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, first_name, last_name, role, preferences, created_at, updated_at
	`

	// Default preferences
	defaultPrefs := map[string]interface{}{"theme": "light"}

	var newUser model.User
	err = s.db.QueryRow(ctx, query, userReq.Email, string(hashedPassword), userReq.FirstName, userReq.LastName, defaultPrefs).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Role,
		&newUser.Preferences,
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
func (s *PostgresUserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, preferences, created_at, updated_at
		FROM users WHERE email = $1
	`
	var user model.User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Preferences,
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

// GetUserByID retrieves a user by their ID.
// Note: This currently selects the password_hash. In a real-world scenario
// you might have a separate function or a different model for public user profiles.
func (s *PostgresUserStore) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, preferences, created_at, updated_at
		FROM users WHERE id = $1
	`
	var user model.User
	err := s.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Preferences,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
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
func (s *PostgresUserStore) MarkLessonAsComplete(ctx context.Context, userID int64, lessonID int64) error {
	query := `
		INSERT INTO user_lesson_progress (user_id, lesson_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, lesson_id) DO NOTHING
	`
	_, err := s.db.Exec(ctx, query, userID, lessonID)
	return err
}

// GetCompletedLessonsForUser retrieves a list of completed lesson IDs for a user.
func (s *PostgresUserStore) GetCompletedLessonsForUser(ctx context.Context, userID int64) ([]int64, error) {
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

// UpdateUserPreferences updates the preferences for a given user.
func (s *PostgresUserStore) UpdateUserPreferences(ctx context.Context, userID int64, prefs map[string]interface{}) error {
	query := `
		UPDATE users
		SET preferences = preferences || $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := s.db.Exec(ctx, query, prefs, userID)
	return err
}

// CreatePasswordResetToken creates a new password reset token.
func (s *PostgresUserStore) CreatePasswordResetToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := s.db.Exec(ctx, query, userID, token, expiresAt)
	return err
}

// GetUserByPasswordResetToken retrieves a user by a password reset token.
func (s *PostgresUserStore) GetUserByPasswordResetToken(ctx context.Context, token string) (*model.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.role, u.preferences, u.created_at, u.updated_at
		FROM users u
		INNER JOIN password_reset_tokens prt ON u.id = prt.user_id
		WHERE prt.token = $1 AND prt.expires_at > NOW()
	`
	var user model.User
	err := s.db.QueryRow(ctx, query, token).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Preferences,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return &user, err
}

// DeletePasswordResetToken deletes a password reset token.
func (s *PostgresUserStore) DeletePasswordResetToken(ctx context.Context, token string) error {
	query := `DELETE FROM password_reset_tokens WHERE token = $1`
	_, err := s.db.Exec(ctx, query, token)
	return err
}

// UpdatePassword updates a user's password.
func (s *PostgresUserStore) UpdatePassword(ctx context.Context, userID int64, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err = s.db.Exec(ctx, query, string(hashedPassword), userID)
	return err
}

// --- Quiz Attempt Storage Functions ---

// CreateQuizAttempt creates a new quiz attempt in the database.
func (s *PostgresUserStore) CreateQuizAttempt(ctx context.Context, attempt *model.CreateQuizAttemptRequest, userID int64) (*model.QuizAttempt, error) {
	query := `
		INSERT INTO quiz_attempts (user_id, quiz_id, score, answers)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, quiz_id, score, answers, created_at
	`
	var newAttempt model.QuizAttempt
	err := s.db.QueryRow(ctx, query, userID, attempt.QuizID, attempt.Score, attempt.Answers).Scan(
		&newAttempt.ID,
		&newAttempt.UserID,
		&newAttempt.QuizID,
		&newAttempt.Score,
		&newAttempt.Answers,
		&newAttempt.CreatedAt,
	)
	return &newAttempt, err
}

// GetQuizAttemptsForUser retrieves all quiz attempts for a given user.
func (s *PostgresUserStore) GetQuizAttemptsForUser(ctx context.Context, userID int64) ([]model.QuizAttempt, error) {
	query := `
		SELECT id, user_id, quiz_id, score, answers, created_at
		FROM quiz_attempts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []model.QuizAttempt
	for rows.Next() {
		var attempt model.QuizAttempt
		if err := rows.Scan(&attempt.ID, &attempt.UserID, &attempt.QuizID, &attempt.Score, &attempt.Answers, &attempt.CreatedAt); err != nil {
			return nil, err
		}
		attempts = append(attempts, attempt)
	}
	return attempts, nil
}
