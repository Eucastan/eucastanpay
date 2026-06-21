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

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) WithTX(ctx context.Context, fn func(tx pgx.Tx) error) error {
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

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (id, email, Phone, first_name, last_name, password_hash, date_of_birth, role, status, email_verified, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING created_at, updated_at, id;
	`

	err := r.db.QueryRow(ctx, query,
		user.ID,
		user.Email,
		user.Phone,
		user.FirstName,
		user.LastName,
		user.Password,
		user.DateOfBirth,
		user.Role,
		user.Status,
		user.EmailVerified,
		user.CreatedAt,
	).Scan(&user.CreatedAt, &user.UpdatedAt, &user.ID)

	return err
}

func (r *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	query := `
        SELECT id, email, phone, first_name, last_name, date_of_birth, role, status, email_verified, created_at, updated_at
        FROM users 
        ORDER BY created_at DESC
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT id, email, Phone, first_name, last_name, password_hash, date_of_birth, status, email_verified, created_at, updated_at
	FROM users 
	WHERE email = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Phone, &user.FirstName, &user.LastName,
		&user.Password, &user.DateOfBirth, &user.Status, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
	    SELECT id, email, Phone, first_name, last_name, date_of_birth, role, status, email_verified, created_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Phone, &user.FirstName, &user.LastName, &user.DateOfBirth,
		&user.Role, &user.Status, &user.EmailVerified, &user.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET
		 first_name = $2, 
		 last_name = $3, 
		 password_hash = $4, 
		 status = $5, 
		 email_verified = $6, 
		 updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, user.ID,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Status,
		user.EmailVerified,
	)
	return err
}

func (r *UserRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
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
