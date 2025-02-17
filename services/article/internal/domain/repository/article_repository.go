package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/article/internal/domain/model"
)

type ArticleRepository interface {
	// Create a new article
	Create(ctx context.Context, article *model.Article) error

	// Get article by ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error)

	// List all articles
	List(ctx context.Context) ([]model.Article, error)

	// Update an article
	Update(ctx context.Context, article *model.Article) error

	// Delete an article
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByStatus gets articles with specific status
	ListByStatus(ctx context.Context, status int) ([]model.Article, error)
}
