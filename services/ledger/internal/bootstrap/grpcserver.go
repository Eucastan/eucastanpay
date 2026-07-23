package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	ledgerpb "github.com/Eucastan/eucastanpay/common/proto/ledger"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/grpcserver"
	"google.golang.org/grpc"
)

func (a *App) initGRPCServer() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.AuthInterceptor(a.cfg.SharedCfg.JWTSecret),
		),
	)

	gs := grpcserver.NewLedgerServiceServer(a.uc)
	ledgerpb.RegisterLedgerServiceServer(grpcServer, gs)

	a.grpcServ = grpcServer
	a.grpcserver = gs
}
