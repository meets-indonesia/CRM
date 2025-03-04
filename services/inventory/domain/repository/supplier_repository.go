package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/inventory/domain/entity"
)

// SupplierRepository mendefinisikan operasi-operasi repository untuk Supplier
type SupplierRepository interface {
	Create(ctx context.Context, supplier *entity.Supplier) error
	FindByID(ctx context.Context, id uint) (*entity.Supplier, error)
	Update(ctx context.Context, supplier *entity.Supplier) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, activeOnly bool, page, limit int) ([]entity.Supplier, int64, error)
	Search(ctx context.Context, query string, activeOnly bool, page, limit int) ([]entity.Supplier, int64, error)
}

// EventPublisher mendefinisikan operasi-operasi untuk mempublikasikan event
type EventPublisher interface {
	PublishLowStockAlert(item *entity.Item, deficit int) error
	PublishStockUpdated(transaction *entity.StockTransaction) error
	Close() error
}
