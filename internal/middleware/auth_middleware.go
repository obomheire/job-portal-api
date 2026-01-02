package middleware

import (
	"net/http"
	"strings"

	"job-portal-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userIdStr, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		c.Set("user_id", userIdStr)

		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}

		if isAdmin, ok := claims["is_admin"].(bool); ok {
			c.Set("is_admin", isAdmin)
		}

		c.Next()
	}
}
