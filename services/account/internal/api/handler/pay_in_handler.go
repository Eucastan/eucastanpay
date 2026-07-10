package handler

import (
	"errors"
	"net/http"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/response"
	"github.com/gin-gonic/gin"
)

// InitiatePayIn godoc
//
// @Summary Pay into owner's account
// @Tags Account
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param request body request.DepositRequest true "Deposit information"
// @Param id path string true "Account ID"
//
// @Success 200 {object} response.MessageResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /account/{id}/pay-in [put]
func (h *AccountHandler) InitiatePayIn(c *gin.Context) {
	ctx := c.Request.Context()

	accID := c.Param("id")

	var req request.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid inputs from client: " + err.Error(),
		})
		return
	}

	if err := h.AccUC.DepositAccount(ctx, accID, &req); err != nil {
		if errors.Is(err, errmessage.ErrNotEligibleForOperation) {
			c.JSON(http.StatusForbidden,
				response.ErrorResponse{
					Error: err.Error(),
				},
			)
			return
		}

		if errors.Is(err, errmessage.ErrAccNotFound) {
			c.JSON(http.StatusNotFound,
				response.ErrorResponse{
					Error: err.Error(),
				},
			)
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to create make payment: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Message: "Credited successfully",
	})
}
