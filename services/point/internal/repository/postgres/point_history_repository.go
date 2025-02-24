package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/user/internal/domain/model"
	"gorm.io/gorm"
)

type pointHistoryRepository struct {
	db *gorm.DB
}

func NewPointHistoryRepository(db *gorm.DB) *pointHistoryRepository {
	return &pointHistoryRepository{db: db}
}

func (r *pointHistoryRepository) Create(ctx context.Context, history *model.PointHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *pointHistoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.PointHistory, error) {
	var histories []model.PointHistory
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *pointHistoryRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]model.PointHistory, error) {
	var histories []model.PointHistory
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND date_earned BETWEEN ? AND ?", userID, startDate, endDate).
		Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}
