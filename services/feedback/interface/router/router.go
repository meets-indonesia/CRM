package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/feedback/config"
	"github.com/kevinnaserwan/crm-be/services/feedback/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/feedback/middleware"
)

// Setup initializes the router
func Setup(cfg *config.Config,
	feedbackHandler *handler.FeedbackHandler,
	qrFeedbackHandler *handler.QRFeedbackHandler) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create a new Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Tambahkan config ke context
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Serve static files untuk gambar feedback
	r.Static("/uploads", cfg.FileStore.UploadDir)

	// Health check endpoint
	r.GET("/health", feedbackHandler.HealthCheck)

	// Feedback endpoints
	feedbacks := r.Group("/feedbacks")
	{
		// Get feedback by ID - no auth required
		feedbacks.GET("/:id", feedbackHandler.GetFeedback)
		// Create feedback - customer only
		feedbacks.POST("", middleware.AuthMiddleware(), middleware.CustomerOnlyMiddleware(), feedbackHandler.CreateFeedback)
		// Respond to feedback - admin only
		feedbacks.PUT("/:id/respond", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), feedbackHandler.RespondFeedback)
		// List all feedback - admin only
		feedbacks.GET("", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), feedbackHandler.ListAllFeedback)
		// List pending feedback - admin only
		feedbacks.GET("/pending", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), feedbackHandler.ListPendingFeedback)
		// List user feedback - can be used by both customer (for own feedback) and admin (for any user's feedback)
		feedbacks.GET("/user", middleware.AuthMiddleware(), feedbackHandler.ListUserFeedback)
		feedbacks.GET("/user/:user_id", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), feedbackHandler.ListUserFeedback)
	}

	// QR Feedback endpoints
	qrFeedbacks := r.Group("/qr-feedbacks")
	{
		// Admin-only endpoints
		qrFeedbacks.POST("", middleware.AuthMiddleware(), qrFeedbackHandler.CreateQRFeedback)
		qrFeedbacks.GET("", middleware.AuthMiddleware(), qrFeedbackHandler.ListQRFeedback)
		qrFeedbacks.GET("/:id", middleware.AuthMiddleware(), qrFeedbackHandler.GetQRFeedback)
		qrFeedbacks.GET("/:id/download", middleware.AuthMiddleware(), qrFeedbackHandler.DownloadQRCode)
	}

	// Public QR scan endpoint - no auth required
	r.GET("/feedback/scan/:qrCode", qrFeedbackHandler.VerifyQRCode)

	return r
}
