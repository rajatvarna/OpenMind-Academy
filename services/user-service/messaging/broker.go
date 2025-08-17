package messaging

import "context"

// MessageBroker defines the interface for a message broker client.
// This allows for easy mocking in tests.
type MessageBroker interface {
	Publish(ctx context.Context, queueName string, eventType string, payload interface{}) error
	Close()
}
