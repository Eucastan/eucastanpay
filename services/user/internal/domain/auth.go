package domain

import "time"

type TypeOfToken string

const (
	EmailToken         TypeOfToken = "email_token"
	RefreshToken       TypeOfToken = "refresh_token"
	PasswordResetToken TypeOfToken = "reset"
)

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Token     string      `json:"token"`
	TokenType TypeOfToken `json:"token_type"`
	ExpiredAt time.Time
	Revoked   bool `json:"revoked"`
	CreatedAt time.Time

	ParentID *string `json:"parent_id"`
}
