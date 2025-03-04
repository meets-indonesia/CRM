package entity

import (
	"time"
)

// Supplier adalah entitas untuk supplier
type Supplier struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	ContactName string    `json:"contact_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSupplierRequest adalah model untuk permintaan pembuatan supplier
type CreateSupplierRequest struct {
	Name        string `json:"name" binding:"required"`
	ContactName string `json:"contact_name"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
}

// UpdateSupplierRequest adalah model untuk permintaan update supplier
type UpdateSupplierRequest struct {
	Name        string `json:"name,omitempty"`
	ContactName string `json:"contact_name,omitempty"`
	Email       string `json:"email,omitempty" binding:"omitempty,email"`
	Phone       string `json:"phone,omitempty"`
	Address     string `json:"address,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// SupplierListResponse adalah model respons untuk daftar supplier
type SupplierListResponse struct {
	Suppliers []Supplier `json:"suppliers"`
	Total     int64      `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}
