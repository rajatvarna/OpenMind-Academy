package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/free-education/user-service/api"
	"github.com/free-education/user-service/messaging"
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

	apiHandler := api.NewAPI(userStore, messageBroker)

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
		v1.POST("/password/forgot", apiHandler.ForgotPasswordHandler)
		v1.POST("/password/reset", apiHandler.ResetPasswordHandler)

		// Authenticated routes are now protected by the API Gateway
		v1.GET("/profile", apiHandler.GetProfileHandler)
		v1.GET("/preferences", apiHandler.GetUserPreferencesHandler)
		v1.PUT("/preferences", apiHandler.UpdateUserPreferencesHandler)
		v1.GET("/users/:userId/progress", apiHandler.GetProgressHandler)
		v1.POST("/users/:userId/progress", apiHandler.MarkLessonCompleteHandler)
		v1.GET("/users/:userId/quiz-attempts", apiHandler.GetQuizAttemptsForUserHandler)
		v1.GET("/users/:userId/full-profile", apiHandler.GetFullProfileHandler)

		// Authenticated routes - specific to the user
		v1.POST("/quiz-attempts", apiHandler.CreateQuizAttemptHandler)
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


