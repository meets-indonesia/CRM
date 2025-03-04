package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/kevinnaserwan/crm-be/services/article/config"
	"github.com/kevinnaserwan/crm-be/services/article/domain/entity"
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
	EventArticleCreated = "article.created"
	EventArticleUpdated = "article.updated"
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

// PublishArticleCreated publishes an article created event
func (r *RabbitMQ) PublishArticleCreated(article *entity.Article) error {
	event := map[string]interface{}{
		"article_id":   article.ID,
		"title":        article.Title,
		"summary":      article.Summary,
		"author_id":    article.AuthorID,
		"created_at":   article.CreatedAt,
		"published_at": article.PublishedAt,
	}

	return r.publishEvent(EventArticleCreated, event)
}

// PublishArticleUpdated publishes an article updated event
func (r *RabbitMQ) PublishArticleUpdated(article *entity.Article) error {
	event := map[string]interface{}{
		"article_id": article.ID,
		"title":      article.Title,
		"summary":    article.Summary,
		"author_id":  article.AuthorID,
		"updated_at": article.UpdatedAt,
	}

	return r.publishEvent(EventArticleUpdated, event)
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
