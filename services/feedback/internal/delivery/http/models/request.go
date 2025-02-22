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

type CreateFeedbackData struct {
	FeedbackTypeID uuid.UUID `json:"feedback_type_id" binding:"required"`
	StationID      uuid.UUID `json:"station_id" binding:"required"`
	Feedback       string    `json:"feedback" binding:"required"`
	Documentation  string    `json:"documentation" binding:"required"`
	Rating         float64   `json:"rating" binding:"required"`
}

type RespondToFeedbackData struct {
	FeedbackID uuid.UUID `json:"feedback_id" binding:"required"`
	Response   string    `json:"response" binding:"required"`
}

type CreateFeedbackTypeData struct {
	Name string `json:"name" binding:"required"`
}

type CreateStationData struct {
	Name string `json:"name" binding:"required"`
}
