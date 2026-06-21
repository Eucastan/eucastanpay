package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errors"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransferRepository struct {
	DB *pgxpool.Pool
}

func NewTransferRepository(db *pgxpool.Pool) *TransferRepository {
	return &TransferRepository{DB: db}
}

func (r *TransferRepository) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *TransferRepository) Create(ctx context.Context, tx pgx.Tx, t *domain.Transfer) error {
	query := `
        INSERT INTO transfers (
            id, user_id, reference, step, from_account_id, from_account_no, 
            to_account_id, to_account_no, amount, description, 
            idempotency_key, type, status, mode, from_balance_after, to_balance_after, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
        RETURNING created_at, updated_at;
    `

	err := tx.QueryRow(ctx, query,
		t.ID, t.UserID, t.Reference, t.Step, t.FromAccID, t.FromAccNo,
		t.ToAccID, t.ToAccNo, t.Amount, t.Description,
		t.IdempotencyKey, t.Type, t.Status, t.Mode,
		t.FromBalanceAfter, t.ToBalanceAfter, t.CreatedAt,
	).Scan(&t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			// Unique constraint violation
			if pgErr.Code == "23505" {
				return errors.ErrDuplicateRequest
			}
		}

		return err
	}
	return nil
}

func (r *TransferRepository) FindAll(ctx context.Context) ([]domain.Transfer, error) {
	query := `
        SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id, 
               to_account_no, amount, description, idempotency_key, type, status, 
               mode, reversal_ref, is_reversed, from_balance_after, to_balance_after, 
               created_at, updated_at
        FROM transfers 
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Transfer])
}

func (r *TransferRepository) FindByIdempotencyKey(ctx context.Context, idemKey string) (*domain.Transfer, error) {
	query := `SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id, 
               to_account_no, amount, description, idempotency_key, type, status, 
               mode, reversal_ref, is_reversed, from_balance_after, to_balance_after, 
               created_at, updated_at
	 FROM transfers 
	 WHERE idempotency_key = $1`

	transfer := &domain.Transfer{}
	err := r.DB.QueryRow(ctx, query, idemKey).Scan(
		&transfer.ID, &transfer.UserID, &transfer.Reference, &transfer.Step, &transfer.FromAccID,
		&transfer.FromAccNo, &transfer.ToAccID, &transfer.ToAccNo, &transfer.Amount,
		&transfer.Description, &transfer.IdempotencyKey, &transfer.Type, &transfer.Status,
		&transfer.Mode, &transfer.ReversalRef, &transfer.IsReversed, &transfer.FromBalanceAfter,
		&transfer.ToBalanceAfter, &transfer.CreatedAt, &transfer.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.ErrTranferNotFound
	}
	return transfer, err
}

func (r *TransferRepository) FindByReference(ctx context.Context, tx pgx.Tx, ref string) (*domain.Transfer, error) {
	query := `SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id, 
		to_account_no, amount, description, idempotency_key, type, status, 
		mode, reversal_ref, is_reversed, from_balance_after, to_balance_after, 
		created_at, updated_at
	    FROM transfers 
	    WHERE reference = $1
	`
	transfer := &domain.Transfer{}
	err := tx.QueryRow(ctx, query, ref).Scan(
		&transfer.ID, &transfer.UserID, &transfer.Reference, &transfer.Step, &transfer.FromAccID,
		&transfer.FromAccNo, &transfer.ToAccID, &transfer.ToAccNo, &transfer.Amount,
		&transfer.Description, &transfer.IdempotencyKey, &transfer.Type, &transfer.Status,
		&transfer.Mode, &transfer.ReversalRef, &transfer.IsReversed, &transfer.FromBalanceAfter,
		&transfer.ToBalanceAfter, &transfer.CreatedAt, &transfer.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.ErrTranferNotFound
	}
	return transfer, err
}

func (r *TransferRepository) FindByReferenceNoTx(ctx context.Context, reference string) (*domain.Transfer, error) {

	query := `SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id, 
		to_account_no, amount, description, idempotency_key, type, status, 
		mode, reversal_ref, is_reversed, from_balance_after, to_balance_after, 
		created_at, updated_at
	    FROM transfers 
	    WHERE reference = $1
	`
	transfer := &domain.Transfer{}
	err := r.DB.QueryRow(ctx, query, reference).Scan(
		&transfer.ID, &transfer.UserID, &transfer.Reference, &transfer.Step, &transfer.FromAccID,
		&transfer.FromAccNo, &transfer.ToAccID, &transfer.ToAccNo, &transfer.Amount,
		&transfer.Description, &transfer.IdempotencyKey, &transfer.Type, &transfer.Status,
		&transfer.Mode, &transfer.ReversalRef, &transfer.IsReversed, &transfer.FromBalanceAfter,
		&transfer.ToBalanceAfter, &transfer.CreatedAt, &transfer.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (r *TransferRepository) FindByID(ctx context.Context, id string) (*domain.Transfer, error) {
	query := `SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id, 
               to_account_no, amount, description, idempotency_key, type, status, 
               mode, reversal_ref, is_reversed, from_balance_after, to_balance_after, 
               created_at, updated_at
        FROM transfers 
        WHERE id = $1`

	transfer := &domain.Transfer{}
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&transfer.ID, &transfer.UserID, &transfer.Reference, &transfer.Step, &transfer.FromAccID,
		&transfer.FromAccNo, &transfer.ToAccID, &transfer.ToAccNo, &transfer.Amount,
		&transfer.Description, &transfer.IdempotencyKey, &transfer.Type, &transfer.Status,
		&transfer.Mode, &transfer.ReversalRef, &transfer.IsReversed, &transfer.FromBalanceAfter,
		&transfer.ToBalanceAfter, &transfer.CreatedAt, &transfer.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errors.ErrTranferNotFound
	}
	return transfer, err
}

func (r *TransferRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, ref, status string) error {
	query := `UPDATE transfers SET status = $2, updated_at = NOW() WHERE reference = $1`
	cmd, err := tx.Exec(ctx, query, ref, status)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.ErrTranferNotFound
	}
	return nil
}

func (r *TransferRepository) UpdateAfterDebit(ctx context.Context, tx pgx.Tx, ref string, FromBalanceAfter int64) error {
	query := `UPDATE transfers SET from_balance_after = $2, updated_at = NOW() WHERE reference = $1`
	cmd, err := tx.Exec(ctx, query, ref, FromBalanceAfter)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.ErrTranferNotFound
	}
	return nil
}

func (r *TransferRepository) UpdateAfterCredit(ctx context.Context, tx pgx.Tx, ref string, ToBalanceAfter int64) error {
	query := `UPDATE transfers SET to_balance_after = $2, updated_at = NOW() WHERE reference = $1`
	cmd, err := tx.Exec(ctx, query, ref, ToBalanceAfter)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.ErrTranferNotFound
	}
	return nil
}

func (r *TransferRepository) UpdateStep(ctx context.Context, tx pgx.Tx, ref string, step string) error {
	query := `UPDATE transfers SET step = $2, updated_at = NOW() WHERE reference = $1`
	cmd, err := tx.Exec(ctx, query, ref, step)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.ErrTranferNotFound
	}
	return nil
}

func (r *TransferRepository) MarkAsReversed(ctx context.Context, tx pgx.Tx, ref string) error {
	query := `UPDATE transfers SET is_reversed = TRUE, status = $2, updated_at = NOW() WHERE reference = $1 AND is_reversed = FALSE`
	cmd, err := tx.Exec(ctx, query, ref, domain.TransferStatusReverse)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.ErrTranferNotFound
	}
	return nil
}

func (r *TransferRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	query := `INSERT INTO outbox (id, topic, key, payload) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ctx, query, uuid.NewString(), topic, key, data)
	return err
}

func (r *TransferRepository) IncrementRecoveryCount(ctx context.Context, tx pgx.Tx, reference string) error {

	_, err := tx.Exec(
		ctx,
		`
        UPDATE transfers
        SET recovery_count = recovery_count + 1,
            last_recovery_at = NOW(),
            updated_at = NOW()
        WHERE reference = $1
        `,
		reference,
	)

	return err
}

func (r *TransferRepository) FindStuckTransfers(ctx context.Context, timeout time.Duration) ([]domain.Transfer, error) {
	query := `
		SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id,
			to_account_no, amount, description, idempotency_key, type, status, mode, reversal_ref,
			is_reversed, recovery_count, from_balance_after, to_balance_after, created_at, updated_at
		FROM transfers
		WHERE status = $1
		  AND updated_at < NOW() - ($2 * INTERVAL '1 second')
		  AND recovery_count < 5
		  AND (
		       last_recovery_at IS NULL
		       OR last_recovery_at < NOW() - INTERVAL '1 minute'
		  )
	`

	rows, err := r.DB.Query(
		ctx,
		query,
		domain.TransferStatusPending,
		int(timeout.Seconds()),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Transfer])
}
