package handler

import (
	"strconv"

	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/application/service"
	auditReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/audit"
	auditResp "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/response/audit"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/httpx"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	*proxy.Base
	auditApp *service.AuditApplication
}

func NewAuditHandler(base *proxy.Base, auditApp *service.AuditApplication) *AuditHandler {
	return &AuditHandler{
		Base:     base,
		auditApp: auditApp,
	}
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
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/audits/search [get]
func (h *AuditHandler) SearchAuditLogs(c *gin.Context) {
	ctx := proxy.Context(c)

	filter := &auditReq.Filter{
		CorrelationID: c.Query("correlation_id"),
		Reference:     c.Query("reference"),
		EventType:     c.Query("event_type"),
		Limit:         50,
		Offset:        0,
	}

	if minStr := c.Query("min_amount"); minStr != "" {
		min, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			httpx.ValidationError(c, err)
			return
		}
		filter.MinAmount = min
	}

	if maxStr := c.Query("max_amount"); maxStr != "" {
		max, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			httpx.ValidationError(c, err)
			return
		}
		filter.MaxAmount = max
	}

	res, err := h.auditApp.SearchAuditLogs(ctx, filter)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, auditResp.SearchResult{
		Count: len(res.Data),
		Data:  res.Data,
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
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/audits/{id} [get]
func (h *AuditHandler) GetAuditRead(c *gin.Context) {
	uri, err := httpx.BindURI[auditReq.AuditIdRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	res, err := h.auditApp.GetAuditByID(proxy.Context(c), uri.AuditID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, res)
}

// GetAllAuditReads godoc
//
// @Summary List AuditReads
// @Tags Audit
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Success 200 {array} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/audits [get]
func (h *AuditHandler) GetAllAuditReads(c *gin.Context) {
	res, err := h.auditApp.GetAllAudits(proxy.Context(c))
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	httpx.Success(c, res)
}
