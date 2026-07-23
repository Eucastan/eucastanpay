package response

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/admin/internal/domain"
)

type AdminResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminLoginResponse struct {
	Message      string        `json:"message"`
	Data         AdminResponse `json:"data"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

func ToAdminResponse(admin *domain.Admin) *AdminResponse {

	return &AdminResponse{
		ID:        admin.ID,
		Email:     admin.Email,
		FirstName: admin.FirstName,
		LastName:  admin.LastName,
		Role:      string(admin.Role),
		Status:    string(admin.Status),
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}
}
