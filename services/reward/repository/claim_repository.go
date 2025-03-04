package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/reward/domain/entity"
	"gorm.io/gorm"
)

// GormClaimRepository implements ClaimRepository with GORM
type GormClaimRepository struct {
	db *gorm.DB
}

// NewGormClaimRepository creates a new GormClaimRepository
func NewGormClaimRepository(db *gorm.DB) *GormClaimRepository {
	return &GormClaimRepository{
		db: db,
	}
}

// Create creates a new claim
func (r *GormClaimRepository) Create(ctx context.Context, claim *entity.RewardClaim) error {
	result := r.db.Create(claim)
	return result.Error
}

// FindByID finds a claim by ID
func (r *GormClaimRepository) FindByID(ctx context.Context, id uint) (*entity.RewardClaim, error) {
	var claim entity.RewardClaim
	result := r.db.Preload("Reward").First(&claim, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &claim, nil
}

// Update updates a claim
func (r *GormClaimRepository) Update(ctx context.Context, claim *entity.RewardClaim) error {
	result := r.db.Save(claim)
	return result.Error
}

// ListByUserID lists claims by user ID
func (r *GormClaimRepository) ListByUserID(ctx context.Context, userID uint, page, limit int) ([]entity.RewardClaim, int64, error) {
	var claims []entity.RewardClaim
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.RewardClaim{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get claims
	if err := r.db.Where("user_id = ?", userID).Preload("Reward").Order("created_at DESC").Offset(offset).Limit(limit).Find(&claims).Error; err != nil {
		return nil, 0, err
	}

	return claims, total, nil
}

// ListAll lists all claims
func (r *GormClaimRepository) ListAll(ctx context.Context, page, limit int) ([]entity.RewardClaim, int64, error) {
	var claims []entity.RewardClaim
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.RewardClaim{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get claims
	if err := r.db.Preload("Reward").Order("created_at DESC").Offset(offset).Limit(limit).Find(&claims).Error; err != nil {
		return nil, 0, err
	}

	return claims, total, nil
}

// ListByStatus lists claims by status
func (r *GormClaimRepository) ListByStatus(ctx context.Context, status entity.ClaimStatus, page, limit int) ([]entity.RewardClaim, int64, error) {
	var claims []entity.RewardClaim
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.RewardClaim{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get claims
	if err := r.db.Where("status = ?", status).Preload("Reward").Order("created_at DESC").Offset(offset).Limit(limit).Find(&claims).Error; err != nil {
		return nil, 0, err
	}

	return claims, total, nil
}
