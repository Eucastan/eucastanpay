package handler

import (
	"errors"
	"net/http"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	Transfer  usecase.TransferUseCase
	telemetry *telemetry.Telemetry
}

func NewTransferHandler(transfer usecase.TransferUseCase, telemetry *telemetry.Telemetry) *TransferHandler {
	return &TransferHandler{Transfer: transfer, telemetry: telemetry}
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
// @Param request body request.TransferRequest true "Transfer Request"
//
// @Success 201 {object} response.UserTransferResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /transfers [post]
func (h *TransferHandler) TransferFromUser(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := h.telemetry.Start(ctx, "TransferHandler.TransferFromUser")
	defer span.End()

	token := c.GetString("token")

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "user_id not found in context",
		})
		return
	}

	userID, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "invalid user_id format",
		})
		return
	}

	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Idempotency-Key header is required",
		})
		return
	}

	ctx = interceptor.AppendJWTToContext(ctx, token)

	var req request.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.Transfer.TransferFromUser(ctx, userID, idemKey, &req)
	if err != nil {
		if errors.Is(err, errmessage.ErrUserNotOwner) {
			c.JSON(http.StatusForbidden, response.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		if errors.Is(err, errmessage.ErrDuplicateRequest) {
			c.JSON(http.StatusForbidden, response.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.UserTransferResponse{
		Message: "Transfer initiated successfully",
		Data:    *resp,
	})
}

// GetAllTransfer godoc
//
// @Summary List transfers
// @Tags Transfer
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Success 200 {array} response.TransferResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /transfers [get]
func (h *TransferHandler) GetAllTransfer(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := h.Transfer.GetAllTransfers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
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
// @Success 200 {object} response.TransferResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /transfers/{id} [get]
func (h *TransferHandler) GetTransfer(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	resp, err := h.Transfer.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errmessage.ErrTranferNotFound) {
			c.JSON(http.StatusNotFound, response.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
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
// @Success 200 {object} response.UserTransferResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /transfers/{reference} [post]
func (h *TransferHandler) ReverseTransfer(c *gin.Context) {
	ctx := c.Request.Context()

	originalRef := c.Param("reference")

	token := c.GetString("token")

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "user_id not found in context",
		})
		return
	}

	userID, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "invalid user_id format",
		})
		return
	}

	ctx = interceptor.AppendJWTToContext(ctx, token)

	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Idempotency-Key header is required",
		})
		return
	}

	resp, err := h.Transfer.ReverseTransfer(ctx, userID, originalRef, idemKey)
	if err != nil {
		if errors.Is(err, errmessage.ErrAlreadyReversed) {
			c.JSON(http.StatusConflict, response.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		if errors.Is(err, errmessage.ErrCannotReverseNonSuccessfulTransfer) {
			c.JSON(http.StatusConflict, response.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.UserTransferResponse{
		Message: "Reversal initiated successfully",
		Data:    *resp,
	})
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
// @Param request body request.ReconciliationRequest true "Reconciliation Request"
// @Param account_id path string true "Account Reconciliation"
//
// @Success 200 {object} response.MessageResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /account/{account_id}/reconcile [post]
func (h *TransferHandler) ReconcileAccount(c *gin.Context) {
	ctx := c.Request.Context()

	token := c.GetString("token")
	accID := c.Param("account_id")

	ctx = interceptor.AppendJWTToContext(ctx, token)

	var req request.ReconciliationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.Transfer.ReconcileAccount(ctx, accID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Message: "Reconciled",
	})
}
