package application

import (
	"context"
	"time"

	"github.com/vigia/vigia-v1/internal/notification/notification"
	"github.com/vigia/vigia-v1/internal/shared/id"
)

type EnqueueNotification struct {
	notifications NotificationRepository
	ids           id.Generator
	clock         func() time.Time
}

func NewEnqueueNotification(notifications NotificationRepository, ids id.Generator, clk func() time.Time) *EnqueueNotification {
	return &EnqueueNotification{notifications: notifications, ids: ids, clock: clk}
}

type EnqueueInput struct {
	Type      notification.Type
	Recipient string
	Payload   notification.Payload
}

func (uc *EnqueueNotification) Execute(ctx context.Context, in EnqueueInput) error {
	if in.Recipient == "" {
		return nil // RN-N006: no recipient → no notification
	}
	n := notification.New(uc.ids(), in.Type, in.Recipient, in.Payload, uc.clock())
	return uc.notifications.Save(ctx, n)
}
