package clients

import (
	commongrpc "github.com/Eucastan/eucastanpay/common/pkg/grpc"

	userpb "github.com/Eucastan/eucastanpay/common/proto/user"
)

func User(manager *commongrpc.Manager) userpb.UserServiceClient {
	return userpb.NewUserServiceClient(
		manager.Get("user"),
	)
}
