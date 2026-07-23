package grpcstatus

import (
	"errors"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToAccountStatus(err error) error {
	if err == nil {
		return nil
	}

	switch {

	case errors.Is(err, errmessage.ErrAccNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, errmessage.ErrAccountDisabled):
		return status.Error(codes.PermissionDenied, err.Error())

	case errors.Is(err, errmessage.ErrNotEligibleForOperation):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, errmessage.ErrInsufficientAmount):
		return status.Error(codes.ResourceExhausted, err.Error())

	case errors.Is(err, errmessage.ErrDuplicateAccNo):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, errmessage.ErrAccAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, errmessage.ErrIdemKeyNotFound):
		return status.Error(codes.NotFound, err.Error())

	default:
		return internal(err)
	}
}
