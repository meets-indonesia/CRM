package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *model.OTP) error
	FindValidOTP(ctx context.Context, userID uuid.UUID, code string) (*model.OTP, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
}
