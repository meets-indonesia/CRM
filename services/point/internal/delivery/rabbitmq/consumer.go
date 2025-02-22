package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel    *amqp.Channel
	exchange   string
	queue      string
	routingKey string
}

func NewConsumer(
	channel *amqp.Channel,
	exchange,
	queue,
	routingKey string,

) (*Consumer, error) {
	// Deklarasi exchange
	err := channel.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Deklarasi queue
	q, err := channel.QueueDeclare(
		queue, // nama queue
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	// Binding queue ke exchange
	err = channel.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		channel:    channel,
		exchange:   exchange,
		queue:      queue,
		routingKey: routingKey,
	}, nil
}

func (c *Consumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Process message
			err := c.processMessage(d)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				// Nack message jika gagal
				d.Nack(false, true)
			} else {
				// Ack message jika berhasil
				d.Ack(false)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

func (c *Consumer) processMessage(d amqp.Delivery) error {
	// Contoh struct untuk message
	var message struct {
		OrderID uint   `json:"order_id"`
		Status  string `json:"status"`
	}

	// Decode message
	if err := json.Unmarshal(d.Body, &message); err != nil {
		return err
	}

	// Process sesuai routing key
	// switch d.RoutingKey {
	// case "order.created":
	// 	return c.useCase.HandleOrderCreated(message.OrderID)
	// case "order.updated":
	// 	return c.useCase.HandleOrderUpdated(message.OrderID, message.Status)
	// default:
	// 	log.Printf("Unknown routing key: %s", d.RoutingKey)
	// }

	return nil
}
