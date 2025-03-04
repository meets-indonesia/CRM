package db

import (
	"fmt"

	"github.com/kevinnaserwan/crm-be/services/user/config"
	"github.com/kevinnaserwan/crm-be/services/user/domain/entity"
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
	if err := db.AutoMigrate(&entity.User{}, &entity.PointTransaction{}, &entity.PointLevel{}); err != nil {
		return nil, err
	}

	// Insert default point levels if they don't exist
	var count int64
	db.Model(&entity.PointLevel{}).Count(&count)
	if count == 0 {
		levels := []entity.PointLevel{
			{Name: "Bronze", MinPoints: 0, Multiplier: 1.0},
			{Name: "Silver", MinPoints: 50, Multiplier: 1.2},
			{Name: "Gold", MinPoints: 100, Multiplier: 1.5},
			{Name: "Platinum", MinPoints: 200, Multiplier: 2.0},
		}

		for _, level := range levels {
			db.Create(&level)
		}
	}

	return db, nil
}
