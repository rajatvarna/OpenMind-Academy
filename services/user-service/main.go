package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/free-education/user-service/api"
	"github.com/free-education/user-service/auth"
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

		// Authenticated routes
		authorized := v1.Group("/")
		authorized.Use(AuthMiddleware("any")) // Accessible to any authenticated user
		{
			authorized.GET("/profile", getProfileHandler)

			// Progress tracking routes
			authorized.GET("/users/:userId/progress", apiHandler.GetProgressHandler)
			authorized.POST("/users/:userId/progress", apiHandler.MarkLessonCompleteHandler)

			// Full profile route
			authorized.GET("/users/:userId/full-profile", apiHandler.GetFullProfileHandler)
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

// AuthMiddleware creates a middleware that validates a JWT and checks for a required role.
func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		// Check if the user has the required role.
		// A more complex system could support multiple roles, e.g., "admin" can do anything a "moderator" can.
		if claims.Role != requiredRole && requiredRole != "any" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		// Set user info in context for downstream handlers to use
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// getProfileHandler is a placeholder for a protected endpoint.
func getProfileHandler(c *gin.Context) {
	// We can retrieve the user ID from the context that the middleware set.
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "This is a protected profile endpoint", "user_id": userID})
}
