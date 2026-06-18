package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	Log *logrus.Logger
	cfg string
}

func New(log *logrus.Logger, secret string) *Middleware {
	return &Middleware{Log: log, cfg: secret}
}

func (m *Middleware) Recovery() gin.HandlerFunc { return Recovery(m.Log) }
func (m *Middleware) Logger() gin.HandlerFunc   { return Logger(m.Log) }
func (m *Middleware) Auth() gin.HandlerFunc     { return Auth(m.cfg) }
func (m *Middleware) Role() gin.HandlerFunc     { return RequireRole(m.cfg) }
