package postgres

import (
	model "github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.FeedbackType{},
		&model.Station{},
		&model.Feedback{},
		&model.FeedbackResponse{},
	)
	if err != nil {
		panic(err)
	}
}
