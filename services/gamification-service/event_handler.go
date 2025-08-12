package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/free-education/gamification-service/storage"
)

// Define points awarded for different events
const pointsForLessonCompleted = 10

// EventHandler holds dependencies for handling events, like the storage layer.
type EventHandler struct {
	store *storage.GamificationStore
}

// NewEventHandler creates a new EventHandler.
func NewEventHandler(store *storage.GamificationStore) *EventHandler {
	return &EventHandler{store: store}
}

// Event represents a generic event received from the message queue.
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// LessonCompletedPayload is the specific structure for the 'lesson_completed' event.
type LessonCompletedPayload struct {
	UserID   int64 `json:"user_id"`
	LessonID int64 `json:"lesson_id"`
	CourseID int64 `json:"course_id"`
}

// HandleEvent parses the event and routes it to the appropriate handler.
func (h *EventHandler) HandleEvent(ctx context.Context, eventBody []byte) error {
	var event Event
	if err := json.Unmarshal(eventBody, &event); err != nil {
		log.Printf("Error unmarshalling event: %v", err)
		return err // Malformed message, don't requeue
	}

	log.Printf("Handling event of type: %s", event.Type)

	switch event.Type {
	case "lesson_completed":
		return h.handleLessonCompleted(ctx, event.Payload)
	// Add cases for other events like "course_completed", "user_registered", etc.
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return nil // Ignore unknown event types
	}
}

func (h *EventHandler) handleLessonCompleted(ctx context.Context, payload json.RawMessage) error {
	var p LessonCompletedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Printf("Error unmarshalling lesson_completed payload: %v", err)
		return err
	}

	log.Printf("Awarding %d points to user %d for completing lesson %d", pointsForLessonCompleted, p.UserID, p.LessonID)

	// Add points to the user
	newScore, err := h.store.AddPointsForUser(ctx, p.UserID, pointsForLessonCompleted)
	if err != nil {
		log.Printf("Failed to add points for user %d: %v", p.UserID, err)
		return fmt.Errorf("failed to process points for user %d: %w", p.UserID, err)
	}

	log.Printf("User %d now has %d points.", p.UserID, newScore)

	// Placeholder for badge logic
	// You could check the user's new score, number of completed lessons, etc.
	// and award badges accordingly.
	// For example: checkAndAwardBadges(ctx, p.UserID, newScore)

	return nil
}
