package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/api-gateway/config"
	"github.com/kevinnaserwan/crm-be/api-gateway/middleware"
	"github.com/kevinnaserwan/crm-be/api-gateway/proxy"
)

// Setup initializes the router
func Setup(cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create a new Gin router
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(cfg.RateLimit.RequestsPerSecond))

	// Initialize proxies
	authProxy := proxy.NewAuthProxy(cfg.Services.AuthURL)
	userProxy := proxy.NewUserProxy(cfg.Services.UserURL)
	// Inisialisasi proxy lainnya...

	// Public routes (no auth required)
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/admin/login", authProxy.AdminLogin)
		auth.POST("/admin/register", authProxy.AdminRegister)
		auth.POST("/admin/reset-password", authProxy.AdminResetPassword)
		auth.POST("/admin/verify-otp", authProxy.AdminVerifyOTP)
		auth.POST("/customer/login", authProxy.CustomerLogin)
		auth.POST("/customer/google", authProxy.CustomerGoogleLogin)
	}

	// Articles public routes
	r.GET("/articles", nil)     // Tambahkan articleProxy.ListArticles nanti
	r.GET("/articles/:id", nil) // Tambahkan articleProxy.GetArticle nanti

	// Routes that require authentication
	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuth(cfg.JWT))

	// Admin routes
	admin := authorized.Group("/")
	admin.Use(middleware.AdminOnly())
	{
		// User management
		admin.GET("/users/admin", userProxy.ListAdmins)
		admin.GET("/users/customer", userProxy.ListCustomers)

		// Tambahkan rute admin lainnya...
	}

	// Customer routes
	customer := authorized.Group("/")
	customer.Use(middleware.CustomerOnly())
	{
		// Profile
		customer.GET("/users/:id", userProxy.GetUser)
		customer.PUT("/users/:id", userProxy.UpdateUser)
		customer.GET("/users/customer/:id/points", userProxy.GetCustomerPoints)

		// Tambahkan rute customer lainnya...
	}

	// Common authenticated routes (accessible by both admin and customer)
	// ...

	return r
}
