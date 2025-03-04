package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/inventory/domain/entity"
	"gorm.io/gorm"
)

// GormSupplierRepository implements SupplierRepository with GORM
type GormSupplierRepository struct {
	db *gorm.DB
}

// NewGormSupplierRepository creates a new GormSupplierRepository
func NewGormSupplierRepository(db *gorm.DB) *GormSupplierRepository {
	return &GormSupplierRepository{
		db: db,
	}
}

// Create creates a new supplier
func (r *GormSupplierRepository) Create(ctx context.Context, supplier *entity.Supplier) error {
	result := r.db.Create(supplier)
	return result.Error
}

// FindByID finds a supplier by ID
func (r *GormSupplierRepository) FindByID(ctx context.Context, id uint) (*entity.Supplier, error) {
	var supplier entity.Supplier
	result := r.db.First(&supplier, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &supplier, nil
}

// Update updates a supplier
func (r *GormSupplierRepository) Update(ctx context.Context, supplier *entity.Supplier) error {
	result := r.db.Save(supplier)
	return result.Error
}

// Delete deletes a supplier
func (r *GormSupplierRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.Delete(&entity.Supplier{}, id)
	return result.Error
}

// List lists suppliers
func (r *GormSupplierRepository) List(ctx context.Context, activeOnly bool, page, limit int) ([]entity.Supplier, int64, error) {
	var suppliers []entity.Supplier
	var total int64

	offset := (page - 1) * limit
	query := r.db

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	// Count total
	if err := query.Model(&entity.Supplier{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get suppliers
	if err := query.Order("name ASC").Offset(offset).Limit(limit).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

// Search searches suppliers by name, contact name, or email
func (r *GormSupplierRepository) Search(ctx context.Context, query string, activeOnly bool, page, limit int) ([]entity.Supplier, int64, error) {
	var suppliers []entity.Supplier
	var total int64

	offset := (page - 1) * limit
	dbQuery := r.db.Where("name ILIKE ? OR contact_name ILIKE ? OR email ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%")

	if activeOnly {
		dbQuery = dbQuery.Where("is_active = ?", true)
	}

	// Count total
	if err := dbQuery.Model(&entity.Supplier{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get suppliers
	if err := dbQuery.Order("name ASC").Offset(offset).Limit(limit).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}
