package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ArticleProxy handles proxying requests to the Article service
type ArticleProxy struct {
	baseURL string
	client  *http.Client
}

// NewArticleProxy creates a new ArticleProxy
func NewArticleProxy(baseURL string) *ArticleProxy {
	return &ArticleProxy{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// CreateArticle handles create article requests
func (p *ArticleProxy) CreateArticle(c *gin.Context) {
	targetURL := p.baseURL + "/articles"

	// Batasi ukuran upload (optional)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // max 10MB

	// Buat request baru ke service article
	req, err := http.NewRequest("POST", targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	// Salin semua header dari request original
	req.Header = c.Request.Header

	// Tambahkan auth headers
	AddAuthHeaders(req)

	// Kirim request ke service article
	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward to article service"})
		return
	}
	defer resp.Body.Close()

	// Salin semua header response dari service
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	c.Status(resp.StatusCode)

	io.Copy(c.Writer, resp.Body)
}

// GetArticle handles get article by ID requests
func (p *ArticleProxy) GetArticle(c *gin.Context) {
	p.proxyRequest(c, "/articles/"+c.Param("id"), nil)
}

// UpdateArticle handles update article requests
func (p *ArticleProxy) UpdateArticle(c *gin.Context) {
	p.proxyRequest(c, "/articles/"+c.Param("id"), nil)
}

// DeleteArticle handles delete article requests
func (p *ArticleProxy) DeleteArticle(c *gin.Context) {
	p.proxyRequest(c, "/articles/"+c.Param("id"), nil)
}

// ListArticles handles list articles requests
func (p *ArticleProxy) ListArticles(c *gin.Context) {
	p.proxyRequest(c, "/articles", nil)
}

// ViewArticle handles viewing an article and incrementing view count
func (p *ArticleProxy) ViewArticle(c *gin.Context) {
	p.proxyRequest(c, "/articles/"+c.Param("id"), nil)
}

// SearchArticles handles search articles requests
func (p *ArticleProxy) SearchArticles(c *gin.Context) {
	p.proxyRequest(c, "/articles/search", nil)
}

// proxyRequest proxies a request to the Article service
func (p *ArticleProxy) proxyRequest(c *gin.Context, path string, transformRequestBody func([]byte) ([]byte, error)) {
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

	// Add auth headers
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to article service"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from article service"})
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

// AccessUploadImages handles accessing uploaded images from article service
func (p *ArticleProxy) AccessUploadImages(c *gin.Context) {
	filepath := c.Param("filepath")
	targetURL := p.baseURL + "/uploads" + filepath

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request to article service"})
		return
	}

	// Copy headers from the original request (optional)
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add auth headers
	AddAuthHeaders(req)

	resp, err := p.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch image from article service"})
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
