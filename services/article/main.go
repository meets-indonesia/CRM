// services/article/main.go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/article/internal/config"
	"github.com/kevinnaserwan/crm-be/services/article/internal/delivery/http"
	"github.com/kevinnaserwan/crm-be/services/article/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/article/internal/infrastructure/messagebroker"
	"github.com/kevinnaserwan/crm-be/services/article/internal/middleware"
	repopostgre "github.com/kevinnaserwan/crm-be/services/article/internal/repository/postgres"
	"github.com/kevinnaserwan/crm-be/services/article/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize Database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize RabbitMQ
	rabbitMQ, err := messagebroker.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitMQ.Close()

	// Auto Migrate
	err = db.AutoMigrate(&model.Article{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Repositories
	articleRepo := repopostgre.NewArticleRepository(db)

	// Initialize Use Cases
	articleUseCase := usecase.NewArticleUseCase(articleRepo, rabbitMQ)

	// Initialize Handlers
	articleHandler := http.NewArticleHandler(articleUseCase)

	// Setup Router
	router := gin.Default()

	// CORS Configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes for mobile app (read-only)
	router.POST("/articles", articleHandler.Handle)

	// Protected routes for admin
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.POST("/admin/articles", articleHandler.Handle)
	}

	// Start server
	log.Printf("Article service starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
