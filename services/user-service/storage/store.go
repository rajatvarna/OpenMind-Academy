package storage

import (
	"context"
	"time"

	"github.com/free-education/user-service/model"
)

// UserStore defines the interface for user storage operations.
type UserStore interface {
	CreateUser(ctx context.Context, userReq *model.RegistrationRequest) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	UpdateUserPreferences(ctx context.Context, userID int64, prefs map[string]interface{}) error
	UpdateProfilePictureURL(ctx context.Context, userID int64, url string) error
	Store2FASecrets(ctx context.Context, userID int64, secret string, recoveryCodes []string) error
	Activate2FA(ctx context.Context, userID int64) error
	Get2FAData(ctx context.Context, userID int64) (secret string, enabled bool, err error)
	CreatePasswordResetToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	GetUserByPasswordResetToken(ctx context.Context, token string) (*model.User, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
	UpdatePassword(ctx context.Context, userID int64, newPassword string) error
	GetCompletedLessonsForUser(ctx context.Context, userID int64) ([]int64, error)
	MarkLessonAsComplete(ctx context.Context, userID int64, lessonID int64) error
	CreateQuizAttempt(ctx context.Context, attempt *model.CreateQuizAttemptRequest, userID int64) (*model.QuizAttempt, error)
	GetQuizAttemptsForUser(ctx context.Context, userID int64) ([]model.QuizAttempt, error)
}
