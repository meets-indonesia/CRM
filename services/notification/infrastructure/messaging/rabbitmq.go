package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kevinnaserwan/crm-be/services/notification/config"
	"github.com/kevinnaserwan/crm-be/services/notification/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/notification/domain/usecase"
	"github.com/streadway/amqp"
)

// RabbitMQ implements EventSubscriber
type RabbitMQ struct {
	conn                *amqp.Connection
	channel             *amqp.Channel
	notificationUsecase usecase.NotificationUsecase
}

// NewRabbitMQ creates a new RabbitMQ instance
func NewRabbitMQ(cfg config.RabbitMQConfig, notificationUsecase usecase.NotificationUsecase) (*RabbitMQ, error) {
	// Connect to RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

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
		conn:                conn,
		channel:             ch,
		notificationUsecase: notificationUsecase,
	}, nil
}

// SubscribeToEvents subscribes to relevant events
func (r *RabbitMQ) SubscribeToEvents() error {
	// Declare exchanges
	exchanges := []string{"article.events", "feedback.events", "reward.events"}
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
			return err
		}
	}

	// Declare queue
	queue, err := r.channel.QueueDeclare(
		"notification.events.queue", // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	if err != nil {
		return err
	}

	// Bind queue to exchanges
	bindings := []struct {
		exchange string
		key      string
	}{
		{"article.events", "article.created"},
		{"feedback.events", "feedback.responded"},
		{"reward.events", "reward.claimed"},
		{"reward.events", "reward.claim_status_updated"},
	}

	for _, binding := range bindings {
		err = r.channel.QueueBind(
			queue.Name,       // queue name
			binding.key,      // routing key
			binding.exchange, // exchange
			false,            // no-wait
			nil,              // arguments
		)
		if err != nil {
			return err
		}
	}

	// Consume messages
	msgs, err := r.channel.Consume(
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
			case "article.created":
				r.handleArticleCreated(msg)
			case "feedback.responded":
				r.handleFeedbackResponded(msg)
			case "reward.claimed", "reward.claim_status_updated":
				r.handleRewardEvent(msg)
			}

			msg.Ack(false)
		}
	}()

	return nil
}

// handleArticleCreated handles article.created events
func (r *RabbitMQ) handleArticleCreated(msg amqp.Delivery) {
	var payload struct {
		ArticleID   uint      `json:"article_id"`
		Title       string    `json:"title"`
		Summary     string    `json:"summary"`
		AuthorID    uint      `json:"author_id"`
		PublishedAt time.Time `json:"published_at"`
	}

	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error unmarshaling article.created event: %v", err)
		return
	}

	err := r.notificationUsecase.HandleArticleCreated(
		payload.ArticleID,
		payload.Title,
		payload.Summary,
		payload.AuthorID,
		payload.PublishedAt,
	)
	if err != nil {
		log.Printf("Error handling article.created event: %v", err)
	}
}

// handleFeedbackResponded handles feedback.responded events
func (r *RabbitMQ) handleFeedbackResponded(msg amqp.Delivery) {
	var payload struct {
		FeedbackID uint   `json:"feedback_id"`
		UserID     uint   `json:"user_id"`
		Title      string `json:"title"`
		Response   string `json:"response"`
	}

	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error unmarshaling feedback.responded event: %v", err)
		return
	}

	err := r.notificationUsecase.HandleFeedbackResponded(
		payload.FeedbackID,
		payload.UserID,
		payload.Title,
		payload.Response,
	)
	if err != nil {
		log.Printf("Error handling feedback.responded event: %v", err)
	}
}

// handleRewardEvent handles reward.claimed and reward.claim_status_updated events
func (r *RabbitMQ) handleRewardEvent(msg amqp.Delivery) {
	var payload struct {
		ClaimID  uint                      `json:"claim_id"`
		UserID   uint                      `json:"user_id"`
		RewardID uint                      `json:"reward_id"`
		Status   entity.NotificationStatus `json:"status"`
	}

	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Error unmarshaling reward event: %v", err)
		return
	}

	err := r.notificationUsecase.HandleRewardClaimed(
		payload.ClaimID,
		payload.UserID,
		payload.RewardID,
		payload.Status,
	)
	if err != nil {
		log.Printf("Error handling reward event: %v", err)
	}
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
