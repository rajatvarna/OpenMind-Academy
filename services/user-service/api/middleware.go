package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a gin middleware for authenticating users.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDHeader := c.GetHeader("X-User-Id")
		if userIDHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID header not provided"})
			c.Abort()
			return
		}

		id, err := strconv.ParseInt(userIDHeader, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
			c.Abort()
			return
		}

		// Set the user ID in the context for downstream handlers to use.
		c.Set("userID", id)

		c.Next()
	}
}
