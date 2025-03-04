package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/auth/config"
	"github.com/kevinnaserwan/crm-be/services/auth/domain/usecase"
	"github.com/kevinnaserwan/crm-be/services/auth/interface/handler"
)

// Setup initializes the router
func Setup(cfg *config.Config, authUsecase usecase.AuthUsecase) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create a new Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)

	// Health check endpoint
	r.GET("/health", authHandler.HealthCheck)

	// Admin auth endpoints
	admin := r.Group("/admin")
	{
		admin.POST("/register", authHandler.RegisterAdmin)
		admin.POST("/login", authHandler.LoginAdmin)
		admin.POST("/reset-password", authHandler.ResetPasswordRequest)
		admin.POST("/verify-otp", authHandler.VerifyOTP)
	}

	// Customer auth endpoints
	customer := r.Group("/customer")
	{
		customer.POST("/login", authHandler.LoginCustomer)
		customer.POST("/google", authHandler.LoginWithGoogle)
	}

	// Token validation endpoint
	r.GET("/validate", authHandler.ValidateToken)

	return r
}
