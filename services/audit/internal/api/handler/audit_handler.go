package handler

import (
	"net/http"
	"strconv"

	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/audit/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	UC usecase.AuditUseCase
}

func NewAuditHandler(uc usecase.AuditUseCase) *AuditHandler {
	return &AuditHandler{UC: uc}
}

func (h *AuditHandler) Search(c *gin.Context) {
	filter := postgres.Filter{
		CorrelationID: c.Query("correlation_id"),
		Reference:     c.Query("reference"),
		EventType:     c.Query("event_type"),
		Limit:         50,
		Offset:        0,
	}

	if minStr := c.Query("min_amount"); minStr != "" {
		filter.MinAmount, _ = strconv.ParseInt(minStr, 10, 64)
	}
	if maxStr := c.Query("max_amount"); maxStr != "" {
		filter.MaxAmount, _ = strconv.ParseInt(maxStr, 10, 64)
	}

	res, err := h.UC.Search(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  res,
		"count": len(res),
	})
}
