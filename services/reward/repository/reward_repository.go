package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/reward/domain/entity"
	"gorm.io/gorm"
)

// GormRewardRepository implements RewardRepository with GORM
type GormRewardRepository struct {
	db *gorm.DB
}

// NewGormRewardRepository creates a new GormRewardRepository
func NewGormRewardRepository(db *gorm.DB) *GormRewardRepository {
	return &GormRewardRepository{
		db: db,
	}
}

// Create creates a new reward
func (r *GormRewardRepository) Create(ctx context.Context, reward *entity.Reward) error {
	result := r.db.Create(reward)
	return result.Error
}

// FindByID finds a reward by ID
func (r *GormRewardRepository) FindByID(ctx context.Context, id uint) (*entity.Reward, error) {
	var reward entity.Reward
	result := r.db.First(&reward, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &reward, nil
}

// Update updates a reward
func (r *GormRewardRepository) Update(ctx context.Context, reward *entity.Reward) error {
	result := r.db.Save(reward)
	return result.Error
}

// Delete deletes a reward
func (r *GormRewardRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.Delete(&entity.Reward{}, id)
	return result.Error
}

// List lists rewards
func (r *GormRewardRepository) List(ctx context.Context, active bool, page, limit int) ([]entity.Reward, int64, error) {
	var rewards []entity.Reward
	var total int64

	offset := (page - 1) * limit
	query := r.db

	if active {
		query = query.Where("is_active = ?", true)
	}

	// Count total
	if err := query.Model(&entity.Reward{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get rewards
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&rewards).Error; err != nil {
		return nil, 0, err
	}

	return rewards, total, nil
}

// CheckStock checks the stock of a reward
func (r *GormRewardRepository) CheckStock(ctx context.Context, id uint) (int, error) {
	var reward entity.Reward

	if err := r.db.Select("stock").First(&reward, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}

	return reward.Stock, nil
}

// DecreaseStock decreases the stock of a reward
func (r *GormRewardRepository) DecreaseStock(ctx context.Context, id uint, amount int) error {
	// Begin transaction
	tx := r.db.Begin()

	// Get current stock with lock
	var reward entity.Reward
	if err := tx.First(&reward, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Check if stock is sufficient
	if reward.Stock < amount {
		tx.Rollback()
		return errors.New("insufficient stock")
	}

	// Update stock
	if err := tx.Model(&entity.Reward{}).Where("id = ?", id).Update("stock", gorm.Expr("stock - ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// IncreaseStock increases the stock of a reward
func (r *GormRewardRepository) IncreaseStock(ctx context.Context, id uint, amount int) error {
	result := r.db.Model(&entity.Reward{}).Where("id = ?", id).Update("stock", gorm.Expr("stock + ?", amount))
	return result.Error
}
