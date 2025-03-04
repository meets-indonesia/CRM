package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/kevinnaserwan/crm-be/services/inventory/config"
	"github.com/kevinnaserwan/crm-be/services/inventory/domain/entity"
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
	EventLowStockAlert = "inventory.low_stock"
	EventStockUpdated  = "inventory.stock_updated"
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

// PublishLowStockAlert publishes a low stock alert event
func (r *RabbitMQ) PublishLowStockAlert(item *entity.Item, deficit int) error {
	event := map[string]interface{}{
		"item_id":          item.ID,
		"name":             item.Name,
		"sku":              item.SKU,
		"current_stock":    item.CurrentStock,
		"minimum_stock":    item.MinimumStock,
		"deficit":          deficit,
		"reorder_quantity": item.ReorderQuantity,
	}

	return r.publishEvent(EventLowStockAlert, event)
}

// PublishStockUpdated publishes a stock updated event
func (r *RabbitMQ) PublishStockUpdated(transaction *entity.StockTransaction) error {
	event := map[string]interface{}{
		"transaction_id": transaction.ID,
		"item_id":        transaction.ItemID,
		"type":           transaction.Type,
		"quantity":       transaction.Quantity,
		"previous_qty":   transaction.PreviousQty,
		"new_qty":        transaction.NewQty,
		"reason":         transaction.Reason,
		"reference_id":   transaction.ReferenceID,
		"performed_by":   transaction.PerformedBy,
		"created_at":     transaction.CreatedAt,
	}

	return r.publishEvent(EventStockUpdated, event)
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
