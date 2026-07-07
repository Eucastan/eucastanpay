package api

import (
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/user/config"
	"github.com/Eucastan/eucastanpay/services/user/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, user *handler.UserHandler, kyc *handler.KYCHandler, cfg *config.Config) {

	public := r.Group("api/v1/public")
	{
		public.POST("/register", user.RegisterUser)
		public.POST("/verify-email", user.VerifyUserEmail)
		public.POST("/forgot-password", user.ForgotPassword)
		public.POST("/reset-password", user.ResetPassword)

		public.POST("/login", user.LoginUser)
	}

	auth := r.Group("api/v1")
	auth.Use(middleware.Auth(cfg.JWTSecret))
	{
		auth.POST("/refresh", middleware.RequireRole("user"), user.RefreshToken)
		auth.PUT("/current-user/:user_id", middleware.RequireRole("admin", "super_admin"), user.UserCurrentStaus)
		auth.POST("/logout-all", middleware.RequireRole("user"), user.LogoutAllUsers)
		auth.POST("/logout", middleware.RequireRole("user"), user.Logout)

		auth.POST("/kyc/:user_id/create", middleware.RequireRole("user"), kyc.CreateKYC)
		auth.GET("/kyc/:user_id", middleware.RequireRole("user"), kyc.GetKYC)
		auth.PATCH("/kyc/:user_id/approve", middleware.RequireRole("user"), kyc.ApprovedKYC)
	}
}
