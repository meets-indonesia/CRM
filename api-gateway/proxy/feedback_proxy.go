package proxy

import (
	"bytes"
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
	// Untuk form-data dengan file upload, kita perlu penanganan khusus
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
	// Meneruskan query parameters
	path := "/feedbacks" + getQueryString(c)
	p.proxyRequest(c, path, nil)
}

// ListPendingFeedback handles list pending feedback requests
func (p *FeedbackProxy) ListPendingFeedback(c *gin.Context) {
	path := "/feedbacks/pending" + getQueryString(c)
	p.proxyRequest(c, path, nil)
}

// AccesUploadImages handles accessing uploaded images
func (p *FeedbackProxy) AccesUploadImages(c *gin.Context) {
	path := "feedbacks/uploads/" + c.Query("filename")
	p.proxyRequest(c, path, nil)
}

// ListUserFeedback handles list user feedback requests
func (p *FeedbackProxy) ListUserFeedback(c *gin.Context) {
	var path string
	if c.Param("user_id") != "" {
		path = "/feedbacks/user/" + c.Param("user_id")
	} else {
		path = "/feedbacks/user"
	}
	// Meneruskan query parameters
	path += getQueryString(c)
	p.proxyRequest(c, path, nil)
}

// proxyMultipartRequest secara khusus menangani request multipart untuk upload file
func (p *FeedbackProxy) proxyMultipartRequest(c *gin.Context, path string) {
	// Target URL
	targetURL := p.baseURL + path

	// Buat form data baru
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read multipart form"})
		return
	}

	// Buat request multipart baru
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Mulai goroutine untuk menulis form data ke pipe
	go func() {
		defer pw.Close()
		defer writer.Close()

		// Salin semua field form
		for key, values := range form.Value {
			for _, value := range values {
				writer.WriteField(key, value)
			}
		}

		// Salin semua file
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

	// Buat HTTP request baru ke service target
	req, err := http.NewRequest(http.MethodPost, targetURL, pr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	// Set content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Salin header authorization
	auth := c.GetHeader("Authorization")
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}

	// Kirim request
	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to feedback service"})
		return
	}
	defer resp.Body.Close()

	// Baca response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from feedback service"})
		return
	}

	// Salin header dari response service target
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set status code dan tulis response body
	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)
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

// getQueryString mengekstrak query string dari gin context
func getQueryString(c *gin.Context) string {
	query := c.Request.URL.RawQuery
	if query != "" {
		return "?" + query
	}
	return ""
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
	path := "/qr-feedbacks" + getQueryString(c)
	p.proxyRequest(c, path, nil)
}

// GenerateQRCodeImage handles QR code image generation requests
func (p *FeedbackProxy) GenerateQRCodeImage(c *gin.Context) {
	path := "/qr-feedbacks/" + c.Param("id") + "/download"

	// Create request
	req, err := http.NewRequest("GET", p.baseURL+path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make request
	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch QR code"})
		return
	}
	defer resp.Body.Close()

	// Copy headers and status
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	c.Status(resp.StatusCode)

	// Stream the image data directly
	io.Copy(c.Writer, resp.Body)
}

// VerifyQRCode handles QR code verification requests
func (p *FeedbackProxy) VerifyQRCode(c *gin.Context) {
	p.proxyRequest(c, "/qr-feedbacks/verify/"+c.Param("code"), nil)
}
