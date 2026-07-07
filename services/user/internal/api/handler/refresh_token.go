package handler

import (
	"net/http"

	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/gin-gonic/gin"
)

// RefreshToken godoc
//
// @Summary Refresh Token
// @Tags Authentication
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Success 200 {object} response.RefreshTokenResponse
//
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	token, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Token not found in context"})
		return
	}

	accessToken, refreshToken, err := h.User.RefreshToken(ctx, token.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to refresh token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
