package util

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/google/uuid"
)

func KYCDbType(userID string, input *request.KYCRequest) *domain.KYC {
	kyc := &domain.KYC{
		ID:        uuid.NewString(),
		UserID:    userID,
		IDType:    input.IDType,
		IDNumber:  input.IDNumber,
		Status:    domain.StatusKycPending,
		CreatedAt: time.Now(),
	}

	return kyc
}
