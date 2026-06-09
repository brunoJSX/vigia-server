// Package application orchestrates the Observability use cases — the
// workflows from spec.yaml plus the Monitor management operations the spec
// implies but does not name explicitly (see the implementation plan's
// "gap identificado" note, governed by RN-037).
package application

import (
	"context"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

// MonitorRepository persists and retrieves Monitor configurations.
type MonitorRepository interface {
	Save(ctx context.Context, m monitor.Monitor) error
	FindByID(ctx context.Context, id string) (monitor.Monitor, error)
	FindActive(ctx context.Context) ([]monitor.Monitor, error)
	FindAllByAccount(ctx context.Context, accountID string) ([]monitor.Monitor, error)
}

// IncidentRepository persists and retrieves Incidents.
type IncidentRepository interface {
	Save(ctx context.Context, i incident.Incident) error
	// FindOpenByMonitorID returns nil when no Incident is Open for the
	// Monitor — RN-002 guarantees there is at most one.
	FindOpenByMonitorID(ctx context.Context, monitorID string) (*incident.Incident, error)
	FindByMonitorAndPeriod(ctx context.Context, monitorID string, from, to time.Time) ([]incident.Incident, error)
	FindAllOpen(ctx context.Context) ([]incident.Incident, error)
	// FindByStatus returns incidents with the given status, ordered by
	// opened_at desc. Pass limit=0 for no limit.
	FindByStatus(ctx context.Context, status incident.Status, limit int) ([]incident.Incident, error)
	// FindByPeriod returns all incidents that overlap [from, to) — i.e.
	// opened before to and not yet resolved or resolved after from.
	FindByPeriod(ctx context.Context, from, to time.Time) ([]incident.Incident, error)
}

// SampleRepository persists and retrieves the raw data the Collector
// produces — operational evidence with no identity of its own.
type SampleRepository interface {
	Save(ctx context.Context, monitorID string, s collector.Sample) error
	// FindRecent returns up to `limit` of the most recent samples, oldest first.
	FindRecent(ctx context.Context, monitorID string, limit int) ([]collector.Sample, error)
	FindByMonitorAndPeriod(ctx context.Context, monitorID string, from, to time.Time) ([]collector.Sample, error)
	// FindLastTimestamps returns the most recent sample timestamp per monitor.
	// MonitorIDs with no samples are absent from the result map.
	FindLastTimestamps(ctx context.Context, monitorIDs []string) (map[string]time.Time, error)
}

// IncidentOpenedPayload carries the data needed to notify about an opened Incident.
type IncidentOpenedPayload struct {
	MonitorName string
}

// IncidentResolvedPayload carries the data needed to notify about a resolved Incident.
type IncidentResolvedPayload struct {
	MonitorName string
	Duration    time.Duration
}

// EventKind names something relevant this context emits — overview.md
// describes the responsibility ("emitir eventos relevantes ... disponibilizar
// informações para alertas e comunicações"), but ownership of "Notification"
// itself is still open (PA-005). Naming this Event, not Notification, keeps
// that boundary honest.
type EventKind string

const (
	EventIncidentOpened    EventKind = "incident_opened"
	EventIncidentResolved  EventKind = "incident_resolved"
	EventDailySummaryReady EventKind = "daily_summary_ready"
)

// Event carries the data needed downstream to communicate something
// relevant. Payload holds kind-specific data (e.g. a DailySummary).
type Event struct {
	Kind       EventKind
	AccountID  string
	MonitorID  string
	IncidentID string
	OccurredAt time.Time
	Payload    any
}

// NotificationPublisher is the boundary to whatever turns an Event into a
// communication to the client — a stub implementation logs it for now.
type NotificationPublisher interface {
	Publish(ctx context.Context, e Event) error
}
