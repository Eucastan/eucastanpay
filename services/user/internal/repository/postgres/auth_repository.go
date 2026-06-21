package postgres

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/errors"
	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Create(ctx context.Context, auth *domain.Token) error {
	query := `
		INSERT INTO tokens (id, user_id, token, token_type, expired_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, token;
	`

	err := r.db.QueryRow(ctx, query,
		auth.ID,
		auth.UserID,
		auth.Token,
		auth.TokenType,
		auth.ExpiredAt,
		auth.Revoked,
		auth.CreatedAt,
	).Scan(&auth.ID, &auth.UserID, &auth.Token)

	return err
}

func (r *AuthRepository) FindToken(ctx context.Context, token, tokenType string) (*domain.Token, error) {
	query := `
	   SELECT id, user_id, token, token_type, expired_at, revoked, created_at
	   FROM tokens 
	   WHERE token = $1 AND token_type = $2
	`

	auth := &domain.Token{}
	err := r.db.QueryRow(ctx, query, token, tokenType).Scan(
		&auth.ID, &auth.UserID, &auth.Token, &auth.TokenType,
		&auth.ExpiredAt, &auth.Revoked, &auth.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}

	return auth, err
}

func (r *AuthRepository) FindByUserID(ctx context.Context, userID, tokenType string) (*domain.Token, error) {
	query := `
	   SELECT id, user_id, token, token_type, expired_at, revoked, created_at
	   FROM tokens 
	   WHERE user_id = $1 AND token_type = $2
	`

	auth := &domain.Token{}
	err := r.db.QueryRow(ctx, query, userID, tokenType).Scan(
		&auth.ID, &auth.UserID, &auth.Token, &auth.TokenType,
		&auth.ExpiredAt, &auth.Revoked, &auth.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}

	return auth, err
}

func (r *AuthRepository) UpdateToken(ctx context.Context, id, token, tokenType string) error {
	query := `
	   UPDATE tokens 
	   SET token = $2, token_type = $3
	   WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id, token, tokenType)
	if err == pgx.ErrNoRows {
		return errors.ErrUserNotFound
	}

	return err
}

func (r *AuthRepository) Revoked(ctx context.Context, token string) error {
	query := `
	    UPDATE tokens SET revoked = True
		WHERE token = $1
	`

	_, err := r.db.Exec(ctx, query, token)
	return err
}

func (r *AuthRepository) RevokeAllByUser(ctx context.Context, userID string) error {
	query := `
		UPDATE tokens SET revoked = true
		WHERE user_id = $1 AND token_type = 'refresh_token'
	`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
