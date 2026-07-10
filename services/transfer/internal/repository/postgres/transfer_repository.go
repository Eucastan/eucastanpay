package postgres

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransferRepository struct {
	DB        *pgxpool.Pool
	telemetry *telemetry.Telemetry
}

func NewTransferRepository(db *pgxpool.Pool, telemetry *telemetry.Telemetry) *TransferRepository {
	return &TransferRepository{
		DB:        db,
		telemetry: telemetry,
	}
}

func (r *TransferRepository) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.WithTX")
	defer span.End()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		span.RecordError(err)
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
		span.RecordError(err)
		return err
	}

	return tx.Commit(ctx)
}

func (r *TransferRepository) Create(ctx context.Context, tx pgx.Tx, t *domain.Transfer) error {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.Create")
	defer span.End()

	r.telemetry.Count(ctx, 1)

	query := `
        INSERT INTO transfers (
            id, user_id, reference, step, from_account_id, from_account_no, to_account_id,
            to_account_no, amount, description, idempotency_key, direction, status, mode, 
			from_balance_after, to_balance_after, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
        RETURNING created_at, updated_at;
    `

	err := tx.QueryRow(ctx, query,
		t.ID, t.UserID, t.Reference, t.Step, t.FromAccID, t.FromAccNo, t.ToAccID,
		t.ToAccNo, t.Amount, t.Description, t.IdempotencyKey, t.Direction, t.Status,
		t.Mode, t.FromBalanceAfter, t.ToBalanceAfter, t.CreatedAt,
	).Scan(&t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			// Unique constraint violation
			if pgErr.Code == "23505" {
				return errmessage.ErrDuplicateRequest
			}
		}

		r.telemetry.RecordError(span, err)

		return err
	}
	return nil
}

func (r *TransferRepository) FindAll(ctx context.Context) ([]domain.Transfer, error) {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.FindAll")
	defer span.End()

	query := `
        SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id,  
               to_account_no, amount, description, idempotency_key, direction, status, 
               mode, reversal_ref, is_reversed, from_balance_after, to_balance_after, 
               created_at, updated_at
        FROM transfers 
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Transfer])
}

func (r *TransferRepository) FindByIdempotencyKey(ctx context.Context, idemKey string) (*domain.Transfer, error) {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.FindByIdempotencyKey")
	defer span.End()

	query := `SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id, 
               to_account_no, amount, description, idempotency_key, direction, status, 
               mode, reversal_ref, is_reversed, from_balance_after, to_balance_after, 
               created_at, updated_at
	 FROM transfers 
	 WHERE idempotency_key = $1`

	transfer := &domain.Transfer{}
	err := r.DB.QueryRow(ctx, query, idemKey).Scan(
		&transfer.ID, &transfer.UserID, &transfer.Reference, &transfer.Step, &transfer.FromAccID,
		&transfer.FromAccNo, &transfer.ToAccID, &transfer.ToAccNo, &transfer.Amount,
		&transfer.Description, &transfer.IdempotencyKey, &transfer.Direction, &transfer.Status,
		&transfer.Mode, &transfer.ReversalRef, &transfer.IsReversed, &transfer.FromBalanceAfter,
		&transfer.ToBalanceAfter, &transfer.CreatedAt, &transfer.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		span.RecordError(err)
		return nil, errmessage.ErrTranferNotFound
	}
	return transfer, err
}

func (r *TransferRepository) FindByReference(ctx context.Context, tx pgx.Tx, ref string) (*domain.Transfer, error) {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.FindByReference")
	defer span.End()

	query := `SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id,  
		to_account_no, amount, description, idempotency_key, direction, status, 
		mode, reversal_ref, is_reversed, from_balance_after, to_balance_after, 
		created_at, updated_at
	    FROM transfers 
	    WHERE reference = $1
	`
	transfer := &domain.Transfer{}
	err := tx.QueryRow(ctx, query, ref).Scan(
		&transfer.ID, &transfer.UserID, &transfer.Reference, &transfer.Step, &transfer.FromAccID,
		&transfer.FromAccNo, &transfer.ToAccID, &transfer.ToAccNo, &transfer.Amount,
		&transfer.Description, &transfer.IdempotencyKey, &transfer.Direction, &transfer.Status,
		&transfer.Mode, &transfer.ReversalRef, &transfer.IsReversed, &transfer.FromBalanceAfter,
		&transfer.ToBalanceAfter, &transfer.CreatedAt, &transfer.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		span.RecordError(err)
		return nil, errmessage.ErrTranferNotFound
	}
	return transfer, err
}

func (r *TransferRepository) FindByReferenceNoTx(ctx context.Context, reference string) (*domain.Transfer, error) {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.FindByReferenceNoTx")
	defer span.End()

	query := `SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id, to_account_no, 
		amount, description, idempotency_key, direction, status, mode, reversal_ref, 
		is_reversed, from_balance_after, to_balance_after, created_at, updated_at
	    FROM transfers 
	    WHERE reference = $1
	`
	transfer := &domain.Transfer{}
	err := r.DB.QueryRow(ctx, query, reference).Scan(
		&transfer.ID, &transfer.UserID, &transfer.Reference, &transfer.Step, &transfer.FromAccID,
		&transfer.FromAccNo, &transfer.ToAccID, &transfer.ToAccNo, &transfer.Amount,
		&transfer.Description, &transfer.IdempotencyKey, &transfer.Direction, &transfer.Status,
		&transfer.Mode, &transfer.ReversalRef, &transfer.IsReversed, &transfer.FromBalanceAfter,
		&transfer.ToBalanceAfter, &transfer.CreatedAt, &transfer.UpdatedAt,
	)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return transfer, nil
}

func (r *TransferRepository) FindByID(ctx context.Context, id string) (*domain.Transfer, error) {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.FindByID")
	defer span.End()

	query := `SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id, to_account_no, 
		amount, description, idempotency_key, direction, status, mode, reversal_ref, 
		is_reversed, from_balance_after, to_balance_after, created_at, updated_at
        FROM transfers 
        WHERE id = $1`

	transfer := &domain.Transfer{}
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&transfer.ID, &transfer.UserID, &transfer.Reference, &transfer.Step, &transfer.FromAccID,
		&transfer.FromAccNo, &transfer.ToAccID, &transfer.ToAccNo, &transfer.Amount,
		&transfer.Description, &transfer.IdempotencyKey, &transfer.Direction, &transfer.Status,
		&transfer.Mode, &transfer.ReversalRef, &transfer.IsReversed, &transfer.FromBalanceAfter,
		&transfer.ToBalanceAfter, &transfer.CreatedAt, &transfer.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		span.RecordError(err)
		return nil, errmessage.ErrTranferNotFound
	}
	return transfer, err
}

func (r *TransferRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, ref, status string) error {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.UpdateStatus")
	defer span.End()

	query := `UPDATE transfers 
	  SET status = $2, updated_at = NOW() 
	  WHERE reference = $1 
	`
	cmd, err := tx.Exec(ctx, query, ref, status)
	if err != nil {
		span.RecordError(err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		span.RecordError(err)
		return errmessage.ErrNothingUpdated
	}
	return nil
}

func (r *TransferRepository) UpdateAfterDebit(ctx context.Context, tx pgx.Tx, ref string, FromBalanceAfter int64) error {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.UpdateAfterDebit")
	defer span.End()

	query := `UPDATE transfers 
	  SET from_balance_after = $2, updated_at = NOW() 
	  WHERE reference = $1 
	`
	cmd, err := tx.Exec(ctx, query, ref, FromBalanceAfter)
	if err != nil {
		span.RecordError(err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		span.RecordError(err)
		return errmessage.ErrNothingUpdated
	}
	return nil
}

func (r *TransferRepository) UpdateAfterCredit(ctx context.Context, tx pgx.Tx, ref string, ToBalanceAfter int64) error {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.UpdateAfterCredit")
	defer span.End()

	query := `UPDATE transfers 
	  SET to_balance_after = $2, updated_at = NOW() 
	  WHERE reference = $1
	`
	cmd, err := tx.Exec(ctx, query, ref, ToBalanceAfter)
	if err != nil {
		span.RecordError(err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		span.RecordError(err)
		return errmessage.ErrNothingUpdated
	}
	return nil
}

func (r *TransferRepository) UpdateStep(ctx context.Context, tx pgx.Tx, ref string, step string) error {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.UpdateStep")
	defer span.End()

	query := `UPDATE transfers 
	  SET step = $2, updated_at = NOW() 
	  WHERE reference = $1 
	`
	cmd, err := tx.Exec(ctx, query, ref, step)
	if err != nil {
		span.RecordError(err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		span.RecordError(err)
		return errmessage.ErrNothingUpdated
	}
	return nil
}

func (r *TransferRepository) MarkAsReversed(ctx context.Context, tx pgx.Tx, ref string) error {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.MarkAsReversed")
	defer span.End()

	query := `
	  UPDATE transfers 
	  SET is_reversed = TRUE, status = $2, updated_at = NOW() 
	  WHERE reference = $1 AND is_reversed = FALSE
	`
	cmd, err := tx.Exec(ctx, query, ref, domain.TransferStatusReversed)
	if err != nil {
		span.RecordError(err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		span.RecordError(err)
		return errmessage.ErrNothingUpdated
	}
	return nil
}

func (r *TransferRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.SaveOutboxEvent")
	defer span.End()

	data, err := producer.Encode(payload)
	if err != nil {
		span.RecordError(err)
		return err
	}
	query := `INSERT INTO outbox (id, topic, key, payload) VALUES ($1, $2, $3, $4)`
	if _, err = tx.Exec(ctx, query, uuid.NewString(), topic, key, data); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *TransferRepository) IncrementRecoveryCount(ctx context.Context, tx pgx.Tx, reference string) error {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.IncrementRecoveryCount")
	defer span.End()

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
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *TransferRepository) FindStuckTransfers(ctx context.Context, timeout time.Duration) ([]domain.Transfer, error) {
	ctx, span := r.telemetry.Start(ctx, "TransferRepository.FindStuckTransfers")
	defer span.End()

	query := `
		SELECT id, user_id, reference, step, from_account_id, from_account_no, to_account_id, to_account_no, 
			amount, description, idempotency_key, direction, status, mode, reversal_ref,
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
		span.RecordError(err)
		return nil, err
	}

	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Transfer])
}
