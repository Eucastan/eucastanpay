package idempotency

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Store interface {
	IsEventProcessedTx(ctx context.Context, tx pgx.Tx, eventID string) (bool, error)
	MarkEventProcessedTx(ctx context.Context, tx pgx.Tx, id, eventID, topic string) error
}
