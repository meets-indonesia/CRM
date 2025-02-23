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

		fmt.Printf("Setting context values - UserID: %v, Email: %v, Role: %v\n",
			claims.UserID, claims.Email, claims.Role)

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email) // Make sure this is set
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
