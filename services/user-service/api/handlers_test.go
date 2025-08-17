package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/free-education/user-service/auth"
	"github.com/free-education/user-service/model"
	"github.com/free-education/user-service/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// MockUserStore is a mock implementation of the UserStore.
type MockUserStore struct{}

func (m *MockUserStore) CreateUser(ctx context.Context, userReq *model.RegistrationRequest) (*model.User, error) {
	return nil, nil
}
func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "test@example.com" {
		return &model.User{ID: 1, Email: "test@example.com", FirstName: "Test"}, nil
	}
	return nil, nil // Simulate user not found
}
func (m *MockUserStore) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	return nil, nil
}
func (m *MockUserStore) UpdateUserPreferences(ctx context.Context, userID int64, prefs map[string]interface{}) error {
	return nil
}
func (m *MockUserStore) CreatePasswordResetToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	return nil
}
func (m *MockUserStore) GetUserByPasswordResetToken(ctx context.Context, token string) (*model.User, error) {
	if token == "valid-token" {
		return &model.User{ID: 1, Email: "test@example.com", FirstName: "Test"}, nil
	}
	return nil, nil // Simulate token not found
}
func (m *MockUserStore) DeletePasswordResetToken(ctx context.Context, token string) error {
	return nil
}
func (m *MockUserStore) UpdatePassword(ctx context.Context, userID int64, newPassword string) error {
	return nil
}
func (m *MockUserStore) GetCompletedLessonsForUser(ctx context.Context, userID int64) ([]int64, error) {
	return nil, nil
}
func (m *MockUserStore) MarkLessonAsComplete(ctx context.Context, userID int64, lessonID int64) error {
	return nil
}

// MockMessageBroker is a mock implementation of the MessageBroker.
type MockMessageBroker struct{}

func (m *MockMessageBroker) Publish(ctx context.Context, queueName string, eventType string, payload interface{}) error {
	return nil
}

func (m *MockMessageBroker) Close() {}

func setupTestKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	if err := os.WriteFile("test_key.pem", pemdata, 0644); err != nil {
		t.Fatalf("Failed to write key to file: %v", err)
	}

	if err := auth.LoadPrivateKey("test_key.pem"); err != nil {
		t.Fatalf("Failed to load private key: %v", err)
	}
}

func TestForgotPasswordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestKey(t)

	t.Run("Successful password reset request", func(t *testing.T) {
		var userStore storage.UserStore = &MockUserStore{}
		mockMessageBroker := &MockMessageBroker{}
		apiHandler := NewAPI(userStore, mockMessageBroker)

		router := gin.Default()
		router.POST("/password/forgot", apiHandler.ForgotPasswordHandler)

		// Create a request
		body := map[string]string{"email": "test@example.com"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/password/forgot", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Check the response
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
		}
	})
}

func TestResetPasswordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestKey(t)

	t.Run("Successful password reset", func(t *testing.T) {
		var userStore storage.UserStore = &MockUserStore{}
		mockMessageBroker := &MockMessageBroker{}
		apiHandler := NewAPI(userStore, mockMessageBroker)

		router := gin.Default()
		router.POST("/password/reset", apiHandler.ResetPasswordHandler)

		// Create a request
		body := map[string]string{"token": "valid-token", "new_password": "new-password"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/password/reset", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Check the response
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
		}
	})
}
