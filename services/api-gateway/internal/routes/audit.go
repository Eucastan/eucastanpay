package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuditRoutes(r *gin.Engine, deps Dependencies, h Handlers) {

	admin := AdminGroup(r, deps.Config)

	{
		admin.GET("/audits", h.Audit.GetAllAuditReads)
		admin.GET("/audits/search", h.Audit.SearchAuditLogs)
		admin.GET("/audits/:id", h.Audit.GetAuditRead)
	}
}
