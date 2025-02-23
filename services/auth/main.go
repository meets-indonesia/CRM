package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/config"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/delivery/http"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/infrastructure/messagebroker"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/middleware"
	repopostgre "github.com/kevinnaserwan/crm-be/services/auth/internal/repository/postgres"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/service"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/usecase"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/util"
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

	// Initialize RabbitMQ
	rabbitMQ, err := messagebroker.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitMQ.Close()

	// Auto Migrate
	err = db.AutoMigrate(&model.Role{}, &model.User{}, &model.OAuthAccount{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Repositories
	userRepo := repopostgre.NewUserRepository(db)
	roleRepo := repopostgre.NewRoleRepository(db)
	otpRepo := repopostgre.NewOTPRepository(db)
	oauthRepo := repopostgre.NewOAuthRepository(db) // Add this

	// Initialize Email Service
	emailService := util.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
	)

	oauthService := service.NewOAuthService(cfg.OAuth) // Add this

	// Initialize Use Cases with RabbitMQ
	authUseCase := usecase.NewAuthUseCase(userRepo, roleRepo, cfg.JWTSecret, cfg.JWTExpiration, rabbitMQ)
	roleUseCase := usecase.NewRoleUseCase(roleRepo)
	forgotPasswordUseCase := usecase.NewForgotPasswordUseCase(userRepo, otpRepo, emailService, rabbitMQ)
	oauthUseCase := usecase.NewOAuthUseCase( // Add this
		userRepo,
		oauthRepo,
		roleRepo,
		oauthService,
		cfg.JWTSecret,
		cfg.JWTExpiration,
	)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize Handlers
	authHandler := http.NewAuthHandler(authUseCase, forgotPasswordUseCase)
	roleHandler := http.NewRoleHandler(roleUseCase)
	oauthHandler := http.NewOAuthHandler(oauthUseCase)
	userHandler := http.NewUserHandler(userUseCase)

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

	// Single endpoint for all auth operations
	router.POST("/auth", authHandler.Handle)

	// Role endpoint remains separate
	router.POST("/roles", roleHandler.CreateRole)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/users/:id", userHandler.GetUserByID)
	}

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// OAuth routes for mobile
	router.POST("/oauth", oauthHandler.Handle)                                   // For meta-based approach
	router.GET("/oauth/google/callback", oauthHandler.GoogleCallback)            // Web flow
	router.GET("/oauth/google/callback/mobile", oauthHandler.GoogleCallbackJSON) // Mobile-friendly JSON response

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
