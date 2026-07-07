package handler

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/gin-gonic/gin"
)

// VerifyUserEmail godoc
//
// @Summary Verify User Email
// @Tags Authentication
//
// @Accept json
// @Produce json
//
// @Param token query string true "Verification Token"
//
// @Success 200 {object} response.MessageResponse
//
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /verify-email [post]
func (h *UserHandler) VerifyUserEmail(c *gin.Context) {
	ctx := c.Request.Context()

	token := c.Query("token")

	if err := h.User.VerifyEmail(ctx, token); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to verify email: " + err.Error()})
		return
	}

	msg := &response.MessageResponse{
		Message: "Email verification successful",
	}

	c.JSON(http.StatusOK, msg)
}
