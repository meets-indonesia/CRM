package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/inventory/config"
	"github.com/kevinnaserwan/crm-be/services/inventory/domain/usecase"
	"github.com/kevinnaserwan/crm-be/services/inventory/infrastructure/db"
	"github.com/kevinnaserwan/crm-be/services/inventory/infrastructure/messaging"
	"github.com/kevinnaserwan/crm-be/services/inventory/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/inventory/interface/router"
	"github.com/kevinnaserwan/crm-be/services/inventory/repository"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup database
	database, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup RabbitMQ
	rabbitMQ, err := messaging.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Setup repositories
	itemRepo := repository.NewGormItemRepository(database)
	stockRepo := repository.NewGormStockRepository(database)
	supplierRepo := repository.NewGormSupplierRepository(database)

	// Setup usecase
	inventoryUsecase := usecase.NewInventoryUsecase(itemRepo, stockRepo, supplierRepo, rabbitMQ)

	// Setup handler
	inventoryHandler := handler.NewInventoryHandler(inventoryUsecase)

	// Setup subscriber
	subscriber, err := messaging.NewRabbitMQSubscriber(cfg.RabbitMQ, inventoryUsecase)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ subscriber: %v", err)
	}
	defer subscriber.Close()

	// Subscribe to events
	if err := subscriber.SubscribeToEvents(); err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}

	// Add config to Gin context
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Setup router
	router := router.Setup(cfg, inventoryHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Inventory Service starting on port %s", port)
	router.Run(":" + port)
}
