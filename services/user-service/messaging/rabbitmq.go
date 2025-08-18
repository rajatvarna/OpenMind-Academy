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

// Consume starts consuming messages from a queue and passes them to the handler.
func (c *RabbitMQClient) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
	q, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			handler(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages in %s. To exit press CTRL+C", q.Name)
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
