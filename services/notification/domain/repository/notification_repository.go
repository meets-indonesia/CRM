package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/notification/domain/entity"
)

// NotificationRepository mendefinisikan operasi-operasi repository untuk Notification
type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	FindByID(ctx context.Context, id uint) (*entity.Notification, error)
	Update(ctx context.Context, notification *entity.Notification) error
	ListByUserID(ctx context.Context, userID uint, page, limit int) ([]entity.Notification, int64, error)
	ListPending(ctx context.Context, limit int) ([]entity.Notification, error)
}

// EmailSender mendefinisikan operasi-operasi untuk mengirim email
type EmailSender interface {
	SendEmail(to, subject, body string, isHTML bool) error
}

// PushNotificationSender mendefinisikan operasi-operasi untuk mengirim push notification
type PushNotificationSender interface {
	SendPushNotification(userID uint, title, message string, data map[string]interface{}) error
}

// EventSubscriber mendefinisikan operasi-operasi untuk berlangganan event
type EventSubscriber interface {
	SubscribeToEvents() error
	Close() error
}
