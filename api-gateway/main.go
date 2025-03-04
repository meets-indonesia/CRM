package main

import (
	"log"
	"os"

	"github.com/kevinnaserwan/crm-be/api-gateway/config"
	"github.com/kevinnaserwan/crm-be/api-gateway/router"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize router
	r := router.Setup(cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	r.Run(":" + port)
}
