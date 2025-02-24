package Postgres

import (
	model "github.com/kevinnaserwan/crm-be/services/user/internal/domain/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.UserPoints{},
		&model.PointHistory{},
		&model.PointResetHistory{},
	)
	if err != nil {
		panic(err)
	}
}
