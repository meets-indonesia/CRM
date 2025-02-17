package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"gorm.io/gorm"
)

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *otpRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(ctx context.Context, otp *model.OTP) error {
	return r.db.WithContext(ctx).Create(otp).Error
}

func (r *otpRepository) FindValidOTP(ctx context.Context, userID uuid.UUID, code string) (*model.OTP, error) {
	var otp model.OTP
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND code = ? AND used = ? AND expires_at > ?",
			userID, code, false, time.Now()).
		First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.OTP{}).
		Where("id = ?", id).
		Update("used", true).Error
}
