package messagebroker

import (
	"encoding/json"
	"log"

	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/infrastructure/email"
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
// Event payloads
type FeedbackCreatedEvent struct {
	FeedbackID   string  `json:"feedback_id"`
	UserID       string  `json:"user_id"`
	UserEmail    string  `json:"user_email"` // Add this
	FeedbackType string  `json:"feedback_type"`
	Station      string  `json:"station"`
	Feedback     string  `json:"feedback"`
	Rating       float64 `json:"rating"`
	CreatedAt    string  `json:"created_at"`
}

type FeedbackRespondedEvent struct {
	FeedbackID  string `json:"feedback_id"`
	UserID      string `json:"user_id"`
	UserEmail   string `json:"user_email"` // Add this
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

func (r *RabbitMQ) PublishFeedbackCreated(feedback *model.Feedback, userEmail string) error {
	event := FeedbackCreatedEvent{
		FeedbackID:   feedback.ID.String(),
		UserID:       feedback.UserID.String(),
		UserEmail:    userEmail, // Add this
		FeedbackType: feedback.FeedbackType.Name,
		Station:      feedback.Station.Name,
		Feedback:     feedback.Feedback,
		Rating:       feedback.Rating,
		CreatedAt:    feedback.FeedbackDate.Format("2006-01-02 15:04:05"),
	}

	return r.publishMessage(FeedbackCreatedQueue, event)
}

func (r *RabbitMQ) PublishFeedbackResponded(response *model.FeedbackResponse, userEmail string) error {
	event := FeedbackRespondedEvent{
		FeedbackID:  response.FeedbackID.String(),
		UserEmail:   userEmail, // Add this
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

func (r *RabbitMQ) ConsumeFeedbackEvents(emailService *email.EmailService) {
	// Consume feedback.created events
	go func() {
		msgs, err := r.channel.Consume(
			FeedbackCreatedQueue,
			"",    // consumer
			true,  // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			log.Printf("Failed to register feedback.created consumer: %v", err)
			return
		}

		for msg := range msgs {
			var event FeedbackCreatedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Error unmarshaling feedback.created event: %v", err)
				continue
			}

			feedbackData := map[string]interface{}{
				"type":    event.FeedbackType,
				"station": event.Station,
				"message": event.Feedback,
				"rating":  event.Rating,
			}

			// Send email to admin
			// Note: You'll need to get admin email from somewhere
			if err := emailService.SendFeedbackCreatedEmail(feedbackData, "admin@example.com"); err != nil {
				log.Printf("Error sending feedback created email: %v", err)
			}
		}
	}()

	// Consume feedback.responded events
	go func() {
		msgs, err := r.channel.Consume(
			FeedbackRespondedQueue,
			"",    // consumer
			true,  // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			log.Printf("Failed to register feedback.responded consumer: %v", err)
			return
		}

		for msg := range msgs {
			var event FeedbackRespondedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Error unmarshaling feedback.responded event: %v", err)
				continue
			}

			responseData := map[string]interface{}{
				"message": event.Response,
				"date":    event.RespondedAt,
			}

			// Send email to user
			if err := emailService.SendFeedbackResponseEmail(responseData, event.UserEmail); err != nil {
				log.Printf("Error sending feedback response email: %v", err)
			}
		}
	}()
}
