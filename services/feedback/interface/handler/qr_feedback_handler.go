package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/feedback/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/feedback/domain/usecase"
)

type QRFeedbackHandler struct {
	qrFeedbackUsecase usecase.QRFeedbackUsecase
}

func NewQRFeedbackHandler(qrFeedbackUsecase usecase.QRFeedbackUsecase) *QRFeedbackHandler {
	return &QRFeedbackHandler{
		qrFeedbackUsecase: qrFeedbackUsecase,
	}
}

// CreateQRFeedback handles the creation of a new QR code for feedback
func (h *QRFeedbackHandler) CreateQRFeedback(c *gin.Context) {
	var req entity.CreateQRFeedbackRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qrFeedback, err := h.qrFeedbackUsecase.CreateQRFeedback(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, qrFeedback)
}

// GetQRFeedback gets a QR feedback by ID
func (h *QRFeedbackHandler) GetQRFeedback(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	qrFeedback, err := h.qrFeedbackUsecase.GetQRFeedback(c, uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrQRFeedbackNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, qrFeedback)
}

// ListQRFeedback lists all QR feedbacks with pagination
func (h *QRFeedbackHandler) ListQRFeedback(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	response, err := h.qrFeedbackUsecase.ListQRFeedback(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DownloadQRCode generates and downloads a QR code image
func (h *QRFeedbackHandler) DownloadQRCode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	// Get the QR feedback first to get the station name for filename
	qrFeedback, err := h.qrFeedbackUsecase.GetQRFeedback(c, uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrQRFeedbackNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Generate QR code image
	qrCodeImage, err := h.qrFeedbackUsecase.GenerateQRCodeImage(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set response headers for file download
	fileName := "feedback-qr-" + qrFeedback.Station + ".png"
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "image/png")
	c.Header("Content-Length", strconv.Itoa(len(qrCodeImage)))

	// Write the QR code image to the response
	c.Writer.Write(qrCodeImage)
}

// VerifyQRCode verifies and returns the QR feedback information
func (h *QRFeedbackHandler) VerifyQRCode(c *gin.Context) {
	qrCode := c.Param("qrCode")

	qrFeedback, err := h.qrFeedbackUsecase.VerifyQRCode(c, qrCode)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrInvalidQRCode {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, qrFeedback)
}
