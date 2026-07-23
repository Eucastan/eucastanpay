package handler

import (
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/application/service"
	adminReq "github.com/Eucastan/eucastanpay/services/api-gateway/internal/dto/request/admin"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/httpx"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	*proxy.Base
	adminApp *service.AdminApplication
}

func NewAdminHandler(base *proxy.Base, adminApp *service.AdminApplication) *AdminHandler {
	return &AdminHandler{
		Base:     base,
		adminApp: adminApp,
	}
}

// CreateBootstrapAdmin godoc
// @Summary Register a new admin
// @Description Creates a new EucastanPay admin account.
// @Tags Admin Auth
//
// @Accept json
// @Produce json
//
// @Param request body adminReq.CreateAdminRequest true "Admin Registration"
//
// @Success 201 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 409 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /auth/register [post]
func (h *AdminHandler) CreateBootstrapAdmin(c *gin.Context) {

	req, err := httpx.BindJSON[adminReq.CreateAdminRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	resp, err := h.adminApp.CreateAdmin(proxy.Context(c), &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Created(c, resp)
}

// CreateAdmin godoc
// @Summary Register a new admin
// @Description Creates a new EucastanPay admin account.
// @Tags Admin
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
//
// @Param request body adminReq.CreateAdminRequest true "Admin Registration"
//
// @Success 201 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 409 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin [post]
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	req, err := httpx.BindJSON[adminReq.CreateAdminRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	resp, err := h.adminApp.CreateAdmin(proxy.Context(c), &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Created(c, resp)
}

// Login godoc
//
// @Summary Login admin
// @Description Authenticates an admin and returns JWT tokens.
//
// @Tags Admin Auth
//
// @Accept json
// @Produce json
//
// @Param request body adminReq.AdminLoginRequest true "Login"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /auth/login [post]
func (h *AdminHandler) Login(c *gin.Context) {
	req, err := httpx.BindJSON[adminReq.AdminLoginRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	resp, err := h.adminApp.Login(proxy.Context(c), &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// GetAdminProfile godoc
//
// @Summary Get Admin
// @Tags Admin
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
// @Router /admin [get]
func (h *AdminHandler) GetAdminProfile(c *gin.Context) {
	adminID := proxy.AdminID(c)

	resp, err := h.adminApp.GetAdmin(proxy.Context(c), adminID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// GetAdmin godoc
//
// @Summary Get Admin
// @Tags Admin
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
// @Param id path string true "Admin ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/{id} [get]
func (h *AdminHandler) GetAdmin(c *gin.Context) {
	uri, err := httpx.BindURI[adminReq.AdminURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	resp, err := h.adminApp.GetAdmin(proxy.Context(c), uri.AdminID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// ListAdmins godoc
//
// @Summary List All Admins
// @Tags Admin
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin [get]
func (h *AdminHandler) ListAdmins(c *gin.Context) {
	query, err := httpx.BindQuery[adminReq.Pagination](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	resp, err := h.adminApp.GetAllAdmins(proxy.Context(c), query.Limit, query.Page)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// UpdateAdmin godoc
//
// @Summary Editing Admin information
// @Tags Admin
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
// @Param id path string true "Admin ID"
// @Param request body adminReq.UpdateAdminRequest true "Update Admin Details"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 404 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/{id} [put]
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	uri, err := httpx.BindURI[adminReq.AdminURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	req, err := httpx.BindJSON[adminReq.UpdateAdminRequest](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	resp, err := h.adminApp.UpdateAdmin(proxy.Context(c), uri.AdminID, &req)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}

// DeleteAdmin godoc
//
// @Summary Deleting Admin information
// @Tags Admin
//
// @Security BearerAuth
//
// @Accept json
// @Produce json
// @Param id path string true "Admin ID"
//
// @Success 200 {object} httpx.APIResponse
//
// @Failure 400 {object} httpx.APIResponse
// @Failure 401 {object} httpx.APIResponse
// @Failure 403 {object} httpx.APIResponse
// @Failure 500 {object} httpx.APIResponse
//
// @Router /admin/{id} [delete]
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	uri, err := httpx.BindURI[adminReq.AdminURI](c)
	if err != nil {
		httpx.ValidationError(c, err)
		return
	}

	resp, err := h.adminApp.DeleteAdmin(proxy.Context(c), uri.AdminID)
	if err != nil {
		httpx.HandleGRPCError(c, err)
		return
	}

	httpx.Success(c, resp)
}
