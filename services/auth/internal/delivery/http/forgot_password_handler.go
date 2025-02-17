package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/usecase"
)

type ForgotPasswordHandler struct {
	forgotPasswordUseCase *usecase.ForgotPasswordUseCase
}

// Add this constructor function
func NewForgotPasswordHandler(forgotPasswordUseCase *usecase.ForgotPasswordUseCase) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		forgotPasswordUseCase: forgotPasswordUseCase,
	}
}

func (h *ForgotPasswordHandler) RequestPasswordReset(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forgotPasswordUseCase.RequestPasswordReset(c.Request.Context(), request.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset OTP has been sent"})
}

func (h *ForgotPasswordHandler) VerifyOTP(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forgotPasswordUseCase.VerifyOTP(c.Request.Context(), request.Email, request.OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

func (h *ForgotPasswordHandler) ResetPassword(c *gin.Context) {
	var request struct {
		Email           string `json:"email" binding:"required,email"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forgotPasswordUseCase.ResetPassword(c.Request.Context(), request.Email, request.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
