package handler

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/services/ledger/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/usecase"
	"github.com/gin-gonic/gin"
)

type LedgerHandler struct {
	ledger usecase.LedgerUseCase
}

func NewLedgerHandler(ledger usecase.LedgerUseCase) *LedgerHandler {
	return &LedgerHandler{ledger: ledger}
}

func (h *LedgerHandler) GetLedger(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	ledger, err := h.ledger.GetTransactionEntry(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ledger entry not found"})
		return
	}

	c.JSON(http.StatusOK, ledger)
}

func (h *LedgerHandler) GetLedgerByEntry(c *gin.Context) {
	ctx := c.Request.Context()
	entryType := c.Query("entry_type")

	if entryType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entry_type query parameter is required"})
		return
	}

	ledgers, err := h.ledger.GetTransactionByEntryType(ctx, &request.EntryTypeRequest{
		EntryType: entryType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ledgers)
}

func (h *LedgerHandler) GetAccountBalance(c *gin.Context) {
	ctx := c.Request.Context()
	accountID := c.Param("account_id")

	balance, err := h.ledger.GetAccountBalance(ctx, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}

	c.JSON(http.StatusOK, balance)
}
