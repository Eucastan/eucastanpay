package handler

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/gin-gonic/gin"
)

// GetKYC godoc
//
// @Summary Get KYC
// @Tags KYC
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param user_id path string true "User ID"
//
// @Success 200 {object} response.KYCResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /kyc/{user_id} [get]
func (h *KYCHandler) GetKYC(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("user_id")

	msg, err := h.Kyc.GetKYC(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to get user kyc: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, msg)
}
