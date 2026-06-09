package application

import (
	"context"

	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

// IncidentView enriches an Incident with the name of the Monitor it belongs to.
type IncidentView struct {
	Incident    incident.Incident
	MonitorName string
}

type QueryIncidentsInput struct {
	AccountID string
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

	// Collect unique monitor IDs, fetch names and filter by account in one pass.
	seen := make(map[string]bool, len(incidents))
	for _, i := range incidents {
		seen[i.MonitorID] = true
	}
	monitorByID := make(map[string]monitor.Monitor, len(seen))
	for monitorID := range seen {
		m, err := uc.monitors.FindByID(ctx, monitorID)
		if err == nil {
			monitorByID[monitorID] = m
		}
	}

	views := make([]IncidentView, 0, len(incidents))
	for _, inc := range incidents {
		m, ok := monitorByID[inc.MonitorID]
		if !ok || m.AccountID != in.AccountID {
			continue
		}
		views = append(views, IncidentView{
			Incident:    inc,
			MonitorName: m.Name,
		})
	}
	return views, nil
}
