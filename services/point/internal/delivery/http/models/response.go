package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPointResponse struct {
	UserID        uuid.UUID `json:"user_id"`
	TotalPoints   int       `json:"total_points"`
	Level         string    `json:"level"`
	NextResetDate time.Time `json:"next_reset_date"`
}

type PointHistoryResponse struct {
	ID           uuid.UUID `json:"id"`
	PointsEarned int       `json:"points_earned"`
	DateEarned   time.Time `json:"date_earned"`
	FeedbackID   uuid.UUID `json:"feedback_id"`
}
