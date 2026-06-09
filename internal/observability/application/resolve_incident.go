package application

import (
	"context"

	"github.com/vigia/vigia-v1/internal/shared/clock"
)

// ResolveIncident closes a Monitor's Open Incident — the "Processo de
// Recuperação de Serviço". It is also how CheckMonitor reacts to a
// DecisionResolveIncident, so it looks the Incident up itself rather than
// requiring the caller to hold one.
type ResolveIncident struct {
	incidents     IncidentRepository
	notifications NotificationPublisher
	clock         clock.Clock
}

func NewResolveIncident(incidents IncidentRepository, notifications NotificationPublisher, c clock.Clock) *ResolveIncident {
	return &ResolveIncident{incidents: incidents, notifications: notifications, clock: c}
}

func (uc *ResolveIncident) Execute(ctx context.Context, monitorID, monitorName string) error {
	open, err := uc.incidents.FindOpenByMonitorID(ctx, monitorID)
	if err != nil {
		return err
	}
	if open == nil {
		return nil
	}

	now := uc.clock()
	open.Resolve(now)

	if err := uc.incidents.Save(ctx, *open); err != nil {
		return err
	}

	return uc.notifications.Publish(ctx, Event{
		Kind:       EventIncidentResolved,
		MonitorID:  open.MonitorID,
		IncidentID: open.ID,
		OccurredAt: now,
		Payload:    IncidentResolvedPayload{MonitorName: monitorName, Duration: open.Duration()},
	})
}
