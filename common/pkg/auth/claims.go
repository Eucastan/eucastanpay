package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"type"`
	jwt.RegisteredClaims
}

type AdminClaims struct {
	AdminID   string `json:"admin_id"`
	Role      string `json:"role"`
	TokenType string `json:"type"`
	jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	if c.UserID == "" || c.Email == "" {
		return jwt.ErrTokenInvalidClaims
	}

	return nil
}
