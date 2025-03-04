package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/reward/domain/entity"
)

// ClaimRepository mendefinisikan operasi-operasi repository untuk RewardClaim
type ClaimRepository interface {
	Create(ctx context.Context, claim *entity.RewardClaim) error
	FindByID(ctx context.Context, id uint) (*entity.RewardClaim, error)
	Update(ctx context.Context, claim *entity.RewardClaim) error
	ListByUserID(ctx context.Context, userID uint, page, limit int) ([]entity.RewardClaim, int64, error)
	ListAll(ctx context.Context, page, limit int) ([]entity.RewardClaim, int64, error)
	ListByStatus(ctx context.Context, status entity.ClaimStatus, page, limit int) ([]entity.RewardClaim, int64, error)
}

// EventPublisher mendefinisikan operasi-operasi untuk mempublikasikan event
type EventPublisher interface {
	PublishRewardClaimed(claim *entity.RewardClaim) error
	PublishClaimStatusUpdated(claim *entity.RewardClaim) error
	Close() error
}

// UserPointService mendefinisikan operasi-operasi untuk mengecek poin user
type UserPointService interface {
	CheckUserPoints(ctx context.Context, userID uint) (int, error)
}
