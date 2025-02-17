// services/auth/internal/usecase/forgot_password_usecase.go
package usecase

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/events"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/repository"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/infrastructure/messagebroker"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type ForgotPasswordUseCase struct {
	userRepo      repository.UserRepository
	otpRepo       repository.OTPRepository
	emailService  *util.EmailService
	messageBroker *messagebroker.RabbitMQ
}

func NewForgotPasswordUseCase(
	userRepo repository.UserRepository,
	otpRepo repository.OTPRepository,
	emailService *util.EmailService,
	messageBroker *messagebroker.RabbitMQ,
) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{
		userRepo:      userRepo,
		otpRepo:       otpRepo,
		emailService:  emailService,
		messageBroker: messageBroker,
	}
}

func (uc *ForgotPasswordUseCase) RequestPasswordReset(ctx context.Context, email string) error {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("email not found")
	}

	// Generate OTP
	otp := &model.OTP{
		UserID:    user.ID,
		Code:      generateOTP(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	// Save OTP to database
	if err := uc.otpRepo.Create(ctx, otp); err != nil {
		return err
	}

	// Print OTP to console for debugging (this is important!)
	log.Printf("\n============================================")
	log.Printf("Generated OTP for %s: %s", email, otp.Code)
	log.Printf("============================================\n")

	// Publish password reset requested event
	event := events.PasswordResetRequestedEvent{
		UserID: user.ID.String(),
		Email:  user.Email,
		OTP:    otp.Code,
	}

	// Publish the password reset event
	if err = uc.messageBroker.PublishMessage(ctx, "auth.events", "password.reset.requested", event); err != nil {
		log.Printf("Warning: Failed to publish password reset event: %v", err)
	}

	// Create and publish email event
	emailEvent := events.EmailEvent{
		To:      user.Email,
		Subject: "Password Reset OTP",
		Body:    fmt.Sprintf("Your OTP for password reset is: %s\nThis code will expire in 15 minutes.", otp.Code),
	}

	if err = uc.messageBroker.PublishMessage(ctx, "email.events", "email.send", emailEvent); err != nil {
		log.Printf("Warning: Failed to publish email event: %v", err)
	}

	// Try to send email directly
	if err := uc.emailService.SendOTP(email, otp.Code); err != nil {
		log.Printf("Warning: Failed to send email directly: %v", err)
		// Don't return error to not reveal if email exists
	}

	return nil
}

func (uc *ForgotPasswordUseCase) VerifyOTP(ctx context.Context, email, code string) error {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("email not found")
	}

	otp, err := uc.otpRepo.FindValidOTP(ctx, user.ID, code)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	return uc.otpRepo.MarkAsUsed(ctx, otp.ID)
}

func (uc *ForgotPasswordUseCase) ResetPassword(ctx context.Context, email, newPassword string) error {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("email not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return uc.userRepo.Update(ctx, user)
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	otp := ""
	for i := 0; i < 6; i++ {
		otp += string(digits[rand.Intn(len(digits))])
	}
	return otp
}
