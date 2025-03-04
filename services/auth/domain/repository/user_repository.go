package repository

import (
	"context"
	"time"

	"github.com/kevinnaserwan/crm-be/services/auth/domain/entity"
)

// UserRepository mendefinisikan operasi-operasi repository untuk User
type UserRepository interface {
	// User operations
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error

	// OTP operations
	CreateOTP(ctx context.Context, otp *entity.OTP) error
	FindOTPByCode(ctx context.Context, code string) (*entity.OTP, error)
	DeleteOTP(ctx context.Context, id uint) error
}

// EventPublisher mendefinisikan operasi-operasi untuk mempublikasikan event
type EventPublisher interface {
	PublishUserCreated(userID uint, email string, role entity.Role) error
	PublishUserLoggedIn(userID uint, email string, role entity.Role) error
	PublishPasswordReset(userID uint, email string) error
}

// EmailSender mendefinisikan operasi-operasi untuk mengirim email
type EmailSender interface {
	SendOTP(email, name, otp string) error
}

// GoogleAuthProvider mendefinisikan operasi-operasi untuk autentikasi Google
type GoogleAuthProvider interface {
	VerifyIDToken(idToken string) (string, string, string, error)
}

// TokenGenerator mendefinisikan operasi-operasi untuk membuat token
type TokenGenerator interface {
	GenerateToken(userID uint, email string, role entity.Role) (string, time.Time, error)
	ValidateToken(tokenString string) (uint, string, entity.Role, error)
}

// OTPGenerator mendefinisikan operasi-operasi untuk membuat OTP
type OTPGenerator interface {
	GenerateOTP() (string, error)
}
