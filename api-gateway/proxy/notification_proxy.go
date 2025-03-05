package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotificationProxy handles proxying requests to the Notification service
type NotificationProxy struct {
	baseURL string
	client  *http.Client
}

// NewNotificationProxy creates a new NotificationProxy
func NewNotificationProxy(baseURL string) *NotificationProxy {
	return &NotificationProxy{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// SendEmail handles send email requests
func (p *NotificationProxy) SendEmail(c *gin.Context) {
	p.proxyRequest(c, "/notifications/email", nil)
}

// SendPushNotification handles send push notification requests
func (p *NotificationProxy) SendPushNotification(c *gin.Context) {
	p.proxyRequest(c, "/notifications/push", nil)
}

// GetNotification handles get notification by ID requests
func (p *NotificationProxy) GetNotification(c *gin.Context) {
	p.proxyRequest(c, "/notifications/"+c.Param("id"), nil)
}

// ListUserNotifications handles list user notifications requests
func (p *NotificationProxy) ListUserNotifications(c *gin.Context) {
	if c.Param("user_id") != "" {
		p.proxyRequest(c, "/notifications/user/"+c.Param("user_id"), nil)
	} else {
		p.proxyRequest(c, "/notifications", nil)
	}
}

// ProcessPendingNotifications handles process pending notifications requests
func (p *NotificationProxy) ProcessPendingNotifications(c *gin.Context) {
	p.proxyRequest(c, "/notifications/process", nil)
}

// proxyRequest proxies a request to the Notification service
func (p *NotificationProxy) proxyRequest(c *gin.Context, path string, transformRequestBody func([]byte) ([]byte, error)) {
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

	// Add query parameters
	req.URL.RawQuery = c.Request.URL.RawQuery

	// Send the request to the target service
	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to notification service"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from notification service"})
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
