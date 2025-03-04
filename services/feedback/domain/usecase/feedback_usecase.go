package usecase

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/feedback/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/feedback/domain/repository"
)

// Errors
var (
	ErrFeedbackNotFound = errors.New("feedback not found")
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrAlreadyResponded = errors.New("feedback already responded to")
)

// FeedbackUsecase mendefinisikan operasi-operasi usecase untuk Feedback
type FeedbackUsecase interface {
	CreateFeedback(ctx context.Context, userID uint, req entity.CreateFeedbackRequest) (*entity.Feedback, error)
	GetFeedback(ctx context.Context, id uint) (*entity.Feedback, error)
	RespondFeedback(ctx context.Context, id uint, req entity.RespondFeedbackRequest) (*entity.Feedback, error)
	ListAllFeedback(ctx context.Context, page, limit int) (*entity.FeedbackListResponse, error)
	ListUserFeedback(ctx context.Context, userID uint, page, limit int) (*entity.FeedbackListResponse, error)
	ListPendingFeedback(ctx context.Context, page, limit int) (*entity.FeedbackListResponse, error)
}

type feedbackUsecase struct {
	feedbackRepo   repository.FeedbackRepository
	eventPublisher repository.EventPublisher
}

// NewFeedbackUsecase membuat instance baru FeedbackUsecase
func NewFeedbackUsecase(feedbackRepo repository.FeedbackRepository, eventPublisher repository.EventPublisher) FeedbackUsecase {
	return &feedbackUsecase{
		feedbackRepo:   feedbackRepo,
		eventPublisher: eventPublisher,
	}
}

// CreateFeedback membuat feedback baru
func (u *feedbackUsecase) CreateFeedback(ctx context.Context, userID uint, req entity.CreateFeedbackRequest) (*entity.Feedback, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	feedback := &entity.Feedback{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Status:  entity.StatusPending,
	}

	if err := u.feedbackRepo.Create(ctx, feedback); err != nil {
		return nil, err
	}

	// Publish event
	if err := u.eventPublisher.PublishFeedbackCreated(feedback); err != nil {
		// Log error but don't fail
	}

	return feedback, nil
}

// GetFeedback mendapatkan feedback berdasarkan ID
func (u *feedbackUsecase) GetFeedback(ctx context.Context, id uint) (*entity.Feedback, error) {
	feedback, err := u.feedbackRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if feedback == nil {
		return nil, ErrFeedbackNotFound
	}

	return feedback, nil
}

// RespondFeedback memberikan respons terhadap feedback
func (u *feedbackUsecase) RespondFeedback(ctx context.Context, id uint, req entity.RespondFeedbackRequest) (*entity.Feedback, error) {
	feedback, err := u.feedbackRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if feedback == nil {
		return nil, ErrFeedbackNotFound
	}

	if feedback.Status == entity.StatusResponded {
		return nil, ErrAlreadyResponded
	}

	feedback.Response = req.Response
	feedback.Status = entity.StatusResponded

	if err := u.feedbackRepo.Update(ctx, feedback); err != nil {
		return nil, err
	}

	// Publish event
	if err := u.eventPublisher.PublishFeedbackResponded(feedback); err != nil {
		// Log error but don't fail
	}

	return feedback, nil
}

// ListAllFeedback mendapatkan semua feedback
func (u *feedbackUsecase) ListAllFeedback(ctx context.Context, page, limit int) (*entity.FeedbackListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	feedbacks, total, err := u.feedbackRepo.ListAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.FeedbackListResponse{
		Feedbacks: feedbacks,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// ListUserFeedback mendapatkan feedback dari user tertentu
func (u *feedbackUsecase) ListUserFeedback(ctx context.Context, userID uint, page, limit int) (*entity.FeedbackListResponse, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	feedbacks, total, err := u.feedbackRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.FeedbackListResponse{
		Feedbacks: feedbacks,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// ListPendingFeedback mendapatkan feedback yang belum direspons
func (u *feedbackUsecase) ListPendingFeedback(ctx context.Context, page, limit int) (*entity.FeedbackListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	feedbacks, total, err := u.feedbackRepo.ListByStatus(ctx, entity.StatusPending, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.FeedbackListResponse{
		Feedbacks: feedbacks,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}
