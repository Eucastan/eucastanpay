package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {

	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"message": "internal server error",
			},
		)
	})
}
