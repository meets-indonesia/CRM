package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/delivery/http/models"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/infrastructure/auth"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/usecase"
)

type FeedbackHandler struct {
	feedbackUseCase *usecase.FeedbackUseCase
	authService     *auth.AuthService
}

func NewFeedbackHandler(feedbackUseCase *usecase.FeedbackUseCase, authService *auth.AuthService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackUseCase: feedbackUseCase,
		authService:     authService,
	}
}

func (h *FeedbackHandler) Handle(c *gin.Context) {
	// Read the raw body
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Restore the body for further processing
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

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
		// Extract user_id directly from the data
		dataMap, ok := baseRequest.Data.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
			return
		}

		userIDStr, ok := dataMap["user_id"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required and must be a string"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		feedbacks, err := h.feedbackUseCase.GetUserFeedbacks(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": feedbacks})
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

	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found"})
		return
	}

	feedback := &model.Feedback{
		UserID:         userID.(uuid.UUID),
		UserEmail:      userEmail.(string), // Tambahkan ini
		FeedbackTypeID: feedbackData.FeedbackTypeID,
		StationID:      feedbackData.StationID,
		Feedback:       feedbackData.Feedback,
		Documentation:  feedbackData.Documentation,
		Rating:         feedbackData.Rating,
		FeedbackDate:   time.Now(),
		Status:         model.FeedbackStatusPending,
	}

	if err := h.feedbackUseCase.CreateFeedback(c.Request.Context(), feedback, userEmail.(string)); err != nil {
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

	// Get the feedback
	feedback, err := h.feedbackUseCase.GetFeedbackByID(c.Request.Context(), responseData.FeedbackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feedback details"})
		return
	}

	// Gunakan email yang tersimpan
	if err := h.feedbackUseCase.RespondToFeedback(
		c.Request.Context(),
		responseData.FeedbackID,
		responseData.Response,
		feedback.UserEmail, // Gunakan email yang tersimpan
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Response submitted successfully"})
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
