package grpc

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"google.golang.org/grpc/metadata"
)

type RequestContext struct {
	UserID        string
	Role          string
	RequestID     string
	CorrelationID string
}

func ContextValues(ctx context.Context) RequestContext {

	md, _ := metadata.FromIncomingContext(ctx)
	rc := RequestContext{}

	if v := md.Get(interceptor.UserIDKey); len(v) > 0 {
		rc.UserID = v[0]
	}

	if v := md.Get(interceptor.RoleKey); len(v) > 0 {
		rc.Role = v[0]
	}

	if v := md.Get(interceptor.RequestIDKey); len(v) > 0 {
		rc.RequestID = v[0]
	}

	if v := md.Get(interceptor.CorrelationKey); len(v) > 0 {
		rc.CorrelationID = v[0]
	}

	return rc

}
