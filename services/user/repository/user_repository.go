package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/user/domain/entity"
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

// Delete deletes a user
func (r *GormUserRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.Delete(&entity.User{}, id)
	return result.Error
}

// List lists users by role
func (r *GormUserRepository) List(ctx context.Context, role entity.Role, page, limit int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.User{}).Where("role = ?", role).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get users
	if err := r.db.Where("role = ?", role).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
