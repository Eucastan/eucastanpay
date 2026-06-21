package postgres

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/services/audit/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Filter struct {
	CorrelationID string
	Reference     string
	EventType     string
	MinAmount     int64
	MaxAmount     int64
	FromDate      *time.Time
	ToDate        *time.Time
	Limit         int
	Offset        int
}

type AuditRepository struct {
	DB *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{DB: db}
}

func (r *AuditRepository) WithTX(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil || err != nil {
			_ = tx.Rollback(ctx)

		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *AuditRepository) Insert(ctx context.Context, tx pgx.Tx, log *domain.AuditLog) error {
	query := `
        INSERT INTO audit_logs (id, event_type, correlation_id, causation_id, reference, payload, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.Exec(ctx, query, log.ID, log.EventType, log.CorrelationID,
		log.CausationID, log.Reference, log.Payload, log.CreatedAt)
	return err
}

func (r *AuditRepository) InsertRead(ctx context.Context, tx pgx.Tx, read *domain.AuditRead) error {
	query := `
        INSERT INTO audit_read (id, event_type, service, correlation_id, causation_id, reference, 
                               account_id, user_id, amount, status, payload, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := tx.Exec(ctx, query, read.ID, read.EventType, read.Service, read.CorrelationID, read.CausationID,
		read.Reference, read.AccountID, read.UserID, read.Amount, read.Status, read.Payload, read.CreatedAt)
	return err
}

func (r *AuditRepository) Search(ctx context.Context, f Filter) ([]domain.AuditRead, error) {
	if f.Limit <= 0 {
		f.Limit = 100
	}

	query := `
        SELECT id, event_type, service, correlation_id, causation_id, reference, account_id, user_id, amount, status, payload, created_at 
        FROM audit_read
        WHERE ($1 = '' OR correlation_id = $1)
          AND ($2 = '' OR reference = $2)
          AND ($3 = '' OR event_type = $3)
          AND ($4 = 0 OR amount >= $4)
          AND ($5 = 0 OR amount <= $5)
          AND ($6::timestamptz IS NULL OR created_at >= $6::timestamptz)
          AND ($7::timestamptz IS NULL OR created_at <= $7::timestamptz)
        ORDER BY created_at DESC
        LIMIT $8 OFFSET $9`

	rows, err := r.DB.Query(ctx, query,
		f.CorrelationID, f.Reference, f.EventType,
		f.MinAmount, f.MaxAmount, f.FromDate, f.ToDate,
		f.Limit, f.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.AuditRead])
}

func (r *AuditRepository) FindByID(ctx context.Context, id string) (*domain.AuditRead, error) {
	query := `
	    SELECT id, event_type, service, correlation_id, causation_id, reference, account_id, user_id, amount, status, payload, created_at
		FROM audit_read
		WHERE id = $1
	`

	auditLog := &domain.AuditRead{}
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&auditLog.ID,
		&auditLog.EventType,
		&auditLog.Service,
		&auditLog.CorrelationID,
		&auditLog.CausationID,
		&auditLog.Reference,
		&auditLog.AccountID,
		&auditLog.UserID,
		&auditLog.Amount,
		&auditLog.Status,
		&auditLog.Payload,
		&auditLog.CreatedAt,
	)
	return auditLog, err
}
