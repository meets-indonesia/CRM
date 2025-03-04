package entity

import (
	"time"
)

// Status mendefinisikan status klaim
type ClaimStatus string

const (
	ClaimStatusPending   ClaimStatus = "PENDING"
	ClaimStatusApproved  ClaimStatus = "APPROVED"
	ClaimStatusRejected  ClaimStatus = "REJECTED"
	ClaimStatusCancelled ClaimStatus = "CANCELLED"
)

// RewardClaim adalah entitas untuk klaim hadiah
type RewardClaim struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	UserID    uint        `json:"user_id" gorm:"index;not null"`
	RewardID  uint        `json:"reward_id" gorm:"index;not null"`
	PointCost int         `json:"point_cost" gorm:"not null"`
	Status    ClaimStatus `json:"status" gorm:"not null;default:PENDING"`
	Notes     string      `json:"notes,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`

	// Eager loading untuk reward
	Reward Reward `json:"reward" gorm:"foreignKey:RewardID"`
}

// ClaimRewardRequest adalah model untuk permintaan klaim hadiah
type ClaimRewardRequest struct {
	RewardID uint `json:"reward_id" binding:"required"`
}

// UpdateClaimStatusRequest adalah model untuk permintaan update status klaim
type UpdateClaimStatusRequest struct {
	Status ClaimStatus `json:"status" binding:"required,oneof=PENDING APPROVED REJECTED CANCELLED"`
	Notes  string      `json:"notes,omitempty"`
}

// ClaimListResponse adalah model respons untuk daftar klaim
type ClaimListResponse struct {
	Claims []RewardClaim `json:"claims"`
	Total  int64         `json:"total"`
	Page   int           `json:"page"`
	Limit  int           `json:"limit"`
}
