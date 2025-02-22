package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/delivery/http/models"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/usecase"
)

type FeedbackHandler struct {
	feedbackUseCase *usecase.FeedbackUseCase
}

func NewFeedbackHandler(feedbackUseCase *usecase.FeedbackUseCase) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackUseCase: feedbackUseCase,
	}
}

func (h *FeedbackHandler) Handle(c *gin.Context) {
	var baseRequest models.BaseRequest
	if err := c.ShouldBindJSON(&baseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch baseRequest.Meta.Action {
	case "create_feedback":
		h.handleCreateFeedback(c, baseRequest.Data)
	case "create_feedback_type":
		h.handleCreateFeedbackType(c, baseRequest.Data)
	case "create_station":
		h.handleCreateStation(c, baseRequest.Data)
	case "respond_to_feedback":
		h.handleRespondToFeedback(c, baseRequest.Data)
	case "get_user_feedbacks":
		h.handleGetUserFeedbacks(c)
	case "get_all_feedbacks":
		h.handleGetAllFeedbacks(c)
	case "get_feedback_types":
		h.handleGetFeedbackTypes(c)
	case "get_stations":
		h.handleGetStations(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
	}
}

func (h *FeedbackHandler) handleCreateFeedback(c *gin.Context, data interface{}) {
	var feedbackData models.CreateFeedbackData
	if err := convertData(data, &feedbackData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Now userID is already UUID type from middleware
	feedback := &model.Feedback{
		ID:             uuid.New(),
		UserID:         userID.(uuid.UUID), // This should work now
		FeedbackTypeID: feedbackData.FeedbackTypeID,
		StationID:      feedbackData.StationID,
		Feedback:       feedbackData.Feedback,
		Documentation:  feedbackData.Documentation,
		Rating:         feedbackData.Rating,
	}

	if err := h.feedbackUseCase.CreateFeedback(c.Request.Context(), feedback); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Feedback submitted successfully",
		"data":    feedback,
	})
}

func (h *FeedbackHandler) handleCreateFeedbackType(c *gin.Context, data interface{}) {
	var typeData models.CreateFeedbackTypeData
	if err := convertData(data, &typeData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedbackType := &model.FeedbackType{
		ID:   uuid.New(),
		Name: typeData.Name,
	}

	if err := h.feedbackUseCase.CreateFeedbackType(c.Request.Context(), feedbackType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Feedback type created successfully",
		"data":    feedbackType,
	})
}

func (h *FeedbackHandler) handleCreateStation(c *gin.Context, data interface{}) {
	var stationData models.CreateStationData
	if err := convertData(data, &stationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	station := &model.Station{
		ID:   uuid.New(),
		Name: stationData.Name,
	}

	if err := h.feedbackUseCase.CreateStation(c.Request.Context(), station); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Station created successfully",
		"data":    station,
	})
}

func (h *FeedbackHandler) handleRespondToFeedback(c *gin.Context, data interface{}) {
	var responseData models.RespondToFeedbackData
	if err := convertData(data, &responseData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.feedbackUseCase.RespondToFeedback(c.Request.Context(), responseData.FeedbackID, responseData.Response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Response submitted successfully"})
}

func (h *FeedbackHandler) handleGetUserFeedbacks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	feedbacks, err := h.feedbackUseCase.GetUserFeedbacks(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedbacks})
}

func (h *FeedbackHandler) handleGetAllFeedbacks(c *gin.Context) {
	feedbacks, err := h.feedbackUseCase.GetAllFeedbacks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedbacks})
}

func (h *FeedbackHandler) handleGetFeedbackTypes(c *gin.Context) {
	types, err := h.feedbackUseCase.GetFeedbackTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": types})
}

func (h *FeedbackHandler) handleGetStations(c *gin.Context) {
	stations, err := h.feedbackUseCase.GetStations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stations})
}

// Helper function to convert interface{} to specific type
func convertData(data interface{}, target interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonStr, target)
}
