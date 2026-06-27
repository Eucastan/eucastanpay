package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/services/account/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type AccountRepository struct {
	DB     *pgxpool.Pool
	logger *logrus.Logger
}

func NewAccountRepository(db *pgxpool.Pool, logger *logrus.Logger) *AccountRepository {
	return &AccountRepository{
		DB:     db,
		logger: logger,
	}
}

func (r *AccountRepository) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
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

	if err = fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *AccountRepository) LockAccount(ctx context.Context, tx pgx.Tx, accID string, accNo int64) (*domain.Account, error) {
	start := time.Now()

	query := `
        SELECT id, user_id, account_no, balance, account_type, currency, status, created_at, updated_at
        FROM accounts
        WHERE id = $1 AND account_no = $2
        FOR UPDATE;
    `

	acc := &domain.Account{}
	err := tx.QueryRow(ctx, query, accID, accNo).Scan(
		&acc.ID, &acc.UserID, &acc.AccountNo, &acc.Balance,
		&acc.AccountType, &acc.Currency, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errmessage.ErrAccNotFound
	}

	duration := time.Since(start)
	r.logger.WithFields(logrus.Fields{
		"query":       "success",
		"duration_ms": duration,
	}).Info("lock operation query")
	return acc, err
}

func (r *AccountRepository) Create(ctx context.Context, tx pgx.Tx, acc *domain.Account) error {
	query := `
        INSERT INTO accounts (id, user_id, balance, account_type, currency, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING account_no, created_at, updated_at;
    `

	return tx.QueryRow(ctx, query,
		acc.ID,
		acc.UserID,
		acc.Balance,
		acc.AccountType,
		acc.Currency,
		acc.Status,
		acc.CreatedAt,
	).Scan(&acc.AccountNo, &acc.CreatedAt, &acc.UpdatedAt)
}

func (r *AccountRepository) Exists(ctx context.Context, userID, accType string) (bool, error) {
	start := time.Now()

	query := `
        SELECT EXISTS (
			SELECT 1
			FROM accounts
			WHERE user_id = $1
			  AND account_type = $2
		)
    `

	var exists bool
	if err := r.DB.QueryRow(ctx, query, userID, accType).Scan(&exists); err != nil {
		return false, err
	}

	duration := time.Since(start)

	r.logger.WithFields(logrus.Fields{
		"query":       "Exists",
		"duration_ms": duration.Milliseconds(),
	}).Info("database query")

	return exists, nil
}

func (r *AccountRepository) FindAll(ctx context.Context) ([]domain.Account, error) {
	query := `
        SELECT id, user_id, account_no, balance, account_type, currency, status, created_at, updated_at 
        FROM accounts 
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Account])
}

func (r *AccountRepository) FindByID(ctx context.Context, accID string) (*domain.Account, error) {
	query := `
        SELECT id, user_id, account_no, balance, account_type, currency, status, created_at, updated_at
        FROM accounts
        WHERE id = $1
    `

	acc := &domain.Account{}
	err := r.DB.QueryRow(ctx, query, accID).Scan(
		&acc.ID, &acc.UserID, &acc.AccountNo, &acc.Balance,
		&acc.AccountType, &acc.Currency, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errmessage.ErrAccNotFound
	}
	return acc, err
}

func (r *AccountRepository) FindByIDTX(ctx context.Context, tx pgx.Tx, accID string, accNo int64) (*domain.Account, error) {
	query := `
        SELECT id, user_id, account_no, balance, account_type, currency, status, created_at, updated_at
        FROM accounts
        WHERE id = $1 AND account_no = $2
    `

	acc := &domain.Account{}
	err := tx.QueryRow(ctx, query, accID, accNo).Scan(
		&acc.ID, &acc.UserID, &acc.AccountNo, &acc.Balance,
		&acc.AccountType, &acc.Currency, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errmessage.ErrAccNotFound
	}
	return acc, err
}

func (r *AccountRepository) FindByUserID(ctx context.Context, userID string) (*domain.Account, error) {
	query := `
        SELECT id, user_id, account_no, balance, account_type, currency, status, created_at, updated_at
        FROM accounts
        WHERE user_id = $1
    `

	acc := &domain.Account{}
	err := r.DB.QueryRow(ctx, query, userID).Scan(
		&acc.ID, &acc.UserID, &acc.AccountNo, &acc.Balance,
		&acc.AccountType, &acc.Currency, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errmessage.ErrAccNotFound
	}
	return acc, err
}

func (r *AccountRepository) FindByAccountIDAndUserID(ctx context.Context, accID, userID string) (*domain.Account, error) {
	query := `
        SELECT id, user_id, account_no, balance, account_type, currency, status, created_at, updated_at
        FROM accounts
        WHERE id = $1 AND user_id = $2
        LIMIT 1
    `

	acc := &domain.Account{}
	err := r.DB.QueryRow(ctx, query, accID, userID).Scan(
		&acc.ID, &acc.UserID, &acc.AccountNo, &acc.Balance,
		&acc.AccountType, &acc.Currency, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errmessage.ErrAccNotFound
	}
	return acc, err
}

func (r *AccountRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, accID, status string) error {

	query := `
		UPDATE accounts 
		SET status = $2, updated_at = NOW()
		WHERE id = $1;
	`

	cmd, err := tx.Exec(ctx, query, accID, status)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errmessage.ErrAccNotFound
	}
	return nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, tx pgx.Tx, accID string, amount int64, isCredit bool) error {
	var query string
	if isCredit {
		query = `
            UPDATE accounts 
            SET balance = balance + $2, 
                updated_at = NOW()
            WHERE id = $1;
        `
	} else {
		query = `
            UPDATE accounts 
            SET balance = balance - $2, 
                updated_at = NOW()
            WHERE id = $1 
              AND balance >= $2;
        `
	}

	cmd, err := tx.Exec(ctx, query, accID, amount)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		if !isCredit {
			return errmessage.ErrInsufficientAmount
		}
		return errmessage.ErrAccNotFound
	}
	return nil
}

func (r *AccountRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
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
