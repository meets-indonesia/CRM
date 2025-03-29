package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthProxy handles proxying requests to the Auth service
type AuthProxy struct {
	baseURL string
	client  *http.Client
}

// NewAuthProxy creates a new AuthProxy
func NewAuthProxy(baseURL string) *AuthProxy {
	return &AuthProxy{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// AdminLogin handles admin login requests
func (p *AuthProxy) AdminLogin(c *gin.Context) {
	p.proxyRequest(c, "/admin/login", nil)
}

// AdminRegister handles admin registration requests
func (p *AuthProxy) AdminRegister(c *gin.Context) {
	p.proxyRequest(c, "/admin/register", nil)
}

// AdminResetPassword handles admin password reset requests
func (p *AuthProxy) AdminResetPassword(c *gin.Context) {
	p.proxyRequest(c, "/admin/reset-password", nil)
}

// AdminVerifyOTP handles admin OTP verification requests
func (p *AuthProxy) AdminVerifyOTP(c *gin.Context) {
	p.proxyRequest(c, "/admin/verify-otp", nil)
}

// CustomerLogin handles customer login requests
func (p *AuthProxy) CustomerLogin(c *gin.Context) {
	p.proxyRequest(c, "/customer/login", nil)
}

// CustomerGoogleLogin handles customer login via Google OAuth
func (p *AuthProxy) CustomerGoogleLogin(c *gin.Context) {
	p.proxyRequest(c, "/customer/google", nil)
}

// ValidateToken handles token validation requests
func (p *AuthProxy) ValidateToken(c *gin.Context) {
	p.proxyRequest(c, "/validate", nil)
}

// proxyRequest proxies a request to the Auth service
func (p *AuthProxy) proxyRequest(c *gin.Context, path string, transformRequestBody func([]byte) ([]byte, error)) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to auth service"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from auth service"})
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
