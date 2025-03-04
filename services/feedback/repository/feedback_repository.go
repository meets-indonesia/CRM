package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/feedback/domain/entity"
	"gorm.io/gorm"
)

// GormFeedbackRepository implements FeedbackRepository with GORM
type GormFeedbackRepository struct {
	db *gorm.DB
}

// NewGormFeedbackRepository creates a new GormFeedbackRepository
func NewGormFeedbackRepository(db *gorm.DB) *GormFeedbackRepository {
	return &GormFeedbackRepository{
		db: db,
	}
}

// Create creates a new feedback
func (r *GormFeedbackRepository) Create(ctx context.Context, feedback *entity.Feedback) error {
	result := r.db.Create(feedback)
	return result.Error
}

// FindByID finds a feedback by ID
func (r *GormFeedbackRepository) FindByID(ctx context.Context, id uint) (*entity.Feedback, error) {
	var feedback entity.Feedback
	result := r.db.First(&feedback, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &feedback, nil
}

// Update updates a feedback
func (r *GormFeedbackRepository) Update(ctx context.Context, feedback *entity.Feedback) error {
	result := r.db.Save(feedback)
	return result.Error
}

// ListAll lists all feedbacks
func (r *GormFeedbackRepository) ListAll(ctx context.Context, page, limit int) ([]entity.Feedback, int64, error) {
	var feedbacks []entity.Feedback
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.Feedback{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get feedbacks
	if err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// ListByUserID lists feedbacks by user ID
func (r *GormFeedbackRepository) ListByUserID(ctx context.Context, userID uint, page, limit int) ([]entity.Feedback, int64, error) {
	var feedbacks []entity.Feedback
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.Feedback{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get feedbacks
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// ListByStatus lists feedbacks by status
func (r *GormFeedbackRepository) ListByStatus(ctx context.Context, status entity.Status, page, limit int) ([]entity.Feedback, int64, error) {
	var feedbacks []entity.Feedback
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.Feedback{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get feedbacks
	if err := r.db.Where("status = ?", status).Order("created_at DESC").Offset(offset).Limit(limit).Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}
