package clients

import (
	commongrpc "github.com/Eucastan/eucastanpay/common/pkg/grpc"

	transferpb "github.com/Eucastan/eucastanpay/common/proto/transfer"
)

func Transfer(manager *commongrpc.Manager) transferpb.TransferServiceClient {
	return transferpb.NewTransferServiceClient(
		manager.Get("transfer"),
	)
}
