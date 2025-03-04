package service

import (
	"crypto/rand"
	"fmt"
	"net/smtp"

	"github.com/kevinnaserwan/crm-be/services/auth/config"
)

// OTPService implements OTPGenerator and EmailSender
type OTPService struct {
	smtpConfig config.SMTPConfig
}

// NewOTPService creates a new OTPService
func NewOTPService(smtpConfig config.SMTPConfig) *OTPService {
	return &OTPService{
		smtpConfig: smtpConfig,
	}
}

// GenerateOTP generates a random 6-digit OTP
func (s *OTPService) GenerateOTP() (string, error) {
	// Generate 6 random digits
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convert to 6 digits
	otp := fmt.Sprintf("%06d", int(b[0])*65536+int(b[1])*256+int(b[2]))

	return otp, nil
}

// SendOTP sends an OTP to the user's email
func (s *OTPService) SendOTP(email, name, otp string) error {
	// SMTP server configuration
	smtpHost := s.smtpConfig.Host
	smtpPort := s.smtpConfig.Port
	smtpUser := s.smtpConfig.User
	smtpPass := s.smtpConfig.Password
	from := s.smtpConfig.From

	// Message
	subject := "Password Reset OTP"
	body := fmt.Sprintf("Hello %s,\n\nYour OTP for password reset is: %s\n\nThis OTP will expire in 15 minutes.\n\nRegards,\nLRT CRM Team", name, otp)
	message := []byte("To: " + email + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// Authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Send mail
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, message)
	if err != nil {
		return err
	}

	return nil
}
