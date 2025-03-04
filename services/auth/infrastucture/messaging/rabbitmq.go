package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/kevinnaserwan/crm-be/services/auth/config"
	"github.com/kevinnaserwan/crm-be/services/auth/domain/entity"
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
	EventUserCreated   = "user.created"
	EventUserLoggedIn  = "user.logged_in"
	EventPasswordReset = "user.password_reset"
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

// PublishUserCreated publishes a user created event
func (r *RabbitMQ) PublishUserCreated(userID uint, email string, role entity.Role) error {
	event := map[string]interface{}{
		"user_id": userID,
		"email":   email,
		"role":    role,
	}

	return r.publishEvent(EventUserCreated, event)
}

// PublishUserLoggedIn publishes a user logged in event
func (r *RabbitMQ) PublishUserLoggedIn(userID uint, email string, role entity.Role) error {
	event := map[string]interface{}{
		"user_id": userID,
		"email":   email,
		"role":    role,
	}

	return r.publishEvent(EventUserLoggedIn, event)
}

// PublishPasswordReset publishes a password reset event
func (r *RabbitMQ) PublishPasswordReset(userID uint, email string) error {
	event := map[string]interface{}{
		"user_id": userID,
		"email":   email,
	}

	return r.publishEvent(EventPasswordReset, event)
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
