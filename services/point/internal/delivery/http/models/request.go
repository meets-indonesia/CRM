package models

import (
	"github.com/google/uuid"
)

type Meta struct {
	Action string `json:"action" binding:"required"`
}

type BaseRequest struct {
	Meta Meta        `json:"meta" binding:"required"`
	Data interface{} `json:"data" binding:"required"`
}

type ProcessFeedbackData struct {
	FeedbackID uuid.UUID `json:"feedback_id" binding:"required"`
}

type GetUserPointsData struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
