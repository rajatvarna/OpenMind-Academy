package api

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/oauth2"

	"github.com/free-education/user-service/auth"
	"github.com/free-education/user-service/messaging"
	"github.com/free-education/user-service/model"
	"github.com/free-education/user-service/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
)

// API holds the dependencies for the API handlers, like the user store.
type API struct {
	UserStore              storage.UserStore
	MessageBroker          messaging.MessageBroker
	FrontendBaseURL        string
	ContentServiceURL      string
	GamificationServiceURL string
	GoogleOAuthConfig      *oauth2.Config
}

// MarkCompleteRequest defines the payload for marking a lesson as complete.
type MarkCompleteRequest struct {
	LessonID int64 `json:"lesson_id" binding:"required"`
}

// NewAPI creates a new API struct with its dependencies.
func NewAPI(userStore storage.UserStore, messageBroker messaging.MessageBroker, frontendBaseURL, contentServiceURL, gamificationServiceURL string, googleOAuthConfig *oauth2.Config) *API {
	return &API{
		UserStore:              userStore,
		MessageBroker:          messageBroker,
		FrontendBaseURL:        frontendBaseURL,
		ContentServiceURL:      contentServiceURL,
		GamificationServiceURL: gamificationServiceURL,
		GoogleOAuthConfig:      googleOAuthConfig,
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

	// Check if the account is deactivated.
	if user.DeactivatedAt != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This account has been deactivated."})
		return
	}

	if !storage.CheckPassword(user.PasswordHash, req.Password) {
		// Incorrect password.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if 2FA is enabled for the user.
	_, twoFactorEnabled, err := a.UserStore.Get2FAData(c.Request.Context(), user.ID)
	if err != nil {
		log.Printf("Error getting 2FA data for user %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process login."})
		return
	}

	if twoFactorEnabled {
		// If 2FA is enabled, issue a temporary token and prompt for 2FA verification.
		tempToken, err := auth.Generate2FATempToken(user.ID)
		if err != nil {
			log.Printf("Error generating 2FA temp token for user %d: %v", user.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate 2FA token."})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":    "2FA token required",
			"temp_token": tempToken,
		})
		return
	}

	// If 2FA is not enabled, issue a full-access token.
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{Token: token})
}

// Login2FARequest represents the payload for the 2FA login request.
type Login2FARequest struct {
	TempToken string `json:"temp_token" binding:"required"`
	Token     string `json:"token" binding:"required"`
}

// Login2FAHandler handles the second step of the 2FA login process.
func (a *API) Login2FAHandler(c *gin.Context) {
	var req Login2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Validate the temporary token
	claims, err := auth.ValidateToken(req.TempToken)
	if err != nil || claims.Type != "2fa_temp" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired temporary token."})
		return
	}

	// Get the user's 2FA secret
	secret, enabled, err := a.UserStore.Get2FAData(c.Request.Context(), claims.UserID)
	if err != nil || !enabled || secret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "2FA is not enabled or user not found."})
		return
	}

	// Validate the TOTP token
	valid := totp.Validate(req.Token, secret)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA token."})
		return
	}

	// Get user role to generate full token
	user, err := a.UserStore.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user details."})
		return
	}

	// If everything is valid, issue a full-access token
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		log.Printf("Error generating JWT for user %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token."})
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

// --- 2FA Handlers ---

// generateRecoveryCodes creates a set of random strings to be used as single-use recovery codes.
func generateRecoveryCodes(count int, length int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		bytes := make([]byte, length)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		codes[i] = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
	}
	return codes, nil
}

// --- OAuth Handlers ---

func (a *API) generateStateOauthCookie(c *gin.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, 3600, "/", "", false, true)
	return state
}

func (a *API) GoogleLoginHandler(c *gin.Context) {
	state := a.generateStateOauthCookie(c)
	url := a.GoogleOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *API) GoogleCallbackHandler(c *gin.Context) {
	// Read oauthState from Cookie
	oauthState, _ := c.Cookie("oauthstate")
	if c.Request.FormValue("state") != oauthState {
		log.Println("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=invalid_state", a.FrontendBaseURL))
		return
	}

	data, err := a.getUserDataFromGoogle(c.Request.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=oauth_failed", a.FrontendBaseURL))
		return
	}

	// Get user info
	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(data, &userInfo); err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=oauth_failed", a.FrontendBaseURL))
		return
	}

	// Check if user already exists
	user, err := a.UserStore.GetUserByOAuthID(c.Request.Context(), "google", userInfo.ID)
	if err != nil {
		// User does not exist, check by email
		user, err = a.UserStore.GetUserByEmail(c.Request.Context(), userInfo.Email)
		if err != nil {
			// User does not exist, create new user
			newUser := &model.User{
				Email:           userInfo.Email,
				FirstName:       userInfo.Name,
				OAuthProvider:   "google",
				OAuthProviderID: userInfo.ID,
			}
			user, err = a.UserStore.CreateOAuthUser(c.Request.Context(), newUser)
			if err != nil {
				log.Printf("Error creating OAuth user: %v", err)
				c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=creation_failed", a.FrontendBaseURL))
				return
			}
		}
	}

	// Generate JWT
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		log.Printf("Error generating JWT for user %d: %v", user.ID, err)
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=token_failed", a.FrontendBaseURL))
		return
	}

	// For simplicity, we'll redirect with the token in the query string.
	// In a real app, you might use a more secure method.
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/callback?token=%s", a.FrontendBaseURL, token))
}

func (a *API) getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := a.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

// Enable2FAHandler begins the process of enabling two-factor authentication.
// It generates a new TOTP secret and recovery codes for the user.
func (a *API) Enable2FAHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)
	user, err := a.UserStore.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "FreeEdu",
		AccountName: user.Email,
	})
	if err != nil {
		log.Printf("Error generating TOTP key for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate 2FA secret."})
		return
	}

	recoveryCodes, err := generateRecoveryCodes(10, 10) // 10 codes, 10 chars each
	if err != nil {
		log.Printf("Error generating recovery codes for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recovery codes."})
		return
	}

	if err := a.UserStore.Store2FASecrets(c.Request.Context(), userID, key.Secret(), recoveryCodes); err != nil {
		log.Printf("Error storing 2FA secrets for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save 2FA configuration."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "2FA setup initiated. Please scan the QR code and verify.",
		"otpauth_url":    key.URL(),
		"recovery_codes": recoveryCodes,
	})
}

// Verify2FARequest represents the payload for the 2FA verification request.
type Verify2FARequest struct {
	Token string `json:"token" binding:"required"`
}

// Verify2FAHandler completes the 2FA setup process by verifying a TOTP token.
func (a *API) Verify2FAHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	secret, _, err := a.UserStore.Get2FAData(c.Request.Context(), userID)
	if err != nil || secret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA not initiated or user not found."})
		return
	}

	valid := totp.Validate(req.Token, secret)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 2FA token."})
		return
	}

	if err := a.UserStore.Activate2FA(c.Request.Context(), userID); err != nil {
		log.Printf("Error activating 2FA for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate 2FA."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully."})
}

// Disable2FAHandler handles disabling 2FA for a user.
func (a *API) Disable2FAHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	if err := a.UserStore.Disable2FA(c.Request.Context(), userID); err != nil {
		log.Printf("Error disabling 2FA for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable 2FA."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully."})
}

// --- Account Deactivation ---

// DeactivateUserHandler handles a user's request to deactivate their own account.
func (a *API) DeactivateUserHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	if err := a.UserStore.DeactivateUser(c.Request.Context(), userID); err != nil {
		log.Printf("Error deactivating user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate account."})
		return
	}

	// Here you might also want to publish an event to a message queue
	// to invalidate all active sessions/tokens for this user.

	c.JSON(http.StatusOK, gin.H{"message": "Account deactivated successfully."})
}

// DeleteUserHandler handles a user's request to permanently delete their own account.
func (a *API) DeleteUserHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	// It's good practice to log this significant event.
	log.Printf("Attempting to permanently delete user %d", userID)

	if err := a.UserStore.DeleteUser(c.Request.Context(), userID); err != nil {
		log.Printf("Error permanently deleting user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account."})
		return
	}

	// Publish an event to notify other services that this user has been deleted.
	// This is crucial for data consistency across the microservices ecosystem.
	payload := map[string]interface{}{"user_id": userID}
	if err := a.MessageBroker.Publish(c.Request.Context(), "user_events", "user_deleted", payload); err != nil {
		// Log the error, but don't fail the request. The primary operation (DB deletion) was successful.
		// A monitoring system should alert on these kinds of failures.
		log.Printf("CRITICAL: Failed to publish user_deleted event for user %d: %v", userID, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account permanently deleted."})
}

// --- User Activity Handlers ---

// GetUserActivityHandler retrieves a log of recent activities for a user.
func (a *API) GetUserActivityHandler(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	// In a real app, you'd add authorization logic here to ensure the requesting user
	// is allowed to see this activity log (e.g., they are the user themselves, or an admin).

	activities, err := a.UserStore.GetUserActivities(c.Request.Context(), targetUserID)
	if err != nil {
		log.Printf("Error getting activities for user %d: %v", targetUserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user activities."})
		return
	}

	c.JSON(http.StatusOK, activities)
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
