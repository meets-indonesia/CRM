// services/article/internal/delivery/http/models/request.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Meta struct {
	Action string `json:"action" binding:"required"`
}

type BaseRequest struct {
	Meta Meta        `json:"meta" binding:"required"`
	Data interface{} `json:"data" binding:"required"`
}

type CreateArticleData struct {
	Title      string    `json:"title" binding:"required"`
	Content    string    `json:"content" binding:"required"`
	Image      string    `json:"image"`
	Status     int       `json:"status" binding:"required"`
	MakingDate time.Time `json:"making_date" binding:"required"`
}

type GetArticleData struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type UpdateArticleData struct {
	ID         uuid.UUID `json:"id" binding:"required"`
	Title      string    `json:"title" binding:"required"`
	Content    string    `json:"content" binding:"required"`
	Image      string    `json:"image"`
	Status     int       `json:"status" binding:"required"`
	MakingDate time.Time `json:"making_date" binding:"required"`
}

type DeleteArticleData struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type ListByStatusData struct {
	Status int `json:"status" binding:"required"`
}
