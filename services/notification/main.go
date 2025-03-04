package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/notification/config"
	"github.com/kevinnaserwan/crm-be/services/notification/domain/usecase"
	"github.com/kevinnaserwan/crm-be/services/notification/infrastructure/db"
	"github.com/kevinnaserwan/crm-be/services/notification/infrastructure/email"
	"github.com/kevinnaserwan/crm-be/services/notification/infrastructure/messaging"
	"github.com/kevinnaserwan/crm-be/services/notification/infrastructure/push"
	"github.com/kevinnaserwan/crm-be/services/notification/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/notification/interface/router"
	"github.com/kevinnaserwan/crm-be/services/notification/repository"
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

	// Setup repositories
	notificationRepo := repository.NewGormNotificationRepository(database)

	// Setup email sender
	emailSender := email.NewSMTPEmailSender(cfg.SMTP)

	// Setup push notification sender
	pushSender, err := push.NewFirebasePushSender(cfg.Push)
	if err != nil {
		log.Fatalf("Failed to create push notification sender: %v", err)
	}

	// Setup usecase
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo, emailSender, pushSender)

	// Setup RabbitMQ
	rabbitMQ, err := messaging.NewRabbitMQ(cfg.RabbitMQ, notificationUsecase)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Subscribe to events
	if err := rabbitMQ.SubscribeToEvents(); err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}

	// Setup handler
	notificationHandler := handler.NewNotificationHandler(notificationUsecase)

	// Setup background processing of pending notifications
	go func() {
		for {
			// Process pending notifications every minute
			time.Sleep(1 * time.Minute)
			if err := notificationUsecase.ProcessPendingNotifications(context.Background()); err != nil {
				log.Printf("Error processing pending notifications: %v", err)
			}
		}
	}()

	// Add config to Gin context
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Setup router
	router := router.Setup(cfg, notificationHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Notification Service starting on port %s", port)
	router.Run(":" + port)
}
