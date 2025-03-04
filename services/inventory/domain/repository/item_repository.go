package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/inventory/domain/entity"
)

// ItemRepository mendefinisikan operasi-operasi repository untuk Item
type ItemRepository interface {
	Create(ctx context.Context, item *entity.Item) error
	FindByID(ctx context.Context, id uint) (*entity.Item, error)
	FindBySKU(ctx context.Context, sku string) (*entity.Item, error)
	Update(ctx context.Context, item *entity.Item) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, activeOnly bool, page, limit int) ([]entity.Item, int64, error)
	ListByCategory(ctx context.Context, category string, activeOnly bool, page, limit int) ([]entity.Item, int64, error)
	Search(ctx context.Context, query string, activeOnly bool, page, limit int) ([]entity.Item, int64, error)
	FindLowStockItems(ctx context.Context) ([]entity.Item, int64, error)
	UpdateStock(ctx context.Context, id uint, quantity int) error
}
