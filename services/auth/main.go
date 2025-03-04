package main

import (
	"log"
	"os"

	"github.com/kevinnaserwan/crm-be/services/auth/config"
	"github.com/kevinnaserwan/crm-be/services/auth/domain/usecase"
	"github.com/kevinnaserwan/crm-be/services/auth/infrastucture/db"
	"github.com/kevinnaserwan/crm-be/services/auth/infrastucture/messaging"
	"github.com/kevinnaserwan/crm-be/services/auth/interface/router"
	"github.com/kevinnaserwan/crm-be/services/auth/repository"
	serviceoauth "github.com/kevinnaserwan/crm-be/services/auth/service/google"
	servicejwt "github.com/kevinnaserwan/crm-be/services/auth/service/jwt"
	serviceotp "github.com/kevinnaserwan/crm-be/services/auth/service/otp"
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

	// Setup repositories and services
	userRepo := repository.NewGormUserRepository(database)
	jwtService := servicejwt.NewJWTService(cfg.JWT)
	otpService := serviceotp.NewOTPService(cfg.SMTP)
	googleAuthService := serviceoauth.NewGoogleAuthService(cfg.GoogleAuth)

	// Setup usecase
	authUsecase := usecase.NewAuthUsecase(
		userRepo,
		jwtService,
		otpService,
		rabbitMQ,
		otpService,
		googleAuthService,
	)

	// Setup router
	r := router.Setup(cfg, authUsecase)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Auth Service starting on port %s", port)
	r.Run(":" + port)
}
