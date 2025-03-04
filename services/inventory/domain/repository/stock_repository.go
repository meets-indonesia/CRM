package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/inventory/domain/entity"
)

// StockRepository mendefinisikan operasi-operasi repository untuk StockTransaction
type StockRepository interface {
	CreateTransaction(ctx context.Context, transaction *entity.StockTransaction) error
	FindTransactionByID(ctx context.Context, id uint) (*entity.StockTransaction, error)
	ListTransactionsByItemID(ctx context.Context, itemID uint, page, limit int) ([]entity.StockTransaction, int64, error)
	ListAllTransactions(ctx context.Context, page, limit int) ([]entity.StockTransaction, int64, error)
}
