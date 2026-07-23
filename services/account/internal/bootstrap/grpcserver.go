package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	accountpb "github.com/Eucastan/eucastanpay/common/proto/account"
	"github.com/Eucastan/eucastanpay/services/account/internal/grpcserver"
	"google.golang.org/grpc"
)

func (a *App) initGRPCServer() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.AuthInterceptor(a.cfg.SharedCfg.JWTSecret),
		),
	)

	gs := grpcserver.NewAccountServiceServer(a.uc)
	accountpb.RegisterAccountServiceServer(grpcServer, gs)

	a.grpcServ = grpcServer
	a.grpcserver = gs
}
