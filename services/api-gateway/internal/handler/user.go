package handler

import (
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/application/service"
	userReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/user"
	"google.golang.org/grpc/metadata"

	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/httpx"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	*proxy.Base
	userApp *service.UserApplication
}

func NewUserHandler(base *proxy.Base, userApp *service.UserApplication) *UserHandler {
	return &UserHandler{
		Base:    base,
		userApp: userApp,
	}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Creates a new EucastanPay user account.
// @Tags Authentication
//
// @Accept json
// @Produce json
//
// @Param request body userReq.RegisterRequest true "User Registration"
//
// @Success 201 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 409 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	req, err := httpx.BindJSON[userReq.RegisterRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	res, err := h.userApp.Register(proxy.Context(c), &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Created(c, res)
}

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
// @Param request body userReq.LoginRequest true "Login"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {

	req, err := httpx.BindJSON[userReq.LoginRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	res, err := h.userApp.Login(proxy.Context(c), &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, res)
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
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /user [get]
func (h *UserHandler) GetUser(c *gin.Context) {

	userID := proxy.UserID(c)

	res, err := h.userApp.GetUserByID(proxy.Context(c), userID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, res)
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
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {

	res, err := h.userApp.GetAllUsers(proxy.Context(c))
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, res)
}

// UpdateUser godoc
//
// @Summary Editing User information
// @Tags User
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
// @Param request body userReq.UpdateRequest true "Update Details"
// @Param id path string true "User ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {

	req, err := httpx.BindJSON[userReq.UpdateRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	uri, err := httpx.BindURI[userReq.UserURI](c)

	res, err := h.userApp.UpdateUser(proxy.Context(c), uri.UserID, &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, res)
}

// DeleteUser godoc
//
// @Summary Deleting User information
// @Tags User
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
// @Param id path string true "User ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {

	uri, err := httpx.BindURI[userReq.UserURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	res, err := h.userApp.DeleteUser(proxy.Context(c), uri.UserID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, res)
}

// CreateKYC godoc
//
// @Summary Submit KYC information
// @Tags KYC
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param user_id path string true "User ID"
// @Param request body userReq.CreateKYCRequest true "KYC Information"
//
// @Success 201 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 409 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /users/{user_id}/kyc [post]
func (h *UserHandler) CreateKYC(c *gin.Context) {
	h.Logger.Printf("CreateKYC handler reached")

	ctx := proxy.Context(c)
	md, _ := metadata.FromOutgoingContext(ctx)
	h.Logger.Printf("Outgoing metadata: %+v\n", md)

	req, err := httpx.BindJSON[userReq.CreateKYCRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	_, err = httpx.BindURI[userReq.UserKYCURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	res, err := h.userApp.CreateKYC(ctx, req.IdNumber, req.IdType)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Created(c, res)
}

// ApproveKYC godoc
//
// @Summary Approve user KYC by admin
// @Tags KYC
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param user_id path string true "User ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/{user_id}/kyc [patch]
func (h *UserHandler) ApproveKYC(c *gin.Context) {

	uri, err := httpx.BindURI[userReq.UserKYCURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	res, err := h.userApp.ApproveKYC(proxy.Context(c), uri.UserID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, res)
}

// GetKYC godoc
//
// @Summary Get KYC
// @Tags KYC
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param user_id path string true "User ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /users/{user_id}/kyc [get]
func (h *UserHandler) GetKYC(c *gin.Context) {

	uri, err := httpx.BindURI[userReq.UserKYCURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	res, err := h.userApp.GetKYC(proxy.Context(c), uri.UserID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, res)
}
