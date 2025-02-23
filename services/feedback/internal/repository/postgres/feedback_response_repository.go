package postgres

import (
	"context"
	"errors"

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

// internal/repository/postgres/feedback_response_repository.go
func (r *feedbackResponseRepository) Create(ctx context.Context, response *model.FeedbackResponse) error {
	// Check if response already exists
	var existing model.FeedbackResponse
	err := r.db.WithContext(ctx).
		Where("feedback_id = ?", response.FeedbackID).
		First(&existing).Error

	if err == nil {
		// Response already exists, update it
		return r.Update(ctx, response)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other error occurred
		return err
	}

	// Create new response
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
