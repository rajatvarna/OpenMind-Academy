package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/free-education/forum-service/api"
	"github.com/free-education/forum-service/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://platform_user:password@db:5432/platform_db"
	}

	dbpool, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	store := storage.NewForumStore(dbpool)
	apiHandler := api.NewAPI(store)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	v1 := router.Group("/api/v1")
	{
		// In a real app, POST routes would be protected
		v1.POST("/threads", apiHandler.CreateThreadHandler)
		v1.POST("/posts", apiHandler.CreatePostHandler)
		v1.GET("/courses/:courseId/threads", apiHandler.GetThreadsForCourseHandler)
		v1.GET("/threads/:threadId/posts", apiHandler.GetPostsForThreadHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3006" // Port for the forum service
	}
	router.Run(":" + port)
}
