package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/feedback/domain/entity"
	"gorm.io/gorm"
)

// GormQRFeedbackRepository implements QRFeedbackRepository with GORM
type GormQRFeedbackRepository struct {
	db *gorm.DB
}

// NewGormQRFeedbackRepository creates a new GormQRFeedbackRepository
func NewGormQRFeedbackRepository(db *gorm.DB) *GormQRFeedbackRepository {
	return &GormQRFeedbackRepository{
		db: db,
	}
}

// Create creates a new QR feedback
func (r *GormQRFeedbackRepository) Create(ctx context.Context, qrFeedback *entity.QRFeedback) error {
	result := r.db.Create(qrFeedback)
	return result.Error
}

// FindByID finds a QR feedback by ID
func (r *GormQRFeedbackRepository) FindByID(ctx context.Context, id uint) (*entity.QRFeedback, error) {
	var qrFeedback entity.QRFeedback
	result := r.db.First(&qrFeedback, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &qrFeedback, nil
}

// FindByQRCode finds a QR feedback by QR code
func (r *GormQRFeedbackRepository) FindByQRCode(ctx context.Context, qrCode string) (*entity.QRFeedback, error) {
	var qrFeedback entity.QRFeedback
	result := r.db.Where("qr_code = ?", qrCode).First(&qrFeedback)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &qrFeedback, nil
}

// ListAll lists all QR feedbacks
func (r *GormQRFeedbackRepository) ListAll(ctx context.Context, page, limit int) ([]entity.QRFeedback, int64, error) {
	var qrFeedbacks []entity.QRFeedback
	var total int64
	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.QRFeedback{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get QR feedbacks
	if err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&qrFeedbacks).Error; err != nil {
		return nil, 0, err
	}

	return qrFeedbacks, total, nil
}

// Delete deletes a QR feedback
func (r *GormQRFeedbackRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.Delete(&entity.QRFeedback{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("qr feedback not found")
	}
	return nil
}
