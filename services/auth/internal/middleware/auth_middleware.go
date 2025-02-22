package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/util"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := util.ValidateJWT(parts[1], jwtSecret)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Explicitly set as string
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		fmt.Printf("DEBUG: Setting role in middleware: %s\n", claims.Role)

		c.Next()
	}
}
