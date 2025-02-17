// services/auth/internal/usecase/forgot_password_usecase.go
package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/repository"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type ForgotPasswordUseCase struct {
	userRepo     repository.UserRepository
	otpRepo      repository.OTPRepository
	emailService *util.EmailService
}

func NewForgotPasswordUseCase(
	userRepo repository.UserRepository,
	otpRepo repository.OTPRepository,
	emailService *util.EmailService,
) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		emailService: emailService,
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

	if err := uc.otpRepo.Create(ctx, otp); err != nil {
		return err
	}

	// Send email with OTP
	if err := uc.emailService.SendOTP(email, otp.Code); err != nil {
		return fmt.Errorf("failed to send OTP: %w", err)
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
