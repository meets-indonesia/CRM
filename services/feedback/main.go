package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/feedback/config"
	"github.com/kevinnaserwan/crm-be/services/feedback/domain/usecase"
	"github.com/kevinnaserwan/crm-be/services/feedback/infrastructure/db"
	"github.com/kevinnaserwan/crm-be/services/feedback/infrastructure/messaging"
	"github.com/kevinnaserwan/crm-be/services/feedback/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/feedback/interface/router"
	"github.com/kevinnaserwan/crm-be/services/feedback/repository"
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
	feedbackRepo := repository.NewGormFeedbackRepository(database)

	// Setup usecase
	feedbackUsecase := usecase.NewFeedbackUsecase(feedbackRepo, rabbitMQ)

	// Setup handler
	feedbackHandler := handler.NewFeedbackHandler(feedbackUsecase)

	// Setup subscriber
	subscriber, err := messaging.NewRabbitMQSubscriber(cfg.RabbitMQ, feedbackUsecase)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ subscriber: %v", err)
	}
	defer subscriber.Close()

	// Subscribe to events
	if err := subscriber.SubscribeToEvents(); err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}

	// Add config to Gin context
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Setup router
	router := router.Setup(cfg, feedbackHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Feedback Service starting on port %s", port)
	router.Run(":" + port)
}
