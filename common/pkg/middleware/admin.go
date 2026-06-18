package middleware

import (
	"net/http"
	"strings"

	"github.com/Eucastan/eucastanpay/common/pkg/auth"
	"github.com/gin-gonic/gin"
)

func AdminAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		claims, err := auth.ValidateAdminToken(parts[1], secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid admin token"})
			return
		}

		// Set admin context
		c.Set("admin_id", claims.AdminID)
		c.Set("admin_role", claims.Role)
		c.Set("is_admin", true)
		c.Set("token", parts[1])

		c.Next()
	}
}

// RequireAdminRole checks for specific admin roles
func RequireAdminRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("admin_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin role not found"})
			return
		}

		currentRole := roleVal.(string)

		for _, allowed := range roles {
			if allowed == currentRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient admin privileges"})
	}
}
