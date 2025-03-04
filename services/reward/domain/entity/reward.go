package entity

import (
	"time"
)

// Reward adalah entitas untuk hadiah yang dapat diklaim
type Reward struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	PointCost   int       `json:"point_cost" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"not null"`
	ImageURL    string    `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateRewardRequest adalah model untuk permintaan pembuatan reward
type CreateRewardRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	PointCost   int    `json:"point_cost" binding:"required,min=1"`
	Stock       int    `json:"stock" binding:"required,min=0"`
	ImageURL    string `json:"image_url,omitempty"`
}

// UpdateRewardRequest adalah model untuk permintaan update reward
type UpdateRewardRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PointCost   int    `json:"point_cost,omitempty" binding:"min=1"`
	Stock       int    `json:"stock,omitempty" binding:"min=0"`
	ImageURL    string `json:"image_url,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// RewardListResponse adalah model respons untuk daftar reward
type RewardListResponse struct {
	Rewards []Reward `json:"rewards"`
	Total   int64    `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}
