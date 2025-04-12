package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	AuthTimestampHeader = "X-Auth-Timestamp"
	AuthSignatureHeader = "X-Auth-Signature"
	SecretKey           = "CRMSUMSEL2025@MEETSIDN" // Your secret key
	TimeWindow          = 5                        // 1 minutes in seconds
)

const (
	APIKeyHeader = "x-api-key"
	StaticAPIKey = "CRMSUMSEL2025@MEETSIDN" // Ganti dengan API key yang diinginkan
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip middleware for health check and auth routes
		if c.Request.URL.Path == "/health" ||
			c.Request.URL.Path == "/qr/verify/:code" ||
			c.Request.URL.Path == "/validate" ||
			c.Request.URL.Path == "/auth/admin/login" ||
			c.Request.URL.Path == "/auth/customer/login" ||
			c.Request.URL.Path == "/article/uploads/*filepath" ||
			c.Request.URL.Path == "/feedbacks/uploads/*filepath" ||
			c.Request.URL.Path == "/uploads/*filepath" {
			c.Next()
			return
		}

		// Get headers
		timestampStr := c.GetHeader(AuthTimestampHeader)
		signature := c.GetHeader(AuthSignatureHeader)

		// Validate headers presence
		if timestampStr == "" || signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "X-Auth-Timestamp and X-Auth-Signature headers are required",
			})
			return
		}

		// Validate timestamp
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid timestamp format",
			})
			return
		}

		// Check if timestamp is within allowed window
		currentTime := time.Now().Unix()
		if timestamp < currentTime-TimeWindow || timestamp > currentTime+TimeWindow {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "timestamp is too old or too new",
			})
			return
		}

		// Verify signature
		expectedSignature := generateSignature(timestampStr, SecretKey)
		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid signature",
			})
			return
		}

		c.Next()
	}
}

func generateSignature(timestamp, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(timestamp))
	return hex.EncodeToString(h.Sum(nil))
}

func SimpleAPIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "X-API-Key header is required",
			})
			return
		}

		if apiKey != StaticAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API key",
			})
			return
		}

		c.Next()
	}
}
