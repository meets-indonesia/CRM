package db

import (
	"fmt"

	"github.com/kevinnaserwan/crm-be/services/inventory/config"
	"github.com/kevinnaserwan/crm-be/services/inventory/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB creates a new PostgreSQL connection
func NewPostgresDB(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&entity.Item{}, &entity.StockTransaction{}, &entity.Supplier{}); err != nil {
		return nil, err
	}

	return db, nil
}
