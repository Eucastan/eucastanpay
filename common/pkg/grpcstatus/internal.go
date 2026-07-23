package grpcstatus

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func internal(err error) error {
	return status.Error(codes.Internal, err.Error())
}
