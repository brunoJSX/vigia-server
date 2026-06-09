package application

import (
	"context"

	"github.com/vigia/vigia-v1/internal/notification/notification"
)

type NotificationRepository interface {
	Save(ctx context.Context, n notification.Notification) error
	FindActionable(ctx context.Context) ([]notification.Notification, error)
}

type WhatsAppProvider interface {
	Send(ctx context.Context, recipient, message string) error
}
