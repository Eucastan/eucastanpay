package handler

import (
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/application/service"
	accountreq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/account"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/httpx"

	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	*proxy.Base
	accountApp *service.AccountApplication
}

func NewAccountHandler(base *proxy.Base, accountApp *service.AccountApplication) *AccountHandler {
	return &AccountHandler{
		Base:       base,
		accountApp: accountApp,
	}
}

// InitiatePayIn godoc
//
// @Summary Pay into owner's account
// @Tags Account
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param request body accountreq.DepositRequest true "Deposit information"
// @Param id path string true "Account ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/account/{id}/pay-in [put]
func (h *AccountHandler) InitiatePayIn(c *gin.Context) {

	uri, err := httpx.BindURI[accountreq.AccountURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	req, err := httpx.BindJSON[accountreq.DepositRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	req.AccountID = uri.AccountID

	resp, err := h.accountApp.Deposit(proxy.Context(c), &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// WithDrawal godoc
//
// @Summary Withdrawing Cash from account
// @Tags Account
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param request body accountreq.DepositRequest true "Deposit information"
// @Param id path string true "Account ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 409 {object} httpx.APIResponse "Insufficient balance"
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/account/{id}/withdraw [put]
func (h *AccountHandler) WithDrawal(c *gin.Context) {

	uri, err := httpx.BindURI[accountreq.AccountURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	req, err := httpx.BindJSON[accountreq.DepositRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	req.AccountID = uri.AccountID

	resp, err := h.accountApp.WithDraw(proxy.Context(c), &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
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
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/account/{id}/balance [get]
func (h *AccountHandler) GetBalance(c *gin.Context) {

	uri, err := httpx.BindURI[accountreq.AccountURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	userID := proxy.UserID(c)

	accReq := accountreq.GetBalanceRequest{
		AccountID: uri.AccountID,
		UserID:    userID,
	}

	resp, err := h.accountApp.GetBalance(proxy.Context(c), &accReq)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
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
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /account/{id} [get]
func (h *AccountHandler) GetUserAccount(c *gin.Context) {

	userID := proxy.UserID(c)

	resp, err := h.accountApp.GetUserAccount(proxy.Context(c), userID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// GetAllUsersAccount godoc
//
// @Summary Get Account Details for all users
// @Tags Account
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/accounts [get]
func (h *AccountHandler) GetAllUsersAccount(c *gin.Context) {
	resp, err := h.accountApp.GetAllAccounts(proxy.Context(c))
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}
