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

	qnaServiceURL := os.Getenv("QNA_SERVICE_URL")
	if qnaServiceURL == "" {
		qnaServiceURL = "http://qna-service:3003/generate-quiz"
		log.Println("QNA_SERVICE_URL not set, using default value.")
	}

	apiHandler := api.NewAPI(contentStore, qnaServiceURL)

	// --- Router Setup ---
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public, read-only routes
		v1.GET("/courses", apiHandler.GetAllCoursesHandler)
		v1.GET("/courses/featured", apiHandler.GetFeaturedCoursesHandler)
		v1.GET("/courses/:courseId", apiHandler.GetCourseHandler)
		v1.GET("/courses/:courseId/reviews", apiHandler.GetReviewsHandler)
		v1.GET("/users/:userId/courses", apiHandler.GetCoursesForUserHandler)
		v1.GET("/paths/:pathId", apiHandler.GetLearningPathHandler)
		v1.GET("/lessons/:lessonId/quiz", apiHandler.GetQuizByLessonIDHandler)

		// Authenticated routes (write operations)
		authRequired := v1.Group("/")
		authRequired.Use(api.AuthMiddleware())
		{
			authRequired.POST("/courses", apiHandler.CreateCourseHandler)
			authRequired.DELETE("/courses/:courseId", apiHandler.DeleteCourseHandler)
			authRequired.POST("/lessons", apiHandler.CreateLessonHandler)
			authRequired.POST("/reviews", apiHandler.CreateReviewHandler)
			authRequired.PATCH("/lessons/:lessonId/transcript", apiHandler.UpdateTranscriptHandler)
			authRequired.POST("/paths", apiHandler.CreateLearningPathHandler)
			authRequired.POST("/quizzes", apiHandler.CreateQuizHandler)
		}
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
