package handler

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/gin-gonic/gin"
)

// ApprovedKYC godoc
//
// @Summary Approved KYC
// @Tags KYC
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param user_id path string true "User ID"
//
// @Success 200 {object} response.MessageResponse
//
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /kyc/{user_id}/approve [patch]
func (h *KYCHandler) ApprovedKYC(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("user_id")

	if err := h.Kyc.ApproveKYC(ctx, userID); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to approve user kyc: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Message: "KYC approved",
	})
}
