package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	APIKeyHeader = "x-api-key"
	ValidAPIKey  = "CRMSUMSEL2025@MEETSIDN"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip middleware untuk health check dan auth routes
		if c.Request.URL.Path == "/health" ||
			c.Request.URL.Path == "/qr/verify/:code" ||
			c.Request.URL.Path == "/validate" ||
			c.Request.URL.Path == "/auth/admin/login" ||
			c.Request.URL.Path == "/auth/customer/login" {
			c.Next()
			return
		}

		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "x-api-key header is required",
			})
			return
		}

		if apiKey != ValidAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid api key",
			})
			return
		}

		c.Next()
	}
}
