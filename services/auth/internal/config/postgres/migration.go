package postgres

import (
	model "github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
	)
	if err != nil {
		panic(err)
	}
}
