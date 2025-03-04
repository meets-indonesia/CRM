package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InventoryProxy handles proxying requests to the Inventory service
type InventoryProxy struct {
	baseURL string
	client  *http.Client
}

// NewInventoryProxy creates a new InventoryProxy
func NewInventoryProxy(baseURL string) *InventoryProxy {
	return &InventoryProxy{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// CreateItem handles create item requests
func (p *InventoryProxy) CreateItem(c *gin.Context) {
	p.proxyRequest(c, "/items", nil)
}

// GetItem handles get item by ID requests
func (p *InventoryProxy) GetItem(c *gin.Context) {
	p.proxyRequest(c, "/items/"+c.Param("id"), nil)
}

// UpdateItem handles update item requests
func (p *InventoryProxy) UpdateItem(c *gin.Context) {
	p.proxyRequest(c, "/items/"+c.Param("id"), nil)
}

// DeleteItem handles delete item requests
func (p *InventoryProxy) DeleteItem(c *gin.Context) {
	p.proxyRequest(c, "/items/"+c.Param("id"), nil)
}

// ListItems handles list items requests
func (p *InventoryProxy) ListItems(c *gin.Context) {
	p.proxyRequest(c, "/items", nil)
}

// ListItemsByCategory handles list items by category requests
func (p *InventoryProxy) ListItemsByCategory(c *gin.Context) {
	p.proxyRequest(c, "/items/category/"+c.Param("category"), nil)
}

// SearchItems handles search items requests
func (p *InventoryProxy) SearchItems(c *gin.Context) {
	p.proxyRequest(c, "/items/search", nil)
}

// GetLowStockItems handles get low stock items requests
func (p *InventoryProxy) GetLowStockItems(c *gin.Context) {
	p.proxyRequest(c, "/items/low-stock", nil)
}

// IncreaseStock handles increase stock requests
func (p *InventoryProxy) IncreaseStock(c *gin.Context) {
	p.proxyRequest(c, "/items/"+c.Param("id")+"/stock/increase", nil)
}

// DecreaseStock handles decrease stock requests
func (p *InventoryProxy) DecreaseStock(c *gin.Context) {
	p.proxyRequest(c, "/items/"+c.Param("id")+"/stock/decrease", nil)
}

// GetStockTransaction handles get stock transaction by ID requests
func (p *InventoryProxy) GetStockTransaction(c *gin.Context) {
	p.proxyRequest(c, "/transactions/"+c.Param("id"), nil)
}

// ListStockTransactionsByItem handles list stock transactions by item requests
func (p *InventoryProxy) ListStockTransactionsByItem(c *gin.Context) {
	p.proxyRequest(c, "/transactions/item/"+c.Param("id"), nil)
}

// ListAllStockTransactions handles list all stock transactions requests
func (p *InventoryProxy) ListAllStockTransactions(c *gin.Context) {
	p.proxyRequest(c, "/transactions", nil)
}

// CreateSupplier handles create supplier requests
func (p *InventoryProxy) CreateSupplier(c *gin.Context) {
	p.proxyRequest(c, "/suppliers", nil)
}

// GetSupplier handles get supplier by ID requests
func (p *InventoryProxy) GetSupplier(c *gin.Context) {
	p.proxyRequest(c, "/suppliers/"+c.Param("id"), nil)
}

// UpdateSupplier handles update supplier requests
func (p *InventoryProxy) UpdateSupplier(c *gin.Context) {
	p.proxyRequest(c, "/suppliers/"+c.Param("id"), nil)
}

// DeleteSupplier handles delete supplier requests
func (p *InventoryProxy) DeleteSupplier(c *gin.Context) {
	p.proxyRequest(c, "/suppliers/"+c.Param("id"), nil)
}

// ListSuppliers handles list suppliers requests
func (p *InventoryProxy) ListSuppliers(c *gin.Context) {
	p.proxyRequest(c, "/suppliers", nil)
}

// SearchSuppliers handles search suppliers requests
func (p *InventoryProxy) SearchSuppliers(c *gin.Context) {
	p.proxyRequest(c, "/suppliers/search", nil)
}

// proxyRequest proxies a request to the Inventory service
func (p *InventoryProxy) proxyRequest(c *gin.Context, path string, transformRequestBody func([]byte) ([]byte, error)) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to inventory service"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from inventory service"})
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
