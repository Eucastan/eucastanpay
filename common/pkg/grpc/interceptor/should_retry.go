package interceptor

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func shouldRetry(err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return true
	}

	switch s.Code() {
	case codes.Unavailable,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Aborted:
		return true

	default:
		return false
	}
}
