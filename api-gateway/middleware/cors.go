package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a middleware for handling CORS
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Consider restricting this in production
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"x-api-key",
			"X-Auth-Timestamp",
			"X-Auth-Signature",
			"X-Requested-With", // Often needed for axios
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Auth-Timestamp", // Explicitly expose custom headers
			"X-Auth-Signature",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache preflight requests
	})
}
