package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/inventory/domain/entity"
	"gorm.io/gorm"
)

// GormItemRepository implements ItemRepository with GORM
type GormItemRepository struct {
	db *gorm.DB
}

// NewGormItemRepository creates a new GormItemRepository
func NewGormItemRepository(db *gorm.DB) *GormItemRepository {
	return &GormItemRepository{
		db: db,
	}
}

// Create creates a new item
func (r *GormItemRepository) Create(ctx context.Context, item *entity.Item) error {
	result := r.db.Create(item)
	return result.Error
}

// FindByID finds an item by ID
func (r *GormItemRepository) FindByID(ctx context.Context, id uint) (*entity.Item, error) {
	var item entity.Item
	result := r.db.First(&item, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &item, nil
}

// FindBySKU finds an item by SKU
func (r *GormItemRepository) FindBySKU(ctx context.Context, sku string) (*entity.Item, error) {
	var item entity.Item
	result := r.db.Where("sku = ?", sku).First(&item)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &item, nil
}

// Update updates an item
func (r *GormItemRepository) Update(ctx context.Context, item *entity.Item) error {
	result := r.db.Save(item)
	return result.Error
}

// Delete deletes an item
func (r *GormItemRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.Delete(&entity.Item{}, id)
	return result.Error
}

// List lists items
func (r *GormItemRepository) List(ctx context.Context, activeOnly bool, page, limit int) ([]entity.Item, int64, error) {
	var items []entity.Item
	var total int64

	offset := (page - 1) * limit
	query := r.db

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	// Count total
	if err := query.Model(&entity.Item{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get items
	if err := query.Order("name ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ListByCategory lists items by category
func (r *GormItemRepository) ListByCategory(ctx context.Context, category string, activeOnly bool, page, limit int) ([]entity.Item, int64, error) {
	var items []entity.Item
	var total int64

	offset := (page - 1) * limit
	query := r.db.Where("category = ?", category)

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	// Count total
	if err := query.Model(&entity.Item{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get items
	if err := query.Order("name ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Search searches items by name, SKU, or description
func (r *GormItemRepository) Search(ctx context.Context, query string, activeOnly bool, page, limit int) ([]entity.Item, int64, error) {
	var items []entity.Item
	var total int64

	offset := (page - 1) * limit
	dbQuery := r.db.Where("name ILIKE ? OR sku ILIKE ? OR description ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%")

	if activeOnly {
		dbQuery = dbQuery.Where("is_active = ?", true)
	}

	// Count total
	if err := dbQuery.Model(&entity.Item{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get items
	if err := dbQuery.Order("name ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindLowStockItems finds items with stock below minimum
func (r *GormItemRepository) FindLowStockItems(ctx context.Context) ([]entity.Item, int64, error) {
	var items []entity.Item
	var total int64

	query := r.db.Where("current_stock < minimum_stock").Where("is_active = ?", true)

	// Count total
	if err := query.Model(&entity.Item{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get items
	if err := query.Order("(minimum_stock - current_stock) DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// UpdateStock updates the stock of an item
func (r *GormItemRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
	result := r.db.Model(&entity.Item{}).Where("id = ?", id).Update("current_stock", quantity)
	return result.Error
}
