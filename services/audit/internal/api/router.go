package api

import (
	"github.com/Eucastan/eucastanpay/services/audit/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, h *handler.AuditHandler) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/audit/search", h.Search)
	}
}
