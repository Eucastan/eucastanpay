package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	transferpb "github.com/Eucastan/eucastanpay/common/proto/transfer"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/grpc/server"
	"google.golang.org/grpc"
)

func (a *App) initGRPCServer() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.AuthInterceptor(a.cfg.SharedCfg.JWTSecret),
		),
	)

	gs := server.NewTransferServiceServer(a.uc)
	transferpb.RegisterTransferServiceServer(grpcServer, gs)

	a.grpcServ = grpcServer
	a.grpcserver = gs
}
