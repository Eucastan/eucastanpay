package domain

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/admin/internal/dto/request"
	"github.com/google/uuid"
)

func ToAdminDB(passwordHash, role string, admin *request.CreateAdminRequest) *Admin {

	return &Admin{
		ID:           uuid.NewString(),
		Email:        admin.Email,
		PasswordHash: passwordHash,
		FirstName:    admin.FirstName,
		LastName:     admin.LastName,
		Role:         AdminRole(role),
		Status:       AdminStatus(StatusActive),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
