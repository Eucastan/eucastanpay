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
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Password      string     `json:"-"`
	DateOfBirth   time.Time  `json:"date_of_birth" binding:"omitempty"`
	Role          string     `json:"role"`
	Status        UserStatus `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
