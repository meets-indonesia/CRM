package entity

import (
	"time"
)

// PointTransaction adalah entitas untuk transaksi poin
type PointTransaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"index;not null"`
	Amount      int       `json:"amount" gorm:"not null"` // Bisa positif (penambahan) atau negatif (pengurangan)
	Type        string    `json:"type" gorm:"not null"`   // 'feedback', 'reward', dll
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// PointLevel adalah entitas untuk level poin
type PointLevel struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"not null"`       // 'Bronze', 'Silver', 'Gold', etc.
	MinPoints  int       `json:"min_points" gorm:"not null"` // Minimum poin untuk level ini
	Multiplier float64   `json:"multiplier" gorm:"not null"` // Multiplier untuk hadiah, diskon, dll.
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PointBalance adalah ringkasan saldo poin pengguna
type PointBalance struct {
	UserID    uint   `json:"user_id"`
	Total     int    `json:"total"`
	Level     string `json:"level"`
	NextLevel string `json:"next_level,omitempty"`
	ToNext    int    `json:"to_next,omitempty"` // Poin yang dibutuhkan untuk naik level
}

// PointTransactionListResponse adalah model respons untuk daftar transaksi poin
type PointTransactionListResponse struct {
	Transactions []PointTransaction `json:"transactions"`
	Total        int64              `json:"total"`
	Page         int                `json:"page"`
	Limit        int                `json:"limit"`
}
