package repository

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/article/domain/entity"
	"gorm.io/gorm"
)

// GormArticleRepository implements ArticleRepository with GORM
type GormArticleRepository struct {
	db *gorm.DB
}

// NewGormArticleRepository creates a new GormArticleRepository
func NewGormArticleRepository(db *gorm.DB) *GormArticleRepository {
	return &GormArticleRepository{
		db: db,
	}
}

// Create creates a new article
func (r *GormArticleRepository) Create(ctx context.Context, article *entity.Article) error {
	result := r.db.Create(article)
	return result.Error
}

// FindByID finds an article by ID
func (r *GormArticleRepository) FindByID(ctx context.Context, id uint) (*entity.Article, error) {
	var article entity.Article
	result := r.db.First(&article, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &article, nil
}

// Update updates an article
func (r *GormArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	result := r.db.Save(article)
	return result.Error
}

// Delete deletes an article
func (r *GormArticleRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.Delete(&entity.Article{}, id)
	return result.Error
}

// List lists articles
func (r *GormArticleRepository) List(ctx context.Context, publishedOnly bool, page, limit int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	offset := (page - 1) * limit
	query := r.db

	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}

	// Count total
	if err := query.Model(&entity.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get articles
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// IncrementViewCount increments the view count of an article
func (r *GormArticleRepository) IncrementViewCount(ctx context.Context, id uint) error {
	result := r.db.Model(&entity.Article{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + ?", 1))
	return result.Error
}

// Search searches articles by title, summary, or content
func (r *GormArticleRepository) Search(ctx context.Context, query string, publishedOnly bool, page, limit int) ([]entity.Article, int64, error) {
	var articles []entity.Article
	var total int64

	offset := (page - 1) * limit
	dbQuery := r.db.Where("title ILIKE ? OR summary ILIKE ? OR content ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%")

	if publishedOnly {
		dbQuery = dbQuery.Where("is_published = ?", true)
	}

	// Count total
	if err := dbQuery.Model(&entity.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get articles
	if err := dbQuery.Order("created_at DESC").Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
