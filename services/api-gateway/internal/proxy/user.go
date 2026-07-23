package proxy

import (
	"github.com/gin-gonic/gin"
)

func UserID(c *gin.Context) string {
	return c.GetString("user_id")
}

func AdminID(c *gin.Context) string {
	return c.GetString("admin_id")
}

func IdemKey(c *gin.Context) string {

	v, _ := c.Get("idempotency_key")
	if v == nil {
		return ""
	}

	return v.(string)
}

func Token(c *gin.Context) string {
	return c.GetString("token")
}

func Email(c *gin.Context) string {
	return c.GetString("email")
}

func Role(c *gin.Context) string {
	return c.GetString("role")
}
