package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/notification/domain/entity"
	"gorm.io/gorm"
)

// GormNotificationRepository implements NotificationRepository with GORM
type GormNotificationRepository struct {
	db *gorm.DB
}

// NewGormNotificationRepository creates a new GormNotificationRepository
func NewGormNotificationRepository(db *gorm.DB) *GormNotificationRepository {
	return &GormNotificationRepository{
		db: db,
	}
}

// Create creates a new notification
func (r *GormNotificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	result := r.db.Create(notification)
	return result.Error
}

// FindByID finds a notification by ID
func (r *GormNotificationRepository) FindByID(ctx context.Context, id uint) (*entity.Notification, error) {
	var notification entity.Notification
	result := r.db.First(&notification, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &notification, nil
}

// Update updates a notification
func (r *GormNotificationRepository) Update(ctx context.Context, notification *entity.Notification) error {
	result := r.db.Save(notification)
	return result.Error
}

// ListByUserID lists notifications by user ID
func (r *GormNotificationRepository) ListByUserID(ctx context.Context, userID uint, page, limit int) ([]entity.Notification, int64, error) {
	var notifications []entity.Notification
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.Notification{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get notifications
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// ListPending lists pending notifications
func (r *GormNotificationRepository) ListPending(ctx context.Context, limit int) ([]entity.Notification, error) {
	var notifications []entity.Notification

	if err := r.db.Where("status = ?", entity.StatusPending).Limit(limit).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}
