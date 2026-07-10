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
	ID        string      `db:"id" json:"id"`
	UserID    string      `db:"user_id" json:"user_id"`
	Token     string      `db:"token" json:"token"`
	TokenType TypeOfToken `db:"token_type" json:"token_type"`
	ExpiredAt time.Time   `db:"expired_at" json:"expire_at"`
	Revoked   bool        `db:"revoked" json:"revoked"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`

	ParentID *string `db:"parent_id" json:"parent_id"`
}
