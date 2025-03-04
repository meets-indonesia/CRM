package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kevinnaserwan/crm-be/services/user/config"
	"github.com/kevinnaserwan/crm-be/services/user/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/user/domain/repository"
	"github.com/streadway/amqp"
)

// Event types
const (
	EventUserCreated     = "user.created"
	EventFeedbackCreated = "feedback.created"
	EventRewardClaimed   = "reward.claimed"
)

// RabbitMQ implements EventSubscriber
type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	eventHandler repository.UserEventProcessor
	pointHandler repository.PointEventProcessor
}

// NewRabbitMQ creates a new RabbitMQ instance
func NewRabbitMQ(
	config config.RabbitMQConfig,
	eventHandler repository.UserEventProcessor,
	pointHandler repository.PointEventProcessor,
) (*RabbitMQ, error) {
	// Connect to RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.User, config.Password, config.Host, config.Port)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		conn:         conn,
		channel:      ch,
		eventHandler: eventHandler,
		pointHandler: pointHandler,
	}, nil
}

// SubscribeToUserEvents subscribes to user-related events
func (r *RabbitMQ) SubscribeToUserEvents() error {
	// Declare auth events exchange
	err := r.channel.ExchangeDeclare(
		"auth.events", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	// Declare feedback events exchange
	err = r.channel.ExchangeDeclare(
		"feedback.events", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}

	// Declare reward events exchange
	err = r.channel.ExchangeDeclare(
		"reward.events", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	// Declare queue for user events
	userQueue, err := r.channel.QueueDeclare(
		"user.events.queue", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return err
	}

	// Bind queue to exchanges with routing keys
	err = r.channel.QueueBind(
		userQueue.Name, // queue name
		"user.created", // routing key
		"auth.events",  // exchange
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	err = r.channel.QueueBind(
		userQueue.Name,     // queue name
		"feedback.created", // routing key
		"feedback.events",  // exchange
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	err = r.channel.QueueBind(
		userQueue.Name,   // queue name
		"reward.claimed", // routing key
		"reward.events",  // exchange
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	// Consume messages
	msgs, err := r.channel.Consume(
		userQueue.Name, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return err
	}

	// Process messages in a goroutine
	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", msg.RoutingKey)

			switch msg.RoutingKey {
			case "user.created":
				var payload struct {
					UserID uint        `json:"user_id"`
					Email  string      `json:"email"`
					Role   entity.Role `json:"role"`
				}

				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					log.Printf("Error unmarshaling user.created event: %v", err)
					msg.Nack(false, true) // Requeue the message
					continue
				}

				err := r.eventHandler.ProcessUserCreated(payload.UserID, payload.Email, payload.Role)
				if err != nil {
					log.Printf("Error processing user.created event: %v", err)
					msg.Nack(false, true) // Requeue the message
					continue
				}

			case "feedback.created":
				var payload struct {
					UserID     uint `json:"user_id"`
					FeedbackID uint `json:"feedback_id"`
				}

				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					log.Printf("Error unmarshaling feedback.created event: %v", err)
					msg.Nack(false, true) // Requeue the message
					continue
				}

				err := r.pointHandler.ProcessFeedbackCreated(payload.UserID, payload.FeedbackID)
				if err != nil {
					log.Printf("Error processing feedback.created event: %v", err)
					msg.Nack(false, true) // Requeue the message
					continue
				}

			case "reward.claimed":
				var payload struct {
					UserID   uint `json:"user_id"`
					RewardID uint `json:"reward_id"`
					Points   int  `json:"points"`
				}

				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					log.Printf("Error unmarshaling reward.claimed event: %v", err)
					msg.Nack(false, true) // Requeue the message
					continue
				}

				err := r.pointHandler.ProcessRewardClaimed(payload.UserID, payload.RewardID, payload.Points)
				if err != nil {
					log.Printf("Error processing reward.claimed event: %v", err)
					msg.Nack(false, true) // Requeue the message
					continue
				}

			default:
				log.Printf("Unknown routing key: %s", msg.RoutingKey)
			}

			msg.Ack(false) // Acknowledge the message
		}
	}()

	return nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
