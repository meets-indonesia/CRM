package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/auth/domain/entity"
	"gorm.io/gorm"
)

// GormUserRepository implements UserRepository with GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		db: db,
	}
}

// FindByID finds a user by ID
func (r *GormUserRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByGoogleID finds a user by Google ID
func (r *GormUserRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("google_id = ?", googleID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// Create creates a new user
func (r *GormUserRepository) Create(ctx context.Context, user *entity.User) error {
	result := r.db.Create(user)
	return result.Error
}

// Update updates a user
func (r *GormUserRepository) Update(ctx context.Context, user *entity.User) error {
	result := r.db.Save(user)
	return result.Error
}

// CreateOTP creates a new OTP
func (r *GormUserRepository) CreateOTP(ctx context.Context, otp *entity.OTP) error {
	result := r.db.Create(otp)
	return result.Error
}

// FindOTPByCode finds an OTP by code
func (r *GormUserRepository) FindOTPByCode(ctx context.Context, code string) (*entity.OTP, error) {
	var otp entity.OTP
	result := r.db.Where("code = ?", code).First(&otp)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &otp, nil
}

// DeleteOTP deletes an OTP
func (r *GormUserRepository) DeleteOTP(ctx context.Context, id uint) error {
	result := r.db.Delete(&entity.OTP{}, id)
	return result.Error
}
