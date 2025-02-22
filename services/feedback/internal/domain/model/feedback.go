package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FeedbackType represents the category of feedback
type FeedbackType struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(255);column:name;not null"`
}

// Station represents a location
type Station struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(255);column:name;not null"`
}

// Status constants
const (
	FeedbackStatusPending = 1
	FeedbackStatusSolved  = 2
)

// Feedback represents the main feedback entity
type Feedback struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID         uuid.UUID `gorm:"type:uuid;column:user_id;not null"`
	FeedbackDate   time.Time `gorm:"type:timestamp;column:feedback_date;not null"`
	FeedbackTypeID uuid.UUID `gorm:"type:uuid;column:feedback_type_id;not null"`
	StationID      uuid.UUID `gorm:"type:uuid;column:station_id;not null"`
	Feedback       string    `gorm:"type:text;column:feedback;not null"`
	Documentation  string    `gorm:"type:text;column:documentation;not null"`
	Rating         float64   `gorm:"type:float;column:rating;not null"`
	Status         int       `gorm:"type:int;column:status;not null;default:1"`

	// Relationships
	FeedbackType FeedbackType      `gorm:"foreignKey:FeedbackTypeID;references:ID"`
	Station      Station           `gorm:"foreignKey:StationID;references:ID"`
	Response     *FeedbackResponse `gorm:"foreignKey:FeedbackID"`
}

// FeedbackResponse represents admin's response to feedback
type FeedbackResponse struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	FeedbackID   uuid.UUID `gorm:"type:uuid;column:feedback_id;not null"`
	Response     string    `gorm:"type:text;column:response;not null"`
	ResponseDate time.Time `gorm:"type:timestamp;column:response_date;not null"`

	// Relationship back to feedback
	Feedback Feedback `gorm:"foreignKey:FeedbackID;references:ID"`
}
