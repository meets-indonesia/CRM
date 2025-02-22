package messagebroker

import (
	"encoding/json"
	"log"

	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"github.com/streadway/amqp"
)

// Queue names
const (
	FeedbackCreatedQueue   = "feedback.created"
	FeedbackRespondedQueue = "feedback.responded"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Event payloads
type FeedbackCreatedEvent struct {
	FeedbackID   string  `json:"feedback_id"`
	UserID       string  `json:"user_id"`
	FeedbackType string  `json:"feedback_type"`
	Station      string  `json:"station"`
	Feedback     string  `json:"feedback"`
	Rating       float64 `json:"rating"`
	CreatedAt    string  `json:"created_at"`
}

type FeedbackRespondedEvent struct {
	FeedbackID  string `json:"feedback_id"`
	UserID      string `json:"user_id"`
	Response    string `json:"response"`
	RespondedAt string `json:"responded_at"`
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare queues
	queues := []string{FeedbackCreatedQueue, FeedbackRespondedQueue}
	for _, queue := range queues {
		_, err = ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return nil, err
		}
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQ) PublishFeedbackCreated(feedback *model.Feedback) error {
	event := FeedbackCreatedEvent{
		FeedbackID:   feedback.ID.String(),
		UserID:       feedback.UserID.String(),
		FeedbackType: feedback.FeedbackType.Name,
		Station:      feedback.Station.Name,
		Feedback:     feedback.Feedback,
		Rating:       feedback.Rating,
		CreatedAt:    feedback.FeedbackDate.Format("2006-01-02 15:04:05"),
	}

	return r.publishMessage(FeedbackCreatedQueue, event)
}

func (r *RabbitMQ) PublishFeedbackResponded(response *model.FeedbackResponse) error {
	event := FeedbackRespondedEvent{
		FeedbackID:  response.FeedbackID.String(),
		Response:    response.Response,
		RespondedAt: response.ResponseDate.Format("2006-01-02 15:04:05"),
	}

	return r.publishMessage(FeedbackRespondedQueue, event)
}

func (r *RabbitMQ) publishMessage(queue string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

// ConsumeMessages is a generic consumer that can be used for testing
func (r *RabbitMQ) ConsumeMessages(queue string, handler func([]byte) error) {
	msgs, err := r.channel.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	<-forever
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
