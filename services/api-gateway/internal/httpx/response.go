package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Created(c *gin.Context, data any) {

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "created",
		Data:    data,
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK,
		APIResponse{
			Success: true,
			Message: "success",
			Data:    data,
		},
	)
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status,
		APIResponse{
			Success: false,
			Message: message,
		},
	)
}

func ValidationError(c *gin.Context, err interface{}) {
	c.JSON(http.StatusBadRequest,
		APIResponse{
			Success: false,
			Message: "validation failed",
			Error:   err,
		},
	)
}
