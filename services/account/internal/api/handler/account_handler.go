package handler

import (
	"fmt"
	"net/http"

	"github.com/Eucastan/eucastanpay/services/account/internal/dto/request"
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

func (h *AccountHandler) OpenAccount(c *gin.Context) {
	ctx := c.Request.Context()

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	userID, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID format"})
		return
	}

	var req request.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Invalid request from client")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request from client"})
		return
	}

	resp, err := h.AccUC.CreateAccount(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open account"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AccountHandler) GetBalance(c *gin.Context) {
	ctx := c.Request.Context()

	accID := c.Param("id")

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	resp, err := h.AccUC.GetBalance(ctx, accID, userID.(string))
	if err != nil {
		fmt.Printf("GET BALANCE ERROR: %+v\n", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
