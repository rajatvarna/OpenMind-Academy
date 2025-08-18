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
	if userID == 1 {
		return &model.User{ID: 1, Email: "test@example.com", FirstName: "Test"}, nil
	}
	return nil, nil // Simulate user not found
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

func (m *MockUserStore) CreateQuizAttempt(ctx context.Context, attempt *model.CreateQuizAttemptRequest, userID int64) (*model.QuizAttempt, error) {
	return nil, nil
}

func (m *MockUserStore) GetQuizAttemptsForUser(ctx context.Context, userID int64) ([]model.QuizAttempt, error) {
	return nil, nil
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
		apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := map[string]string{"email": "test@example.com"}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/password/forgot", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.ForgotPasswordHandler(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("User not found", func(t *testing.T) {
		var userStore storage.UserStore = &MockUserStore{}
		mockMessageBroker := &MockMessageBroker{}
		apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := map[string]string{"email": "not-found@example.com"}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/password/forgot", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.ForgotPasswordHandler(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
		}
	})
}

func TestGetFullProfileHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock Gamification Service
	gamificationServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"score": "100", "rank": "gold"})
	}))
	defer gamificationServer.Close()

	// Mock Content Service (configured to fail)
	contentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer contentServer.Close()

	var userStore storage.UserStore = &MockUserStore{}
	mockMessageBroker := &MockMessageBroker{}
	apiHandler := NewAPI(userStore, mockMessageBroker, "", contentServer.URL, gamificationServer.URL)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "userId", Value: "1"}}

	c.Request, _ = http.NewRequest(http.MethodGet, "/users/1/full-profile", nil)

	apiHandler.GetFullProfileHandler(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d; got %d", http.StatusInternalServerError, w.Code)
	}

	var responseBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedError := "Failed to retrieve full user profile due to an error with a downstream service."
	if responseBody["error"] != expectedError {
		t.Errorf("expected error message '%s'; got '%s'", expectedError, responseBody["error"])
	}
}

func TestResetPasswordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestKey(t)

	t.Run("Successful password reset", func(t *testing.T) {
		var userStore storage.UserStore = &MockUserStore{}
		mockMessageBroker := &MockMessageBroker{}
		apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := map[string]string{"token": "valid-token", "new_password": "new-password"}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/password/reset", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.ResetPasswordHandler(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		var userStore storage.UserStore = &MockUserStore{}
		mockMessageBroker := &MockMessageBroker{}
		apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := map[string]string{"token": "invalid-token", "new_password": "new-password"}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/password/reset", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.ResetPasswordHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d; got %d", http.StatusBadRequest, w.Code)
		}
	})
}
