package postgres

import (
	"context"

	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/notification/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	DB        *pgxpool.Pool
	telemetry *telemetry.Telemetry
}

func NewNotificationRepository(db *pgxpool.Pool, telemetry *telemetry.Telemetry) *NotificationRepository {
	return &NotificationRepository{DB: db, telemetry: telemetry}
}

func (r *NotificationRepository) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.WithTx")
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

func (r *NotificationRepository) CreateTx(ctx context.Context, tx pgx.Tx, n *domain.Notification) error {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.Create")
	defer span.End()

	query := `
		INSERT INTO notifications (id, user_id, title, message, channel, type, priority, reference, metadata, status, scheduled_for)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at
	`

	metadata, err := producer.Encode(n.Metadata)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return tx.QueryRow(ctx, query,
		n.ID, n.UserID, n.Title, n.Message, n.Channel, n.Type, n.Priority,
		n.Reference, metadata, n.Status, n.ScheduledFor,
	).Scan(&n.CreatedAt)
}

func (r *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.Create")
	defer span.End()

	query := `
		INSERT INTO notifications (id, user_id, title, message, channel, type, priority, reference, metadata, status, scheduled_for)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at
	`

	metadata, err := producer.Encode(n.Metadata)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return r.DB.QueryRow(ctx, query,
		n.ID, n.UserID, n.Title, n.Message, n.Channel, n.Type, n.Priority,
		n.Reference, metadata, n.Status, n.ScheduledFor,
	).Scan(&n.CreatedAt)
}

func (r *NotificationRepository) UpdateStatusTx(ctx context.Context, tx pgx.Tx, id string, status string) error {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.UpdateStatus")
	defer span.End()

	query := `
	  UPDATE notifications 
	  SET status = $1, sent_at = NOW() 
	  WHERE id = $2
	`
	if _, err := tx.Exec(ctx, query, status, id); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (r *NotificationRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.UpdateStatus")
	defer span.End()

	query := `
	  UPDATE notifications 
	  SET status = $1, sent_at = NOW() 
	  WHERE id = $2
	`
	if _, err := r.DB.Exec(ctx, query, status, id); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (r *NotificationRepository) FindByUserID(ctx context.Context, userID string, limit int) ([]domain.Notification, error) {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.FindByUserID")
	defer span.End()

	query := `
		SELECT id, user_id, title, message, channel, type, priority, reference, metadata, status, created_at, sent_at, scheduled_for FROM notifications 
		WHERE user_id = $1 
		ORDER BY created_at DESC LIMIT $2
	`

	rows, err := r.DB.Query(ctx, query, userID, limit)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Notification])
}

func (r *NotificationRepository) FindByReference(ctx context.Context, reference string) ([]domain.Notification, error) {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.FindByReference")
	defer span.End()

	query := `
		SELECT id, user_id, title, message, channel, type, priority, reference, metadata, status, created_at, sent_at, scheduled_for FROM notifications 
		WHERE reference = $1 
		ORDER BY created_at DESC LIMIT 2
	`

	rows, err := r.DB.Query(ctx, query, reference)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Notification])
}

func (r *NotificationRepository) SaveTemplate(ctx context.Context, temp *domain.NotificationTemplate) error {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.SaveTemplate")
	defer span.End()

	query := `
		INSERT INTO notification_templates (id, name, subject, body, channel)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	if err := r.DB.QueryRow(ctx, query, &temp.ID, &temp.Name,
		&temp.Subject, &temp.Body, &temp.Channel,
	).Scan(&temp.ID); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *NotificationRepository) FindTemplateByName(ctx context.Context, name string) (*domain.NotificationTemplate, error) {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.FindTemplateByName")
	defer span.End()

	query := `
		SELECT id, name, subject, body, channel FROM notification_templates
		WHERE name = $1 
		ORDER BY created_at
	`

	temp := &domain.NotificationTemplate{}
	if err := r.DB.QueryRow(ctx, query, name).Scan(&temp.ID, &temp.Name, &temp.Subject,
		&temp.Body, &temp.Channel,
	); err != nil {
		span.RecordError(err)
		return nil, err
	}

	return temp, nil
}

func (r *NotificationRepository) SaveOutboxEvent(ctx context.Context, tx pgx.Tx, topic, key string, payload interface{}) error {
	ctx, span := r.telemetry.Start(ctx, "NotificationRepository.SaveOutboxEvent")
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
