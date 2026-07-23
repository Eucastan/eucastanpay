package domain

import (
	"time"
)

type AdminRole string
type AdminStatus string

const (
	RoleSuperAdmin AdminRole = "super_admin"
	RoleAdmin      AdminRole = "admin"
	RoleModerator  AdminRole = "moderator"
)
const (
	StatusActive   AdminStatus = "active"
	StatusDisabled AdminStatus = "disabled"
)

type Admin struct {
	ID           string      `json:"id"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"-"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Role         AdminRole   `json:"role"`
	Status       AdminStatus `json:"status"`
	TwoFAEnabled bool        `json:"two_fa_enabled"`
	TwoFASecret  *string     `json:"-"`
	LastLoginAt  *time.Time  `json:"last_login_at,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type AdminAction struct {
	ID         string    `json:"id"`
	AdminID    string    `json:"admin_id"`
	Action     string    `json:"action"`
	TargetType string    `json:"target_type"`
	TargetID   string    `json:"target_id"`
	Reason     string    `json:"reason"`
	Payload    any       `json:"payload,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
