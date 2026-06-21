package handler

import (
	"errors"
	"net/http"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase"
	"github.com/gin-gonic/gin"
)

type KYCHandler struct {
	Kyc usecase.KYCUseCase
}

func NewKYCHandler(kyc usecase.KYCUseCase) *KYCHandler {
	return &KYCHandler{Kyc: kyc}
}

func (h *KYCHandler) CreateKYC(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("user_id")

	var req request.KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	if err := h.Kyc.CreateKYC(ctx, userID, &req); err != nil {
		if errors.Is(err, errmessage.ErrKYCAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Failed to create user kyc: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user kyc: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "KYC created successfully",
	})
}

func (h *KYCHandler) ApprovedKYC(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("user_id")

	if err := h.Kyc.ApproveKYC(ctx, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to approve user kyc: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "KYC approved",
	})
}

func (h *KYCHandler) GetKYC(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("user_id")

	msg, err := h.Kyc.GetKYC(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user kyc: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
	})
}
