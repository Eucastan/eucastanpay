package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleVal, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "no role found"})
			c.Abort()
			return
		}

		currentRole := roleVal.(string)

		for _, r := range roles {
			if r == currentRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}
