package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/free-education/user-service/auth"
	"github.com/free-education/user-service/messaging"
	"github.com/free-education/user-service/model"
	"github.com/free-education/user-service/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
)

// API holds the dependencies for the API handlers, like the user store.
type API struct {
	UserStore             storage.UserStore
	MessageBroker         messaging.MessageBroker
	FrontendBaseURL       string
	ContentServiceURL     string
	GamificationServiceURL string
}

// MarkCompleteRequest defines the payload for marking a lesson as complete.
type MarkCompleteRequest struct {
	LessonID int64 `json:"lesson_id" binding:"required"`
}

// NewAPI creates a new API struct with its dependencies.
func NewAPI(userStore storage.UserStore, messageBroker messaging.MessageBroker, frontendBaseURL, contentServiceURL, gamificationServiceURL string) *API {
	return &API{
		UserStore:             userStore,
		MessageBroker:         messageBroker,
		FrontendBaseURL:       frontendBaseURL,
		ContentServiceURL:     contentServiceURL,
		GamificationServiceURL: gamificationServiceURL,
	}
}

// RegisterUserHandler handles new user registration.
// It expects a JSON payload with the user's email, password, and name.
// On success, it returns the newly created user object with a 201 status code.
// If the email already exists, it returns a 409 Conflict error.
func (a *API) RegisterUserHandler(c *gin.Context) {
	var req model.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	newUser, err := a.UserStore.CreateUser(c.Request.Context(), &req)
	if err != nil {
		var pgErr *pgconn.PgError
		// Check if the error is a PostgreSQL error and if it's a unique violation (code 23505).
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "An account with this email already exists."})
			return
		}
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// LoginUserHandler handles user authentication.
// It expects an email and password, and upon successful validation,
// returns a JWT token for use in subsequent authenticated requests.
// Returns a 401 Unauthorized error for invalid credentials.
func (a *API) LoginUserHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	user, err := a.UserStore.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		// User not found. Return a generic error to avoid revealing user existence.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !storage.CheckPassword(user.PasswordHash, req.Password) {
		// Incorrect password.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{Token: token})
}

// ForgotPasswordHandler initiates the password reset process.
// It generates a secure, single-use token, stores it, and publishes an event
// for the notifications service to send an email with the reset link.
// To prevent user enumeration, it always returns a 200 OK response.
func (a *API) ForgotPasswordHandler(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	user, err := a.UserStore.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil || user == nil {
		// Don't reveal if the user exists or not for security reasons.
		c.JSON(http.StatusOK, gin.H{"message": "If a user with that email exists, a password reset link has been sent."})
		return
	}

	token, err := auth.GenerateSecureToken(32)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request."})
		return
	}

	expiresAt := time.Now().Add(time.Hour * 1) // Token valid for 1 hour
	if err := a.UserStore.CreatePasswordResetToken(c.Request.Context(), user.ID, token, expiresAt); err != nil {
		log.Printf("Error creating password reset token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request."})
		return
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", a.FrontendBaseURL, token)

	// Publish an event to the message broker. The notifications service will consume this
	// and send the actual email. This decouples the services.
	payload := map[string]interface{}{
		"email":     user.Email,
		"name":      user.FirstName,
		"resetLink": resetLink,
	}
	if err := a.MessageBroker.Publish(c.Request.Context(), "notifications_events", "password_reset_requested", payload); err != nil {
		log.Printf("Error publishing password reset event: %v", err)
		// We still return a success response to the user even if the notification fails.
		// The operation should be idempotent and can be retried by the user.
	}

	c.JSON(http.StatusOK, gin.H{"message": "If a user with that email exists, a password reset link has been sent."})
}

// ResetPasswordHandler completes the password reset process.
// It requires a valid, non-expired token and a new password.
// Upon success, it updates the user's password and deletes the token to prevent reuse.
func (a *API) ResetPasswordHandler(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	user, err := a.UserStore.GetUserByPasswordResetToken(c.Request.Context(), req.Token)
	if err != nil || user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token."})
		return
	}

	if err := a.UserStore.UpdatePassword(c.Request.Context(), user.ID, req.NewPassword); err != nil {
		log.Printf("Error updating password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password."})
		return
	}

	// Clean up the used token to ensure it cannot be used again.
	if err := a.UserStore.DeletePasswordResetToken(c.Request.Context(), req.Token); err != nil {
		log.Printf("Error deleting password reset token: %v", err)
		// Don't fail the main request if cleanup fails, but log it as it's important.
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully."})
}

// GetProfileHandler retrieves the profile for the currently authenticated user.
// The user ID is injected by the AuthMiddleware.
func (a *API) GetProfileHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	user, err := a.UserStore.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error fetching profile for user %d: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	// The User model already omits the password hash, so it's safe to return.
	c.JSON(http.StatusOK, user)
}

// GetUserPreferencesHandler retrieves the preferences for the currently authenticated user.
// The user ID is injected by the AuthMiddleware.
func (a *API) GetUserPreferencesHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	user, err := a.UserStore.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error fetching preferences for user %d: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.Preferences)
}

// UpdateUserPreferencesHandler updates the preferences for the currently authenticated user.
// The user ID is injected by the AuthMiddleware. It expects a JSON object
// containing the preferences to be updated.
func (a *API) UpdateUserPreferencesHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	var prefs map[string]interface{}
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := a.UserStore.UpdateUserPreferences(c.Request.Context(), userID, prefs); err != nil {
		log.Printf("Error updating preferences for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	c.Status(http.StatusNoContent)
}

// uploadToCloudStorage is a placeholder for a real cloud storage upload function.
// In a real application, this would use the AWS, GCS, or Azure SDK to upload the file
// and would require proper error handling and configuration.
func uploadToCloudStorage(fileHeader *multipart.FileHeader) (string, error) {
	// For this example, we'll just simulate an upload and return a fake URL.
	// We'll use the filename to make the URL unique.
	// In a real app, you would generate a unique ID (e.g., a UUID) for the filename
	// to prevent collisions.
	log.Printf("Simulating upload for file: %s", fileHeader.Filename)
	fakeURL := fmt.Sprintf("https://storage.example.com/profiles/%d-%s", time.Now().UnixNano(), fileHeader.Filename)
	return fakeURL, nil
}

// UploadProfilePictureHandler handles the profile picture upload process.
func (a *API) UploadProfilePictureHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	file, err := c.FormFile("picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided or invalid."})
		return
	}

	// In a real application, you would add more validation here:
	// - Check file size
	// - Check file type (e.g., only allow jpeg, png)

	// Upload the file to cloud storage (using our mock function)
	url, err := uploadToCloudStorage(file)
	if err != nil {
		log.Printf("Error uploading file for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file."})
		return
	}

	// Update the user's profile picture URL in the database
	if err := a.UserStore.UpdateProfilePictureURL(c.Request.Context(), userID, url); err != nil {
		log.Printf("Error updating profile picture URL for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile picture URL."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile picture updated successfully.", "url": url})
}

// --- Quiz Attempt Handlers ---

// CreateQuizAttemptHandler handles saving a user's quiz attempt.
// The user ID is injected by the AuthMiddleware.
func (a *API) CreateQuizAttemptHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	var req model.CreateQuizAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	attempt, err := a.UserStore.CreateQuizAttempt(c.Request.Context(), &req, userID)
	if err != nil {
		log.Printf("Error creating quiz attempt for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save quiz attempt"})
		return
	}

	c.JSON(http.StatusCreated, attempt)
}

// GetQuizAttemptsForUserHandler retrieves all quiz attempts for a specific user.
// Authorization should be handled by the API Gateway to ensure only the user
// themselves or an authorized role (e.g., admin) can access this.
func (a *API) GetQuizAttemptsForUserHandler(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	attempts, err := a.UserStore.GetQuizAttemptsForUser(c.Request.Context(), targetUserID)
	if err != nil {
		log.Printf("Error getting quiz attempts for user %d: %v", targetUserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quiz attempts"})
		return
	}

	c.JSON(http.StatusOK, attempts)
}

// GetProgressHandler retrieves the list of completed lesson IDs for a user.
// Authorization should be handled by the API Gateway.
func (a *API) GetProgressHandler(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	completed, err := a.UserStore.GetCompletedLessonsForUser(c.Request.Context(), targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"completed_lessons": completed})
}

// MarkLessonCompleteHandler marks a lesson as complete for a user.
// Authorization should be handled by the API Gateway.
func (a *API) MarkLessonCompleteHandler(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	var req MarkCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	err = a.UserStore.MarkLessonAsComplete(c.Request.Context(), targetUserID, req.LessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark lesson as complete"})
		return
	}

	c.Status(http.StatusNoContent)
}

// --- Full Profile Aggregation ---

// FullProfileResponse defines the aggregated data for a user profile.
type FullProfileResponse struct {
	User             *model.User       `json:"user"`
	GamificationStats map[string]string `json:"gamification_stats"`
	CreatedCourses   []interface{}     `json:"created_courses"` // Using interface{} for simplicity
}

// GetFullProfileHandler demonstrates the aggregator pattern. It fetches data from
// multiple services to construct a complete user profile.
// It concurrently calls the gamification and content services.
// NOTE: This approach has trade-offs. While it simplifies the frontend, it creates
// coupling between services and can be a performance bottleneck. In a real-world
// scenario, other patterns like event-driven data replication might be preferable.
func (a *API) GetFullProfileHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 1. Get base user data from our own DB
	user, err := a.UserStore.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error getting user for full profile %d: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Create a channel to receive results from concurrent API calls
	type apiResult struct {
		data interface{}
		err  error
		from string
	}
	ch := make(chan apiResult, 2)

	// 2. Fetch gamification stats concurrently
	go func() {
		resp, err := http.Get(fmt.Sprintf("%s/users/%d/stats", a.GamificationServiceURL, userID))
		if err != nil {
			ch <- apiResult{err: err, from: "gamification"}
			return
		}
		defer resp.Body.Close()
		var stats map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			ch <- apiResult{err: err, from: "gamification"}
			return
		}
		ch <- apiResult{data: stats, from: "gamification"}
	}()

	// 3. Fetch user's created courses concurrently
	go func() {
		// This endpoint doesn't exist yet, we'd need to add it to the Content Service
		resp, err := http.Get(fmt.Sprintf("%s/users/%d/courses", a.ContentServiceURL, userID))
		if err != nil {
			ch <- apiResult{err: err, from: "content"}
			return
		}
		defer resp.Body.Close()
		var courses []interface{}
		if err := json.NewDecoder(resp.Body).Decode(&courses); err != nil {
			ch <- apiResult{err: err, from: "content"}
			return
		}
		ch <- apiResult{data: courses, from: "content"}
	}()

	// 4. Aggregate results
	response := FullProfileResponse{User: user}
	for i := 0; i < 2; i++ {
		result := <-ch
		if result.err != nil {
			log.Printf("Error fetching from %s service: %v", result.from, result.err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve full user profile due to an error with a downstream service."})
			return
		}
		switch result.from {
		case "gamification":
			response.GamificationStats = result.data.(map[string]string)
		case "content":
			response.CreatedCourses = result.data.([]interface{})
		}
	}

	c.JSON(http.StatusOK, response)
}
