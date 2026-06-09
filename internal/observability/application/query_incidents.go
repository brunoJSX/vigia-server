package application

import (
	"context"

	"github.com/vigia/vigia-v1/internal/observability/incident"
)

// IncidentView enriches an Incident with the name of the Monitor it belongs to.
type IncidentView struct {
	Incident    incident.Incident
	MonitorName string
}

type QueryIncidentsInput struct {
	// Status filters by Incident status. Empty string returns all incidents.
	Status incident.Status
	// Limit caps the number of results. 0 means no limit.
	Limit int
}

type QueryIncidents struct {
	incidents IncidentRepository
	monitors  MonitorRepository
}

func NewQueryIncidents(incidents IncidentRepository, monitors MonitorRepository) *QueryIncidents {
	return &QueryIncidents{incidents: incidents, monitors: monitors}
}

func (uc *QueryIncidents) Execute(ctx context.Context, in QueryIncidentsInput) ([]IncidentView, error) {
	incidents, err := uc.incidents.FindByStatus(ctx, in.Status, in.Limit)
	if err != nil {
		return nil, err
	}
	if len(incidents) == 0 {
		return nil, nil
	}

	// Collect unique monitor IDs, then fetch names in one pass.
	seen := make(map[string]bool, len(incidents))
	for _, i := range incidents {
		seen[i.MonitorID] = true
	}
	names := make(map[string]string, len(seen))
	for monitorID := range seen {
		m, err := uc.monitors.FindByID(ctx, monitorID)
		if err == nil {
			names[monitorID] = m.Name
		}
	}

	views := make([]IncidentView, len(incidents))
	for i, inc := range incidents {
		views[i] = IncidentView{
			Incident:    inc,
			MonitorName: names[inc.MonitorID],
		}
	}
	return views, nil
}
