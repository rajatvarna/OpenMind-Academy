package messaging

import "context"

// MessageHandler defines the function signature for handling consumed messages.
type MessageHandler func(body []byte)

// MessageBroker defines the interface for a message broker client.
// This allows for easy mocking in tests.
type MessageBroker interface {
	Publish(ctx context.Context, queueName string, eventType string, payload interface{}) error
	Consume(ctx context.Context, queueName string, handler MessageHandler) error
	Close()
}
