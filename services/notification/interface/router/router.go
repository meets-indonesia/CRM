package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/notification/config"
	"github.com/kevinnaserwan/crm-be/services/notification/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/notification/middleware"
)

// Setup initializes the router
func Setup(cfg *config.Config, notificationHandler *handler.NotificationHandler) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create a new Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add config to Gin context
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", notificationHandler.HealthCheck)

	// Notification endpoints
	notifications := r.Group("/notifications")
	{
		// Auth required endpoints
		auth := notifications.Group("")
		auth.Use(middleware.AuthMiddleware())

		// User notifications
		auth.GET("", notificationHandler.ListUserNotifications)
		auth.GET("/:id", notificationHandler.GetNotification)

		// Admin only endpoints
		admin := notifications.Group("")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())

		// Send notifications
		admin.POST("/email", notificationHandler.SendEmail)
		admin.POST("/push", notificationHandler.SendPushNotification)

		// Process pending notifications
		admin.POST("/process", notificationHandler.ProcessPendingNotifications)

		// List user notifications (admin can see any user's notifications)
		admin.GET("/user/:user_id", notificationHandler.ListUserNotifications)
	}

	return r
}
