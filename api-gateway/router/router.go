package router

import (
	"io"
	"net/http"

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

	// Serve static files from Article service
	r.Use(middleware.APIKeyAuth()).GET("/article/uploads/*filepath", articleProxy.AccessUploadImages)

	r.Use(middleware.APIKeyAuth()).GET("/feedbacks/uploads/*filepath", func(c *gin.Context) {
		// Ambil hanya path file
		filepath := c.Param("filepath") // contoh: /17429xxx.png

		// Bangun ulang path yang sesuai di service backend
		targetURL := cfg.Services.FeedbackURL + "/uploads" + filepath

		resp, err := http.Get(targetURL)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch image from feedback service"})
			return
		}
		defer resp.Body.Close()

		// Salin content-type & status
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}
		c.Status(resp.StatusCode)

		// Stream response
		io.Copy(c.Writer, resp.Body)
	})

	// Serve reward uploads directly from reward service
	r.Use(middleware.APIKeyAuth()).GET("/uploads/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		targetURL := cfg.Services.RewardURL + "/uploads" + filepath

		resp, err := http.Get(targetURL)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch image from reward service"})
			return
		}
		defer resp.Body.Close()

		// Copy headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}
		c.Status(resp.StatusCode)

		// Stream content
		io.Copy(c.Writer, resp.Body)
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.Use(middleware.SimpleAPIKeyAuth()).POST("/admin/login", authProxy.AdminLogin)
		auth.Use(middleware.SimpleAPIKeyAuth()).POST("/admin/register", authProxy.AdminRegister)
		auth.Use(middleware.SimpleAPIKeyAuth()).POST("/admin/reset-password", authProxy.AdminResetPassword)
		auth.Use(middleware.SimpleAPIKeyAuth()).POST("/admin/verify-otp", authProxy.AdminVerifyOTP)
		auth.Use(middleware.APIKeyAuth()).POST("/customer/login", authProxy.CustomerLogin)
		auth.Use(middleware.APIKeyAuth()).POST("/customer/google", authProxy.CustomerGoogleLogin)
	}

	// validate user token
	r.Use(middleware.SimpleAPIKeyAuth()).GET("/validate", authProxy.ValidateToken)

	// Articles public routes
	r.Use(middleware.APIKeyAuth()).GET("/articles", articleProxy.ListArticles)
	r.Use(middleware.APIKeyAuth()).GET("/articles/:id", articleProxy.ViewArticle)
	r.Use(middleware.APIKeyAuth()).GET("/articles/search", articleProxy.SearchArticles)

	// Public rewards routes
	r.Use(middleware.APIKeyAuth()).GET("/rewards", rewardProxy.ListRewards)
	r.Use(middleware.APIKeyAuth()).GET("/rewards/:id", rewardProxy.GetReward)

	// Qr Scan
	r.Use(middleware.APIKeyAuth()).GET("/qr/verify/:code", feedbackProxy.VerifyQRCode)

	// Routes that require authentication
	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuth(cfg.JWT))

	// Admin routes
	admin := authorized.Group("/")
	admin.Use(middleware.AdminOnly())
	admin.Use(middleware.SimpleAPIKeyAuth())
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

		// QR Feedback Management
		admin.POST("/qr-feedbacks", feedbackProxy.CreateQRFeedback)
		admin.GET("/qr-feedbacks", feedbackProxy.ListQRFeedbacks)
		admin.GET("/qr-feedbacks/:id", feedbackProxy.GetQRFeedback)
		admin.GET("/qr-feedbacks/:id/download", feedbackProxy.GenerateQRCodeImage)

		// Tambahkan rute admin lainnya...
	}

	// Customer routes
	customer := authorized.Group("/")
	customer.Use(middleware.CustomerOnly())
	{
		// Profile
		// customer.GET("/users/:id", userProxy.GetUser)  // hapus command jika ingin get user data hanya dilakukan oleh user
		customer.Use(middleware.APIKeyAuth()).PUT("/users/:id", userProxy.UpdateUser)

		// Feedback
		customer.Use(middleware.APIKeyAuth()).POST("/feedbacks", feedbackProxy.CreateFeedback)

		// Reward
		customer.Use(middleware.APIKeyAuth()).POST("/claims", rewardProxy.ClaimReward)
		customer.Use(middleware.SimpleAPIKeyAuth()).GET("/claims/user", rewardProxy.ListUserClaims)

		// Notifications
		customer.Use(middleware.APIKeyAuth()).GET("/notifications", notificationProxy.ListUserNotifications)

		// Tambahkan rute customer lainnya...
	}

	// Common authenticated routes (accessible by both admin and customer)
	authorized.GET("/users/:id", userProxy.GetUser)
	authorized.GET("/feedbacks/:id", feedbackProxy.GetFeedback)
	authorized.GET("/notifications/:id", notificationProxy.GetNotification)
	authorized.GET("/feedbacks/user/:user_id", feedbackProxy.ListUserFeedback)
	authorized.GET("/feedbacks/user", feedbackProxy.ListUserFeedback)
	authorized.GET("/users/customer/:id/points", userProxy.GetCustomerPoints)

	// Inventory public routes
	r.GET("/items", inventoryProxy.ListItems)
	r.GET("/items/:id", inventoryProxy.GetItem)
	r.GET("/items/category/:category", inventoryProxy.ListItemsByCategory)
	r.GET("/items/search", inventoryProxy.SearchItems)

	return r
}
