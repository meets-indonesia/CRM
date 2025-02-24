package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/user/internal/domain/model"
)

type UserPointsRepository interface {
	Create(ctx context.Context, points *model.UserPoints) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserPoints, error)
	List(ctx context.Context) ([]model.UserPoints, error)
	UpdatePoints(ctx context.Context, userID uuid.UUID, points int) error
	UpdateLastPointEarned(ctx context.Context, userID uuid.UUID, date time.Time) error
	UpdateResetDates(ctx context.Context, userID uuid.UUID, lastReset, nextReset time.Time) error
}

type PointHistoryRepository interface {
	Create(ctx context.Context, history *model.PointHistory) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.PointHistory, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]model.PointHistory, error)
}

type PointResetHistoryRepository interface {
	Create(ctx context.Context, resetHistory *model.PointResetHistory) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.PointResetHistory, error)
	GetLastReset(ctx context.Context, userID uuid.UUID) (*model.PointResetHistory, error)
}
