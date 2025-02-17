// services/auth/internal/delivery/http/auth_handler.go
package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/delivery/http/models"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/usecase"
)

type AuthHandler struct {
	authUseCase           *usecase.AuthUseCase
	forgotPasswordUseCase *usecase.ForgotPasswordUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase, forgotPasswordUseCase *usecase.ForgotPasswordUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase:           authUseCase,
		forgotPasswordUseCase: forgotPasswordUseCase,
	}
}

func (h *AuthHandler) Handle(c *gin.Context) {
	var baseRequest models.BaseRequest
	if err := c.ShouldBindJSON(&baseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch baseRequest.Meta.Action {
	case "login":
		h.handleLogin(c, baseRequest.Data)
	case "register":
		h.handleRegister(c, baseRequest.Data)
	case "forgot_password":
		h.handleForgotPassword(c, baseRequest.Data)
	case "verify_otp":
		h.handleVerifyOTP(c, baseRequest.Data)
	case "reset_password":
		h.handleResetPassword(c, baseRequest.Data)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
	}
}

func (h *AuthHandler) handleLogin(c *gin.Context, data interface{}) {
	var loginData models.LoginData
	if err := convertData(data, &loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authUseCase.Login(c.Request.Context(), loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) handleRegister(c *gin.Context, data interface{}) {
	var registerData models.RegisterData
	if err := convertData(data, &registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		FirstName: registerData.FirstName,
		LastName:  registerData.LastName,
		Email:     registerData.Email,
		Password:  registerData.Password,
		RoleID:    registerData.RoleID,
	}

	if err := h.authUseCase.Register(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) handleForgotPassword(c *gin.Context, data interface{}) {
	var forgotPasswordData models.ForgotPasswordData
	if err := convertData(data, &forgotPasswordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forgotPasswordUseCase.RequestPasswordReset(c.Request.Context(), forgotPasswordData.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset OTP has been sent"})
}

func (h *AuthHandler) handleVerifyOTP(c *gin.Context, data interface{}) {
	var verifyOTPData models.VerifyOTPData
	if err := convertData(data, &verifyOTPData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forgotPasswordUseCase.VerifyOTP(c.Request.Context(), verifyOTPData.Email, verifyOTPData.OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

func (h *AuthHandler) handleResetPassword(c *gin.Context, data interface{}) {
	var resetPasswordData models.ResetPasswordData
	if err := convertData(data, &resetPasswordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.forgotPasswordUseCase.ResetPassword(c.Request.Context(), resetPasswordData.Email, resetPasswordData.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func convertData(data interface{}, target interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, target)
}
