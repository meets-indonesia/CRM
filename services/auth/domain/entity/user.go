package entity

import (
	"time"
)

// Role mendefinisikan peran pengguna
type Role string

const (
	RoleAdmin    Role = "ADMIN"
	RoleCustomer Role = "CUSTOMER"
)

// User adalah entitas pengguna
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"` // Password tidak dikembalikan dalam JSON
	Name      string    `json:"name" gorm:"not null"`
	Role      Role      `json:"role" gorm:"not null"`
	GoogleID  string    `json:"google_id,omitempty" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OTP adalah entitas untuk One-Time Password
type OTP struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Code      string    `json:"code" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest adalah model untuk permintaan login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest adalah model untuk permintaan registrasi
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// ResetPasswordRequest adalah model untuk permintaan reset password
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyOTPRequest adalah model untuk permintaan verifikasi OTP
type VerifyOTPRequest struct {
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// GoogleLoginRequest adalah model untuk permintaan login dengan Google
type GoogleLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

// LoginResponse adalah model untuk respons login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
