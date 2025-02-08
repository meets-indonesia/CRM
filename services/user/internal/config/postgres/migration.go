package postgres

import (
	model "github.com/kevinnaserwan/crm-be/services/user/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Point{},
	)
	if err != nil {
		panic(err)
	}
}
