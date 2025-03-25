package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/feedback/config"
	"github.com/kevinnaserwan/crm-be/services/feedback/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/feedback/middleware"
)

// Setup initializes the router
// Setup initializes the router
func Setup(cfg *config.Config, feedbackHandler *handler.FeedbackHandler) *gin.Engine {
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

	return r
}
