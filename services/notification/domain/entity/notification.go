package entity

import (
	"time"
)

// NotificationType mendefinisikan tipe notifikasi
type NotificationType string

const (
	TypeEmail            NotificationType = "EMAIL"
	TypePushNotification NotificationType = "PUSH"
)

// NotificationStatus mendefinisikan status notifikasi
type NotificationStatus string

const (
	StatusPending NotificationStatus = "PENDING"
	StatusSent    NotificationStatus = "SENT"
	StatusFailed  NotificationStatus = "FAILED"
)

// Notification adalah entitas untuk notifikasi
type Notification struct {
	ID        uint               `json:"id" gorm:"primaryKey"`
	UserID    uint               `json:"user_id" gorm:"index;not null"`
	Type      NotificationType   `json:"type" gorm:"not null"`
	Title     string             `json:"title" gorm:"not null"`
	Content   string             `json:"content" gorm:"not null"`
	Data      string             `json:"data,omitempty"` // JSON string untuk data tambahan
	Status    NotificationStatus `json:"status" gorm:"not null;default:PENDING"`
	SentAt    *time.Time         `json:"sent_at,omitempty"`
	Error     string             `json:"error,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// CreateEmailRequest adalah model untuk permintaan pembuatan notifikasi email
type CreateEmailRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	EmailTo string `json:"email_to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
	IsHTML  bool   `json:"is_html"`
}

// CreatePushNotificationRequest adalah model untuk permintaan pembuatan push notification
type CreatePushNotificationRequest struct {
	UserID  uint                   `json:"user_id" binding:"required"`
	Title   string                 `json:"title" binding:"required"`
	Message string                 `json:"message" binding:"required"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// NotificationListResponse adalah model respons untuk daftar notifikasi
type NotificationListResponse struct {
	Notifications []Notification `json:"notifications"`
	Total         int64          `json:"total"`
	Page          int            `json:"page"`
	Limit         int            `json:"limit"`
}
