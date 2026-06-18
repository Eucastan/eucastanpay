package clients

import (
	"time"
)

type Config struct {
	UserServiceAddr     string
	AccountServiceAddr  string
	TransferServiceAddr string
	LedgerServiceAddr   string
	AuditServiceAddr    string
	Timeout             time.Duration
	MaxRetries          int
	Insecure            bool
}
