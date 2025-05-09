package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/article/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/article/domain/usecase"
)

// ArticleHandler handles article requests
type ArticleHandler struct {
	articleUsecase usecase.ArticleUsecase
}

// NewArticleHandler creates a new ArticleHandler
func NewArticleHandler(articleUsecase usecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		articleUsecase: articleUsecase,
	}
}

// CreateArticle handles create article requests
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Multipart form
	title := c.PostForm("title")
	content := c.PostForm("content")
	summary := c.PostForm("summary")
	isPublished, _ := strconv.ParseBool(c.PostForm("is_published"))

	// File upload
	file, err := c.FormFile("image")
	var imageURL string
	if err == nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		imageURL = "/" + path
	}

	req := entity.CreateArticleRequest{
		Title:       title,
		Content:     content,
		Summary:     summary,
		ImageURL:    imageURL,
		IsPublished: isPublished,
	}

	article, err := h.articleUsecase.CreateArticle(c, authorID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

// GetArticle handles get article by ID requests
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	article, err := h.articleUsecase.GetArticle(c, uint(id))
	if err != nil {
		if err == usecase.ErrArticleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

// UpdateArticle handles update article requests
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")
	summary := c.PostForm("summary")
	isPublishedStr := c.PostForm("is_published")
	var isPublished *bool
	if isPublishedStr != "" {
		val, err := strconv.ParseBool(isPublishedStr)
		if err == nil {
			isPublished = &val
		}
	}

	file, err := c.FormFile("image")
	var imageURL *string
	if err == nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		fullPath := "/" + path
		imageURL = &fullPath
	}

	var req entity.UpdateArticleRequest
	if title != "" {
		req.Title = &title
	}
	if content != "" {
		req.Content = &content
	}
	if summary != "" {
		req.Summary = &summary
	}
	if imageURL != nil {
		req.ImageURL = imageURL
	}
	if isPublished != nil {
		req.IsPublished = isPublished
	}

	article, err := h.articleUsecase.UpdateArticle(c, uint(id), req)
	if err != nil {
		if err == usecase.ErrArticleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

// DeleteArticle handles delete article requests
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	err = h.articleUsecase.DeleteArticle(c, uint(id))
	if err != nil {
		if err == usecase.ErrArticleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

// ListArticles handles list articles requests
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Check if user is admin
	role, exists := c.Get("role")
	publishedOnly := true

	if exists && role == "ADMIN" {
		// Admin can see both published and unpublished articles
		publishedOnly, _ = strconv.ParseBool(c.DefaultQuery("published_only", "false"))
	}

	articles, err := h.articleUsecase.ListArticles(c, publishedOnly, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// ViewArticle handles viewing an article and incrementing view count
func (h *ArticleHandler) ViewArticle(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	article, err := h.articleUsecase.ViewArticle(c, uint(id))
	if err != nil {
		if err == usecase.ErrArticleNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

// SearchArticles handles searching articles
func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	// Parse query parameters
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Check if user is admin
	role, exists := c.Get("role")
	publishedOnly := true

	if exists && role == "ADMIN" {
		// Admin can see both published and unpublished articles
		publishedOnly, _ = strconv.ParseBool(c.DefaultQuery("published_only", "false"))
	}

	articles, err := h.articleUsecase.SearchArticles(c, query, publishedOnly, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// HealthCheck handles health check requests
func (h *ArticleHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
