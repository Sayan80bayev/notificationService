package repository

import (
	"context"
	"database/sql"
	"errors"
	"notificationService/internal/model"
	"time"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *model.Notification) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Notification, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID, readAt time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type notificationRepo struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &notificationRepo{db: db}
}

// Create inserts a new notification securely
func (r *notificationRepo) Create(ctx context.Context, n *model.Notification) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	query := `
		INSERT INTO notifications 
		    (id, user_id, message, is_read, created_at, read_at)
		VALUES 
		    ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		n.ID,
		n.UserID,
		n.Message,
		n.IsRead,
		n.CreatedAt,
		n.ReadAt,
	)
	return err
}

// FindByID retrieves a notification by its ID
func (r *notificationRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Notification, error) {
	query := `
		SELECT id, user_id, message, is_read, created_at, read_at
		FROM notifications
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var n model.Notification
	if err := row.Scan(
		&n.ID,
		&n.UserID,
		&n.Message,
		&n.IsRead,
		&n.CreatedAt,
		&n.ReadAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

// FindByUserID retrieves notifications for a user with pagination
func (r *notificationRepo) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, error) {
	query := `
		SELECT id, user_id, message, is_read, created_at, read_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Message,
			&n.IsRead,
			&n.CreatedAt,
			&n.ReadAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

// MarkAsRead sets a notification as read
func (r *notificationRepo) MarkAsRead(ctx context.Context, id uuid.UUID, readAt time.Time) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, readAt, id)
	return err
}

// Delete removes a notification
func (r *notificationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
