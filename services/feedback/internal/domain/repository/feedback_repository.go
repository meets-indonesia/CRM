package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
)

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *model.Feedback) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Feedback, error)
	List(ctx context.Context) ([]model.Feedback, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Feedback, error)
	Update(ctx context.Context, feedback *model.Feedback) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status int) error
}

type FeedbackResponseRepository interface {
	Create(ctx context.Context, response *model.FeedbackResponse) error
	GetByFeedbackID(ctx context.Context, feedbackID uuid.UUID) (*model.FeedbackResponse, error)
	Update(ctx context.Context, response *model.FeedbackResponse) error
}

type FeedbackTypeRepository interface {
	Create(ctx context.Context, feedbackType *model.FeedbackType) error
	List(ctx context.Context) ([]model.FeedbackType, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.FeedbackType, error)
}

type StationRepository interface {
	Create(ctx context.Context, station *model.Station) error
	List(ctx context.Context) ([]model.Station, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Station, error)
}
