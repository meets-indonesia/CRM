package messagebroker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Exchange names
const (
	ArticleExchange      = "article.events"
	NotificationExchange = "notification.events"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Initialize RabbitMQ instance
	rmq := &RabbitMQ{
		conn:    conn,
		channel: ch,
	}

	// Setup exchanges
	if err := rmq.setupExchanges(); err != nil {
		return nil, err
	}

	return rmq, nil
}

func (r *RabbitMQ) setupExchanges() error {
	// Declare exchanges
	exchanges := []string{ArticleExchange, NotificationExchange}
	for _, exchange := range exchanges {
		err := r.channel.ExchangeDeclare(
			exchange, // name
			"topic",  // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", exchange, err)
		}
		log.Printf("Declared exchange: %s", exchange)
	}
	return nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, exchange, routingKey string, message interface{}) error {
	// Convert message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish message
	err = r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to exchange %s with routing key %s", exchange, routingKey)
	return nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
