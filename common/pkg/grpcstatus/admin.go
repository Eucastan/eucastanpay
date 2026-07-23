package grpcstatus

import (
	"errors"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToAdminStatus(err error) error {

	switch {

	case errors.Is(err, errmessage.ErrDuplicateEmail):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, errmessage.ErrAdminNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, errmessage.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, errmessage.ErrBootstrapClosed):
		return status.Error(codes.PermissionDenied, err.Error())

	default:
		return status.Error(codes.Internal, err.Error())
	}
}
