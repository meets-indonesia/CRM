package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/kevinnaserwan/crm-be/services/reward/config"
	"github.com/kevinnaserwan/crm-be/services/reward/domain/entity"
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
	EventRewardClaimed      = "reward.claimed"
	EventClaimStatusUpdated = "reward.claim_status_updated"
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

// PublishRewardClaimed publishes a reward claimed event
func (r *RabbitMQ) PublishRewardClaimed(claim *entity.RewardClaim) error {
	event := map[string]interface{}{
		"claim_id":   claim.ID,
		"user_id":    claim.UserID,
		"reward_id":  claim.RewardID,
		"points":     claim.PointCost,
		"created_at": claim.CreatedAt,
	}

	return r.publishEvent(EventRewardClaimed, event)
}

// PublishClaimStatusUpdated publishes a claim status updated event
func (r *RabbitMQ) PublishClaimStatusUpdated(claim *entity.RewardClaim) error {
	event := map[string]interface{}{
		"claim_id":   claim.ID,
		"user_id":    claim.UserID,
		"reward_id":  claim.RewardID,
		"status":     claim.Status,
		"updated_at": claim.UpdatedAt,
	}

	return r.publishEvent(EventClaimStatusUpdated, event)
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
