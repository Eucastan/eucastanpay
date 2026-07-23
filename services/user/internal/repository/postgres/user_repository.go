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

type UserRepository struct {
	db        *pgxpool.Pool
	telemetry *telemetry.Telemetry
}

func NewUserRepository(db *pgxpool.Pool, telemetry *telemetry.Telemetry) *UserRepository {
	return &UserRepository{
		db:        db,
		telemetry: telemetry,
	}
}

func (r *UserRepository) WithTX(ctx context.Context, fn func(tx pgx.Tx) error) error {
	ctx, span := r.telemetry.Start(ctx, "UserRepository.WithTX")
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

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	ctx, span := r.telemetry.Start(ctx, "UserRepository.Create")
	defer span.End()

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
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	ctx, span := r.telemetry.Start(ctx, "UserRepository.FindAll")
	defer span.End()

	query := `
        SELECT id, email, phone, first_name, last_name, date_of_birth, role, status, email_verified, created_at, updated_at
        FROM users 
        ORDER BY created_at DESC
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := r.telemetry.Start(ctx, "UserRepository.FindByEmail")
	defer span.End()

	query := `
	SELECT id, email, Phone, first_name, last_name, password_hash, date_of_birth, role, status, email_verified, created_at, updated_at
	FROM users 
	WHERE email = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Phone, &user.FirstName, &user.LastName,
		&user.Password, &user.DateOfBirth, &user.Role, &user.Status, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		span.RecordError(err)
		return nil, errmessage.ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	ctx, span := r.telemetry.Start(ctx, "UserRepository.FindByID")
	defer span.End()

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
		span.RecordError(err)
		return nil, errmessage.ErrUserNotFound
	}

	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	ctx, span := r.telemetry.Start(ctx, "UserRepository.Update")
	defer span.End()

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
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userID string) error {
	ctx, span := r.telemetry.Start(ctx, "UserRepository.Delete")
	defer span.End()

	query := `
	     DELETE FROM users 
		 WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err == pgx.ErrNoRows {
		return errmessage.ErrUserNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
	ctx, span := r.telemetry.Start(ctx, "UserRepository.SaveOutboxEvent")
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
		return err
	}

	return nil
}
