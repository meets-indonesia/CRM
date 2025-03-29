package entity

import (
	"time"
)

// Category mendefinisikan kategori feedback
type Category string

const (
	CategorySuggestionCriticism Category = "SARAN_KRITIK"
	CategoryTicketPayment       Category = "PEMBAYARAN_TIKET"
	CategoryFacilityIssue       Category = "MASALAH_FASILITAS"
	CategoryServiceComplaint    Category = "KELUHAN_PELAYANAN"
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
	Category  Category  `json:"category" gorm:"not null"`
	Station   string    `json:"station" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	Rating    uint      `json:"rating" gorm:"not null"`
	ImagePath string    `json:"image_path,omitempty" gorm:"default:null"`
	Status    Status    `json:"status" gorm:"not null;default:PENDING"`
	Response  string    `json:"response,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateFeedbackRequest adalah model untuk permintaan pembuatan feedback
type CreateFeedbackRequest struct {
	Category Category `json:"category" form:"category" binding:"required"`
	Station  string   `json:"station" form:"station" binding:"required"`
	Title    string   `json:"title" form:"title" binding:"required"`
	Content  string   `json:"content" form:"content" binding:"required"`
	Rating   uint     `json:"rating" form:"rating" binding:"required,min=1,max=5"`
	// Image akan dihandle melalui form-data
}

type QRFeedback struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Station   string    `json:"station" gorm:"not null"`
	QRCode    string    `json:"qr_code" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateQRFeedbackRequest struct {
	Station string `json:"station" form:"station" binding:"required"`
}

type RespondQRFeedbackRequest struct {
	Response string `json:"response" binding:"required"`
}

// QRFeedbackResponse is the response model for QR feedback operations
type QRFeedbackResponse struct {
	ID        uint      `json:"id"`
	QRCode    string    `json:"qr_code"`
	Station   string    `json:"station"`
	CreatedAt time.Time `json:"created_at"`
}

// In your entity package
type QRFeedbackListResponse struct {
	QRFeedbacks []QRFeedbackResponse `json:"qr_feedbacks"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
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
