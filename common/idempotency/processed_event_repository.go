package idempotency

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PostgresStore struct{}

func NewPostgresStore() *PostgresStore {
	return &PostgresStore{}
}

func (s *PostgresStore) IsEventProcessedTx(ctx context.Context, tx pgx.Tx, eventID string) (bool, error) {
	query := `SELECT 1 FROM processed_events WHERE event_id = $1 LIMIT 1`

	var exists int
	err := tx.QueryRow(ctx, query, eventID).Scan(&exists)

	if err == pgx.ErrNoRows {
		return false, nil
	}
	return true, err
}

func (s *PostgresStore) MarkEventProcessedTx(ctx context.Context, tx pgx.Tx, id, eventID, topic string) error {
	query := `
	INSERT INTO processed_events (id, event_id, topic)
	VALUES ($1, $2, $3)
	ON CONFLICT (event_id) DO NOTHING
	`

	_, err := tx.Exec(ctx, query, id, eventID, topic)
	return err
}
