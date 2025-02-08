package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	FirstName      string    `gorm:"type:varchar(255);column:first_name;not null"`
	LastName       string    `gorm:"type:varchar(255);column:last_name;not null"`
	ProfilePicture string    `gorm:"type:varchar(255);column:profile_picture"`
	Email          string    `gorm:"type:varchar(255);column:email;unique;not null"`
	Password       string    `gorm:"type:varchar(255);column:password;not null"`
	RoleID         uuid.UUID `gorm:"type:uuid;column:role_id;not null"`
	Role           Role      `gorm:"foreignKey:RoleID;references:ID"`
}

type Role struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(255);column:name;not null"`
}

type Point struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;column:user_id;not null"`
	Point  int       `gorm:"type:int;column:point;not null"`
	Tier   int       `gorm:"type:int;column:tier;not null"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`
}
