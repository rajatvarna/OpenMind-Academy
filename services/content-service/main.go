package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/free-education/content-service/api"
	"github.com/free-education/content-service/storage"
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
	contentStore := storage.NewContentStore(dbpool)
	apiHandler := api.NewAPI(contentStore)

	// --- Router Setup ---
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// API v1 routes
	// Note: In a real scenario, POST/PUT/DELETE endpoints would be protected by an auth middleware.
	v1 := router.Group("/api/v1")
	{
		v1.GET("/courses", apiHandler.GetAllCoursesHandler)
		v1.GET("/courses/featured", apiHandler.GetFeaturedCoursesHandler)
		v1.POST("/courses", apiHandler.CreateCourseHandler)
		v1.GET("/courses/:courseId", apiHandler.GetCourseHandler)
		v1.DELETE("/courses/:courseId", apiHandler.DeleteCourseHandler) // New route
		v1.POST("/lessons", apiHandler.CreateLessonHandler)

		// Review routes
		v1.POST("/reviews", apiHandler.CreateReviewHandler)
		v1.GET("/courses/:courseId/reviews", apiHandler.GetReviewsHandler)

		// Transcript route
		v1.PATCH("/lessons/:lessonId/transcript", apiHandler.UpdateTranscriptHandler)

		// User-specific routes
		v1.GET("/users/:userId/courses", apiHandler.GetCoursesForUserHandler)
	}

	// --- Start Server ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001" // Default port for the content-service
	}
	log.Printf("Content service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
