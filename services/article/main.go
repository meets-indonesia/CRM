package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/article/config"
	"github.com/kevinnaserwan/crm-be/services/article/domain/usecase"
	"github.com/kevinnaserwan/crm-be/services/article/infrastructure/db"
	"github.com/kevinnaserwan/crm-be/services/article/infrastructure/messaging"
	"github.com/kevinnaserwan/crm-be/services/article/interface/handler"
	"github.com/kevinnaserwan/crm-be/services/article/interface/router"
	"github.com/kevinnaserwan/crm-be/services/article/repository"
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
	articleRepo := repository.NewGormArticleRepository(database)

	// Setup usecase
	articleUsecase := usecase.NewArticleUsecase(articleRepo, rabbitMQ)

	// Setup handler
	articleHandler := handler.NewArticleHandler(articleUsecase)

	// Add config to Gin context
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// Setup router
	router := router.Setup(cfg, articleHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Article Service starting on port %s", port)
	router.Run(":" + port)
}
