package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/reward/config"
	"github.com/kevinnaserwan/crm-be/services/reward/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/reward/middleware"
)

// Setup initializes the router
func Setup(cfg *config.Config, rewardHandler *handler.RewardHandler) *gin.Engine {
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
	r.GET("/health", rewardHandler.HealthCheck)
	r.Static("/uploads", "./uploads")

	// Reward endpoints
	rewards := r.Group("/rewards")
	{
		// Public endpoints
		rewards.GET("", rewardHandler.ListRewards)
		rewards.GET("/:id", rewardHandler.GetReward)

		// Admin only endpoints
		admin := rewards.Group("")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
		admin.POST("", rewardHandler.CreateReward)
		admin.PUT("/:id", rewardHandler.UpdateReward)
		admin.DELETE("/:id", rewardHandler.DeleteReward)
	}

	// Claim endpoints
	claims := r.Group("/claims")
	{
		// Customer endpoints
		customer := claims.Group("")
		customer.Use(middleware.AuthMiddleware(), middleware.CustomerOnlyMiddleware())
		customer.POST("", rewardHandler.ClaimReward)
		customer.GET("/user", rewardHandler.ListUserClaims)

		// Public endpoints
		claims.GET("/:id", middleware.AuthMiddleware(), rewardHandler.GetClaim)

		// Admin only endpoints
		admin := claims.Group("")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
		admin.PUT("/:id/status", rewardHandler.UpdateClaimStatus)
		admin.GET("", rewardHandler.ListAllClaims)
		admin.GET("/status/:status", rewardHandler.ListClaimsByStatus)
		admin.GET("/user/:user_id", rewardHandler.ListUserClaims)
	}

	return r
}
