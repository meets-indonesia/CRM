package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Point rules constants
const (
	// Point thresholds
	SilverThreshold   = 40
	GoldThreshold     = 100
	PlatinumThreshold = 200

	// Point rules
	PointsPerFeedback = 1
	ResetMonths       = 12
)

// UserPoints stores the user's point balance and reset dates
type UserPoints struct {
	gorm.Model
	ID                  uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID              uuid.UUID  `gorm:"type:uuid;column:user_id;not null"`
	TotalPoints         int        `gorm:"column:total_points;default:0"`
	LastPointEarnedDate *time.Time `gorm:"column:last_point_earned_date"`
	RegistrationDate    time.Time  `gorm:"column:registration_date;not null"`
	LastResetDate       time.Time  `gorm:"column:last_reset_date;not null"`
	NextResetDate       time.Time  `gorm:"column:next_reset_date;not null"`
}

// PointHistory tracks each point transaction
type PointHistory struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID       uuid.UUID `gorm:"type:uuid;column:user_id;not null"`
	PointsEarned int       `gorm:"column:points_earned;not null"`
	DateEarned   time.Time `gorm:"column:date_earned;not null"`
	FeedbackID   uuid.UUID `gorm:"type:uuid;column:feedback_id;not null"`
}

// PointResetHistory tracks point reset events
type PointResetHistory struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID            uuid.UUID `gorm:"type:uuid;column:user_id;not null"`
	ResetDate         time.Time `gorm:"column:reset_date;not null"`
	PointsBeforeReset int       `gorm:"column:points_before_reset;not null"`
	NextResetDate     time.Time `gorm:"column:next_reset_date;not null"`
}
