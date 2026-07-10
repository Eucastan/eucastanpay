package handler

import (
	"net/http"
	"strconv"

	"github.com/Eucastan/eucastanpay/services/audit/internal/dto/response"
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

// SearchAuditLogs godoc
//
// @Summary Search AuditRead
// @Tags Audit
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param event_type query string false "Event type"
// @Param correlation_id query string false "Correlation ID"
// @Param reference query string false "Reference"
// @Param min_amount query integer false "Minimum amount"
// @Param max_amount query integer false "Maximum amount"
//
// @Success 200 {array} response.ReadResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /audit/search [get]
func (h *AuditHandler) SearchAuditLogs(c *gin.Context) {
	filter := postgres.Filter{
		CorrelationID: c.Query("correlation_id"),
		Reference:     c.Query("reference"),
		EventType:     c.Query("event_type"),
		Limit:         50,
		Offset:        0,
	}

	if minStr := c.Query("min_amount"); minStr != "" {
		min, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Error: "invalid min_amount",
			})
			return
		}
		filter.MinAmount = min
	}

	if maxStr := c.Query("max_amount"); maxStr != "" {
		max, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Error: "invalid max_amount",
			})
			return
		}
		filter.MaxAmount = max
	}

	res, err := h.UC.Search(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.ReadResponse{
		Data:  res,
		Count: len(res),
	})
}

// GetAuditRead godoc
//
// @Summary Get AuditRead by ID
// @Tags Audit
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Audit Read ID"
//
// @Success 200 {object} response.AuditReadResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /audit/{id} [get]
func (h *AuditHandler) GetAuditRead(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "invalid query, must provide id",
		})
		return
	}

	res, err := h.UC.GetAuditReadByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
