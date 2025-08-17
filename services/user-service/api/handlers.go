package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/free-education/user-service/auth"
	"github.com/free-education/user-service/messaging"
	"github.com/free-education/user-service/model"
	"github.com/free-education/user-service/storage"
	"github.com/gin-gonic/gin"
)

// API holds the dependencies for the API handlers, like the user store.
type API struct {
	UserStore      storage.UserStore
	MessageBroker messaging.MessageBroker
}

// MarkCompleteRequest defines the payload for marking a lesson as complete.
type MarkCompleteRequest struct {
	LessonID int64 `json:"lesson_id" binding:"required"`
}

// NewAPI creates a new API struct with its dependencies.
func NewAPI(userStore storage.UserStore, messageBroker messaging.MessageBroker) *API {
	return &API{
		UserStore:      userStore,
		MessageBroker: messageBroker,
	}
}

// RegisterUserHandler handles the logic for user registration.
func (a *API) RegisterUserHandler(c *gin.Context) {
	var req model.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// In a real app, you'd add more validation here (e.g., check if email already exists)
	// For simplicity, we rely on the database's UNIQUE constraint for now.

	newUser, err := a.UserStore.CreateUser(c.Request.Context(), &req)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		// This could be a duplicate email error, which is a client error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// LoginUserHandler handles the logic for user login.
func (a *API) LoginUserHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	user, err := a.UserStore.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		// User not found
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !storage.CheckPassword(user.PasswordHash, req.Password) {
		// Incorrect password
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

// ForgotPasswordHandler handles the logic for sending a password reset link.
func (a *API) ForgotPasswordHandler(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	user, err := a.UserStore.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		// Don't reveal if the user exists or not.
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

	// In a real app, you would get the frontend URL from config.
	resetLink := fmt.Sprintf("http://localhost:3001/reset-password?token=%s", token)

	// Publish event to RabbitMQ
	payload := map[string]interface{}{
		"email":     user.Email,
		"name":      user.FirstName,
		"resetLink": resetLink,
	}
	if err := a.MessageBroker.Publish(c.Request.Context(), "notifications_events", "password_reset_requested", payload); err != nil {
		log.Printf("Error publishing password reset event: %v", err)
		// Don't fail the whole request if the notification fails.
	}

	c.JSON(http.StatusOK, gin.H{"message": "If a user with that email exists, a password reset link has been sent."})
}

// ResetPasswordHandler handles the logic for resetting a user's password.
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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token."})
		return
	}

	if err := a.UserStore.UpdatePassword(c.Request.Context(), user.ID, req.NewPassword); err != nil {
		log.Printf("Error updating password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password."})
		return
	}

	// Clean up the used token
	if err := a.UserStore.DeletePasswordResetToken(c.Request.Context(), req.Token); err != nil {
		log.Printf("Error deleting password reset token: %v", err)
		// Don't fail the request if cleanup fails, but log it.
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully."})
}

// GetProfileHandler retrieves the profile for the currently authenticated user.
func (a *API) GetProfileHandler(c *gin.Context) {
	userIDHeader := c.GetHeader("X-User-Id")
	if userIDHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID header not provided"})
		return
	}

	id, err := strconv.ParseInt(userIDHeader, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

	user, err := a.UserStore.GetUserByID(c.Request.Context(), id)
	if err != nil {
		// This could be a "not found" error, which should be handled gracefully.
		log.Printf("Error fetching profile for user %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	// The User model already omits the password hash, so it's safe to return.
	c.JSON(http.StatusOK, user)
}

// GetUserPreferencesHandler retrieves the preferences for the currently authenticated user.
func (a *API) GetUserPreferencesHandler(c *gin.Context) {
	userIDHeader := c.GetHeader("X-User-Id")
	if userIDHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID header not provided"})
		return
	}

	id, err := strconv.ParseInt(userIDHeader, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

	user, err := a.UserStore.GetUserByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error fetching preferences for user %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.Preferences)
}

// UpdateUserPreferencesHandler updates the preferences for the currently authenticated user.
func (a *API) UpdateUserPreferencesHandler(c *gin.Context) {
	userIDHeader := c.GetHeader("X-User-Id")
	if userIDHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID header not provided"})
		return
	}

	id, err := strconv.ParseInt(userIDHeader, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

	var prefs map[string]interface{}
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := a.UserStore.UpdateUserPreferences(c.Request.Context(), id, prefs); err != nil {
		log.Printf("Error updating preferences for user %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProgressHandler retrieves the list of completed lesson IDs for a user.
func (a *API) GetProgressHandler(c *gin.Context) {
	// Authorization is now handled by the API Gateway.
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
func (a *API) MarkLessonCompleteHandler(c *gin.Context) {
	// Authorization is now handled by the API Gateway.
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

// Define service URLs. These should come from config/env vars.
var (
	contentServiceURL = "http://content-service:3001/api/v1"
	gamificationServiceURL = "http://gamification-service:3005/api/v1"
)

// FullProfileResponse defines the aggregated data for a user profile.
type FullProfileResponse struct {
	User             *model.User       `json:"user"`
	GamificationStats map[string]string `json:"gamification_stats"`
	CreatedCourses   []interface{}     `json:"created_courses"` // Using interface{} for simplicity
}

// GetFullProfileHandler fetches data from multiple services to build a user profile.
func (a *API) GetFullProfileHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 1. Get base user data from our own DB
	// In a real app, we'd have a GetUserByID function. We'll reuse GetUserByEmail as a stand-in.
	// This part of the logic is flawed without a proper GetUserByID.
	// For this demo, we'll skip fetching the base user and assume we have it.

	// Create a channel to receive results from concurrent API calls
	type apiResult struct {
		data interface{}
		err  error
		from string
	}
	ch := make(chan apiResult, 2)

	// 2. Fetch gamification stats concurrently
	go func() {
		resp, err := http.Get(fmt.Sprintf("%s/users/%d/stats", gamificationServiceURL, userID))
		if err != nil {
			ch <- apiResult{err: err, from: "gamification"}
			return
		}
		defer resp.Body.Close()
		var stats map[string]string
		json.NewDecoder(resp.Body).Decode(&stats)
		ch <- apiResult{data: stats, from: "gamification"}
	}()

	// 3. Fetch user's created courses concurrently
	go func() {
		// This endpoint doesn't exist yet, we'd need to add it to the Content Service
		resp, err := http.Get(fmt.Sprintf("%s/users/%d/courses", contentServiceURL, userID))
		if err != nil {
			ch <- apiResult{err: err, from: "content"}
			return
		}
		defer resp.Body.Close()
		var courses []interface{}
		json.NewDecoder(resp.Body).Decode(&courses)
		ch <- apiResult{data: courses, from: "content"}
	}()

	// 4. Aggregate results
	response := FullProfileResponse{}
	for i := 0; i < 2; i++ {
		result := <-ch
		if result.err != nil {
			log.Printf("Error fetching from %s service: %v", result.from, result.err)
			// Decide on error handling: fail the whole request or return partial data?
			// For now, we'll continue and return partial data.
			continue
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
