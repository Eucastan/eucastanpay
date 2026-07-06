package postgres

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KYCRepository struct {
	db        *pgxpool.Pool
	telemetry *telemetry.Telemetry
}

func NewKYCRepository(db *pgxpool.Pool, telemetry *telemetry.Telemetry) *KYCRepository {
	return &KYCRepository{
		db:        db,
		telemetry: telemetry,
	}
}

func (r *KYCRepository) WithTX(ctx context.Context, fn func(tx pgx.Tx) error) error {
	ctx, span := r.telemetry.Start(ctx, "KYCRepository.WithTX")
	defer span.End()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return err
	}

	defer func() {
		if p := recover(); p != nil || err != nil {
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

func (r *KYCRepository) Create(ctx context.Context, kyc *domain.KYC) error {
	ctx, span := r.telemetry.Start(ctx, "KYCRepository.Create")
	defer span.End()

	query := `
		INSERT INTO kycs (id, user_id, id_type, id_number, status, verified_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, id, user_id;
	`

	err := r.db.QueryRow(ctx, query,
		kyc.ID,
		kyc.UserID,
		kyc.IDType,
		kyc.IDNumber,
		kyc.Status,
		kyc.VerifiedAt,
		kyc.CreatedAt,
	).Scan(&kyc.CreatedAt, &kyc.ID, &kyc.UserID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *KYCRepository) FindByID(ctx context.Context, id string) (*domain.KYC, error) {
	ctx, span := r.telemetry.Start(ctx, "KYCRepository.FindByID")
	defer span.End()

	query := `
	   SELECT id, user_id, id_type, id_number, status, verified_at, created_at
	   FROM kycs 
	   WHERE id = $1
	`

	kyc := &domain.KYC{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&kyc.ID, &kyc.UserID, &kyc.IDType, &kyc.IDNumber, &kyc.Status,
		&kyc.VerifiedAt, &kyc.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		span.RecordError(err)
		return nil, errmessage.ErrUserNotFound
	}

	return kyc, err
}

func (r *KYCRepository) FindByUserID(ctx context.Context, userID string) (*domain.KYC, error) {
	ctx, span := r.telemetry.Start(ctx, "KYCRepository.FindByUserID")
	defer span.End()

	query := `
	   SELECT id, user_id, id_type, id_number, status, verified_at, created_at
	   FROM kycs 
	   WHERE user_id = $1
	`

	kyc := &domain.KYC{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&kyc.ID, &kyc.UserID, &kyc.IDType, &kyc.IDNumber, &kyc.Status,
		&kyc.VerifiedAt, &kyc.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		span.RecordError(err)
		return nil, errmessage.ErrUserNotFound
	}

	return kyc, err
}

func (r *KYCRepository) Update(ctx context.Context, kyc *domain.KYC) error {
	ctx, span := r.telemetry.Start(ctx, "KYCRepository.Update")
	defer span.End()

	query := `
	    UPDATE kycs 
		SET status = $2
		WHERE id = $1
	`

	if _, err := r.db.Exec(ctx, query, kyc.ID, kyc.Status); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (r *KYCRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
	ctx, span := r.telemetry.Start(ctx, "KYCRepository.SaveOutboxEvent")
	defer span.End()

	data, err := producer.Encode(payload)
	if err != nil {
		span.RecordError(err)
		return err
	}

	query := `
        INSERT INTO outbox (id, topic, key, payload)
        VALUES ($1, $2, $3, $4)
    `

	if _, err = tx.Exec(ctx, query, uuid.NewString(), topic, key, data); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
