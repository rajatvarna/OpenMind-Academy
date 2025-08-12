package api

import (
	"log"
	"net/http"

	"github.com/free-education/user-service/auth"
	"github.com/free-education/user-service/model"
	"github.com/free-education/user-service/storage"
	"github.com/gin-gonic/gin"
)

// API holds the dependencies for the API handlers, like the user store.
type API struct {
	UserStore *storage.UserStore
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

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{Token: token})
}
