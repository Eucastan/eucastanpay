package grpcstatus

import (
	"errors"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToTransferStatus(err error) error {
	if err == nil {
		return nil
	}

	switch {

	case errors.Is(err, errmessage.ErrTranferNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, errmessage.ErrLedgerNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, errmessage.ErrUserNotOwner):
		return status.Error(codes.PermissionDenied, err.Error())

	case errors.Is(err, errmessage.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, errmessage.ErrCanNotTransferToSelf):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, errmessage.ErrInvalidTransferMode):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, errmessage.ErrNothingUpdated):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, errmessage.ErrPreviousTransferFailure):
		return status.Error(codes.Aborted, err.Error())

	case errors.Is(err, errmessage.ErrCannotReverseNonSuccessfulTransfer):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, errmessage.ErrAlreadyReversed):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, errmessage.ErrDuplicateRequest):
		return status.Error(codes.AlreadyExists, err.Error())

	default:
		return internal(err)
	}
}
