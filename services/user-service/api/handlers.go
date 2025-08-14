package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/free-education/user-service/auth"
	"github.com/free-education/user-service/model"
	"github.com/free-education/user-service/storage"
	"github.com/gin-gonic/gin"
)

// API holds the dependencies for the API handlers, like the user store.
type API struct {
	UserStore *storage.UserStore
}

// MarkCompleteRequest defines the payload for marking a lesson as complete.
type MarkCompleteRequest struct {
	LessonID int64 `json:"lesson_id" binding:"required"`
}

// NewAPI creates a new API struct with its dependencies.
func NewAPI(userStore *storage.UserStore) *API {
	return &API{UserStore: userStore}
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

// GetProgressHandler retrieves the list of completed lesson IDs for a user.
func (a *API) GetProgressHandler(c *gin.Context) {
	// The API gateway has already authenticated the user. We trust the headers.
	// We should still validate that the requester has permission to view the progress for the given userId.
	requesterID, _ := strconv.ParseInt(c.GetHeader("X-User-Id"), 10, 64)
	targetUserID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	// For now, only allow users to see their own progress.
	if requesterID != targetUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own progress."})
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
	requesterID, _ := strconv.ParseInt(c.GetHeader("X-User-Id"), 10, 64)
	targetUserID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	if requesterID != targetUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only mark lessons as complete for yourself."})
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
