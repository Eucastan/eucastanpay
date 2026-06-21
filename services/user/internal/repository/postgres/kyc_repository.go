package postgres

import (
	"context"
	"encoding/json"

	"github.com/Eucastan/eucastanpay/common/pkg/errors"
	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KYCRepository struct {
	db *pgxpool.Pool
}

func NewKYCRepository(db *pgxpool.Pool) *KYCRepository {
	return &KYCRepository{db: db}
}

func (r *KYCRepository) WithTX(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
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
		return err
	}

	return tx.Commit(ctx)
}

func (r *KYCRepository) Create(ctx context.Context, kyc *domain.KYC) error {
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

	return err
}

func (r *KYCRepository) FindByID(ctx context.Context, id string) (*domain.KYC, error) {
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
		return nil, errors.ErrUserNotFound
	}

	return kyc, err
}

func (r *KYCRepository) FindByUserID(ctx context.Context, userID string) (*domain.KYC, error) {
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
		return nil, errors.ErrUserNotFound
	}

	return kyc, err
}

func (r *KYCRepository) Update(ctx context.Context, kyc *domain.KYC) error {
	query := `
	    UPDATE kycs 
		SET status = $2
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, kyc.ID, kyc.Status)
	return err
}

func (r *KYCRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO outbox (id, topic, key, payload)
        VALUES ($1, $2, $3, $4)
    `

	_, err = tx.Exec(ctx, query, uuid.NewString(), topic, key, data)
	return err
}
