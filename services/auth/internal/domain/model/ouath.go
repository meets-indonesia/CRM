// services/auth/internal/domain/model/oauth.go
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OAuthProvider represents supported OAuth providers
type OAuthProvider string

const (
	GoogleProvider OAuthProvider = "google"
	AppleProvider  OAuthProvider = "apple"
)

// OAuthAccount represents a user's OAuth connection
type OAuthAccount struct {
	gorm.Model
	ID           uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID       uuid.UUID     `gorm:"type:uuid;not null"`
	User         User          `gorm:"foreignKey:UserID;references:ID"`
	Provider     OAuthProvider `gorm:"type:varchar(20);not null"`
	ProviderID   string        `gorm:"type:varchar(255);not null"`
	Email        string        `gorm:"type:varchar(255);not null"`
	AccessToken  string        `gorm:"type:text"`
	RefreshToken string        `gorm:"type:text"`
	ExpiresAt    time.Time
}
