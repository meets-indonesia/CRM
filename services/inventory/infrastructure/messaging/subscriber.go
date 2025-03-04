package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kevinnaserwan/crm-be/services/inventory/config"
	"github.com/kevinnaserwan/crm-be/services/inventory/domain/usecase"
	"github.com/streadway/amqp"
)

// RabbitMQSubscriber implements event subscriber
type RabbitMQSubscriber struct {
	conn             *amqp.Connection
	channel          *amqp.Channel
	inventoryUsecase usecase.InventoryUsecase
}

// NewRabbitMQSubscriber creates a new RabbitMQSubscriber
func NewRabbitMQSubscriber(config config.RabbitMQConfig, inventoryUsecase usecase.InventoryUsecase) (*RabbitMQSubscriber, error) {
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
		conn:             conn,
		channel:          ch,
		inventoryUsecase: inventoryUsecase,
	}, nil
}

// SubscribeToEvents subscribes to events
func (s *RabbitMQSubscriber) SubscribeToEvents() error {
	// Declare reward exchange
	err := s.channel.ExchangeDeclare(
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

	// Declare queue
	queue, err := s.channel.QueueDeclare(
		"inventory.reward.queue", // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return err
	}

	// Bind queue to exchange
	err = s.channel.QueueBind(
		queue.Name,       // queue name
		"reward.claimed", // routing key
		"reward.events",  // exchange
		false,            // no-wait
		nil,              // arguments
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

			if msg.RoutingKey == "reward.claimed" {
				s.handleRewardClaimed(msg)
			}

			msg.Ack(false)
		}
	}()

	return nil
}

// handleRewardClaimed handles reward.claimed events
func (s *RabbitMQSubscriber) handleRewardClaimed(msg amqp.Delivery) {
	var payload struct {
		ClaimID  uint `json:"claim_id"`
		UserID   uint `json:"user_id"`
		RewardID uint `json:"reward_id"`
		Quantity int  `json:"quantity"`
	}

	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error unmarshaling reward.claimed event: %v", err)
		return
	}

	// If quantity is not provided, default to 1
	if payload.Quantity <= 0 {
		payload.Quantity = 1
	}

	// Process reward claimed event
	err := s.inventoryUsecase.ProcessRewardClaimed(payload.ClaimID, payload.RewardID, payload.Quantity)
	if err != nil {
		log.Printf("Error processing reward.claimed event: %v", err)
	}
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
