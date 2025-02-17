package events

import (
	"time"
)

// Event types constants
const (
	EventTypeArticleCreated = "article.created"
	EventTypeArticleUpdated = "article.updated"
	EventTypeArticleDeleted = "article.deleted"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
	EventType  string    `json:"event_type"`
	OccurredAt time.Time `json:"occurred_at"`
}

// ArticleCreatedEvent is published when a new article is created
type ArticleCreatedEvent struct {
	BaseEvent
	ArticleID  string    `json:"article_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Image      string    `json:"image"`
	Status     int       `json:"status"`
	MakingDate time.Time `json:"making_date"`
	CreatedBy  string    `json:"created_by"`
}

// ArticleUpdatedEvent is published when an article is updated
type ArticleUpdatedEvent struct {
	BaseEvent
	ArticleID     string    `json:"article_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Image         string    `json:"image"`
	Status        int       `json:"status"`
	MakingDate    time.Time `json:"making_date"`
	UpdatedBy     string    `json:"updated_by"`
	UpdatedFields []string  `json:"updated_fields"`
}

// ArticleDeletedEvent is published when an article is deleted
type ArticleDeletedEvent struct {
	BaseEvent
	ArticleID string `json:"article_id"`
	DeletedBy string `json:"deleted_by"`
}
