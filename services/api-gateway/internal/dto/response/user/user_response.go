package user

import (
	"time"
)

type UserResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone,omitempty"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}
