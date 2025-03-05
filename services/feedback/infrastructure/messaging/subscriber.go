package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kevinnaserwan/crm-be/services/feedback/config"
	"github.com/kevinnaserwan/crm-be/services/feedback/domain/usecase"
	"github.com/streadway/amqp"
)

// RabbitMQSubscriber implements event subscriber
type RabbitMQSubscriber struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	feedbackUsecase usecase.FeedbackUsecase
}

// NewRabbitMQSubscriber creates a new RabbitMQSubscriber
func NewRabbitMQSubscriber(config config.RabbitMQConfig, feedbackUsecase usecase.FeedbackUsecase) (*RabbitMQSubscriber, error) {
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

	return &RabbitMQSubscriber{
		conn:            conn,
		channel:         ch,
		feedbackUsecase: feedbackUsecase,
	}, nil
}

// SubscribeToEvents subscribes to events
func (s *RabbitMQSubscriber) SubscribeToEvents() error {
	// Declare user exchange
	err := s.channel.ExchangeDeclare(
		"user.events", // name
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

	// Declare queue for user events
	queue, err := s.channel.QueueDeclare(
		"feedback.user.queue", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return err
	}

	// Bind queue to exchange for user events
	err = s.channel.QueueBind(
		queue.Name,     // queue name
		"user.updated", // routing key
		"user.events",  // exchange
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	// Also bind to user.created events
	err = s.channel.QueueBind(
		queue.Name,     // queue name
		"user.created", // routing key
		"user.events",  // exchange
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	// Consume messages
	msgs, err := s.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
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
				s.handleUserCreated(msg)
			case "user.updated":
				s.handleUserUpdated(msg)
			}

			msg.Ack(false)
		}
	}()

	return nil
}

// handleUserCreated handles user.created events
func (s *RabbitMQSubscriber) handleUserCreated(msg amqp.Delivery) {
	var payload struct {
		UserID uint   `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
	}

	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error unmarshaling user.created event: %v", err)
		return
	}

	log.Printf("User created: %d, %s", payload.UserID, payload.Name)

	// Here you would implement your logic for handling user creation
	// This might involve caching user data or other operations specific
	// to your feedback service's needs
}

// handleUserUpdated handles user.updated events
func (s *RabbitMQSubscriber) handleUserUpdated(msg amqp.Delivery) {
	var payload struct {
		UserID uint   `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
	}

	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error unmarshaling user.updated event: %v", err)
		return
	}

	log.Printf("User updated: %d, %s", payload.UserID, payload.Name)

	// Here you would implement your logic for handling user updates
	// This might involve updating cached user data or performing
	// operations on feedback associated with this user
}

// Close closes the RabbitMQ connection
func (s *RabbitMQSubscriber) Close() error {
	if s.channel != nil {
		s.channel.Close()
	}

	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}
