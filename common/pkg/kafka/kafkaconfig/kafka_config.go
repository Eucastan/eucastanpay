package kafkaconfig

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/config"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func NewMechanism(cfg config.KafkaConfig) *plain.Mechanism {
	var mechanism plain.Mechanism

	if cfg.SASL {

		mechanism = plain.Mechanism{
			Username: cfg.Username,
			Password: cfg.Password,
		}

	}

	return &mechanism
}

func NewTransport(cfg config.KafkaConfig) *kafka.Transport {
	transport := &kafka.Transport{}
	if cfg.SASL {
		transport.SASL = NewMechanism(cfg)
	}

	if cfg.TLS {

		caCert, err := os.ReadFile(cfg.CaCert)
		if err != nil {
			panic(err)
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caCert) {
			panic("failed to load kafka CA certificate")
		}

		transport.TLS = &tls.Config{
			RootCAs: pool,
		}
	}

	return transport
}

func NewDialer(cfg config.KafkaConfig) *kafka.Dialer {
	dialer := &kafka.Dialer{}

	if cfg.SASL {
		dialer.Timeout = 10 * time.Second
		dialer.DualStack = true
		dialer.SASLMechanism = NewMechanism(cfg)
	}

	if cfg.TLS {

		caCert, err := os.ReadFile(cfg.CaCert)
		if err != nil {
			panic(err)
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caCert) {
			panic("failed to load kafka CA certificate")
		}

		dialer.TLS = &tls.Config{
			RootCAs: pool,
		}
	}

	return dialer
}
