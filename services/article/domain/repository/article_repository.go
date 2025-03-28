package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/article/domain/entity"
)

// ArticleRepository mendefinisikan operasi-operasi repository untuk Article
type ArticleRepository interface {
	Create(ctx context.Context, article *entity.Article) error
	FindByID(ctx context.Context, id uint) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, publishedOnly bool, page, limit int) ([]entity.Article, int64, error)
	IncrementViewCount(ctx context.Context, id uint) error
	Search(ctx context.Context, query string, publishedOnly bool, page, limit int) ([]entity.Article, int64, error)
	PartialUpdate(ctx context.Context, id uint, data map[string]interface{}) error
}

// EventPublisher mendefinisikan operasi-operasi untuk mempublikasikan event
type EventPublisher interface {
	PublishArticleCreated(article *entity.Article) error
	PublishArticleUpdated(article *entity.Article) error
	Close() error
}
