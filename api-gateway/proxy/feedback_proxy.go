package proxy

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

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
	if c.ContentType() == "multipart/form-data" {
		p.proxyMultipartRequest(c, "/feedbacks")
		return
	}
	p.proxyRequest(c, "/feedbacks", nil)
}

// GetFeedback handles get feedback by ID requests
func (p *FeedbackProxy) GetFeedback(c *gin.Context) {
	p.proxyRequest(c, "/feedbacks/"+c.Param("id"), nil)
}

// RespondToFeedback handles respond to feedback requests
func (p *FeedbackProxy) RespondToFeedback(c *gin.Context) {
	p.proxyRequest(c, "/feedbacks/"+c.Param("id")+"/respond", nil)
}

// ListAllFeedback handles list all feedback requests
func (p *FeedbackProxy) ListAllFeedback(c *gin.Context) {
	p.proxyRequest(c, "/feedbacks"+getQueryString(c), nil)
}

// ListPendingFeedback handles list pending feedback requests
func (p *FeedbackProxy) ListPendingFeedback(c *gin.Context) {
	p.proxyRequest(c, "/feedbacks/pending"+getQueryString(c), nil)
}

// AccessUploadImages handles accessing uploaded images (tanpa auth headers)
func (p *FeedbackProxy) AccessUploadImages(c *gin.Context) {
	filepath := c.Param("filepath")
	targetURL := p.baseURL + "/uploads" + filepath

	// Buat request tanpa menambahkan auth headers
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy query parameters
	req.URL.RawQuery = c.Request.URL.RawQuery

	// Kirim request
	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch image from feedback service"})
		return
	}
	defer resp.Body.Close()

	// Copy headers dari response
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Stream response
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// ListUserFeedback handles list user feedback requests
func (p *FeedbackProxy) ListUserFeedback(c *gin.Context) {
	path := "/feedbacks/user"
	if userId := c.Param("user_id"); userId != "" {
		path += "/" + userId
	}
	p.proxyRequest(c, path+getQueryString(c), nil)
}

// CreateQRFeedback handles QR feedback creation requests
func (p *FeedbackProxy) CreateQRFeedback(c *gin.Context) {
	p.proxyRequest(c, "/qr-feedbacks", nil)
}

// GetQRFeedback handles get QR feedback by ID requests
func (p *FeedbackProxy) GetQRFeedback(c *gin.Context) {
	p.proxyRequest(c, "/qr-feedbacks/"+c.Param("id"), nil)
}

// ListQRFeedbacks handles list all QR feedbacks requests
func (p *FeedbackProxy) ListQRFeedbacks(c *gin.Context) {
	p.proxyRequest(c, "/qr-feedbacks"+getQueryString(c), nil)
}

// GenerateQRCodeImage handles QR code image generation requests
func (p *FeedbackProxy) GenerateQRCodeImage(c *gin.Context) {
	path := "/qr-feedbacks/" + c.Param("id") + "/download"
	p.proxyRequest(c, path, nil)
}

// VerifyQRCode handles QR code verification requests
func (p *FeedbackProxy) VerifyQRCode(c *gin.Context) {
	code := c.Param("code")
	fmt.Printf("Proxying QR verification to: %s/feedback/scan/%s\n", p.baseURL, code)
	p.proxyRequest(c, "/feedback/scan/"+code, nil)
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

	// Add auth headers using the shared function
	AddAuthHeaders(req)

	// Set content type if it's not already set
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add query parameters
	req.URL.RawQuery = c.Request.URL.RawQuery

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

// proxyMultipartRequest handles multipart requests for file uploads
func (p *FeedbackProxy) proxyMultipartRequest(c *gin.Context, path string) {
	targetURL := p.baseURL + path

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read multipart form"})
		return
	}

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()

		// Copy form fields
		for key, values := range form.Value {
			for _, value := range values {
				writer.WriteField(key, value)
			}
		}

		// Copy files
		for _, fileHeaders := range form.File {
			for _, fileHeader := range fileHeaders {
				formFile, err := fileHeader.Open()
				if err != nil {
					continue
				}
				defer formFile.Close()

				part, err := writer.CreateFormFile("image", filepath.Base(fileHeader.Filename))
				if err != nil {
					continue
				}

				_, err = io.Copy(part, formFile)
				if err != nil {
					continue
				}
			}
		}
	}()

	req, err := http.NewRequest(c.Request.Method, targetURL, pr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	AddAuthHeaders(req)

	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to feedback service"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from feedback service"})
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)
}

// getQueryString extracts query string from gin context
func getQueryString(c *gin.Context) string {
	if query := c.Request.URL.RawQuery; query != "" {
		return "?" + query
	}
	return ""
}
