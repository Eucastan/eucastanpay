package service

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/services/user/config"
	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/user/internal/repository"
	"github.com/Eucastan/eucastanpay/services/user/internal/util"
	"github.com/jackc/pgx/v5"
)

type KYCUseCase struct {
	Kyc repository.KYCRepository
	cfg *config.Config
}

func NewKYCUseCase(
	kyc repository.KYCRepository,
	cfg *config.Config,
) *KYCUseCase {
	return &KYCUseCase{
		Kyc: kyc,
		cfg: cfg,
	}
}

func (u *KYCUseCase) CreateKYC(ctx context.Context, userID string, input *request.KYCRequest) error {
	_, err := u.Kyc.FindByUserID(ctx, userID)
	if err == nil {
		return errmessage.ErrKYCAlreadyExists
	}

	kyc := util.KYCDbType(userID, input)
	if err := u.Kyc.Create(ctx, kyc); err != nil {
		return err
	}

	return u.Kyc.WithTX(ctx, func(tx pgx.Tx) error {

		kycEvent := events.UserKYCVerifiedEvent{
			UserID:    kyc.UserID,
			KYCStatus: string(kyc.Status),
			Timestamp: time.Now().Unix(),
		}

		return u.Kyc.SaveOutboxEvent(ctx, tx, events.TopicUserKYCVerified, kyc.ID, kycEvent)
	})
}

func (u *KYCUseCase) GetKYC(ctx context.Context, userID string) (*response.KYCResponse, error) {
	kyc, err := u.Kyc.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &response.KYCResponse{
		Message: "kyc ready",
		Status:  string(kyc.Status),
	}, nil
}

func (u *KYCUseCase) ApproveKYC(ctx context.Context, userID string) error {
	kyc, err := u.Kyc.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	kyc.Status = domain.StatusApproved
	kyc.VerifiedAt = &now

	return u.Kyc.Update(ctx, kyc)
}
