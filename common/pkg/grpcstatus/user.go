package grpcstatus

import (
	"errors"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToUserStatus(err error) error {
	if err == nil {
		return nil
	}

	switch {

	case errors.Is(err, errmessage.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, errmessage.ErrDuplicateEmail):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, errmessage.ErrKYCAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, errmessage.ErrKYCAlreadySubmitted):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, errmessage.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, errmessage.ErrInvalidToken):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, errmessage.ErrExpiredToken):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, errmessage.ErrPasswordNotConfirmed):
		return status.Error(codes.InvalidArgument, err.Error())

	default:
		return internal(err)
	}
}
