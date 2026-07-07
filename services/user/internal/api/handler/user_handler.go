package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase"
)

type UserHandler struct {
	User usecase.UserUseCaseInterface
}

func NewUserHandler(user usecase.UserUseCaseInterface) *UserHandler {
	return &UserHandler{
		User: user,
	}
}

// UserCurrentStatus godoc
//
// @Summary Update a user's account status
// @Tags Authentication
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param request body request.CurrentStatusRequest true "User Current Status"
// @Param user_id path string true "User ID"
//
// @Success 200 {object} response.MessageResponse
//
// @Failure 401 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /current-user/{user_id} [put]
func (h *UserHandler) UserCurrentStaus(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("user_id")

	var req request.CurrentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid client request: " + err.Error()})
		return
	}

	msg, err := h.User.UserCurrentStatus(ctx, userID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Message: msg,
	})
}
