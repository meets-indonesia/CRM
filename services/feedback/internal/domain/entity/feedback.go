package model

import (
	"time"

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

type FeedbackType struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(255);column:name;not null"`
}

type Station struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(255);column:name;not null"`
}

type Feedback struct {
	gorm.Model
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID         uuid.UUID    `gorm:"type:uuid;column:user_id;not null"`
	FeedbackDate   time.Time    `gorm:"type:timestamp;column:feedback_date;not null"`
	FeedbackTypeID uuid.UUID    `gorm:"type:uuid;column:feedback_type_id;not null"`
	StationID      uuid.UUID    `gorm:"type:uuid;column:station_id;not null"`
	Feedback       string       `gorm:"type:text;column:feedback;not null"`
	Documentation  string       `gorm:"type:text;column:documentation;not null"`
	Rating         float64      `gorm:"type:int;column:rating;not null"`
	Status         int          `gorm:"type:int;column:status;not null"`
	User           User         `gorm:"foreignKey:UserID;references:ID"`
	FeedbackType   FeedbackType `gorm:"foreignKey:FeedbackTypeID;references:ID"`
	Station        Station      `gorm:"foreignKey:StationID;references:ID"`
}
