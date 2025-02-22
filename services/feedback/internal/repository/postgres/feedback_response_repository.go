package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"gorm.io/gorm"
)

type feedbackResponseRepository struct {
	db *gorm.DB
}

func NewFeedbackResponseRepository(db *gorm.DB) *feedbackResponseRepository {
	return &feedbackResponseRepository{db: db}
}

func (r *feedbackResponseRepository) Create(ctx context.Context, response *model.FeedbackResponse) error {
	return r.db.WithContext(ctx).Create(response).Error
}

func (r *feedbackResponseRepository) GetByFeedbackID(ctx context.Context, feedbackID uuid.UUID) (*model.FeedbackResponse, error) {
	var response model.FeedbackResponse
	err := r.db.WithContext(ctx).Where("feedback_id = ?", feedbackID).First(&response).Error
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *feedbackResponseRepository) Update(ctx context.Context, response *model.FeedbackResponse) error {
	return r.db.WithContext(ctx).Save(response).Error
}
