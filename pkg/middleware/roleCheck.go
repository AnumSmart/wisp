package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleCheckMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем пользователя из контекста (предполагаем, что он был добавлен в auth middleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role information missing"})
			return
		}

		if requiredRole != userRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		c.Next()
	}
}
