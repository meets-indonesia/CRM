package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/config"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/delivery/http"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/middleware"
	repopostgre "github.com/kevinnaserwan/crm-be/services/auth/internal/repository/postgres"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/usecase"
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
	err = db.AutoMigrate(&model.Role{}, &model.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Repositories
	userRepo := repopostgre.NewUserRepository(db)
	roleRepo := repopostgre.NewRoleRepository(db)

	// Initialize Use Cases
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	roleUseCase := usecase.NewRoleUseCase(roleRepo)

	// Initialize Handlers
	authHandler := http.NewAuthHandler(authUseCase)
	roleHandler := http.NewRoleHandler(roleUseCase)

	// Setup Router
	router := gin.Default()

	// Public routes
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/roles", roleHandler.CreateRole) // New endpoint for creating roles

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
