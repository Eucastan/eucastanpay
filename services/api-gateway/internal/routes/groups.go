package routes

import (
	commonmw "github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/api-gateway/config"
	"github.com/gin-gonic/gin"
)

func UserGroup(r *gin.Engine, cfg *config.Config) *gin.RouterGroup {

	g := r.Group("/")
	g.Use(commonmw.Auth(cfg.SharedCfg.JWTSecret))
	g.Use(commonmw.RequireRole("user"))

	return g
}

func AdminGroup(r *gin.Engine, cfg *config.Config) *gin.RouterGroup {

	g := r.Group("/admin")
	g.Use(commonmw.AdminAuth(cfg.SharedCfg.JWTSecret))
	g.Use(commonmw.RequireAdminRole("admin"))

	return g
}
