package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FeedbackProxy handles proxying requests to the Feedback service
type FeedbackProxy struct {
	baseURL string
	client  *http.Client
}

// NewFeedbackProxy creates a new FeedbackProxy
func NewFeedbackProxy(baseURL string) *FeedbackProxy {
	return &FeedbackProxy{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// CreateFeedback handles feedback creation requests
func (p *FeedbackProxy) CreateFeedback(c *gin.Context) {
	p.proxyRequest(c, "", nil)
}

// GetFeedback handles get feedback by ID requests
func (p *FeedbackProxy) GetFeedback(c *gin.Context) {
	p.proxyRequest(c, "/"+c.Param("id"), nil)
}

// RespondToFeedback handles respond to feedback requests
func (p *FeedbackProxy) RespondToFeedback(c *gin.Context) {
	p.proxyRequest(c, "/"+c.Param("id")+"/respond", nil)
}

// ListAllFeedback handles list all feedback requests
func (p *FeedbackProxy) ListAllFeedback(c *gin.Context) {
	p.proxyRequest(c, "", nil)
}

// ListPendingFeedback handles list pending feedback requests
func (p *FeedbackProxy) ListPendingFeedback(c *gin.Context) {
	p.proxyRequest(c, "/pending", nil)
}

// ListUserFeedback handles list user feedback requests
func (p *FeedbackProxy) ListUserFeedback(c *gin.Context) {
	if c.Param("user_id") != "" {
		p.proxyRequest(c, "/user/"+c.Param("user_id"), nil)
	} else {
		p.proxyRequest(c, "/user", nil)
	}
}

// proxyRequest proxies a request to the Feedback service
func (p *FeedbackProxy) proxyRequest(c *gin.Context, path string, transformRequestBody func([]byte) ([]byte, error)) {
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

	// Set content type if it's not already set
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Send the request to the target service
	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to feedback service"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from feedback service"})
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
