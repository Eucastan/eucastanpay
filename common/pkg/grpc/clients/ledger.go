package clients

import (
	commongrpc "github.com/Eucastan/eucastanpay/common/pkg/grpc"

	ledgerpb "github.com/Eucastan/eucastanpay/common/proto/ledger"
)

func Ledger(manager *commongrpc.Manager) ledgerpb.LedgerServiceClient {
	return ledgerpb.NewLedgerServiceClient(
		manager.Get("ledger"),
	)
}
