package messaging

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient handles the connection and publishing to RabbitMQ.
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQClient creates and returns a new RabbitMQClient.
func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{conn: conn, channel: ch}, nil
}

// Publish sends a message to a specific queue.
func (c *RabbitMQClient) Publish(ctx context.Context, queueName string, eventType string, payload interface{}) error {
	body, err := json.Marshal(map[string]interface{}{
		"eventType": eventType,
		"payload":   payload,
	})
	if err != nil {
		return err
	}

	err = c.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s", body)
	return nil
}

// Close closes the RabbitMQ connection and channel.
func (c *RabbitMQClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
