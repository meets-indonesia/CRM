package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/user/internal/domain/model"
	"gorm.io/gorm"
)

type pointResetHistoryRepository struct {
	db *gorm.DB
}

func NewPointResetHistoryRepository(db *gorm.DB) *pointResetHistoryRepository {
	return &pointResetHistoryRepository{db: db}
}

func (r *pointResetHistoryRepository) Create(ctx context.Context, resetHistory *model.PointResetHistory) error {
	return r.db.WithContext(ctx).Create(resetHistory).Error
}

func (r *pointResetHistoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.PointResetHistory, error) {
	var histories []model.PointResetHistory
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *pointResetHistoryRepository) GetLastReset(ctx context.Context, userID uuid.UUID) (*model.PointResetHistory, error) {
	var history model.PointResetHistory
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("reset_date DESC").
		First(&history).Error; err != nil {
		return nil, err
	}
	return &history, nil
}
