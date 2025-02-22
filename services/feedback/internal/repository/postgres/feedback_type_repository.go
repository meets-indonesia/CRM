package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"gorm.io/gorm"
)

type feedbackTypeRepository struct {
	db *gorm.DB
}

func NewFeedbackTypeRepository(db *gorm.DB) *feedbackTypeRepository {
	return &feedbackTypeRepository{db: db}
}

func (r *feedbackTypeRepository) List(ctx context.Context) ([]model.FeedbackType, error) {
	var types []model.FeedbackType
	err := r.db.WithContext(ctx).Find(&types).Error
	return types, err
}

func (r *feedbackTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.FeedbackType, error) {
	var feedbackType model.FeedbackType
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&feedbackType).Error
	if err != nil {
		return nil, err
	}
	return &feedbackType, nil
}

func (r *feedbackTypeRepository) Create(ctx context.Context, feedbackType *model.FeedbackType) error {
	return r.db.WithContext(ctx).Create(feedbackType).Error
}
