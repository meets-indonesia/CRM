package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/feedback/config"
	"github.com/kevinnaserwan/crm-be/services/feedback/domain/usecase"
	"github.com/kevinnaserwan/crm-be/services/feedback/infrastructure/db"
	"github.com/kevinnaserwan/crm-be/services/feedback/infrastructure/email"
	"github.com/kevinnaserwan/crm-be/services/feedback/infrastructure/filestore"
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

	// Setup file service
	fileService := filestore.NewLocalFileService(cfg.FileStore.UploadDir, cfg.FileStore.MaxSize)

	// Setup email service
	emailService := email.NewSMTPEmailService(
		cfg.Email.Host,
		cfg.Email.Port,
		cfg.Email.Username,
		cfg.Email.Password,
		cfg.Email.From,
		cfg.Email.AdminEmail,
	)

	// Setup repositories
	feedbackRepo := repository.NewGormFeedbackRepository(database)

	// Setup usecase
	feedbackUsecase := usecase.NewFeedbackUsecase(feedbackRepo, rabbitMQ, fileService, emailService)

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

	// Create upload directory if it doesn't exist
	if _, err := os.Stat(cfg.FileStore.UploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.FileStore.UploadDir, 0755); err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
	}

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
