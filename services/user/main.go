package main

import (
	"log"
	"os"

	"github.com/kevinnaserwan/crm-be/services/user/config"
	"github.com/kevinnaserwan/crm-be/services/user/domain/usecase"
	"github.com/kevinnaserwan/crm-be/services/user/infrastructure/db"
	"github.com/kevinnaserwan/crm-be/services/user/infrastructure/messaging"
	"github.com/kevinnaserwan/crm-be/services/user/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/user/interface/router"
	"github.com/kevinnaserwan/crm-be/services/user/repository"
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
	userRepo := repository.NewGormUserRepository(database)
	pointRepo := repository.NewGormPointRepository(database)

	// Setup usecase
	userUsecase := usecase.NewUserUsecase(userRepo, pointRepo)

	// Setup RabbitMQ
	rabbitMQ, err := messaging.NewRabbitMQ(cfg.RabbitMQ, userUsecase, userUsecase)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Subscribe to events
	if err := rabbitMQ.SubscribeToUserEvents(); err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}

	// Setup handler
	userHandler := handler.NewUserHandler(userUsecase)

	// Setup router
	r := router.Setup(cfg, userHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("User Service starting on port %s", port)
	r.Run(":" + port)
}
