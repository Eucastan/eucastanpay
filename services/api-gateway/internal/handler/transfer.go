package handler

import (
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/application/service"
	transferReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/transfer"

	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/httpx"

	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	*proxy.Base
	transferApp *service.TransferApplication
}

func NewTransferHandler(base *proxy.Base, transferApp *service.TransferApplication) *TransferHandler {
	return &TransferHandler{
		Base:        base,
		transferApp: transferApp,
	}
}

// TransferFromUser godoc
// @Summary TransferFromUser Initiate transfer
// @Description Initiates transfer from one user to the other (sender and receiver).
// @Tags Transfer
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param Idempotency-Key header string true "Unique idempotency key"
// @Param request body transferReq.TransferRequest true "Transfer Request"
//
// @Success 201 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 409 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /transfers [post]
func (h *TransferHandler) TransferFromUser(c *gin.Context) {
	ctx := proxy.Context(c)

	ctx, span := h.Telemetry.Start(ctx, "TransferHandler.TransferFromUser")
	defer span.End()

	token := proxy.Token(c)

	userID := proxy.UserID(c)

	idemKey := proxy.IdemKey(c)

	ctx = interceptor.AppendJWTToContext(ctx, token)

	req, err := httpx.BindJSON[transferReq.TransferRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	req.UserID = userID
	req.IdemKey = idemKey

	resp, err := h.transferApp.Transfer(ctx, &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Created(c, resp)
}

// GetAllTransfers godoc
//
// @Summary List transfers
// @Tags Transfer
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/transfers [get]
func (h *TransferHandler) GetAllTransfers(c *gin.Context) {

	resp, err := h.transferApp.GetAllTransfers(proxy.Context(c))
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// GetTransferByUserID godoc
//
// @Summary Get transfer by User ID
// @Tags Transfer
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "User ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /transfers/me [get]
func (h *TransferHandler) GetTransferByUserID(c *gin.Context) {

	userID := proxy.UserID(c)

	resp, err := h.transferApp.GetTransfer(proxy.Context(c), userID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// GetTransfer godoc
//
// @Summary Get transfer by ID
// @Tags Transfer
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Transfer ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/transfers/{id} [get]
func (h *TransferHandler) GetTransfer(c *gin.Context) {

	uri, err := httpx.BindURI[transferReq.TransferURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	resp, err := h.transferApp.GetTransfer(proxy.Context(c), uri.TranferID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// ReverseTransfer godoc
// @Summary Reverse transfer
// @Description Reverses transfer from receiver back to sender (receiver -> sender).
// @Tags Transfer
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param Idempotency-Key header string true "Unique idempotency key"
// @Param reference path string true "Transfer Reference"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 409 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/transfers/{reference} [post]
func (h *TransferHandler) ReverseTransfer(c *gin.Context) {
	ctx := proxy.Context(c)

	uri, err := httpx.BindURI[transferReq.ReverseURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	token := proxy.Token(c)

	adminID := proxy.AdminID(c)

	ctx = interceptor.AppendJWTToContext(ctx, token)

	idemKey := proxy.IdemKey(c)

	resp, err := h.transferApp.ReverseTransfer(ctx, adminID, uri.Reference, idemKey)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// ReconcileAccount godoc
// @Summary Reconcile account
// @Description Reconciliation Account from Transfer service.
// @Tags Transfer
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param request body transferReq.ReconciliationRequest true "Reconciliation Request"
// @Param account_id path string true "Account Reconciliation"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 409 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/accounts/{account_id}/reconcile [post]
func (h *TransferHandler) ReconcileAccount(c *gin.Context) {
	ctx := proxy.Context(c)

	token := proxy.Token(c)

	uri, err := httpx.BindURI[transferReq.AccountIdURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	ctx = interceptor.AppendJWTToContext(ctx, token)

	req, err := httpx.BindJSON[transferReq.ReconciliationRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	resp, err := h.transferApp.ReconcileAccount(ctx, uri.AccountID, req.AccountNo)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}
