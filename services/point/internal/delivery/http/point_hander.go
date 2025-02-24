package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/user/internal/delivery/http/models"
	"github.com/kevinnaserwan/crm-be/services/user/internal/usecase"
)

type PointHandler struct {
	pointUseCase *usecase.PointUseCase
}

func NewPointHandler(pointUseCase *usecase.PointUseCase) *PointHandler {
	return &PointHandler{
		pointUseCase: pointUseCase,
	}
}

func (h *PointHandler) Handle(c *gin.Context) {
	var baseRequest models.BaseRequest
	if err := c.ShouldBindJSON(&baseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received request with action: %s", baseRequest.Meta.Action)

	switch baseRequest.Meta.Action {
	case "process_feedback_point":
		h.handleProcessFeedbackPoint(c, baseRequest.Data)
	case "get_user_points":
		h.handleGetUserPoints(c)
	case "get_point_history":
		h.handleGetPointHistory(c)
	case "get_all_points":
		h.handleGetAllPoints(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
	}
}

func (h *PointHandler) handleProcessFeedbackPoint(c *gin.Context, data interface{}) {
	var feedbackData models.ProcessFeedbackData
	if err := convertData(data, &feedbackData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := h.pointUseCase.ProcessFeedbackPoint(c.Request.Context(), userID.(uuid.UUID), feedbackData.FeedbackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Points processed successfully"})
}

func (h *PointHandler) handleGetUserPoints(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Make sure user has points record
	if err := h.pointUseCase.InitializeIfNotExists(c.Request.Context(), userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	points, level, err := h.pointUseCase.GetUserPointsAndLevel(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.UserPointResponse{
		UserID:        points.UserID,
		TotalPoints:   points.TotalPoints,
		Level:         level,
		NextResetDate: points.NextResetDate,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *PointHandler) handleGetPointHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Ensure user has points record
	if err := h.pointUseCase.InitializeIfNotExists(c.Request.Context(), userID.(uuid.UUID)); err != nil {
		log.Printf("Error initializing points: %v", err)
		// Continue anyway, just return empty history
	}

	history, err := h.pointUseCase.GetPointHistory(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []models.PointHistoryResponse
	for _, h := range history {
		response = append(response, models.PointHistoryResponse{
			ID:           h.ID,
			PointsEarned: h.PointsEarned,
			DateEarned:   h.DateEarned,
			FeedbackID:   h.FeedbackID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *PointHandler) handleGetAllPoints(c *gin.Context) {
	points, err := h.pointUseCase.GetAllUserPoints(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": points})
}

func convertData(data interface{}, target interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonStr, target)
}
