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
// @Router /users/{user_id}/status [put]
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

// GetAllUsers godoc
//
// @Summary List All Users
// @Tags User
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Success 200 {array} response.UserResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.User.GetAllUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
//
// @Summary Get User
// @Tags User
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Success 200 {object} response.UserResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /user [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "user ID not found in context",
		})
		return
	}

	userID, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "invalid user ID format",
		})
		return
	}

	user, err := h.User.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
