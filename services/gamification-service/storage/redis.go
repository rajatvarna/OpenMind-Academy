package storage

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const (
	leaderboardKey = "leaderboard:global"
	userProfileKeyPrefix = "user:"
)

// GamificationStore handles all database operations with Redis.
type GamificationStore struct {
	client *redis.Client
}

// NewGamificationStore creates a new store with a Redis client.
func NewGamificationStore(client *redis.Client) *GamificationStore {
	return &GamificationStore{client: client}
}

// AddPointsForUser adds a specified number of points to a user.
// It also updates the leaderboard with the new total score.
func (s *GamificationStore) AddPointsForUser(ctx context.Context, userID int64, points int) (int64, error) {
	userKey := fmt.Sprintf("%s%d", userProfileKeyPrefix, userID)

	// Increment the user's score in their profile hash
	newScore, err := s.client.HIncrBy(ctx, userKey, "score", int64(points)).Result()
	if err != nil {
		return 0, err
	}

	// Update the user's score in the global leaderboard sorted set
	err = s.client.ZAdd(ctx, leaderboardKey, &redis.Z{
		Score:  float64(newScore),
		Member: fmt.Sprintf("%d", userID),
	}).Err()
	if err != nil {
		// Attempt to roll back the score increment if the leaderboard update fails
		s.client.HIncrBy(ctx, userKey, "score", int64(-points))
		return 0, err
	}

	return newScore, nil
}

// GetUserRank finds the rank of a user in the leaderboard (0-indexed).
func (s *GamificationStore) GetUserRank(ctx context.Context, userID int64) (int64, error) {
	// ZRevRank is used because higher scores are better.
	return s.client.ZRevRank(ctx, leaderboardKey, fmt.Sprintf("%d", userID)).Result()
}

// GetTopUsers retrieves the top N users from the leaderboard.
func (s *GamificationStore) GetTopUsers(ctx context.Context, count int64) ([]redis.Z, error) {
	// ZRevRangeWithScores returns members from highest to lowest score.
	return s.client.ZRevRangeWithScores(ctx, leaderboardKey, 0, count-1).Result()
}

// GetUserStats retrieves all stats for a given user.
func (s *GamificationStore) GetUserStats(ctx context.Context, userID int64) (map[string]string, error) {
	userKey := fmt.Sprintf("%s%d", userProfileKeyPrefix, userID)
	return s.client.HGetAll(ctx, userKey).Result()
}
