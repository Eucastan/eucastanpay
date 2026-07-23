package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindURI[T any](c *gin.Context) (T, error) {

	var req T

	if err := c.ShouldBindUri(&req); err != nil {

		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return req, err
	}

	return req, nil
}
