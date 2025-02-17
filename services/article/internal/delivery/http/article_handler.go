// services/article/internal/delivery/http/article_handler.go
package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/article/internal/delivery/http/models"
	"github.com/kevinnaserwan/crm-be/services/article/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/article/internal/usecase"
)

type ArticleHandler struct {
	articleUseCase *usecase.ArticleUseCase
}

func NewArticleHandler(articleUseCase *usecase.ArticleUseCase) *ArticleHandler {
	return &ArticleHandler{
		articleUseCase: articleUseCase,
	}
}

func (h *ArticleHandler) Handle(c *gin.Context) {
	var baseRequest models.BaseRequest
	if err := c.ShouldBindJSON(&baseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch baseRequest.Meta.Action {
	case "create":
		h.handleCreate(c, baseRequest.Data)
	case "get":
		h.handleGet(c, baseRequest.Data)
	case "list":
		h.handleList(c)
	case "list_by_status":
		h.handleListByStatus(c, baseRequest.Data)
	case "update":
		h.handleUpdate(c, baseRequest.Data)
	case "delete":
		h.handleDelete(c, baseRequest.Data)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
	}
}

func (h *ArticleHandler) handleCreate(c *gin.Context, data interface{}) {
	var articleData models.CreateArticleData
	if err := convertData(data, &articleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	article := &model.Article{
		Title:      articleData.Title,
		Content:    articleData.Content,
		Image:      articleData.Image,
		Status:     articleData.Status,
		MakingDate: articleData.MakingDate,
	}

	if err := h.articleUseCase.CreateArticle(c.Request.Context(), article, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Article created successfully",
		"article": article,
	})
}

func (h *ArticleHandler) handleGet(c *gin.Context, data interface{}) {
	var getData models.GetArticleData
	if err := convertData(data, &getData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := h.articleUseCase.GetArticle(c.Request.Context(), getData.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

func (h *ArticleHandler) handleList(c *gin.Context) {
	articles, err := h.articleUseCase.ListArticles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

func (h *ArticleHandler) handleListByStatus(c *gin.Context, data interface{}) {
	var statusData models.ListByStatusData
	if err := convertData(data, &statusData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articles, err := h.articleUseCase.ListArticlesByStatus(c.Request.Context(), statusData.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

func (h *ArticleHandler) handleUpdate(c *gin.Context, data interface{}) {
	var updateData models.UpdateArticleData
	if err := convertData(data, &updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	article := &model.Article{
		ID:         updateData.ID,
		Title:      updateData.Title,
		Content:    updateData.Content,
		Image:      updateData.Image,
		Status:     updateData.Status,
		MakingDate: updateData.MakingDate,
	}

	if err := h.articleUseCase.UpdateArticle(c.Request.Context(), article, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
		"article": article,
	})
}

func (h *ArticleHandler) handleDelete(c *gin.Context, data interface{}) {
	var deleteData models.DeleteArticleData
	if err := convertData(data, &deleteData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := h.articleUseCase.DeleteArticle(c.Request.Context(), deleteData.ID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func convertData(data interface{}, target interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, target)
}
