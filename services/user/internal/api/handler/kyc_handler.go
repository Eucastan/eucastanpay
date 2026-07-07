package handler

import (
	"errors"
	"net/http"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase"
	"github.com/gin-gonic/gin"
)

type KYCHandler struct {
	Kyc usecase.KYCUseCase
}

func NewKYCHandler(kyc usecase.KYCUseCase) *KYCHandler {
	return &KYCHandler{Kyc: kyc}
}

// CreateKYC godoc
//
// @Summary Submit KYC information
// @Tags KYC
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param user_id path string true "User ID"
// @Param request body request.KYCRequest true "KYC Information"
//
// @Success 201 {object} response.MessageResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /kyc/{user_id}/create [post]
func (h *KYCHandler) CreateKYC(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("user_id")

	var req request.KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "invalid request: " + err.Error(),
		})
		return
	}

	if err := h.Kyc.CreateKYC(ctx, userID, &req); err != nil {
		if errors.Is(err, errmessage.ErrKYCAlreadyExists) {
			c.JSON(http.StatusConflict, response.ErrorResponse{
				Error: "Failed to create user kyc: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to create user kyc: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.MessageResponse{
		Message: "KYC created successfully",
	})
}
