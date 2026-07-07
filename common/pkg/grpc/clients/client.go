package clients

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Eucastan/eucastanpay/common/proto/account"
	"github.com/Eucastan/eucastanpay/common/proto/audit"
	"github.com/Eucastan/eucastanpay/common/proto/ledger"
	"github.com/Eucastan/eucastanpay/common/proto/transfer"
	"github.com/Eucastan/eucastanpay/common/proto/user"
)

type Clients struct {
	User     user.UserServiceClient
	Account  account.AccountServiceClient
	Transfer transfer.TransferServiceClient
	Ledger   ledger.LedgerServiceClient
	Audit    audit.AuditServiceClient

	ConnUser     *grpc.ClientConn
	ConnAccount  *grpc.ClientConn
	ConnTransfer *grpc.ClientConn
	ConnLedger   *grpc.ClientConn
	ConnAudit    *grpc.ClientConn
}

// var log *logrus.Logger

func NewClients(cfg Config, log *logrus.Logger) (*Clients, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	connUser, err := newConnection(cfg.UserServiceAddr, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	connAccount, err := newConnection(cfg.AccountServiceAddr, cfg)
	if err != nil {
		connUser.Close()
		return nil, fmt.Errorf("failed to connect to account service: %w", err)
	}

	connTransfer, err := newConnection(cfg.TransferServiceAddr, cfg)
	if err != nil {
		connUser.Close()
		connAccount.Close()
		return nil, fmt.Errorf("failed to connect to transfer service: %w", err)
	}

	connLedger, err := newConnection(cfg.LedgerServiceAddr, cfg)
	if err != nil {
		connUser.Close()
		connAccount.Close()
		connTransfer.Close()
		return nil, fmt.Errorf("failed to connect to ledger service: %w", err)
	}

	connAudit, err := newConnection(cfg.AuditServiceAddr, cfg)
	if err != nil {
		connUser.Close()
		connAccount.Close()
		connTransfer.Close()
		connLedger.Close()
		return nil, fmt.Errorf("failed to connect to ledger service: %w", err)
	}

	log.Info("All gRPC clients connected successfully")

	return &Clients{
		User:         user.NewUserServiceClient(connUser),
		Account:      account.NewAccountServiceClient(connAccount),
		Transfer:     transfer.NewTransferServiceClient(connTransfer),
		Ledger:       ledger.NewLedgerServiceClient(connLedger),
		Audit:        audit.NewAuditServiceClient(connAudit),
		ConnUser:     connUser,
		ConnAccount:  connAccount,
		ConnTransfer: connTransfer,
		ConnLedger:   connLedger,
		ConnAudit:    connAudit,
	}, nil
}

func newConnection(target string, cfg Config) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10 * 1024 * 1024),
		),
	}

	if cfg.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("gRPC client failed: %w", err)
	}

	return conn, nil
}

func (c *Clients) Close() error {
	if c.ConnUser != nil {
		c.ConnUser.Close()
	}

	if c.ConnAccount != nil {
		c.ConnAccount.Close()
	}

	if c.ConnTransfer != nil {
		c.ConnTransfer.Close()
	}

	if c.ConnLedger != nil {
		c.ConnLedger.Close()
	}

	if c.ConnAudit != nil {
		c.ConnAudit.Close()
	}

	// log.Info("All gRPC clients closed")

	return nil
}
