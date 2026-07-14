package clients

import (
	commongrpc "github.com/Eucastan/eucastanpay/common/pkg/grpc"

	accountpb "github.com/Eucastan/eucastanpay/common/proto/account"
)

func Account(manager *commongrpc.Manager) accountpb.AccountServiceClient {
	return accountpb.NewAccountServiceClient(
		manager.Get("account"),
	)
}
