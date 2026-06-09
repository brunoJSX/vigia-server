package incident

import "time"

// Status represents the lifecycle state of an Incident (RN-038).
type Status string

const (
	StatusOpen     Status = "open"
	StatusResolved Status = "resolved"
)

// Incident represents an operational situation relevant to the client —
// unavailability, slowness, dependency failure. It carries the duration of
// the problem it represents (RN-012).
//
// SequenceNumber is assigned by the persistence layer (DB serial) and used
// to produce human-readable identifiers (e.g. INC-42). Zero until persisted.
type Incident struct {
	ID             string
	MonitorID      string
	Status         Status
	OpenedAt       time.Time
	ResolvedAt     *time.Time
	SequenceNumber int
}

// Open creates an Incident in state Open (RN-038).
func Open(id, monitorID string, openedAt time.Time) Incident {
	return Incident{
		ID:        id,
		MonitorID: monitorID,
		Status:    StatusOpen,
		OpenedAt:  openedAt,
	}
}

// Resolve transitions the Incident to Resolved and fixes its ResolvedAt.
// There is deliberately no transition back to Open (RN-038: an Incident in
// state Resolved never returns to Open) — calling Resolve again is a no-op.
func (i *Incident) Resolve(resolvedAt time.Time) {
	if i.Status == StatusResolved {
		return
	}
	i.Status = StatusResolved
	i.ResolvedAt = &resolvedAt
}

// Duration returns how long the problem represented by this Incident lasted
// (RN-012). It is zero while the Incident remains Open.
func (i Incident) Duration() time.Duration {
	if i.ResolvedAt == nil {
		return 0
	}
	return i.ResolvedAt.Sub(i.OpenedAt)
}
