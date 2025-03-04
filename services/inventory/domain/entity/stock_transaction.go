package entity

import (
	"time"
)

// TransactionType mendefinisikan tipe transaksi stok
type TransactionType string

const (
	TransactionTypeIncrease TransactionType = "INCREASE"
	TransactionTypeDecrease TransactionType = "DECREASE"
)

// StockTransaction adalah entitas untuk transaksi stok
type StockTransaction struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	ItemID      uint            `json:"item_id" gorm:"index;not null"`
	Type        TransactionType `json:"type" gorm:"not null"`
	Quantity    int             `json:"quantity" gorm:"not null"`
	PreviousQty int             `json:"previous_qty" gorm:"not null"`
	NewQty      int             `json:"new_qty" gorm:"not null"`
	Reason      string          `json:"reason" gorm:"not null"`
	ReferenceID string          `json:"reference_id,omitempty"`
	PerformedBy uint            `json:"performed_by" gorm:"not null"` // User ID yang melakukan
	CreatedAt   time.Time       `json:"created_at"`

	// Relasi
	Item Item `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// StockTransactionListResponse adalah model respons untuk daftar transaksi stok
type StockTransactionListResponse struct {
	Transactions []StockTransaction `json:"transactions"`
	Total        int64              `json:"total"`
	Page         int                `json:"page"`
	Limit        int                `json:"limit"`
}
