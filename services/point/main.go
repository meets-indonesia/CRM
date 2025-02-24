package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/user/internal/config"
	"github.com/kevinnaserwan/crm-be/services/user/internal/delivery/http"
	"github.com/kevinnaserwan/crm-be/services/user/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/user/internal/infrastructure/messagebroker"
	"github.com/kevinnaserwan/crm-be/services/user/internal/middleware"
	repopostgre "github.com/kevinnaserwan/crm-be/services/user/internal/repository/postgres"
	"github.com/kevinnaserwan/crm-be/services/user/internal/usecase"
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

	// Auto Migrate
	err = db.AutoMigrate(&model.UserPoints{}, &model.PointHistory{}, &model.PointResetHistory{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize RabbitMQ
	rabbitMQ, err := messagebroker.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitMQ.Close()

	// Initialize Repositories
	userPointsRepo := repopostgre.NewUserPointsRepository(db)
	pointHistoryRepo := repopostgre.NewPointHistoryRepository(db)
	pointResetRepo := repopostgre.NewPointResetHistoryRepository(db)

	// Initialize Use Cases
	pointUseCase := usecase.NewPointUseCase(
		userPointsRepo,
		pointHistoryRepo,
		pointResetRepo,
	)

	// Initialize handlers
	pointHandler := http.NewPointHandler(pointUseCase) // No auth service needed

	// Setup RabbitMQ consumer
	go rabbitMQ.ConsumeFeedbackEvents(pointUseCase)

	// Process existing feedback (if needed)
	go func() {
		// Sleep to ensure service is fully started
		time.Sleep(5 * time.Second)

		// Process any existing feedbacks for users
		ctx := context.Background()
		users, err := userPointsRepo.List(ctx)
		if err != nil {
			log.Printf("Error getting users: %v", err)
			return
		}

		log.Printf("Found %d users with point records", len(users))
	}()

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

	// Routes
	router.POST("/points", middleware.AuthMiddleware(cfg.JWTSecret), pointHandler.Handle)

	// Admin routes
	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	adminRoutes.Use(middleware.AdminMiddleware())
	{
		adminRoutes.POST("/points", pointHandler.Handle)
	}

	// Start server
	log.Printf("Point service starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
