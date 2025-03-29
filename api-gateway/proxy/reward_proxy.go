package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RewardProxy handles proxying requests to the Reward service
type RewardProxy struct {
	baseURL string
	client  *http.Client
}

// NewRewardProxy creates a new RewardProxy
func NewRewardProxy(baseURL string) *RewardProxy {
	return &RewardProxy{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// CreateReward handles create reward requests
func (p *RewardProxy) CreateReward(c *gin.Context) {
	p.proxyRequest(c, "/rewards", nil)
}

// GetReward handles get reward by ID requests
func (p *RewardProxy) GetReward(c *gin.Context) {
	p.proxyRequest(c, "/rewards/"+c.Param("id"), nil)
}

// UpdateReward handles update reward requests
func (p *RewardProxy) UpdateReward(c *gin.Context) {
	p.proxyRequest(c, "/rewards/"+c.Param("id"), nil)
}

// DeleteReward handles delete reward requests
func (p *RewardProxy) DeleteReward(c *gin.Context) {
	p.proxyRequest(c, "/rewards/"+c.Param("id"), nil)
}

// ListRewards handles list rewards requests
func (p *RewardProxy) ListRewards(c *gin.Context) {
	p.proxyRequest(c, "/rewards", nil)
}

// ClaimReward handles claim reward requests
func (p *RewardProxy) ClaimReward(c *gin.Context) {
	p.proxyRequest(c, "/claims", nil)
}

// GetClaim handles get claim by ID requests
func (p *RewardProxy) GetClaim(c *gin.Context) {
	p.proxyRequest(c, "/claims/"+c.Param("id"), nil)
}

// UpdateClaimStatus handles update claim status requests
func (p *RewardProxy) UpdateClaimStatus(c *gin.Context) {
	p.proxyRequest(c, "/claims/"+c.Param("id")+"/status", nil)
}

// ListUserClaims handles list user claims requests
func (p *RewardProxy) ListUserClaims(c *gin.Context) {
	if c.Param("user_id") != "" {
		p.proxyRequest(c, "/claims/user/"+c.Param("user_id"), nil)
	} else {
		p.proxyRequest(c, "/claims/user", nil)
	}
}

// ListAllClaims handles list all claims requests
func (p *RewardProxy) ListAllClaims(c *gin.Context) {
	p.proxyRequest(c, "/claims", nil)
}

// ListClaimsByStatus handles list claims by status requests
func (p *RewardProxy) ListClaimsByStatus(c *gin.Context) {
	p.proxyRequest(c, "/claims/status/"+c.Param("status"), nil)
}

// proxyRequest proxies a request to the Reward service
func (p *RewardProxy) proxyRequest(c *gin.Context, path string, transformRequestBody func([]byte) ([]byte, error)) {
	targetURL := p.baseURL + path

	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Transform the request body if needed
	if transformRequestBody != nil {
		body, err = transformRequestBody(body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transform request body"})
			return
		}
	}

	// Create a new request to the target service
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	// Copy headers from the original request
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Tambahkan bagian ini untuk meneruskan x-api-key ke service backend
	apiKey := c.GetHeader("x-api-key")
	if apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
	}

	// Set content type if it's not already set
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Send the request to the target service
	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to reward service"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from reward service"})
		return
	}

	// Copy headers from the target service response
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set the status code and write the response body
	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)
}
