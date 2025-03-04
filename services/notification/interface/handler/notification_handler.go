package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/notification/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/notification/domain/usecase"
)

// NotificationHandler handles notification requests
type NotificationHandler struct {
	notificationUsecase usecase.NotificationUsecase
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(notificationUsecase usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{
		notificationUsecase: notificationUsecase,
	}
}

// SendEmail handles send email requests
func (h *NotificationHandler) SendEmail(c *gin.Context) {
	var req entity.CreateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.notificationUsecase.SendEmail(c, req)
	if err != nil {
		if err == usecase.ErrInvalidUserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		if err == usecase.ErrSendingNotification {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification", "details": notification})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// SendPushNotification handles send push notification requests
func (h *NotificationHandler) SendPushNotification(c *gin.Context) {
	var req entity.CreatePushNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.notificationUsecase.SendPushNotification(c, req)
	if err != nil {
		if err == usecase.ErrInvalidUserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		if err == usecase.ErrSendingNotification {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification", "details": notification})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// GetNotification handles get notification by ID requests
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	notification, err := h.notificationUsecase.GetNotification(c, uint(id))
	if err != nil {
		if err == usecase.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// ListUserNotifications handles list user notifications requests
func (h *NotificationHandler) ListUserNotifications(c *gin.Context) {
	// Get user ID from token or parameter
	var userID uint
	userIDParam := c.Param("user_id")

	if userIDParam != "" {
		// Admin is checking another user's notifications
		id, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userID = uint(id)
	} else {
		// User is checking their own notifications
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

	notifications, err := h.notificationUsecase.ListUserNotifications(c, userID, page, limit)
	if err != nil {
		if err == usecase.ErrInvalidUserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// ProcessPendingNotifications handles process pending notifications requests
func (h *NotificationHandler) ProcessPendingNotifications(c *gin.Context) {
	err := h.notificationUsecase.ProcessPendingNotifications(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pending notifications processed"})
}

// HealthCheck handles health check requests
func (h *NotificationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
