package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/feedback/domain/entity"
)

// FeedbackRepository mendefinisikan operasi-operasi repository untuk Feedback
type FeedbackRepository interface {
	Create(ctx context.Context, feedback *entity.Feedback) error
	FindByID(ctx context.Context, id uint) (*entity.Feedback, error)
	Update(ctx context.Context, feedback *entity.Feedback) error
	ListAll(ctx context.Context, page, limit int) ([]entity.Feedback, int64, error)
	ListByUserID(ctx context.Context, userID uint, page, limit int) ([]entity.Feedback, int64, error)
	ListByStatus(ctx context.Context, status entity.Status, page, limit int) ([]entity.Feedback, int64, error)
}

// QRFeedbackRepository defines operations for QRFeedback
type QRFeedbackRepository interface {
	Create(ctx context.Context, qrFeedback *entity.QRFeedback) error
	FindByID(ctx context.Context, id uint) (*entity.QRFeedback, error)
	FindByQRCode(ctx context.Context, qrCode string) (*entity.QRFeedback, error)
	ListAll(ctx context.Context, page, limit int) ([]entity.QRFeedback, int64, error)
	Delete(ctx context.Context, id uint) error
}

// EventPublisher mendefinisikan operasi-operasi untuk mempublikasikan event
type EventPublisher interface {
	PublishFeedbackCreated(feedback *entity.Feedback) error
	PublishFeedbackResponded(feedback *entity.Feedback) error
	Close() error
}
