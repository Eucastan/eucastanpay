package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/Eucastan/eucastanpay/common/pkg/auth"
	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "auth header empty"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, err := auth.ValidateToken(parts[1], secret)
		if err != nil {
			log.Println("TOKEN ERROR:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		log.Printf("CLAIMS: %+v\n", claims)

		c.Set("token", parts[1])
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		log.Printf(
			"user=%s email=%s role=%s\n",
			claims.UserID,
			claims.Email,
			claims.Role,
		)

		c.Next()
	}
}
