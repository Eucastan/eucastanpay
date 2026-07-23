package handler

import (
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/application/service"
	ledgerReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/ledger"

	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/httpx"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

type LedgerHandler struct {
	*proxy.Base
	ledgerApp *service.LedgerApplication
}

func NewLedgerHandler(base *proxy.Base, ledgerApp *service.LedgerApplication) *LedgerHandler {
	return &LedgerHandler{
		Base:      base,
		ledgerApp: ledgerApp,
	}
}

// GetAllLedgers godoc
//
// @Summary Get All Ledger Entries
// @Tags Ledger
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param entry_type query string false "Entry Type" Enums(credit,debit)
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/ledgers [get]
func (h *LedgerHandler) GetAllLedgers(c *gin.Context) {

	uri, err := httpx.BindQuery[ledgerReq.EntryTypeURI](c)
	if err == nil && uri.EntryType != "" {
		entry := &ledgerReq.EntryTypeRequest{
			EntryType: uri.EntryType,
		}

		ledgers, err := h.ledgerApp.GetLedgersByEntryType(
			proxy.Context(c),
			entry,
		)

		if err != nil {
			httpx.HandleGRPCError(c, err)
			return
		}

		httpx.Success(c, ledgers)
		return
	}

	ledgers, err := h.ledgerApp.GetAllLedgers(proxy.Context(c))
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, ledgers)
}

// GetLedgerByUserID godoc
//
// @Summary Get Ledger by User ID
// @Tags Ledger
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Ledger UserID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /ledgers/me [get]
func (h *LedgerHandler) GetLedgerByUserID(c *gin.Context) {

	userID := proxy.UserID(c)

	ledger, err := h.ledgerApp.GetLedger(proxy.Context(c), userID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, ledger)
}

// GetLedger godoc
//
// @Summary Get Ledger by ID
// @Tags Ledger
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Ledger ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/ledgers/{id} [get]
func (h *LedgerHandler) GetLedger(c *gin.Context) {

	uri, err := httpx.BindURI[ledgerReq.LedgerURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	ledger, err := h.ledgerApp.GetLedger(proxy.Context(c), uri.LedgerID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, ledger)
}

// GetLedgerByAccountID godoc
//
// @Summary Get Ledger by Account ID
// @Tags Ledger
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Ledger ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/ledgers/{account_id} [get]
func (h *LedgerHandler) GetLedgerByAccountID(c *gin.Context) {

	uri, err := httpx.BindURI[ledgerReq.AccountURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	ledger, err := h.ledgerApp.GetLedgerByAccountID(proxy.Context(c), uri.AccountID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, ledger)
}

// GetAccountBalance godoc
//
// @Summary Get Account Balance
// @Description Get Account Balance from Ledger
// @Tags Ledger
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param account_id path string true "Ledger Account ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/accounts/{account_id}/balance [get]
func (h *LedgerHandler) GetAccountBalance(c *gin.Context) {
	uri, err := httpx.BindURI[ledgerReq.AccountURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	balance, err := h.ledgerApp.GetLedgerBalance(proxy.Context(c), uri.AccountID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, balance)
}

// ReconciliationResult godoc
//
// @Summary Get Reconciliation Result
// @Tags Ledger
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param account_id path string true "Ledger Account ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/accounts/{account_id}/reconcile [get]
func (h *LedgerHandler) ReconciliationResult(c *gin.Context) {
	ctx := proxy.Context(c)

	uri, err := httpx.BindURI[ledgerReq.AccountURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	token := c.GetString("token")

	ctx = interceptor.AppendJWTToContext(ctx, token)

	result, err := h.ledgerApp.GetLedgerReconciliation(ctx, uri.AccountID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, result)
}
