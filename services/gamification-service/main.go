package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/free-education/gamification-service/storage"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

const (
	eventQueueName = "gamification_events"
)

func main() {
	ctx := context.Background()

	// --- Connect to Redis ---
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis.")

	// --- Dependency Injection ---
	gamificationStore := storage.NewGamificationStore(rdb)
	eventHandler := NewEventHandler(gamificationStore)

	// --- Connect to RabbitMQ and start consumer ---
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	for {
		log.Println("Attempting to connect to RabbitMQ...")
		conn, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer conn.Close()
		log.Println("Successfully connected to RabbitMQ.")

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Failed to open a channel: %v", err)
			continue
		}
		defer ch.Close()

		q, err := ch.QueueDeclare(
			eventQueueName, // name
			true,           // durable
			false,          // delete when unused
			false,          // exclusive
			false,          // no-wait
			nil,            // arguments
		)
		if err != nil {
			log.Printf("Failed to declare a queue: %v", err)
			continue
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack is false, we will manually ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Printf("Failed to register a consumer: %v", err)
			continue
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
				err := eventHandler.HandleEvent(ctx, d.Body)
				if err != nil {
					log.Printf("Error handling event: %v. Message will be Nacked.", err)
					// Negative acknowledgement, don't requeue
					d.Nack(false, false)
				} else {
					log.Printf("Event handled successfully. Acknowledging message.")
					// Acknowledge the message
					d.Ack(false)
				}
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	}
}
