package postgres

import (
	"context"
	"encoding/json"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/services/admin/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil || err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *AdminRepository) Create(ctx context.Context, admin *domain.Admin) error {
	query := `
        INSERT INTO admins (id, email, password_hash, first_name, last_name, role, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		admin.ID, admin.Email, admin.PasswordHash,
		admin.FirstName, admin.LastName, admin.Role, admin.Status,
	).Scan(&admin.CreatedAt, &admin.UpdatedAt)
}

func (r *AdminRepository) FindByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	query := `
        SELECT id, email, password_hash, first_name, last_name, role, 
		  status, last_login_at, created_at, updated_at
        FROM admins WHERE email = $1`

	admin := &domain.Admin{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&admin.ID, &admin.Email, &admin.PasswordHash, &admin.FirstName,
		&admin.LastName, &admin.Role, &admin.Status, &admin.LastLoginAt,
		&admin.CreatedAt, &admin.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errmessage.ErrAdminNotFound
	}
	return admin, err
}

func (r *AdminRepository) FindByID(ctx context.Context, id string) (*domain.Admin, error) {

	query := `
		SELECT id, email, password_hash, first_name, last_name, role, 
		  status, last_login_at, created_at, updated_at
		FROM admins
		WHERE id = $1
	`

	admin := &domain.Admin{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.FirstName,
		&admin.LastName,
		&admin.Role,
		&admin.Status,
		&admin.LastLoginAt,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errmessage.ErrAdminNotFound
	}

	return admin, err
}

func (r *AdminRepository) Update(ctx context.Context, admin *domain.Admin) error {

	query := `
		UPDATE admins
		SET email = $1, password_hash = $2, first_name = $3, last_name = $4, 
		  role = $5, status = $6, last_login_at = $7, updated_at = NOW()
		WHERE id = $8
	`

	_, err := r.db.Exec(ctx, query,
		admin.Email,
		admin.PasswordHash,
		admin.FirstName,
		admin.LastName,
		admin.Role,
		admin.Status,
		admin.LastLoginAt,
		admin.UpdatedAt,
		admin.ID,
	)

	return err
}

func (r *AdminRepository) Delete(ctx context.Context, id string) error {

	query := `DELETE FROM admins WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *AdminRepository) List(ctx context.Context, limit, offset int) ([]domain.Admin, error) {

	query := `
		SELECT id, email, password_hash, first_name, last_name, role, status, 
		  last_login_at, created_at, updated_at
		FROM admins
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []domain.Admin

	for rows.Next() {
		var admin domain.Admin

		err := rows.Scan(
			&admin.ID,
			&admin.Email,
			&admin.PasswordHash,
			&admin.FirstName,
			&admin.LastName,
			&admin.Role,
			&admin.Status,
			&admin.LastLoginAt,
			&admin.CreatedAt,
			&admin.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		admins = append(admins, admin)
	}

	return admins, nil
}

func (r *AdminRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	query := `INSERT INTO outbox (id, topic, key, payload) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ctx, query, uuid.NewString(), topic, key, data)
	return err
}
