package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/feedback/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/feedback/domain/usecase"
)

// FeedbackHandler handles feedback requests
type FeedbackHandler struct {
	feedbackUsecase usecase.FeedbackUsecase
}

// NewFeedbackHandler creates a new FeedbackHandler
func NewFeedbackHandler(feedbackUsecase usecase.FeedbackUsecase) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackUsecase: feedbackUsecase,
	}
}

// CreateFeedback handles create feedback requests
func (h *FeedbackHandler) CreateFeedback(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req entity.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedback, err := h.feedbackUsecase.CreateFeedback(c, userID.(uint), req)
	if err != nil {
		if err == usecase.ErrInvalidUserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

// GetFeedback handles get feedback by ID requests
func (h *FeedbackHandler) GetFeedback(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}

	feedback, err := h.feedbackUsecase.GetFeedback(c, uint(id))
	if err != nil {
		if err == usecase.ErrFeedbackNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// RespondFeedback handles respond feedback requests
func (h *FeedbackHandler) RespondFeedback(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}

	var req entity.RespondFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedback, err := h.feedbackUsecase.RespondFeedback(c, uint(id), req)
	if err != nil {
		if err == usecase.ErrFeedbackNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
			return
		}
		if err == usecase.ErrAlreadyResponded {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Feedback already responded to"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// ListAllFeedback handles list all feedback requests
func (h *FeedbackHandler) ListAllFeedback(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	feedbacks, err := h.feedbackUsecase.ListAllFeedback(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

// ListUserFeedback handles list user feedback requests
func (h *FeedbackHandler) ListUserFeedback(c *gin.Context) {
	// Get user ID from token or parameter
	userIDParam := c.Param("user_id")
	var userID uint

	if userIDParam != "" {
		id, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userID = uint(id)
	} else {
		// Get user ID from token
		tokenUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userID = tokenUserID.(uint)
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	feedbacks, err := h.feedbackUsecase.ListUserFeedback(c, userID, page, limit)
	if err != nil {
		if err == usecase.ErrInvalidUserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

// ListPendingFeedback handles list pending feedback requests
func (h *FeedbackHandler) ListPendingFeedback(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	feedbacks, err := h.feedbackUsecase.ListPendingFeedback(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

// HealthCheck handles health check requests
func (h *FeedbackHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
