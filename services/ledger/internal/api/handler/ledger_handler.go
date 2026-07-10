package handler

import (
	"errors"
	"net/http"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/usecase"
	"github.com/gin-gonic/gin"
)

type LedgerHandler struct {
	ledger usecase.LedgerUseCase
}

func NewLedgerHandler(ledger usecase.LedgerUseCase) *LedgerHandler {
	return &LedgerHandler{ledger: ledger}
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
// @Success 200 {array} response.LedgerResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /ledgers [get]
func (h *LedgerHandler) GetAllLedgers(c *gin.Context) {
	ctx := c.Request.Context()

	if entryType := c.Query("entry_type"); entryType != "" {
		ledgers, err := h.ledger.GetTransactionByEntryType(
			ctx,
			&request.EntryTypeRequest{
				EntryType: entryType,
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, ledgers)
		return
	}

	ledgers, err := h.ledger.GetAllLedgers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ledgers)
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
// @Success 200 {object} response.LedgerResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /ledgers/{id} [get]
func (h *LedgerHandler) GetLedger(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "ledger id is required",
		})
		return
	}

	ledger, err := h.ledger.GetTransactionEntry(ctx, id)
	if err != nil {
		if errors.Is(err, errmessage.ErrLedgerNotFound) {
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

	c.JSON(http.StatusOK, ledger)
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
// @Success 200 {object} response.AccountBalanceResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /accounts/{account_id}/balance [get]
func (h *LedgerHandler) GetAccountBalance(c *gin.Context) {
	ctx := c.Request.Context()
	accountID := c.Param("account_id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Must provide account id",
		})
		return
	}

	balance, err := h.ledger.GetAccountBalance(ctx, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "failed to get balance: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.AccountBalanceResponse{
		AccountID: accountID,
		Balance:   balance,
	})
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
// @Success 200 {object} response.ReconciliationResultResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /accounts/{account_id}/reconcile [get]
func (h *LedgerHandler) ReconciliationResult(c *gin.Context) {
	ctx := c.Request.Context()

	accountID := c.Param("account_id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Must provide account id",
		})
		return
	}

	token := c.GetString("token")

	ctx = interceptor.AppendJWTToContext(ctx, token)

	result, err := h.ledger.ReconcileAccount(ctx, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.ReconciliationResultResponse{
		Data: *result,
	})
}
