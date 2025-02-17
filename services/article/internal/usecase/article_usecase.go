package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/article/internal/domain/events"
	"github.com/kevinnaserwan/crm-be/services/article/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/article/internal/domain/repository"
	"github.com/kevinnaserwan/crm-be/services/article/internal/infrastructure/messagebroker"
)

type ArticleUseCase struct {
	articleRepo   repository.ArticleRepository
	messageBroker *messagebroker.RabbitMQ
}

func NewArticleUseCase(
	articleRepo repository.ArticleRepository,
	messageBroker *messagebroker.RabbitMQ,
) *ArticleUseCase {
	return &ArticleUseCase{
		articleRepo:   articleRepo,
		messageBroker: messageBroker,
	}
}

func (uc *ArticleUseCase) CreateArticle(ctx context.Context, article *model.Article, userID string) error {
	// Set a new UUID if not provided
	if article.ID == uuid.Nil {
		article.ID = uuid.New()
	}

	// Save article to database
	if err := uc.articleRepo.Create(ctx, article); err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}

	// Create and publish event
	event := events.ArticleCreatedEvent{
		BaseEvent: events.BaseEvent{
			EventType:  events.EventTypeArticleCreated,
			OccurredAt: article.CreatedAt,
		},
		ArticleID:  article.ID.String(),
		Title:      article.Title,
		Content:    article.Content,
		Image:      article.Image,
		Status:     article.Status,
		MakingDate: article.MakingDate,
		CreatedBy:  userID,
	}

	// Publish event
	if err := uc.messageBroker.PublishMessage(
		ctx,
		messagebroker.ArticleExchange,
		events.EventTypeArticleCreated,
		event,
	); err != nil {
		log.Printf("Failed to publish article created event: %v", err)
	}

	return nil
}

func (uc *ArticleUseCase) GetArticle(ctx context.Context, id uuid.UUID) (*model.Article, error) {
	article, err := uc.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}
	return article, nil
}

func (uc *ArticleUseCase) ListArticles(ctx context.Context) ([]model.Article, error) {
	articles, err := uc.articleRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list articles: %w", err)
	}
	return articles, nil
}

func (uc *ArticleUseCase) UpdateArticle(ctx context.Context, article *model.Article, userID string) error {
	// Get existing article to track changes
	existing, err := uc.articleRepo.GetByID(ctx, article.ID)
	if err != nil {
		return fmt.Errorf("article not found: %w", err)
	}

	// Track which fields were updated
	var updatedFields []string
	if existing.Title != article.Title {
		updatedFields = append(updatedFields, "title")
	}
	if existing.Content != article.Content {
		updatedFields = append(updatedFields, "content")
	}
	if existing.Image != article.Image {
		updatedFields = append(updatedFields, "image")
	}
	if existing.Status != article.Status {
		updatedFields = append(updatedFields, "status")
	}
	if !existing.MakingDate.Equal(article.MakingDate) {
		updatedFields = append(updatedFields, "making_date")
	}

	// Update in database
	if err := uc.articleRepo.Update(ctx, article); err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	// Create and publish event
	event := events.ArticleUpdatedEvent{
		BaseEvent: events.BaseEvent{
			EventType:  events.EventTypeArticleUpdated,
			OccurredAt: article.UpdatedAt,
		},
		ArticleID:     article.ID.String(),
		Title:         article.Title,
		Content:       article.Content,
		Image:         article.Image,
		Status:        article.Status,
		MakingDate:    article.MakingDate,
		UpdatedBy:     userID,
		UpdatedFields: updatedFields,
	}

	// Publish event
	if err := uc.messageBroker.PublishMessage(
		ctx,
		messagebroker.ArticleExchange,
		events.EventTypeArticleUpdated,
		event,
	); err != nil {
		log.Printf("Failed to publish article updated event: %v", err)
	}

	return nil
}

func (uc *ArticleUseCase) DeleteArticle(ctx context.Context, id uuid.UUID, userID string) error {
	// Check if article exists
	_, err := uc.articleRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("article not found: %w", err)
	}

	// Delete from database
	if err := uc.articleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	// Create and publish event
	event := events.ArticleDeletedEvent{
		BaseEvent: events.BaseEvent{
			EventType:  events.EventTypeArticleDeleted,
			OccurredAt: time.Now(),
		},
		ArticleID: id.String(),
		DeletedBy: userID,
	}

	// Publish event
	if err := uc.messageBroker.PublishMessage(
		ctx,
		messagebroker.ArticleExchange,
		events.EventTypeArticleDeleted,
		event,
	); err != nil {
		log.Printf("Failed to publish article deleted event: %v", err)
	}

	return nil
}

func (uc *ArticleUseCase) ListArticlesByStatus(ctx context.Context, status int) ([]model.Article, error) {
	articles, err := uc.articleRepo.ListByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list articles by status: %w", err)
	}
	return articles, nil
}
