package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/inventory/domain/entity"
	"gorm.io/gorm"
)

// GormStockRepository implements StockRepository with GORM
type GormStockRepository struct {
	db *gorm.DB
}

// NewGormStockRepository creates a new GormStockRepository
func NewGormStockRepository(db *gorm.DB) *GormStockRepository {
	return &GormStockRepository{
		db: db,
	}
}

// CreateTransaction creates a new stock transaction
func (r *GormStockRepository) CreateTransaction(ctx context.Context, transaction *entity.StockTransaction) error {
	result := r.db.Create(transaction)
	return result.Error
}

// FindTransactionByID finds a stock transaction by ID
func (r *GormStockRepository) FindTransactionByID(ctx context.Context, id uint) (*entity.StockTransaction, error) {
	var transaction entity.StockTransaction
	result := r.db.Preload("Item").First(&transaction, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &transaction, nil
}

// ListTransactionsByItemID lists stock transactions by item ID
func (r *GormStockRepository) ListTransactionsByItemID(ctx context.Context, itemID uint, page, limit int) ([]entity.StockTransaction, int64, error) {
	var transactions []entity.StockTransaction
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.StockTransaction{}).Where("item_id = ?", itemID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get transactions
	if err := r.db.Where("item_id = ?", itemID).Preload("Item").Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// ListAllTransactions lists all stock transactions
func (r *GormStockRepository) ListAllTransactions(ctx context.Context, page, limit int) ([]entity.StockTransaction, int64, error) {
	var transactions []entity.StockTransaction
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&entity.StockTransaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get transactions
	if err := r.db.Preload("Item").Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}
