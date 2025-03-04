package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/user/config"
	"github.com/kevinnaserwan/crm-be/services/user/interface/handler"
)

// Setup initializes the router
func Setup(cfg *config.Config, userHandler *handler.UserHandler) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create a new Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/health", userHandler.HealthCheck)

	// User endpoints
	users := r.Group("/users")
	{
		// Get user by ID
		users.GET("/:id", userHandler.GetUser)

		// Update user
		users.PUT("/:id", userHandler.UpdateUser)

		// List admins
		users.GET("/admin", userHandler.ListAdmins)

		// List customers
		users.GET("/customer", userHandler.ListCustomers)

		// Get customer points
		users.GET("/customer/:id/points", userHandler.GetCustomerPoints)

		// Get point transactions
		users.GET("/customer/:id/points/transactions", userHandler.GetPointTransactions)
	}

	return r
}
