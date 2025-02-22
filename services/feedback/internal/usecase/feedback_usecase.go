package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/repository"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/infrastructure/messagebroker"
)

type FeedbackUseCase struct {
	feedbackRepo         repository.FeedbackRepository
	feedbackResponseRepo repository.FeedbackResponseRepository
	feedbackTypeRepo     repository.FeedbackTypeRepository
	stationRepo          repository.StationRepository
	rabbitMQ             *messagebroker.RabbitMQ
}

func NewFeedbackUseCase(
	feedbackRepo repository.FeedbackRepository,
	feedbackResponseRepo repository.FeedbackResponseRepository,
	feedbackTypeRepo repository.FeedbackTypeRepository,
	stationRepo repository.StationRepository,
	rabbitMQ *messagebroker.RabbitMQ,
) *FeedbackUseCase {
	return &FeedbackUseCase{
		feedbackRepo:         feedbackRepo,
		feedbackResponseRepo: feedbackResponseRepo,
		feedbackTypeRepo:     feedbackTypeRepo,
		stationRepo:          stationRepo,
		rabbitMQ:             rabbitMQ,
	}
}

// CreateFeedback creates a new feedback
func (uc *FeedbackUseCase) CreateFeedback(ctx context.Context, feedback *model.Feedback) error {
	// Validate feedback type exists
	if _, err := uc.feedbackTypeRepo.GetByID(ctx, feedback.FeedbackTypeID); err != nil {
		return errors.New("invalid feedback type")
	}

	// Validate station exists
	if _, err := uc.stationRepo.GetByID(ctx, feedback.StationID); err != nil {
		return errors.New("invalid station")
	}

	feedback.Status = model.FeedbackStatusPending
	feedback.FeedbackDate = time.Now()

	// Save feedback
	if err := uc.feedbackRepo.Create(ctx, feedback); err != nil {
		return err
	}

	// Publish event to RabbitMQ
	return uc.rabbitMQ.PublishFeedbackCreated(feedback)
}

// RespondToFeedback allows admin to respond to a feedback
func (uc *FeedbackUseCase) RespondToFeedback(ctx context.Context, feedbackID uuid.UUID, responseText string) error {
	feedback, err := uc.feedbackRepo.GetByID(ctx, feedbackID)
	if err != nil {
		return errors.New("feedback not found")
	}

	if feedback.Status == model.FeedbackStatusSolved {
		return errors.New("feedback is already solved")
	}

	response := &model.FeedbackResponse{
		ID:           uuid.New(),
		FeedbackID:   feedbackID,
		Response:     responseText,
		ResponseDate: time.Now(),
	}

	if err := uc.feedbackResponseRepo.Create(ctx, response); err != nil {
		return err
	}

	// Update feedback status
	if err := uc.feedbackRepo.UpdateStatus(ctx, feedbackID, model.FeedbackStatusSolved); err != nil {
		return err
	}

	// Publish event to RabbitMQ
	return uc.rabbitMQ.PublishFeedbackResponded(response)
}

// GetUserFeedbacks gets all feedbacks for a user
func (uc *FeedbackUseCase) GetUserFeedbacks(ctx context.Context, userID uuid.UUID) ([]model.Feedback, error) {
	return uc.feedbackRepo.ListByUserID(ctx, userID)
}

// GetAllFeedbacks gets all feedbacks (for admin)
func (uc *FeedbackUseCase) GetAllFeedbacks(ctx context.Context) ([]model.Feedback, error) {
	return uc.feedbackRepo.List(ctx)
}

// GetFeedbackTypes gets all feedback types
func (uc *FeedbackUseCase) GetFeedbackTypes(ctx context.Context) ([]model.FeedbackType, error) {
	return uc.feedbackTypeRepo.List(ctx)
}

// GetStations gets all stations
func (uc *FeedbackUseCase) GetStations(ctx context.Context) ([]model.Station, error) {
	return uc.stationRepo.List(ctx)
}

// CreateFeedbackType creates a new feedback type (admin only)
func (uc *FeedbackUseCase) CreateFeedbackType(ctx context.Context, feedbackType *model.FeedbackType) error {
	return uc.feedbackTypeRepo.Create(ctx, feedbackType)
}

// CreateStation creates a new station (admin only)
func (uc *FeedbackUseCase) CreateStation(ctx context.Context, station *model.Station) error {
	return uc.stationRepo.Create(ctx, station)
}
