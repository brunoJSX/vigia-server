// Package notification provides a stub NotificationPublisher — ownership of
// "Notification" is still an open question (PA-005); this implementation
// only keeps the application boundary explicit by logging emitted events.
package notification

import (
	"context"
	"log"

	"github.com/vigia/vigia-v1/internal/observability/application"
)

type StubPublisher struct {
	logger *log.Logger
}

func NewStubPublisher(logger *log.Logger) *StubPublisher {
	if logger == nil {
		logger = log.Default()
	}
	return &StubPublisher{logger: logger}
}

func (p *StubPublisher) Publish(ctx context.Context, e application.Event) error {
	p.logger.Printf("observability event: kind=%s monitor=%s incident=%s at=%s",
		e.Kind, e.MonitorID, e.IncidentID, e.OccurredAt.Format("2006-01-02T15:04:05Z07:00"))
	return nil
}
