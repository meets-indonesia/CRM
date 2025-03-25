package email

import (
	"fmt"
	"log"
	"net/smtp"
)

// EmailService mendefinisikan service untuk mengirim email
type EmailService interface {
	SendNotification(subject string, body string) error
}

// SMTPEmailService implementasi email service dengan SMTP
type SMTPEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
	to       string
}

// NewSMTPEmailService membuat instance baru SMTPEmailService
func NewSMTPEmailService(host, port, username, password, from, to string) EmailService {
	return &SMTPEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		to:       to,
	}
}

// SendNotification mengirim notifikasi email
func (s *SMTPEmailService) SendNotification(subject string, body string) error {
	// Jika kredensial kosong, hanya log pesan
	if s.username == "" || s.password == "" {
		log.Printf("Email notification would be sent: Subject: %s, Body: %s", subject, body)
		return nil
	}

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s", s.from, s.to, subject, body)

	err := smtp.SendMail(addr, auth, s.from, []string{s.to}, []byte(msg))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email notification sent: Subject: %s", subject)
	return nil
}
