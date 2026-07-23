package grpcstatus

import (
	"errors"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToLedgerStatus(err error) error {
	if err == nil {
		return nil
	}

	switch {

	case errors.Is(err, errmessage.ErrLedgerNotFound):
		return status.Error(codes.NotFound, err.Error())

	default:
		return internal(err)
	}
}
