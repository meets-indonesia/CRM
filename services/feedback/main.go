package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/config"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/delivery/http"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/infrastructure/auth"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/infrastructure/email"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/infrastructure/messagebroker"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/middleware"
	repopostgre "github.com/kevinnaserwan/crm-be/services/feedback/internal/repository/postgres"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
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
	err = db.AutoMigrate(&model.FeedbackType{}, &model.Station{}, &model.Feedback{}, &model.FeedbackResponse{})
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
	feedbackRepo := repopostgre.NewFeedbackRepository(db)
	feedbackResponseRepo := repopostgre.NewFeedbackResponseRepository(db)
	feedbackTypeRepo := repopostgre.NewFeedbackTypeRepository(db)
	stationRepo := repopostgre.NewStationRepository(db)

	// Initialize Use Cases
	feedbackUseCase := usecase.NewFeedbackUseCase(
		feedbackRepo,
		feedbackResponseRepo,
		feedbackTypeRepo,
		stationRepo,
		rabbitMQ,
	)

	// Initialize auth service
	authService := auth.NewAuthService(cfg.AuthServiceURL)

	// Initialize handlers with auth service
	feedbackHandler := http.NewFeedbackHandler(feedbackUseCase, authService)

	// Initialize email service
	emailConfig := email.EmailConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
	}
	emailService := email.NewEmailService(emailConfig)

	// Start RabbitMQ consumers
	go rabbitMQ.ConsumeFeedbackEvents(emailService)

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
	router.POST("/feedback", middleware.AuthMiddleware(cfg.JWTSecret), feedbackHandler.Handle)

	// Admin routes
	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	adminRoutes.Use(middleware.AdminMiddleware())
	{
		adminRoutes.POST("/feedback", feedbackHandler.Handle)
	}

	// Start server
	log.Printf("Feedback service starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
