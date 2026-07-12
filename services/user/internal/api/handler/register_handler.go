package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description Creates a new EucastanPay user account.
// @Tags Authentication
//
// @Accept json
// @Produce json
//
// @Param request body request.RegisterRequest true "User Registration"
//
// @Success 201 {object} response.RegisterResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /public/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "invalid client request: " + err.Error(),
		})
		return
	}

	user, err := h.User.Register(ctx, &req)
	if err != nil {
		if errors.Is(err, errmessage.ErrDuplicateEmail) {
			c.JSON(http.StatusConflict, response.ErrorResponse{Error: "Duplicate email: " + err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, response.RegisterResponse{
		Message: "User created successfully",
		User: response.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Phone:     user.Phone,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		},
	})
}
