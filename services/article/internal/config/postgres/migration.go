package postgres

import (
	model "github.com/kevinnaserwan/crm-be/services/article/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.Article{},
	)
	if err != nil {
		panic(err)
	}
}
