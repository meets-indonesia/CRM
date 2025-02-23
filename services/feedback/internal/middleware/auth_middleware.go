package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/util"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"` // Add this
	Role   string `json:"role"`
	jwt.StandardClaims
}

// internal/middleware/auth_middleware.go
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

		// Debug print
		fmt.Printf("Debug - Setting user data in AuthMiddleware: ID=%v, Role=%v, Email=%v\n",
			claims.UserID, claims.Role, claims.Email)

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by AuthMiddleware)
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized: Role not found"})
			c.Abort()
			return
		}

		// Check if role is admin
		if role.(string) != "admin" {
			c.JSON(403, gin.H{"error": "Forbidden: Admin access required"})
			c.Abort()
			return
		}

		// Add debug log
		// fmt.Printf("Admin access granted for role: %v\n", role)

		c.Next()
	}
}
