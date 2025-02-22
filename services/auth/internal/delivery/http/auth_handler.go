// services/auth/internal/delivery/http/auth_handler.go
package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	case "login_oauth":
		h.handleLoginOAuth(c, baseRequest.Data)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
	}
}

func (h *AuthHandler) handleRegister(c *gin.Context, data interface{}) {
	var registerData models.RegisterData
	if err := convertData(data, &registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var password *string
	if registerData.Password != "" {
		password = &registerData.Password
	}

	var googleID *string
	if registerData.GoogleID != "" {
		googleID = &registerData.GoogleID
	}

	// Handle role ID
	var roleID uuid.UUID
	var err error

	// If it's a mobile user (has GoogleID) and no roleID provided
	if googleID != nil && registerData.RoleID == "" {
		defaultRole, err := h.authUseCase.GetDefaultRole(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get default role"})
			return
		}
		roleID = defaultRole.ID
	} else if registerData.RoleID != "" {
		// Convert string to UUID if provided
		roleID, err = uuid.Parse(registerData.RoleID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID is required for admin registration"})
		return
	}

	user := &model.User{
		FirstName: registerData.FirstName,
		LastName:  registerData.LastName,
		Email:     registerData.Email,
		Password:  password,
		GoogleID:  googleID,
		RoleID:    roleID,
	}

	if err := h.authUseCase.Register(c.Request.Context(), user); err != nil {
		// Check if it's an existing user error with Google ID
		if autoLoginErr, ok := err.(*usecase.AutoLoginError); ok {
			// Return special response for auto-login
			c.JSON(http.StatusConflict, gin.H{
				"error":   "user_exists",
				"message": autoLoginErr.Message,
				"action":  "auto_login",
				"data": gin.H{
					"email":     autoLoginErr.Email,
					"google_id": autoLoginErr.GoogleID,
				},
			})
			return
		}

		// Handle regular email exists error
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "email_exists",
				"message": "Email already registered",
			})
			return
		}

		// Handle other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) handleLoginOAuth(c *gin.Context, data interface{}) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		GoogleID string `json:"google_id" binding:"required"`
	}

	if err := convertData(data, &loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authUseCase.LoginOAuth(c.Request.Context(), loginData.Email, loginData.GoogleID)
	if err != nil {
		if err.Error() == "account uses password authentication" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_auth_type",
				"message": "This account uses password authentication",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) handleLogin(c *gin.Context, data interface{}) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := convertData(data, &loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authUseCase.Login(c.Request.Context(), loginData.Email, loginData.Password)
	if err != nil {
		if err.Error() == "this account uses Google authentication" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_auth_type",
				"message": "This account uses Google authentication",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
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
