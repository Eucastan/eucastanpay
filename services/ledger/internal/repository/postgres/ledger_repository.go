package postgres

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerRepository struct {
	DB        *pgxpool.Pool
	telemetry *telemetry.Telemetry
}

func NewLedgerRepository(db *pgxpool.Pool, telemetry *telemetry.Telemetry) *LedgerRepository {
	return &LedgerRepository{
		DB:        db,
		telemetry: telemetry,
	}
}

func (r *LedgerRepository) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	ctx, span := r.telemetry.Start(ctx, "LedgerRepository.WithTx")
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

	if err = fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		span.RecordError(err)
		return err
	}

	return tx.Commit(ctx)
}

func (r *LedgerRepository) CreateLedgerEntry(ctx context.Context, tx pgx.Tx, entry *domain.Ledger) error {
	ctx, span := r.telemetry.Start(ctx, "LedgerRepository.CreateLedgerEntry")
	defer span.End()

	query := `
        INSERT INTO ledgers (id, user_id, account_id, amount, entry_type, reference, balance_after, description)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at;
    `

	return tx.QueryRow(ctx, query,
		entry.ID,
		entry.UserID,
		entry.AccountID,
		entry.Amount,
		entry.EntryType,
		entry.Reference,
		entry.BalanceAfter,
		entry.Description,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
}

func (r *LedgerRepository) FindByID(ctx context.Context, id string) (*domain.Ledger, error) {
	ctx, span := r.telemetry.Start(ctx, "LedgerRepository.FindByID")
	defer span.End()

	query := `
        SELECT id, user_id, account_id, amount, entry_type, reference, balance_after, 
               description, created_at, updated_at 
        FROM ledgers 
        WHERE id = $1
    `

	entry := &domain.Ledger{}
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&entry.ID, &entry.UserID, &entry.AccountID, &entry.Amount, &entry.EntryType,
		&entry.Reference, &entry.BalanceAfter, &entry.Description,
		&entry.CreatedAt, &entry.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		span.RecordError(err)
		return nil, errmessage.ErrLedgerNotFound
	}
	return entry, err
}

func (r *LedgerRepository) FindByReference(ctx context.Context, reference string) (*domain.Ledger, error) {
	ctx, span := r.telemetry.Start(ctx, "LedgerRepository.FindByReference")
	defer span.End()

	query := `
        SELECT id, user_id, account_id, amount, entry_type, reference, balance_after, 
               description, created_at, updated_at 
        FROM ledgers 
        WHERE id = $1
    `

	entry := &domain.Ledger{}
	err := r.DB.QueryRow(ctx, query, reference).Scan(
		&entry.ID, &entry.UserID, &entry.AccountID, &entry.Amount, &entry.EntryType,
		&entry.Reference, &entry.BalanceAfter, &entry.Description,
		&entry.CreatedAt, &entry.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		span.RecordError(err)
		return nil, errmessage.ErrLedgerNotFound
	}
	return entry, err
}

func (r *LedgerRepository) FindAll(ctx context.Context) ([]domain.Ledger, error) {
	ctx, span := r.telemetry.Start(ctx, "LedgerRepository.FindAll")
	defer span.End()

	query := `
        SELECT id, user_id, account_id, amount, entry_type, reference, balance_after, 
               description, created_at, updated_at 
        FROM ledgers 
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Ledger])
}

func (r *LedgerRepository) FindByEntryType(ctx context.Context, entryType string) ([]domain.Ledger, error) {
	ctx, span := r.telemetry.Start(ctx, "LedgerRepository.FindByEntryType")
	defer span.End()

	query := `
        SELECT id, user_id, account_id, amount, entry_type, reference, balance_after, 
               description, created_at, updated_at 
        FROM ledgers 
        WHERE entry_type = $1 
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(ctx, query, entryType)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Ledger])
}

func (r *LedgerRepository) SumByAccountID(ctx context.Context, accID string) (int64, error) {
	ctx, span := r.telemetry.Start(ctx, "LedgerRepository.SumByAccountID")
	defer span.End()

	query := `
        SELECT COALESCE(SUM(
            CASE 
                WHEN entry_type = 'credit' THEN amount 
                WHEN entry_type = 'debit'  THEN -amount 
            END), 0) AS balance
        FROM ledgers 
        WHERE account_id = $1
    `

	var balance int64
	err := r.DB.QueryRow(ctx, query, accID).Scan(&balance)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}
	return balance, nil
}

func (r *LedgerRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
	ctx, span := r.telemetry.Start(ctx, "LedgerRepository.SaveOutboxEvent")
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
