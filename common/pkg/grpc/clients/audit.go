package clients

import (
	commongrpc "github.com/Eucastan/eucastanpay/common/pkg/grpc"

	auditpb "github.com/Eucastan/eucastanpay/common/proto/audit"
)

func Audit(manager *commongrpc.Manager) auditpb.AuditServiceClient {
	return auditpb.NewAuditServiceClient(
		manager.Get("audit"),
	)
}
