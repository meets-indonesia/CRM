package entity

import (
	"time"
)

// Item adalah entitas untuk item dalam inventaris
type Item struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null"`
	Description     string    `json:"description"`
	SKU             string    `json:"sku" gorm:"uniqueIndex;not null"`
	ImageURL        string    `json:"image_url,omitempty"`
	Category        string    `json:"category"`
	CurrentStock    int       `json:"current_stock" gorm:"not null;default:0"`
	MinimumStock    int       `json:"minimum_stock" gorm:"not null;default:5"`
	IsActive        bool      `json:"is_active" gorm:"not null;default:true"`
	SupplierID      *uint     `json:"supplier_id,omitempty" gorm:"index"`
	ReorderQuantity int       `json:"reorder_quantity" gorm:"not null;default:10"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateItemRequest adalah model untuk permintaan pembuatan item
type CreateItemRequest struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description"`
	SKU             string `json:"sku" binding:"required"`
	ImageURL        string `json:"image_url,omitempty"`
	Category        string `json:"category"`
	InitialStock    int    `json:"initial_stock" binding:"min=0"`
	MinimumStock    int    `json:"minimum_stock" binding:"min=0"`
	SupplierID      *uint  `json:"supplier_id,omitempty"`
	ReorderQuantity int    `json:"reorder_quantity" binding:"min=1"`
}

// UpdateItemRequest adalah model untuk permintaan update item
type UpdateItemRequest struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	SKU             string `json:"sku,omitempty"`
	ImageURL        string `json:"image_url,omitempty"`
	Category        string `json:"category,omitempty"`
	MinimumStock    *int   `json:"minimum_stock,omitempty" binding:"omitempty,min=0"`
	IsActive        *bool  `json:"is_active,omitempty"`
	SupplierID      *uint  `json:"supplier_id,omitempty"`
	ReorderQuantity *int   `json:"reorder_quantity,omitempty" binding:"omitempty,min=1"`
}

// StockUpdateRequest adalah model untuk permintaan update stok
type StockUpdateRequest struct {
	Quantity    int    `json:"quantity" binding:"required,nonzero"`
	Reason      string `json:"reason" binding:"required"`
	ReferenceID string `json:"reference_id,omitempty"`
}

// ItemListResponse adalah model respons untuk daftar item
type ItemListResponse struct {
	Items []Item `json:"items"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

// LowStockItem adalah model untuk item dengan stok rendah
type LowStockItem struct {
	Item          Item `json:"item"`
	Deficit       int  `json:"deficit"`
	ReorderAmount int  `json:"reorder_amount"`
}

// LowStockResponse adalah model respons untuk daftar item dengan stok rendah
type LowStockResponse struct {
	Items []LowStockItem `json:"items"`
	Total int64          `json:"total"`
}
