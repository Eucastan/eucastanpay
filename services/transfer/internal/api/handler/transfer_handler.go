package handler

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/dto/request"
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

func (h *TransferHandler) TransferFromUser(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := h.telemetry.Start(ctx, "TransferHandler.TransferFromUser")
	defer span.End()

	token := c.GetString("token")

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

	ctx = interceptor.AppendJWTToContext(ctx, token)

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

	token := c.GetString("token")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	ctx = interceptor.AppendJWTToContext(ctx, token)

	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Idempotency-Key header is required",
		})
		return
	}

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
