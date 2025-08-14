package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/free-education/user-service/api"
	"github.com/free-education/user-service/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
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
	apiHandler := api.NewAPI(userStore)

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

		// Authenticated routes are now protected by the API Gateway
		v1.GET("/profile", apiHandler.GetProfileHandler)
		v1.GET("/users/:userId/progress", apiHandler.GetProgressHandler)
		v1.POST("/users/:userId/progress", apiHandler.MarkLessonCompleteHandler)
		v1.GET("/users/:userId/full-profile", apiHandler.GetFullProfileHandler)
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


