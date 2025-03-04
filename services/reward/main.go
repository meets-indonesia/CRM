package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/reward/config"
	"github.com/kevinnaserwan/crm-be/services/reward/domain/usecase"
	"github.com/kevinnaserwan/crm-be/services/reward/infrastructure/db"
	"github.com/kevinnaserwan/crm-be/services/reward/infrastructure/messaging"
	"github.com/kevinnaserwan/crm-be/services/reward/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/reward/interface/router"
	"github.com/kevinnaserwan/crm-be/services/reward/repository"
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

	// Setup repositories and services
	rewardRepo := repository.NewGormRewardRepository(database)
	claimRepo := repository.NewGormClaimRepository(database)
	userPointService := repository.NewHTTPUserPointService(cfg.Services)

	// Setup usecase
	rewardUsecase := usecase.NewRewardUsecase(rewardRepo, claimRepo, rabbitMQ, userPointService)

	// Setup handler
	rewardHandler := handler.NewRewardHandler(rewardUsecase)

	// Add config to Gin context
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Setup router
	router := router.Setup(cfg, rewardHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Reward Service starting on port %s", port)
	router.Run(":" + port)
}
