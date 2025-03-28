package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/reward/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/reward/domain/usecase"
)

// RewardHandler handles reward requests
type RewardHandler struct {
	rewardUsecase usecase.RewardUsecase
}

// NewRewardHandler creates a new RewardHandler
func NewRewardHandler(rewardUsecase usecase.RewardUsecase) *RewardHandler {
	return &RewardHandler{
		rewardUsecase: rewardUsecase,
	}
}

// CreateReward handles create reward requests
func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req entity.CreateRewardRequest

	// Baca form field
	req.Name = c.PostForm("name")
	req.Description = c.PostForm("description")
	pointCost, _ := strconv.Atoi(c.PostForm("point_cost"))
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	req.PointCost = pointCost
	req.Stock = stock

	// Handle image upload
	file, err := c.FormFile("image")
	if err == nil {
		// Simpan file ke folder tertentu di server
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		// Set URL/Path gambar
		req.ImageURL = "/" + filename
	}

	reward, err := h.rewardUsecase.CreateReward(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reward)
}

// GetReward handles get reward by ID requests
func (h *RewardHandler) GetReward(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reward ID"})
		return
	}

	reward, err := h.rewardUsecase.GetReward(c, uint(id))
	if err != nil {
		if err == usecase.ErrRewardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reward not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reward)
}

// UpdateReward handles update reward requests
// UpdateReward handles update reward requests
func (h *RewardHandler) UpdateReward(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reward ID"})
		return
	}

	// Ambil form field
	name := c.PostForm("name")
	description := c.PostForm("description")
	pointCostStr := c.PostForm("point_cost")
	stockStr := c.PostForm("stock")
	isActiveStr := c.PostForm("is_active")

	var pointCost, stock int
	if pointCostStr != "" {
		pointCost, _ = strconv.Atoi(pointCostStr)
	}
	if stockStr != "" {
		stock, _ = strconv.Atoi(stockStr)
	}

	var isActive *bool
	if isActiveStr != "" {
		val, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			isActive = &val
		}
	}

	// Handle image (jika ada file)
	file, err := c.FormFile("image")
	var imageURL string
	if err == nil {
		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		imageURL = "/" + filename
	}

	// Bangun request
	req := entity.UpdateRewardRequest{}

	if name != "" {
		req.Name = name
	}
	if description != "" {
		req.Description = description
	}
	if pointCost > 0 {
		req.PointCost = pointCost
	}
	if stock >= 0 {
		req.Stock = stock
	}
	if imageURL != "" {
		req.ImageURL = imageURL
	}
	if isActive != nil {
		req.IsActive = isActive
	}

	// Panggil usecase
	reward, err := h.rewardUsecase.UpdateReward(c, uint(id), req)
	if err != nil {
		if err == usecase.ErrRewardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reward not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reward)
}

// DeleteReward handles delete reward requests
func (h *RewardHandler) DeleteReward(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reward ID"})
		return
	}

	err = h.rewardUsecase.DeleteReward(c, uint(id))
	if err != nil {
		if err == usecase.ErrRewardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reward not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reward deleted successfully"})
}

// ListRewards handles list rewards requests
func (h *RewardHandler) ListRewards(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active_only", "false"))

	rewards, err := h.rewardUsecase.ListRewards(c, activeOnly, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// ClaimReward handles claim reward requests
func (h *RewardHandler) ClaimReward(c *gin.Context) {
	// Get user ID from token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req entity.ClaimRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claim, err := h.rewardUsecase.ClaimReward(c, userID.(uint), req)
	if err != nil {
		switch err {
		case usecase.ErrRewardNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Reward not found or inactive"})
		case usecase.ErrInsufficientStock:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		case usecase.ErrInsufficientPoints:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient points"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, claim)
}

// GetClaim handles get claim by ID requests
func (h *RewardHandler) GetClaim(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	claim, err := h.rewardUsecase.GetClaim(c, uint(id))
	if err != nil {
		if err == usecase.ErrClaimNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, claim)
}

// UpdateClaimStatus handles update claim status requests
func (h *RewardHandler) UpdateClaimStatus(c *gin.Context) {
	// Parse ID parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}

	var req entity.UpdateClaimStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claim, err := h.rewardUsecase.UpdateClaimStatus(c, uint(id), req)
	if err != nil {
		switch err {
		case usecase.ErrClaimNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		case usecase.ErrInvalidClaimStatus:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status transition"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, claim)
}

// ListUserClaims handles list user claims requests
func (h *RewardHandler) ListUserClaims(c *gin.Context) {
	// Get user ID from token or parameter
	var userID uint
	userIDParam := c.Param("user_id")

	if userIDParam != "" {
		// Admin is checking another user's claims
		id, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userID = uint(id)
	} else {
		// User is checking their own claims
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

	claims, err := h.rewardUsecase.ListUserClaims(c, userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, claims)
}

// ListAllClaims handles list all claims requests
func (h *RewardHandler) ListAllClaims(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	claims, err := h.rewardUsecase.ListAllClaims(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, claims)
}

// ListClaimsByStatus handles list claims by status requests
func (h *RewardHandler) ListClaimsByStatus(c *gin.Context) {
	// Parse status parameter
	status := entity.ClaimStatus(c.Param("status"))

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	claims, err := h.rewardUsecase.ListClaimsByStatus(c, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, claims)
}

// HealthCheck handles health check requests
func (h *RewardHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
