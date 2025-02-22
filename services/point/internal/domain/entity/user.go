package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Point struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;column:user_id;not null"`
	Point  int       `gorm:"type:int;column:point;not null"`
	Tier   int       `gorm:"type:int;column:tier;not null"`
}
