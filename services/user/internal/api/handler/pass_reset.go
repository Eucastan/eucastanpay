package handler

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/gin-gonic/gin"
)

// ForgotPassword godoc
//
// @Summary Forgot Password
// @Tags Authentication
//
// @Accept json
// @Produce json
//
// @Param request body request.ForgotPasswordRequest true "Forgot Password"
//
// @Success 200 {object} response.MessageResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /forgot-password [post]
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid client request: " + err.Error()})
		return
	}

	if err := h.User.ForgotPassword(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	msg := &response.MessageResponse{
		Message: "Reset link sent to your email",
	}

	c.JSON(http.StatusOK, msg)

}

// ResetPassword godoc
//
// @Summary Reset Password
// @Tags AUTH
//
// @Accept json
// @Produce json
//
// @Param request body request.ResetPasswordRequest true "Reset Password"
//
// @Success 200 {object} response.MessageResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid client request: " + err.Error()})
		return
	}

	if err := h.User.ResetPassword(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to change password: " + err.Error(),
		})
		return
	}

	msg := &response.MessageResponse{
		Message: "Password reset successful",
	}

	c.JSON(http.StatusOK, msg)
}
