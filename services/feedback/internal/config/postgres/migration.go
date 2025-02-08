package postgres

import (
	model "github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.FeedbackType{},
		&model.Station{},
		&model.Feedback{},
	)
	if err != nil {
		panic(err)
	}
}
