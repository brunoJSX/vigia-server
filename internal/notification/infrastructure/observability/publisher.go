package observability

import (
	"context"
	"fmt"
	"log"
	"time"

	obsapp "github.com/vigia/vigia-v1/internal/observability/application"

	notifapp "github.com/vigia/vigia-v1/internal/notification/application"
	"github.com/vigia/vigia-v1/internal/notification/notification"
)

// Publisher adapts observability events to notification enqueue calls.
type Publisher struct {
	enqueue   *notifapp.EnqueueNotification
	recipient string
}

func NewPublisher(enqueue *notifapp.EnqueueNotification, recipient string) *Publisher {
	return &Publisher{enqueue: enqueue, recipient: recipient}
}

func (p *Publisher) Publish(ctx context.Context, e obsapp.Event) error {
	log.Printf("notification publisher: received event kind=%s recipient=%q", e.Kind, p.recipient)

	var input notifapp.EnqueueInput
	input.Recipient = p.recipient

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
		log.Printf("notification publisher: ignoring unknown event kind=%s", e.Kind)
		return nil
	}

	err := p.enqueue.Execute(ctx, input)
	log.Printf("notification publisher: enqueue result type=%s err=%v", input.Type, err)
	return err
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
