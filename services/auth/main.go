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
	"github.com/kevinnaserwan/crm-be/services/auth/internal/middleware"
	repopostgre "github.com/kevinnaserwan/crm-be/services/auth/internal/repository/postgres"
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

	// Auto Migrate
	err = db.AutoMigrate(&model.Role{}, &model.User{}, &model.OTP{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Repositories
	userRepo := repopostgre.NewUserRepository(db)
	roleRepo := repopostgre.NewRoleRepository(db)
	otpRepo := repopostgre.NewOTPRepository(db)

	// Initialize Email Service
	emailService := util.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
	)

	// Initialize Use Cases
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	roleUseCase := usecase.NewRoleUseCase(roleRepo)
	forgotPasswordUseCase := usecase.NewForgotPasswordUseCase(userRepo, otpRepo, emailService)

	// Initialize Handlers
	authHandler := http.NewAuthHandler(authUseCase)
	roleHandler := http.NewRoleHandler(roleUseCase)
	forgotPasswordHandler := http.NewForgotPasswordHandler(forgotPasswordUseCase)

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

	// Public routes
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/roles", roleHandler.CreateRole)

	// Forgot password routes
	router.POST("/auth/forgot-password", forgotPasswordHandler.RequestPasswordReset)
	router.POST("/auth/verify-otp", forgotPasswordHandler.VerifyOTP)
	router.POST("/auth/reset-password", forgotPasswordHandler.ResetPassword)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Add protected routes here
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
