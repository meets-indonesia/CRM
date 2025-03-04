package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/inventory/config"
	"github.com/kevinnaserwan/crm-be/services/inventory/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/inventory/middleware"
)

// Setup initializes the router
func Setup(cfg *config.Config, inventoryHandler *handler.InventoryHandler) *gin.Engine {
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
	r.GET("/health", inventoryHandler.HealthCheck)

	// Items endpoints
	items := r.Group("/items")
	{
		// Public endpoints
		items.GET("", inventoryHandler.ListItems)
		items.GET("/:id", inventoryHandler.GetItem)
		items.GET("/category/:category", inventoryHandler.ListItemsByCategory)
		items.GET("/search", inventoryHandler.SearchItems)
		items.GET("/low-stock", inventoryHandler.GetLowStockItems)

		// Admin only endpoints
		admin := items.Group("")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
		admin.POST("", inventoryHandler.CreateItem)
		admin.PUT("/:id", inventoryHandler.UpdateItem)
		admin.DELETE("/:id", inventoryHandler.DeleteItem)

		// Stock endpoints
		stock := items.Group("/:id/stock")
		stock.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
		stock.POST("/increase", inventoryHandler.IncreaseStock)
		stock.POST("/decrease", inventoryHandler.DecreaseStock)
	}

	// Stock transactions endpoints
	transactions := r.Group("/transactions")
	transactions.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
	{
		transactions.GET("", inventoryHandler.ListAllStockTransactions)
		transactions.GET("/:id", inventoryHandler.GetStockTransaction)
		transactions.GET("/item/:id", inventoryHandler.ListStockTransactionsByItem)
	}

	// Suppliers endpoints
	suppliers := r.Group("/suppliers")
	suppliers.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
	{
		suppliers.GET("", inventoryHandler.ListSuppliers)
		suppliers.GET("/:id", inventoryHandler.GetSupplier)
		suppliers.GET("/search", inventoryHandler.SearchSuppliers)
		suppliers.POST("", inventoryHandler.CreateSupplier)
		suppliers.PUT("/:id", inventoryHandler.UpdateSupplier)
		suppliers.DELETE("/:id", inventoryHandler.DeleteSupplier)
	}

	return r
}
