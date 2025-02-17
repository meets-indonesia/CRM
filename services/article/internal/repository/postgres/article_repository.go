package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/article/internal/domain/model"
	"gorm.io/gorm"
)

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *articleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) Create(ctx context.Context, article *model.Article) error {
	if article.ID == uuid.Nil {
		article.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) List(ctx context.Context) ([]model.Article, error) {
	var articles []model.Article
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("making_date DESC").
		Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *articleRepository) Update(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *articleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&model.Article{}, "id = ?", id).Error
}

func (r *articleRepository) ListByStatus(ctx context.Context, status int) ([]model.Article, error) {
	var articles []model.Article
	err := r.db.WithContext(ctx).
		Where("status = ? AND deleted_at IS NULL", status).
		Order("making_date DESC").
		Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}
