package memory

import (
	"context"
	"sync"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/incident"
)

type IncidentRepository struct {
	mu        sync.Mutex
	incidents map[string]incident.Incident
	seq       int
}

func NewIncidentRepository() *IncidentRepository {
	return &IncidentRepository{incidents: make(map[string]incident.Incident)}
}

func (r *IncidentRepository) Save(ctx context.Context, i incident.Incident) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if i.SequenceNumber == 0 {
		r.seq++
		i.SequenceNumber = r.seq
	}
	r.incidents[i.ID] = i
	return nil
}

// FindOpenByMonitorID returns nil when there is none — RN-002 guarantees
// at most one Open Incident per Monitor, so the first match is the answer.
func (r *IncidentRepository) FindOpenByMonitorID(ctx context.Context, monitorID string) (*incident.Incident, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, i := range r.incidents {
		if i.MonitorID == monitorID && i.Status == incident.StatusOpen {
			open := i
			return &open, nil
		}
	}
	return nil, nil
}

func (r *IncidentRepository) FindByMonitorAndPeriod(ctx context.Context, monitorID string, from, to time.Time) ([]incident.Incident, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []incident.Incident
	for _, i := range r.incidents {
		if i.MonitorID == monitorID && !i.OpenedAt.Before(from) && i.OpenedAt.Before(to) {
			out = append(out, i)
		}
	}
	return out, nil
}

func (r *IncidentRepository) FindAllOpen(ctx context.Context) ([]incident.Incident, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []incident.Incident
	for _, i := range r.incidents {
		if i.Status == incident.StatusOpen {
			out = append(out, i)
		}
	}
	return out, nil
}

func (r *IncidentRepository) FindByPeriod(ctx context.Context, from, to time.Time) ([]incident.Incident, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []incident.Incident
	for _, i := range r.incidents {
		if i.OpenedAt.Before(to) && (i.ResolvedAt == nil || i.ResolvedAt.After(from)) {
			out = append(out, i)
		}
	}
	return out, nil
}

func (r *IncidentRepository) FindByStatus(ctx context.Context, status incident.Status, limit int) ([]incident.Incident, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []incident.Incident
	for _, i := range r.incidents {
		if status == "" || i.Status == status {
			out = append(out, i)
		}
	}
	// sort newest first
	for a := 0; a < len(out)-1; a++ {
		for b := a + 1; b < len(out); b++ {
			if out[b].OpenedAt.After(out[a].OpenedAt) {
				out[a], out[b] = out[b], out[a]
			}
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
