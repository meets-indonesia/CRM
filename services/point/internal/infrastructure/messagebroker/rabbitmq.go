package messagebroker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const (
	FeedbackCreatedQueue = "feedback.created"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type FeedbackCreatedEvent struct {
	FeedbackID string `json:"feedback_id"`
	UserID     string `json:"user_id"`
	CreatedAt  string `json:"created_at"`
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	log.Printf("Connecting to RabbitMQ at %s", url)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare queues
	log.Printf("Declaring queue: %s", FeedbackCreatedQueue)
	_, err = ch.QueueDeclare(
		FeedbackCreatedQueue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	log.Printf("RabbitMQ connection established successfully")
	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQ) ConsumeFeedbackEvents(pointUseCase interface{}) {
	log.Printf("Starting consumer for queue: %s", FeedbackCreatedQueue)
	msgs, err := r.channel.Consume(
		FeedbackCreatedQueue,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf("Failed to register consumer: %v", err)
		return
	}

	log.Printf("Consumer registered successfully. Waiting for messages...")

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", string(d.Body))

			var event FeedbackCreatedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			log.Printf("Parsed event: UserID=%s, FeedbackID=%s", event.UserID, event.FeedbackID)

			userID, err := uuid.Parse(event.UserID)
			if err != nil {
				log.Printf("Error parsing user ID: %v", err)
				continue
			}

			feedbackID, err := uuid.Parse(event.FeedbackID)
			if err != nil {
				log.Printf("Error parsing feedback ID: %v", err)
				continue
			}

			// Handle the feedback with the process function passed in
			if processor, ok := pointUseCase.(interface {
				ProcessFeedbackPoint(ctx context.Context, userID, feedbackID uuid.UUID) error
			}); ok {
				ctx := context.Background()
				if err := processor.ProcessFeedbackPoint(ctx, userID, feedbackID); err != nil {
					log.Printf("Error processing feedback points: %v", err)
				} else {
					log.Printf("Successfully processed points for feedback %s", feedbackID)
				}
			} else {
				log.Printf("Point use case doesn't implement ProcessFeedbackPoint method")
			}
		}
	}()
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	log.Printf("RabbitMQ connection closed")
}
