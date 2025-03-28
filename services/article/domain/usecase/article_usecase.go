package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/kevinnaserwan/crm-be/services/article/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/article/domain/repository"
)

// Errors
var (
	ErrArticleNotFound = errors.New("article not found")
)

// ArticleUsecase mendefinisikan operasi-operasi usecase untuk Article
type ArticleUsecase interface {
	CreateArticle(ctx context.Context, authorID uint, req entity.CreateArticleRequest) (*entity.Article, error)
	GetArticle(ctx context.Context, id uint) (*entity.Article, error)
	UpdateArticle(ctx context.Context, id uint, req entity.UpdateArticleRequest) (*entity.Article, error)
	DeleteArticle(ctx context.Context, id uint) error
	ListArticles(ctx context.Context, publishedOnly bool, page, limit int) (*entity.ArticleListResponse, error)
	ViewArticle(ctx context.Context, id uint) (*entity.Article, error)
	SearchArticles(ctx context.Context, query string, publishedOnly bool, page, limit int) (*entity.ArticleListResponse, error)
}

type articleUsecase struct {
	articleRepo    repository.ArticleRepository
	eventPublisher repository.EventPublisher
}

// NewArticleUsecase membuat instance baru ArticleUsecase
func NewArticleUsecase(
	articleRepo repository.ArticleRepository,
	eventPublisher repository.EventPublisher,
) ArticleUsecase {
	return &articleUsecase{
		articleRepo:    articleRepo,
		eventPublisher: eventPublisher,
	}
}

// CreateArticle membuat artikel baru
func (u *articleUsecase) CreateArticle(ctx context.Context, authorID uint, req entity.CreateArticleRequest) (*entity.Article, error) {
	article := &entity.Article{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		ImageURL:    req.ImageURL,
		AuthorID:    authorID,
		IsPublished: req.IsPublished,
	}

	if article.IsPublished {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := u.articleRepo.Create(ctx, article); err != nil {
		return nil, err
	}

	// Publish event jika artikel dipublish
	if article.IsPublished {
		if err := u.eventPublisher.PublishArticleCreated(article); err != nil {
			// Log error but don't fail
		}
	}

	return article, nil
}

// GetArticle mendapatkan artikel berdasarkan ID
func (u *articleUsecase) GetArticle(ctx context.Context, id uint) (*entity.Article, error) {
	article, err := u.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, ErrArticleNotFound
	}

	return article, nil
}

// UpdateArticle memperbarui artikel
func (u *articleUsecase) UpdateArticle(ctx context.Context, id uint, req entity.UpdateArticleRequest) (*entity.Article, error) {
	article, err := u.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, ErrArticleNotFound
	}

	updateData := make(map[string]interface{})

	if req.Title != nil {
		updateData["title"] = *req.Title
	}
	if req.Content != nil {
		updateData["content"] = *req.Content
	}
	if req.Summary != nil {
		updateData["summary"] = *req.Summary
	}
	if req.ImageURL != nil {
		updateData["image_url"] = *req.ImageURL
	}
	if req.IsPublished != nil {
		updateData["is_published"] = *req.IsPublished
		if *req.IsPublished {
			now := time.Now()
			updateData["published_at"] = &now
		} else {
			updateData["published_at"] = nil
		}
	}

	if len(updateData) == 0 {
		return article, nil // Tidak ada yang perlu di-update
	}

	if err := u.articleRepo.PartialUpdate(ctx, id, updateData); err != nil {
		return nil, err
	}

	return u.articleRepo.FindByID(ctx, id)
}

// DeleteArticle menghapus artikel
func (u *articleUsecase) DeleteArticle(ctx context.Context, id uint) error {
	article, err := u.articleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if article == nil {
		return ErrArticleNotFound
	}

	return u.articleRepo.Delete(ctx, id)
}

// ListArticles mendapatkan daftar artikel
func (u *articleUsecase) ListArticles(ctx context.Context, publishedOnly bool, page, limit int) (*entity.ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	articles, total, err := u.articleRepo.List(ctx, publishedOnly, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

// ViewArticle membuka artikel dan menambah view count
func (u *articleUsecase) ViewArticle(ctx context.Context, id uint) (*entity.Article, error) {
	article, err := u.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, ErrArticleNotFound
	}

	// Hanya menambah view count jika artikel dipublish
	if article.IsPublished {
		if err := u.articleRepo.IncrementViewCount(ctx, id); err != nil {
			return nil, err
		}

		// Reload article untuk mendapatkan view count terbaru
		article, err = u.articleRepo.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return article, nil
}

// SearchArticles mencari artikel berdasarkan kata kunci
func (u *articleUsecase) SearchArticles(ctx context.Context, query string, publishedOnly bool, page, limit int) (*entity.ArticleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	articles, total, err := u.articleRepo.Search(ctx, query, publishedOnly, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.ArticleListResponse{
		Articles: articles,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}
