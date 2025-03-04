package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/reward/domain/entity"
)

// RewardRepository mendefinisikan operasi-operasi repository untuk Reward
type RewardRepository interface {
	Create(ctx context.Context, reward *entity.Reward) error
	FindByID(ctx context.Context, id uint) (*entity.Reward, error)
	Update(ctx context.Context, reward *entity.Reward) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, active bool, page, limit int) ([]entity.Reward, int64, error)

	// Inventory operations
	CheckStock(ctx context.Context, id uint) (int, error)
	DecreaseStock(ctx context.Context, id uint, amount int) error
	IncreaseStock(ctx context.Context, id uint, amount int) error
}
