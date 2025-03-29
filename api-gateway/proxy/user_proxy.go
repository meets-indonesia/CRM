package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserProxy handles proxying requests to the User service
type UserProxy struct {
	baseURL string
	client  *http.Client
}

// NewUserProxy creates a new UserProxy
func NewUserProxy(baseURL string) *UserProxy {
	return &UserProxy{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// GetUser handles get user by ID requests
func (p *UserProxy) GetUser(c *gin.Context) {
	userID := c.Param("id")
	p.proxyRequest(c, "/users/"+userID, nil)
}

// UpdateUser handles update user requests
func (p *UserProxy) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	p.proxyRequest(c, "/users/"+userID, nil)
}

// ListAdmins handles list admins requests
func (p *UserProxy) ListAdmins(c *gin.Context) {
	p.proxyRequest(c, "/users/admin", nil)
}

// ListCustomers handles list customers requests
func (p *UserProxy) ListCustomers(c *gin.Context) {
	p.proxyRequest(c, "/users/customer", nil)
}

// GetCustomerPoints handles get customer points requests
func (p *UserProxy) GetCustomerPoints(c *gin.Context) {
	customerID := c.Param("id")
	p.proxyRequest(c, "/users/customer/"+customerID+"/points", nil)
}

// proxyRequest proxies a request to the User service
func (p *UserProxy) proxyRequest(c *gin.Context, path string, transformRequestBody func([]byte) ([]byte, error)) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to user service"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from user service"})
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
