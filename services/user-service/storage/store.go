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
	CreatePasswordResetToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	GetUserByPasswordResetToken(ctx context.Context, token string) (*model.User, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
	UpdatePassword(ctx context.Context, userID int64, newPassword string) error
	GetCompletedLessonsForUser(ctx context.Context, userID int64) ([]int64, error)
	MarkLessonAsComplete(ctx context.Context, userID int64, lessonID int64) error
}
