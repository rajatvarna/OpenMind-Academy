package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/free-education/user-service/api"
	"github.com/free-education/user-service/messaging"
	"github.com/free-education/user-service/model"
	"github.com/free-education/user-service/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/free-education/user-service/auth"
)

func main() {
	// --- Key Loading ---
	privateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "../secrets/jwtRS256.key"
		log.Println("JWT_PRIVATE_KEY_PATH not set, using default value.")
	}
	if err := auth.LoadPrivateKey(privateKeyPath); err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// --- Database Connection ---
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://platform_user:password@db:5432/platform_db"
		log.Println("DATABASE_URL not set, using default Docker Compose value.")
	}

	dbpool, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()
	log.Println("Successfully connected to the database.")

	// --- Dependency Injection ---
	userStore := storage.NewUserStore(dbpool)

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@rabbitmq:5672/"
		log.Println("RABBITMQ_URL not set, using default Docker Compose value.")
	}
	messageBroker, err := messaging.NewRabbitMQClient(rabbitMQURL)
	if err != nil {
		log.Fatalf("Unable to connect to RabbitMQ: %v\n", err)
	}
	defer messageBroker.Close()

	// --- Start RabbitMQ Consumer for User Activities ---
	activityHandler := func(body []byte) {
		var event struct {
			EventType string             `json:"eventType"`
			Payload   model.UserActivity `json:"payload"`
		}
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("Error unmarshalling activity event: %v", err)
			return
		}

		// We only care about the payload, which should be a UserActivity
		if err := userStore.CreateUserActivity(context.Background(), &event.Payload); err != nil {
			log.Printf("Error creating user activity: %v", err)
		}
	}
	if err := messageBroker.Consume(context.Background(), "user_activity_events", activityHandler); err != nil {
		log.Fatalf("Failed to start user activity consumer: %v", err)
	}

	// --- OAuth2 Config ---
	googleOAuthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/users/login/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	frontendBaseURL := os.Getenv("FRONTEND_BASE_URL")
	if frontendBaseURL == "" {
		frontendBaseURL = "http://localhost:3001"
		log.Println("FRONTEND_BASE_URL not set, using default value.")
	}
	contentServiceURL := os.Getenv("CONTENT_SERVICE_URL")
	if contentServiceURL == "" {
		contentServiceURL = "http://content-service:3001/api/v1"
		log.Println("CONTENT_SERVICE_URL not set, using default value.")
	}
	gamificationServiceURL := os.Getenv("GAMIFICATION_SERVICE_URL")
	if gamificationServiceURL == "" {
		gamificationServiceURL = "http://gamification-service:3005/api/v1"
		log.Println("GAMIFICATION_SERVICE_URL not set, using default value.")
	}

	apiHandler := api.NewAPI(userStore, messageBroker, frontendBaseURL, contentServiceURL, gamificationServiceURL, googleOAuthConfig)

	// --- Router Setup ---
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Unauthenticated routes
		v1.POST("/register", apiHandler.RegisterUserHandler)
		v1.POST("/login", apiHandler.LoginUserHandler)
		v1.POST("/login/2fa", apiHandler.Login2FAHandler)
		v1.GET("/login/google", apiHandler.GoogleLoginHandler)
		v1.GET("/login/google/callback", apiHandler.GoogleCallbackHandler)
		v1.POST("/password/forgot", apiHandler.ForgotPasswordHandler)
		v1.POST("/password/reset", apiHandler.ResetPasswordHandler)

		// Authenticated routes are now protected by the API Gateway
		authenticated := v1.Group("/")
		authenticated.Use(api.AuthMiddleware())
		{
			authenticated.GET("/profile", apiHandler.GetProfileHandler)
			authenticated.DELETE("/profile", apiHandler.DeactivateUserHandler) // Kept for deactivation
			authenticated.DELETE("/account", apiHandler.DeleteUserHandler)     // New route for permanent deletion
			authenticated.POST("/profile/picture", apiHandler.UploadProfilePictureHandler)

			// 2FA routes
			authenticated.POST("/2fa/enable", apiHandler.Enable2FAHandler)
			authenticated.POST("/2fa/verify", apiHandler.Verify2FAHandler)
			authenticated.POST("/2fa/disable", apiHandler.Disable2FAHandler)

			authenticated.GET("/preferences", apiHandler.GetUserPreferencesHandler)
			authenticated.PUT("/preferences", apiHandler.UpdateUserPreferencesHandler)
			authenticated.GET("/users/:userId/progress", apiHandler.GetProgressHandler)
			authenticated.POST("/users/:userId/progress", apiHandler.MarkLessonCompleteHandler)
			authenticated.GET("/users/:userId/quiz-attempts", apiHandler.GetQuizAttemptsForUserHandler)
			authenticated.GET("/users/:userId/activity", apiHandler.GetUserActivityHandler)
			authenticated.GET("/users/:userId/full-profile", apiHandler.GetFullProfileHandler)

			// Authenticated routes - specific to the user
			authenticated.POST("/quiz-attempts", apiHandler.CreateQuizAttemptHandler)
		}
	}

	// --- Start Server ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("User service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}


