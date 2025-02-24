package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/user/internal/domain/model"
	"gorm.io/gorm"
)

type userPointsRepository struct {
	db *gorm.DB
}

func NewUserPointsRepository(db *gorm.DB) *userPointsRepository {
	return &userPointsRepository{db: db}
}

func (r *userPointsRepository) Create(ctx context.Context, points *model.UserPoints) error {
	return r.db.WithContext(ctx).Create(points).Error
}

func (r *userPointsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserPoints, error) {
	var points model.UserPoints
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&points).Error; err != nil {
		return nil, err
	}
	return &points, nil
}

func (r *userPointsRepository) List(ctx context.Context) ([]model.UserPoints, error) {
	var points []model.UserPoints
	if err := r.db.WithContext(ctx).Find(&points).Error; err != nil {
		return nil, err
	}
	return points, nil
}

func (r *userPointsRepository) UpdatePoints(ctx context.Context, userID uuid.UUID, points int) error {
	return r.db.WithContext(ctx).
		Model(&model.UserPoints{}).
		Where("user_id = ?", userID).
		Update("total_points", points).
		Error
}

func (r *userPointsRepository) UpdateLastPointEarned(ctx context.Context, userID uuid.UUID, date time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.UserPoints{}).
		Where("user_id = ?", userID).
		Update("last_point_earned_date", date).
		Error
}

func (r *userPointsRepository) UpdateResetDates(ctx context.Context, userID uuid.UUID, lastReset, nextReset time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.UserPoints{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"last_reset_date": lastReset,
			"next_reset_date": nextReset,
		}).Error
}
