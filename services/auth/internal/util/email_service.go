package util

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	host     string
	port     int
	username string
	password string
}

func NewEmailService(host string, port int, username, password string) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *EmailService) SendOTP(toEmail, otp string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	subject := "Password Reset OTP"
	body := fmt.Sprintf("Your OTP for password reset is: %s\nThis code will expire in 15 minutes.", otp)
	msg := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", toEmail, subject, body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.host, s.port),
		auth,
		s.username,
		[]string{toEmail},
		[]byte(msg),
	)

	return err
}
