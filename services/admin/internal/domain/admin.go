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
	ID           string      `db:"id" json:"id"`
	Email        string      `db:"email" json:"email"`
	PasswordHash string      `db:"password_hash" json:"password"`
	FirstName    string      `db:"first_name" json:"first_name"`
	LastName     string      `db:"last_name" json:"last_name"`
	Role         AdminRole   `db:"role" json:"role"`
	Status       AdminStatus `db:"status" json:"status"`
	LastLoginAt  *time.Time  `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updated_at"`
}

type AdminAction struct {
	ID         string    `db:"id" json:"id"`
	AdminID    string    `db:"admin_id" json:"admin_id"`
	Action     string    `db:"action" json:"action"`
	TargetType string    `db:"target_type" json:"target_type"`
	TargetID   string    `db:"target_id" json:"target_id"`
	Reason     string    `db:"reason" json:"reason"`
	Payload    any       `db:"payload" json:"payload,omitempty"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
