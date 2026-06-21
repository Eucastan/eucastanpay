package handler

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransferHandler struct {
	Transfer usecase.TransferUseCase
}

func NewTransferHandler(transfer usecase.TransferUseCase) *TransferHandler {
	return &TransferHandler{Transfer: transfer}
}

func (h *TransferHandler) TransferFromUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Idempotency-Key header is required",
		})
		return
	}

	var req request.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Transfer.TransferFromUser(ctx, userID.(string), idemKey, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transfer initiated successfully",
		"data":    resp,
	})
}

func (h *TransferHandler) GetTransfer(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	resp, err := h.Transfer.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

func (h *TransferHandler) ReverseTransfer(c *gin.Context) {
	ctx := c.Request.Context()
	originalRef := c.Param("reference")
	idemKey := uuid.NewString()

	userID, _ := c.Get("user_id")

	resp, err := h.Transfer.ReverseTransfer(ctx, userID.(string), originalRef, idemKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reversal initiated successfully",
		"data":    resp,
	})
}
