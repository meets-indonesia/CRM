package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kevinnaserwan/crm-be/services/feedback/domain/entity"
	"github.com/streadway/amqp"
)

type rabbitMQEventPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// NewRabbitMQEventPublisher creates a new RabbitMQ event publisher
func NewRabbitMQEventPublisher(url string, exchange string) (*rabbitMQEventPublisher, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare the exchange
	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	return &rabbitMQEventPublisher{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
	}, nil
}

// PublishFeedbackCreated publishes an event when a feedback is created
func (p *rabbitMQEventPublisher) PublishFeedbackCreated(feedback *entity.Feedback) error {
	// Create event data
	event := map[string]interface{}{
		"id":         feedback.ID,
		"user_id":    feedback.UserID,
		"title":      feedback.Title,
		"content":    feedback.Content,
		"status":     feedback.Status,
		"created_at": feedback.CreatedAt,
	}

	// Marshal the event to JSON
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish the event
	err = p.channel.Publish(
		p.exchange,         // exchange
		"feedback.created", // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published feedback.created event for feedback ID: %d", feedback.ID)
	return nil
}

// PublishFeedbackResponded publishes an event when a feedback is responded to
func (p *rabbitMQEventPublisher) PublishFeedbackResponded(feedback *entity.Feedback) error {
	// Create event data
	event := map[string]interface{}{
		"id":         feedback.ID,
		"user_id":    feedback.UserID,
		"title":      feedback.Title,
		"content":    feedback.Content,
		"status":     feedback.Status,
		"response":   feedback.Response,
		"created_at": feedback.CreatedAt,
		"updated_at": feedback.UpdatedAt,
	}

	// Marshal the event to JSON
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish the event
	err = p.channel.Publish(
		p.exchange,           // exchange
		"feedback.responded", // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published feedback.responded event for feedback ID: %d", feedback.ID)
	return nil
}

// Close closes the connection and channel
func (p *rabbitMQEventPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
