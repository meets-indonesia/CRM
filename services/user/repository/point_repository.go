package repository

import (
	"context"
	"errors"
	"time"

	"github.com/kevinnaserwan/crm-be/services/user/domain/entity"
	"gorm.io/gorm"
)

// GormPointRepository implements PointRepository with GORM
type GormPointRepository struct {
	db *gorm.DB
}

// NewGormPointRepository creates a new GormPointRepository
func NewGormPointRepository(db *gorm.DB) *GormPointRepository {
	return &GormPointRepository{
		db: db,
	}
}

// GetPointBalance gets the point balance for a user
func (r *GormPointRepository) GetPointBalance(ctx context.Context, userID uint) (int, error) {
	var total int

	// Sum all point transactions for the user
	err := r.db.Model(&entity.PointTransaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ?", userID).
		Scan(&total).Error

	return total, err
}

// AddPoints adds points to a user's balance
func (r *GormPointRepository) AddPoints(ctx context.Context, transaction *entity.PointTransaction) error {
	// Begin transaction
	tx := r.db.Begin()

	// Create point transaction
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// GetPointTransactions gets the point transactions for a user
func (r *GormPointRepository) GetPointTransactions(ctx context.Context, userID uint, page, limit int) ([]entity.PointTransaction, int64, error) {
	var transactions []entity.PointTransaction
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.PointTransaction{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get transactions
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetDailyPoints gets the points earned by a user on a specific date
func (r *GormPointRepository) GetDailyPoints(ctx context.Context, userID uint, date time.Time) (int, error) {
	var total int

	// Get the start and end of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Sum all point transactions for the user on the specified date
	err := r.db.Model(&entity.PointTransaction{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("user_id = ? AND created_at >= ? AND created_at < ? AND amount > 0 AND type = 'feedback'", userID, startOfDay, endOfDay).
		Scan(&total).Error

	return total, err
}

// GetPointLevels gets all point levels
func (r *GormPointRepository) GetPointLevels(ctx context.Context) ([]entity.PointLevel, error) {
	var levels []entity.PointLevel

	err := r.db.Order("min_points ASC").Find(&levels).Error

	return levels, err
}

// GetUserLevel gets the level for a user based on their points
func (r *GormPointRepository) GetUserLevel(ctx context.Context, points int) (*entity.PointLevel, error) {
	var level entity.PointLevel

	// Get the highest level that the user qualifies for
	err := r.db.Where("min_points <= ?", points).
		Order("min_points DESC").
		First(&level).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no level found, get the lowest level
			err = r.db.Order("min_points ASC").First(&level).Error
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &level, nil
}

// GetNextLevel gets the next level after the current level
func (r *GormPointRepository) GetNextLevel(ctx context.Context, currentLevel *entity.PointLevel) (*entity.PointLevel, error) {
	var nextLevel entity.PointLevel

	// Get the level with the next highest min_points
	err := r.db.Where("min_points > ?", currentLevel.MinPoints).
		Order("min_points ASC").
		First(&nextLevel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No next level
			return nil, nil
		}
		return nil, err
	}

	return &nextLevel, nil
}
