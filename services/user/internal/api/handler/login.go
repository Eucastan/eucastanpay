package handler

import (
	"errors"
	"net/http"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/gin-gonic/gin"
)

// LoginUser godoc
//
// @Summary Login user
// @Description Authenticates a user and returns JWT tokens.
//
// @Tags Authentication
//
// @Accept json
// @Produce json
//
// @Param request body request.LoginRequest true "Login"
//
// @Success 200 {object} response.LoginResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /public/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid client request"})
		return
	}

	user, err := h.User.Login(ctx, &req)
	if err != nil {
		if errors.Is(err, errmessage.ErrInvalidCredentials) {
			c.JSON(http.StatusConflict, response.ErrorResponse{Error: "Invalid credentials: " + err.Error()})
			return
		}

		if errors.Is(err, errmessage.ErrPasswordNotConfirmed) {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Password mis-match: " + err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Failed to login user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		Message:  "User login successful",
		Response: *user,
	})
}
