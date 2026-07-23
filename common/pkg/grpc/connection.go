package grpc

import (
	"crypto/tls"
	"fmt"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewConnection(cfg ServiceConfig, log *logrus.Logger, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),

		grpc.WithChainUnaryInterceptor(
			interceptor.CorrelationClientInterceptor(),
		),

		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10 * 1024 * 1024),
		),
	}

	if cfg.Insecure {
		options = append(options, opts...)
	} else {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	log.Infof("Connecting to %s", cfg.Address)

	conn, err := grpc.NewClient(cfg.Address, options...)
	if err != nil {
		return nil, fmt.Errorf("%s connection failed: %w", cfg.Name, err)
	}

	log.Infof("Connected to %s", cfg.Address)

	return conn, nil
}
