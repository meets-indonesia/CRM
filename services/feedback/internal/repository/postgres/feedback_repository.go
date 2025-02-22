package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"gorm.io/gorm"
)

type feedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) *feedbackRepository {
	return &feedbackRepository{db: db}
}

func (r *feedbackRepository) Create(ctx context.Context, feedback *model.Feedback) error {
	return r.db.WithContext(ctx).Create(feedback).Error
}

func (r *feedbackRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Feedback, error) {
	var feedback model.Feedback
	err := r.db.WithContext(ctx).
		Preload("FeedbackType").
		Preload("Station").
		Preload("Response").
		Where("id = ?", id).
		First(&feedback).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *feedbackRepository) List(ctx context.Context) ([]model.Feedback, error) {
	var feedbacks []model.Feedback
	err := r.db.WithContext(ctx).
		Preload("FeedbackType").
		Preload("Station").
		Preload("Response").
		Find(&feedbacks).Error
	return feedbacks, err
}

func (r *feedbackRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Feedback, error) {
	var feedbacks []model.Feedback
	err := r.db.WithContext(ctx).
		Preload("FeedbackType").
		Preload("Station").
		Preload("Response").
		Where("user_id = ?", userID).
		Find(&feedbacks).Error
	return feedbacks, err
}

func (r *feedbackRepository) Update(ctx context.Context, feedback *model.Feedback) error {
	return r.db.WithContext(ctx).Save(feedback).Error
}

func (r *feedbackRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status int) error {
	return r.db.WithContext(ctx).
		Model(&model.Feedback{}).
		Where("id = ?", id).
		Update("status", status).Error
}
