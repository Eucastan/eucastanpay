package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	auditpb "github.com/Eucastan/eucastanpay/common/proto/audit"
	"github.com/Eucastan/eucastanpay/services/audit/internal/grpc/server"
	"google.golang.org/grpc"
)

func (a *App) initGRPCServer() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(a.cfg.SharedCfg.JWTSecret)),
	)

	gs := server.NewAuditServiceServer(a.uc)
	auditpb.RegisterAuditServiceServer(grpcServer, gs)

	a.grpcServ = grpcServer
	a.grpcserver = gs
}
