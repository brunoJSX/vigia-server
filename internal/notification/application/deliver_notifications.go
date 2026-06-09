package application

import (
	"context"
	"time"

	"github.com/vigia/vigia-v1/internal/notification/notification"
)

type DeliverNotifications struct {
	notifications NotificationRepository
	provider      WhatsAppProvider
	clock         func() time.Time
}

func NewDeliverNotifications(notifications NotificationRepository, provider WhatsAppProvider, clk func() time.Time) *DeliverNotifications {
	return &DeliverNotifications{notifications: notifications, provider: provider, clock: clk}
}

func (uc *DeliverNotifications) Execute(ctx context.Context) error {
	pending, err := uc.notifications.FindActionable(ctx)
	if err != nil {
		return err
	}
	for _, n := range pending {
		if err := uc.deliver(ctx, n); err != nil {
			return err
		}
	}
	return nil
}

func (uc *DeliverNotifications) deliver(ctx context.Context, n notification.Notification) error {
	err := uc.provider.Send(ctx, n.Recipient, notification.Render(n))
	if err != nil {
		n.RecordFailure()
	} else {
		n.MarkDelivered(uc.clock())
	}
	return uc.notifications.Save(ctx, n)
}
