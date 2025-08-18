package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/free-education/user-service/auth"
	"github.com/free-education/user-service/messaging"
	"github.com/free-education/user-service/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// MockUserStore is a mock implementation of the UserStore with an in-memory map.
type MockUserStore struct {
	users               map[int64]*model.User
	emailToID           map[string]int64
	passwordResetTokens map[string]int64 // token -> userID
	nextID              int64
}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		users:               make(map[int64]*model.User),
		emailToID:           make(map[string]int64),
		passwordResetTokens: make(map[string]int64),
		nextID:              1,
	}
}

func (m *MockUserStore) CreateUser(ctx context.Context, userReq *model.RegistrationRequest) (*model.User, error) {
	// In a real scenario, you'd use a proper salt and hashing cost.
	// We use a low cost here to make tests run faster.
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.MinCost)
	newUser := &model.User{
		ID:           m.nextID,
		Email:        userReq.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    userReq.FirstName,
		LastName:     userReq.LastName,
	}
	m.users[newUser.ID] = newUser
	m.emailToID[newUser.Email] = newUser.ID
	m.nextID++
	return newUser, nil
}

func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	id, ok := m.emailToID[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	user, _ := m.users[id]
	return user, nil
}

func (m *MockUserStore) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	user, ok := m.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
func (m *MockUserStore) UpdateUserPreferences(ctx context.Context, userID int64, prefs map[string]interface{}) error {
	return nil
}
func (m *MockUserStore) UpdateProfilePictureURL(ctx context.Context, userID int64, url string) error {
	return nil
}
func (m *MockUserStore) Store2FASecrets(ctx context.Context, userID int64, secret string, recoveryCodes []string) error {
	return nil
}
func (m *MockUserStore) Activate2FA(ctx context.Context, userID int64) error {
	return nil
}
func (m *MockUserStore) DeactivateUser(ctx context.Context, userID int64) error {
	return nil
}
func (m *MockUserStore) Get2FAData(ctx context.Context, userID int64) (string, bool, error) {
	return "", false, nil
}
func (m *MockUserStore) CreatePasswordResetToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	m.passwordResetTokens[token] = userID
	return nil
}
func (m *MockUserStore) GetUserByPasswordResetToken(ctx context.Context, token string) (*model.User, error) {
	userID, ok := m.passwordResetTokens[token]
	if !ok {
		return nil, errors.New("token not found")
	}
	return m.GetUserByID(ctx, userID)
}
func (m *MockUserStore) DeletePasswordResetToken(ctx context.Context, token string) error {
	delete(m.passwordResetTokens, token)
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
func (m *MockUserStore) CreateUserActivity(ctx context.Context, activity *model.UserActivity) error {
	return nil
}
func (m *MockUserStore) GetUserActivities(ctx context.Context, userID int64) ([]*model.UserActivity, error) {
	return []*model.UserActivity{}, nil
}

// MockMessageBroker is a mock implementation of the MessageBroker.
type MockMessageBroker struct{}

func (m *MockMessageBroker) Publish(ctx context.Context, queueName string, eventType string, payload interface{}) error {
	return nil
}

func (m *MockMessageBroker) Consume(ctx context.Context, queueName string, handler messaging.MessageHandler) error {
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
		userStore := NewMockUserStore()
		userStore.CreateUser(context.Background(), &model.RegistrationRequest{Email: "test@example.com", Password: "password"})
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
		userStore := NewMockUserStore()
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

func TestGetUserActivityHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userStore := NewMockUserStore()
	mockMessageBroker := &MockMessageBroker{}
	apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "userId", Value: "1"}}

	c.Request, _ = http.NewRequest(http.MethodGet, "/users/1/activity", nil)

	apiHandler.GetUserActivityHandler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
	}
}

func TestDeactivateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userStore := NewMockUserStore()
	mockMessageBroker := &MockMessageBroker{}
	apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", int64(1)) // Set user ID in context

	c.Request, _ = http.NewRequest(http.MethodDelete, "/profile", nil)

	apiHandler.DeactivateUserHandler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
	}
}

func TestLoginUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestKey(t)

	t.Run("Successful login", func(t *testing.T) {
		userStore := NewMockUserStore()
		userStore.CreateUser(context.Background(), &model.RegistrationRequest{Email: "test@example.com", Password: "password"})
		mockMessageBroker := &MockMessageBroker{}
		apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := map[string]string{"email": "test@example.com", "password": "password"}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.LoginUserHandler(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Deactivated user login", func(t *testing.T) {
		userStore := NewMockUserStore()
		user, _ := userStore.CreateUser(context.Background(), &model.RegistrationRequest{Email: "deactivated@example.com", Password: "password"})
		now := time.Now()
		user.DeactivatedAt = &now

		mockMessageBroker := &MockMessageBroker{}
		apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := map[string]string{"email": "deactivated@example.com", "password": "password"}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.LoginUserHandler(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d; got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestEnable2FAHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userStore := NewMockUserStore()
	userStore.CreateUser(context.Background(), &model.RegistrationRequest{Email: "test@example.com", Password: "password"})
	mockMessageBroker := &MockMessageBroker{}
	apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", int64(1)) // Set user ID in context

	c.Request, _ = http.NewRequest(http.MethodPost, "/2fa/enable", nil)

	apiHandler.Enable2FAHandler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if _, ok := responseBody["otpauth_url"]; !ok {
		t.Error("expected 'otpauth_url' in response")
	}
	if _, ok := responseBody["recovery_codes"]; !ok {
		t.Error("expected 'recovery_codes' in response")
	}
}

func TestVerify2FAHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Invalid token", func(t *testing.T) {
		userStore := NewMockUserStore()
		mockMessageBroker := &MockMessageBroker{}
		apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", int64(1))

		body := map[string]string{"token": "invalid-token"}
		jsonBody, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest(http.MethodPost, "/2fa/verify", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		apiHandler.Verify2FAHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d; got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestUploadProfilePictureHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userStore := NewMockUserStore()
	mockMessageBroker := &MockMessageBroker{}
	apiHandler := NewAPI(userStore, mockMessageBroker, "", "", "")

	// Create a buffer to store our request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a new form-data header
	part, err := writer.CreateFormFile("picture", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}

	// Write a fake file to the form-data header
	_, err = io.Copy(part, strings.NewReader("fake image data"))
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", int64(1)) // Set user ID in context

	c.Request, _ = http.NewRequest(http.MethodPost, "/profile/picture", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	apiHandler.UploadProfilePictureHandler(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
	}

	var responseBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if responseBody["message"] != "Profile picture updated successfully." {
		t.Errorf("expected success message; got '%s'", responseBody["message"])
	}

	if !strings.HasPrefix(responseBody["url"], "https://storage.example.com/profiles/") {
		t.Errorf("expected URL with prefix; got '%s'", responseBody["url"])
	}
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

	userStore := NewMockUserStore()
	userStore.CreateUser(context.Background(), &model.RegistrationRequest{Email: "test@example.com", Password: "password"})
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
		userStore := NewMockUserStore()
		user, _ := userStore.CreateUser(context.Background(), &model.RegistrationRequest{Email: "test@example.com", Password: "password"})
		userStore.CreatePasswordResetToken(context.Background(), user.ID, "valid-token", time.Now().Add(time.Hour))
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
		userStore := NewMockUserStore()
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
