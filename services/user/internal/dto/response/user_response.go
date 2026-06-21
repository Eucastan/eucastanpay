package response

import (
	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
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

// Helper function to convert domain to response
func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Phone:         u.Phone,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Status:        string(u.Status),
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
	}
}
