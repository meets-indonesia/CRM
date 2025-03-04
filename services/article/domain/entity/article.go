package entity

import (
	"time"
)

// Article adalah entitas untuk artikel/berita
type Article struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Content     string     `json:"content" gorm:"not null"`
	Summary     string     `json:"summary" gorm:"not null"`
	ImageURL    string     `json:"image_url,omitempty"`
	AuthorID    uint       `json:"author_id" gorm:"index;not null"` // ID admin yang membuat
	IsPublished bool       `json:"is_published" gorm:"not null;default:false"`
	ViewCount   int        `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// CreateArticleRequest adalah model untuk permintaan pembuatan artikel
type CreateArticleRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Summary     string `json:"summary" binding:"required"`
	ImageURL    string `json:"image_url,omitempty"`
	IsPublished bool   `json:"is_published"`
}

// UpdateArticleRequest adalah model untuk permintaan update artikel
type UpdateArticleRequest struct {
	Title       string `json:"title,omitempty"`
	Content     string `json:"content,omitempty"`
	Summary     string `json:"summary,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	IsPublished *bool  `json:"is_published,omitempty"`
}

// ArticleListResponse adalah model respons untuk daftar artikel
type ArticleListResponse struct {
	Articles []Article `json:"articles"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}
