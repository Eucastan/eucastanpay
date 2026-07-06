package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
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

func (h *UserHandler) RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid client request: " + err.Error(),
		})
		return
	}

	user, err := h.User.Register(ctx, &req)
	if err != nil {
		if errors.Is(err, errmessage.ErrDuplicateEmail) {
			c.JSON(http.StatusConflict, gin.H{"error": "Duplicate email: " + err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "User created successfully",
		"response": user,
	})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client request"})
		return
	}

	user, err := h.User.Login(ctx, &req)
	if err != nil {
		if errors.Is(err, errmessage.ErrInvalidCredentials) {
			c.JSON(http.StatusConflict, gin.H{"error": "Invalid credentials: " + err.Error()})
			return
		}

		if errors.Is(err, errmessage.ErrPasswordNotConfirmed) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password mis-match: " + err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User login successful",
		"response": user,
	})
}

func (h *UserHandler) VerifyUserEmail(c *gin.Context) {
	ctx := c.Request.Context()

	token := c.Query("token")

	if err := h.User.VerifyEmail(ctx, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email: " + err.Error()})
		return
	}

	msg := &response.MessageResponse{
		Message: "Email verification successful",
	}

	c.JSON(http.StatusOK, msg)
}

func (h *UserHandler) UserCurrentStaus(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("id")
	role, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found in context"})
		return
	}

	var req request.CurrentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client request: " + err.Error()})
		return
	}

	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized Access"})
		return
	}

	msg, err := h.User.UserCurrentStatus(ctx, userID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, msg)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	token, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found in context"})
		return
	}

	accessToken, refreshToken, err := h.User.RefreshToken(ctx, token.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refresh token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access":  accessToken,
		"refresh": refreshToken,
	})
}

func (h *UserHandler) LogoutAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	if err := h.User.LogoutAllUsers(ctx, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User logout successful",
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	token, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found in context"})
		return
	}

	if err := h.User.Logout(ctx, token.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User logout successful",
	})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client request: " + err.Error()})
		return
	}

	if err := h.User.ForgotPassword(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := &response.MessageResponse{
		Message: "Reset link sent to your email",
	}

	c.JSON(http.StatusOK, msg)

}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client request: " + err.Error()})
		return
	}

	if err := h.User.ResetPassword(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to change password: " + err.Error(),
		})
		return
	}

	msg := &response.MessageResponse{
		Message: "Password reset successful",
	}

	c.JSON(http.StatusOK, msg)
}
