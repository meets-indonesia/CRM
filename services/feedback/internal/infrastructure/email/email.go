package email

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	host     string
	port     string
	username string
	password string
}

func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		host:     config.Host,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
	}
}

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (s *EmailService) SendFeedbackCreatedEmail(feedback map[string]interface{}, adminEmail string) error {
	subject := "New Feedback Received"
	body := fmt.Sprintf(`
        New feedback has been received:

        Type: %s
        Station: %s
        Feedback: %s
        Rating: %.1f

        Please log in to respond to this feedback.
    `, feedback["type"], feedback["station"], feedback["message"], feedback["rating"])

	return s.sendEmail(adminEmail, subject, body)
}

func (s *EmailService) SendFeedbackResponseEmail(response map[string]interface{}, userEmail string) error {
	subject := "Response to Your Feedback"
	body := fmt.Sprintf(`
        Your feedback has received a response:

        Response: %s
        Date: %s

        Thank you for your feedback!
    `, response["message"], response["date"])

	return s.sendEmail(userEmail, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"\r\n"+
		"%s", to, subject, body))

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", s.host, s.port),
		auth,
		s.username,
		[]string{to},
		message,
	)
}
