package observability

import (
	"context"
	"fmt"
	"time"

	obsapp "github.com/vigia/vigia-v1/internal/observability/application"

	notifapp "github.com/vigia/vigia-v1/internal/notification/application"
	"github.com/vigia/vigia-v1/internal/notification/notification"
)

// RecipientResolver resolves the WhatsApp recipient number for an account.
// Implemented by account.ResolveRecipient — defined here to avoid cross-context import.
type RecipientResolver interface {
	Execute(ctx context.Context, accountID string) (string, error)
}

// Publisher adapts observability events to notification enqueue calls.
type Publisher struct {
	enqueue  *notifapp.EnqueueNotification
	resolver RecipientResolver
}

func NewPublisher(enqueue *notifapp.EnqueueNotification, resolver RecipientResolver) *Publisher {
	return &Publisher{enqueue: enqueue, resolver: resolver}
}

func (p *Publisher) Publish(ctx context.Context, e obsapp.Event) error {
	recipient, err := p.resolver.Execute(ctx, e.AccountID)
	if err != nil {
		return err
	}

	var input notifapp.EnqueueInput
	input.Recipient = recipient

	switch e.Kind {
	case obsapp.EventIncidentOpened:
		payload, _ := e.Payload.(obsapp.IncidentOpenedPayload)
		input.Type = notification.TypeIncidentOpened
		input.Payload = notification.Payload{MonitorName: payload.MonitorName}
	case obsapp.EventIncidentResolved:
		payload, _ := e.Payload.(obsapp.IncidentResolvedPayload)
		input.Type = notification.TypeIncidentResolved
		input.Payload = notification.Payload{
			MonitorName: payload.MonitorName,
			Duration:    formatDuration(payload.Duration),
		}
	default:
		return nil
	}

	return p.enqueue.Execute(ctx, input)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return ""
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%02dm", h, m)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
