package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/kevinnaserwan/crm-be/services/feedback/config"
	"github.com/kevinnaserwan/crm-be/services/feedback/domain/entity"
	"github.com/streadway/amqp"
)

// RabbitMQ implements EventPublisher
type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// Event types
const (
	EventFeedbackCreated   = "feedback.created"
	EventFeedbackResponded = "feedback.responded"
)

// NewRabbitMQ creates a new RabbitMQ instance
func NewRabbitMQ(config config.RabbitMQConfig) (*RabbitMQ, error) {
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

	// Declare exchange
	err = ch.ExchangeDeclare(
		config.Exchange, // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		conn:     conn,
		channel:  ch,
		exchange: config.Exchange,
	}, nil
}

// PublishFeedbackCreated publishes a feedback created event
func (r *RabbitMQ) PublishFeedbackCreated(feedback *entity.Feedback) error {
	event := map[string]interface{}{
		"feedback_id": feedback.ID,
		"user_id":     feedback.UserID,
		"title":       feedback.Title,
		"created_at":  feedback.CreatedAt,
	}

	return r.publishEvent(EventFeedbackCreated, event)
}

// PublishFeedbackResponded publishes a feedback responded event
func (r *RabbitMQ) PublishFeedbackResponded(feedback *entity.Feedback) error {
	event := map[string]interface{}{
		"feedback_id":  feedback.ID,
		"user_id":      feedback.UserID,
		"title":        feedback.Title,
		"response":     feedback.Response,
		"responded_at": feedback.UpdatedAt,
	}

	return r.publishEvent(EventFeedbackResponded, event)
}

// publishEvent publishes an event to RabbitMQ
func (r *RabbitMQ) publishEvent(eventType string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		r.exchange, // exchange
		eventType,  // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
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
