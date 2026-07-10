package domain

import (
	"time"
)

type UserStatus string

const (
	StatusPending   UserStatus = "pending"
	StatusActive    UserStatus = "active"
	StatusSuspended UserStatus = "suspended"
	StatusClosed    UserStatus = "closed"
)

type User struct {
	ID            string     `db:"id" json:"id"`
	Email         string     `db:"email" json:"email"`
	Phone         string     `db:"phone" json:"phone"`
	FirstName     string     `db:"first_name" json:"first_name"`
	LastName      string     `db:"last_name" json:"last_name"`
	Password      string     `db:"-" json:"-"`
	DateOfBirth   string     `db:"date_of_birth" json:"date_of_birth,omitempty"`
	Role          string     `db:"role" json:"role"`
	Status        UserStatus `db:"status" json:"status"`
	EmailVerified bool       `db:"email_verified" json:"email_verified"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}
