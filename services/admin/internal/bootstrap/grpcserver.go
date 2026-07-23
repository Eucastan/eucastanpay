package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	adminpb "github.com/Eucastan/eucastanpay/common/proto/admin"
	"github.com/Eucastan/eucastanpay/services/admin/internal/grpcserver"
	"google.golang.org/grpc"
)

func (a *App) initGRPCServer() {
	server := grpc.NewServer(grpc.UnaryInterceptor(
		interceptor.AuthInterceptor(a.cfg.JWTSecret),
	))

	gs := grpcserver.NewAdminServiceServer(a.uc)
	adminpb.RegisterAdminServiceServer(server, gs)

	a.grpcServ = server
	a.grpcserver = gs
}
