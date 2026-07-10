package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/account/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/account/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	AccUC usecase.AccountUseCase
}

func NewAccountHandler(acc usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{
		AccUC: acc,
	}
}

// OpenAccount godoc
// @Summary Account creation for a new user
// @Description Creates a new EucastanPay account.
// @Tags Account
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param request body request.CreateAccountRequest true "Account Creation"
//
// @Success 201 {object} response.AccountResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse "Account already exists"
// @Failure 500 {object} response.ErrorResponse
//
// @Router /account [post]
func (h *AccountHandler) OpenAccount(c *gin.Context) {
	ctx := c.Request.Context()

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "user ID not found in context",
		})
		return
	}

	userID, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Invalid ID format"})
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "email not found in context",
		})
		return
	}

	em, ok := email.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Invalid email format",
		})
		return
	}

	var req request.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Invalid request from client")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid request from client: " + err.Error(),
		})
		return
	}

	resp, err := h.AccUC.CreateAccount(ctx, userID, em, &req)
	if err != nil {
		if errors.Is(err, errmessage.ErrAccAlreadyExists) {
			c.JSON(http.StatusConflict, response.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to open account: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetBalance godoc
//
// @Summary Get Account owner's balance
// @Tags Account
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Account ID"
//
// @Success 200 {object} response.AccountResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /account/{id}/balance [get]
func (h *AccountHandler) GetBalance(c *gin.Context) {
	ctx := c.Request.Context()

	accID := c.Param("id")

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "user ID not found in context",
		})
		return
	}

	userID, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	resp, err := h.AccUC.GetBalance(ctx, accID, userID)
	if err != nil {
		if errors.Is(err, errmessage.ErrAccNotFound) {
			c.JSON(http.StatusNotFound,
				response.ErrorResponse{
					Error: err.Error(),
				},
			)
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserAccount godoc
//
// @Summary Get Account Details with accountID and userID
// @Tags Account
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Account ID"
//
// @Success 200 {object} response.AccountResponse
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
//
// @Router /account/{id} [get]
func (h *AccountHandler) GetUserAccount(c *gin.Context) {
	ctx := c.Request.Context()

	accID := c.Param("id")

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "user ID not found in context",
		})
		return
	}

	userID, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Invalid ID format",
		})
		return
	}

	resp, err := h.AccUC.GetByAccountIDAndUserID(ctx, accID, userID)
	if err != nil {
		if errors.Is(err, errmessage.ErrAccNotFound) {
			c.JSON(http.StatusNotFound,
				response.ErrorResponse{
					Error: err.Error(),
				},
			)
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
