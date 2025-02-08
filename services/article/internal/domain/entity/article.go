package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Title      string    `gorm:"type:varchar(255);column:title;not null"`
	Content    string    `gorm:"type:text;column:content;not null"`
	Image      string    `gorm:"type:varchar(255);column:image"`
	Status     int       `gorm:"type:int;column:status;not null"`
	MakingDate time.Time `gorm:"type:timestamp;column:making_date;not null"`
}
