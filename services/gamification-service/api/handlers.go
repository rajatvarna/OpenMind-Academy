package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/free-education/gamification-service/storage"
	"github.com/gin-gonic/gin"
)

// API holds dependencies for the API handlers.
type API struct {
	store *storage.GamificationStore
}

// NewAPI creates a new API struct.
func NewAPI(store *storage.GamificationStore) *API {
	return &API{store: store}
}

// GetStatsHandler handles fetching gamification stats for a user.
func (a *API) GetStatsHandler(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := a.store.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Failed to get stats for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user stats"})
		return
	}

	// If no stats are found, return an empty object instead of an error.
	if stats == nil {
		stats = make(map[string]string)
	}

	c.JSON(http.StatusOK, stats)
}

// GetLeaderboardHandler handles fetching the top users for the leaderboard.
func (a *API) GetLeaderboardHandler(c *gin.Context) {
	// Get top 10 users by default. Could be a query param.
	topUsers, err := a.store.GetTopUsers(c.Request.Context(), 10)
	if err != nil {
		log.Printf("Failed to get leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaderboard"})
		return
	}

	// The data from redis.Z needs to be formatted into a more friendly JSON structure.
	type LeaderboardEntry struct {
		UserID string  `json:"user_id"`
		Score  float64 `json:"score"`
	}

	response := make([]LeaderboardEntry, len(topUsers))
	for i, user := range topUsers {
		response[i] = LeaderboardEntry{
			UserID: user.Member.(string),
			Score:  user.Score,
		}
	}

	c.JSON(http.StatusOK, response)
}
