package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
)

// LogoutUser godoc
//
// @Summary Logout user
// @Tags Authentication
//
// @Security BearerAuth
//
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	token, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Token not found in context"})
		return
	}

	if err := h.User.Logout(ctx, token.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to logout: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Message: "User logout successful",
	})
}

// LogoutAllUser godoc
//
// @Summary Logout user from all devices
// @Tags Authentication
//
// @Security BearerAuth
//
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /logout-all [post]
func (h *UserHandler) LogoutAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "user ID not found in context"})
		return
	}

	if err := h.User.LogoutAllUsers(ctx, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to logout: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Message: "User logout successful",
	})
}
