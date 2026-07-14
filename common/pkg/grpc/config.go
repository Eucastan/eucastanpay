package grpc

import "time"

type ServiceConfig struct {
	Name     string
	Address  string
	Insecure bool
	Timeout  time.Duration
	Retries  int
}
