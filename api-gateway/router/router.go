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
	feedbackProxy := proxy.NewFeedbackProxy(cfg.Services.FeedbackURL)
	rewardProxy := proxy.NewRewardProxy(cfg.Services.RewardURL)
	inventoryProxy := proxy.NewInventoryProxy(cfg.Services.InventoryURL)
	articleProxy := proxy.NewArticleProxy(cfg.Services.ArticleURL)
	notificationProxy := proxy.NewNotificationProxy(cfg.Services.NotificationURL)
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

	// validate user token
	r.GET("/validate", authProxy.ValidateToken)

	// Articles public routes
	r.GET("/articles", articleProxy.ListArticles)
	r.GET("/articles/:id", articleProxy.ViewArticle)
	r.GET("/articles/search", articleProxy.SearchArticles)

	// Public rewards routes
	r.GET("/rewards", rewardProxy.ListRewards)
	r.GET("/rewards/:id", rewardProxy.GetReward)

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

		// Feedback management
		admin.GET("/feedbacks", feedbackProxy.ListAllFeedback)
		admin.GET("/feedbacks/pending", feedbackProxy.ListPendingFeedback)
		admin.PUT("/feedbacks/:id/respond", feedbackProxy.RespondToFeedback)

		// Reward management
		admin.POST("/rewards", rewardProxy.CreateReward)
		admin.PUT("/rewards/:id", rewardProxy.UpdateReward)
		admin.DELETE("/rewards/:id", rewardProxy.DeleteReward)
		admin.GET("/claims", rewardProxy.ListAllClaims)
		admin.GET("/claims/:id", rewardProxy.GetClaim)
		admin.PUT("/claims/:id/status", rewardProxy.UpdateClaimStatus)
		admin.GET("/claims/user/:user_id", rewardProxy.ListUserClaims)
		admin.GET("/claims/status/:status", rewardProxy.ListClaimsByStatus)

		// Inventory management
		admin.POST("/items", inventoryProxy.CreateItem)
		admin.PUT("/items/:id", inventoryProxy.UpdateItem)
		admin.DELETE("/items/:id", inventoryProxy.DeleteItem)
		admin.GET("/items/low-stock", inventoryProxy.GetLowStockItems)
		admin.POST("/items/:id/stock/increase", inventoryProxy.IncreaseStock)
		admin.POST("/items/:id/stock/decrease", inventoryProxy.DecreaseStock)
		admin.GET("/transactions", inventoryProxy.ListAllStockTransactions)
		admin.GET("/transactions/:id", inventoryProxy.GetStockTransaction)
		admin.GET("/transactions/item/:id", inventoryProxy.ListStockTransactionsByItem)
		admin.POST("/suppliers", inventoryProxy.CreateSupplier)
		admin.GET("/suppliers", inventoryProxy.ListSuppliers)
		admin.GET("/suppliers/:id", inventoryProxy.GetSupplier)
		admin.PUT("/suppliers/:id", inventoryProxy.UpdateSupplier)
		admin.DELETE("/suppliers/:id", inventoryProxy.DeleteSupplier)
		admin.GET("/suppliers/search", inventoryProxy.SearchSuppliers)

		// Article management
		admin.POST("/articles", articleProxy.CreateArticle)
		admin.PUT("/articles/:id", articleProxy.UpdateArticle)
		admin.DELETE("/articles/:id", articleProxy.DeleteArticle)

		// Notification management
		admin.POST("/notifications/email", notificationProxy.SendEmail)
		admin.POST("/notifications/push", notificationProxy.SendPushNotification)
		admin.POST("/notifications/process", notificationProxy.ProcessPendingNotifications)
		admin.GET("/notifications/user/:user_id", notificationProxy.ListUserNotifications)

		// Tambahkan rute admin lainnya...
	}

	// Customer routes
	customer := authorized.Group("/")
	customer.Use(middleware.CustomerOnly())
	{
		// Profile
		// customer.GET("/users/:id", userProxy.GetUser)  // hapus command jika ingin get user data hanya dilakukan oleh user
		customer.PUT("/users/:id", userProxy.UpdateUser)

		// Feedback
		customer.POST("/feedbacks", feedbackProxy.CreateFeedback)

		// Reward
		customer.POST("/claims", rewardProxy.ClaimReward)
		customer.GET("/claims/user", rewardProxy.ListUserClaims)

		// Notifications
		customer.GET("/notifications", notificationProxy.ListUserNotifications)

		// Tambahkan rute customer lainnya...
	}

	// Common authenticated routes (accessible by both admin and customer)
	authorized.GET("/users/:id", userProxy.GetUser)
	authorized.GET("/feedbacks/:id", feedbackProxy.GetFeedback)
	authorized.GET("/notifications/:id", notificationProxy.GetNotification)
	authorized.GET("/feedbacks/user/:user_id", feedbackProxy.ListUserFeedback)
	authorized.GET("/feedbacks/user", feedbackProxy.ListUserFeedback)
	authorized.GET("/feedbacks/uploads/:filename", feedbackProxy.AccesUploadImages)
	authorized.GET("/users/customer/:id/points", userProxy.GetCustomerPoints)

	// Inventory public routes
	r.GET("/items", inventoryProxy.ListItems)
	r.GET("/items/:id", inventoryProxy.GetItem)
	r.GET("/items/category/:category", inventoryProxy.ListItemsByCategory)
	r.GET("/items/search", inventoryProxy.SearchItems)

	return r
}
