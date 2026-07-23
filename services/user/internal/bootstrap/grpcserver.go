package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	userpb "github.com/Eucastan/eucastanpay/common/proto/user"
	"github.com/Eucastan/eucastanpay/services/user/internal/grpcserver"
	"google.golang.org/grpc"
)

func (a *App) initGRPCServer() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.AuthInterceptor(a.cfg.SharedCfg.JWTSecret),
		),
	)

	gs := grpcserver.NewUserServiceServer(a.userUC, a.kycUC)
	userpb.RegisterUserServiceServer(grpcServer, gs)

	a.grpcServ = grpcServer
	a.grpcserver = gs
}
