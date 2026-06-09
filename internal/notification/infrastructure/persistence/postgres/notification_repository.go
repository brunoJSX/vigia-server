package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vigia/vigia-v1/internal/notification/notification"
)

type NotificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{pool: pool}
}

func (r *NotificationRepository) Save(ctx context.Context, n notification.Notification) error {
	payload, err := json.Marshal(n.Payload)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO notifications (id, type, recipient, payload, status, attempts, created_at, delivered_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			status       = EXCLUDED.status,
			attempts     = EXCLUDED.attempts,
			delivered_at = EXCLUDED.delivered_at
	`, n.ID, string(n.Type), n.Recipient, payload, string(n.Status), n.Attempts, n.CreatedAt, n.DeliveredAt)
	return err
}

func (r *NotificationRepository) FindActionable(ctx context.Context) ([]notification.Notification, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, type, recipient, payload, status, attempts, created_at, delivered_at
		FROM notifications
		WHERE (status = 'pending' OR status = 'failed') AND attempts < $1
		ORDER BY created_at
	`, notification.MaxAttempts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []notification.Notification
	for rows.Next() {
		var (
			n           notification.Notification
			t           string
			status      string
			payloadJSON []byte
			deliveredAt *time.Time
		)
		if err := rows.Scan(&n.ID, &t, &n.Recipient, &payloadJSON, &status, &n.Attempts, &n.CreatedAt, &deliveredAt); err != nil {
			return nil, err
		}
		n.Type = notification.Type(t)
		n.Status = notification.Status(status)
		n.DeliveredAt = deliveredAt
		if err := json.Unmarshal(payloadJSON, &n.Payload); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}
