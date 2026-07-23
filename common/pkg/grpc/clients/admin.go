package clients

import (
	commongrpc "github.com/Eucastan/eucastanpay/common/pkg/grpc"

	adminpb "github.com/Eucastan/eucastanpay/common/proto/admin"
)

func Admin(manager *commongrpc.Manager) adminpb.AdminServiceClient {
	return adminpb.NewAdminServiceClient(
		manager.Get("admin"),
	)
}
