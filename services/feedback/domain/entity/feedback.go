package entity

import (
	"time"
)

// Status mendefinisikan status feedback
type Status string

const (
	StatusPending   Status = "PENDING"
	StatusResponded Status = "RESPONDED"
)

// Feedback adalah entitas untuk feedback dari customer
type Feedback struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	Status    Status    `json:"status" gorm:"not null;default:PENDING"`
	Response  string    `json:"response,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateFeedbackRequest adalah model untuk permintaan pembuatan feedback
type CreateFeedbackRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// RespondFeedbackRequest adalah model untuk permintaan respons feedback
type RespondFeedbackRequest struct {
	Response string `json:"response" binding:"required"`
}

// FeedbackListResponse adalah model respons untuk daftar feedback
type FeedbackListResponse struct {
	Feedbacks []Feedback `json:"feedbacks"`
	Total     int64      `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}
