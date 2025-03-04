package email

import (
	"fmt"
	"net/smtp"

	"github.com/kevinnaserwan/crm-be/services/notification/config"
)

// SMTPEmailSender implements EmailSender with SMTP
type SMTPEmailSender struct {
	config config.SMTPConfig
}

// NewSMTPEmailSender creates a new SMTPEmailSender
func NewSMTPEmailSender(config config.SMTPConfig) *SMTPEmailSender {
	return &SMTPEmailSender{
		config: config,
	}
}

// SendEmail sends an email using SMTP
func (s *SMTPEmailSender) SendEmail(to, subject, body string, isHTML bool) error {
	// Combine host and port
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	// Set up authentication information
	auth := smtp.PlainAuth("", s.config.User, s.config.Password, s.config.Host)

	// Set content type
	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}

	// Construct message
	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: %s; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.config.From, to, subject, contentType, body))

	// Send email
	return smtp.SendMail(addr, auth, s.config.From, []string{to}, message)
}
