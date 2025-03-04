package repository

import (
	"context"
	"time"

	"github.com/kevinnaserwan/crm-be/services/user/domain/entity"
)

// PointRepository mendefinisikan operasi-operasi repository untuk Point
type PointRepository interface {
	// Point Transaction
	GetPointBalance(ctx context.Context, userID uint) (int, error)
	AddPoints(ctx context.Context, transaction *entity.PointTransaction) error
	GetPointTransactions(ctx context.Context, userID uint, page, limit int) ([]entity.PointTransaction, int64, error)
	GetDailyPoints(ctx context.Context, userID uint, date time.Time) (int, error)

	// Point Level
	GetPointLevels(ctx context.Context) ([]entity.PointLevel, error)
	GetUserLevel(ctx context.Context, points int) (*entity.PointLevel, error)
	GetNextLevel(ctx context.Context, currentLevel *entity.PointLevel) (*entity.PointLevel, error)
}

// PointEventProcessor mendefinisikan operasi-operasi untuk memproses event point
type PointEventProcessor interface {
	ProcessFeedbackCreated(userID uint, feedbackID uint) error
	ProcessRewardClaimed(userID uint, rewardID uint, points int) error
}
