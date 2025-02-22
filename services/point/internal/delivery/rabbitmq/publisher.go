package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Publisher struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

func NewPublisher(channel *amqp.Channel, exchange, routingKey string) (*Publisher, error) {
	// Deklarasi exchange
	err := channel.ExchangeDeclare(
		exchange, // nama exchange
		"topic",  // tipe exchange
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return nil, err
	}

	return &Publisher{
		channel:    channel,
		exchange:   exchange,
		routingKey: routingKey,
	}, nil
}

func (p *Publisher) Publish(message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		p.exchange,   // exchange
		p.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // pesan akan disimpan ke disk
		},
	)
}
