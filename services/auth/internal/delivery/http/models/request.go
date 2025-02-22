package models

type Meta struct {
	Action string `json:"action" binding:"required"`
}

type BaseRequest struct {
	Meta Meta        `json:"meta" binding:"required"`
	Data interface{} `json:"data" binding:"required"`
}

// Auth request structures
type LoginData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// services/auth/internal/delivery/http/models/request.go
type RegisterData struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password"`
	GoogleID  string `json:"google_id"`
	RoleID    string `json:"role_id"` // Change to string to handle empty value
}

type ForgotPasswordData struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPData struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type ResetPasswordData struct {
	Email           string `json:"email" binding:"required,email"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
