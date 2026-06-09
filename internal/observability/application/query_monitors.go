package application

import (
	"context"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

// CurrentState is the operational health derived from a Monitor's open Incident
// (if any). It is a read-side concept — not persisted.
type CurrentState string

const (
	CurrentStateFunctioning CurrentState = "functioning"
	CurrentStateSlow        CurrentState = "slow"
	CurrentStateDown        CurrentState = "down"
)

// MonitorView enriches a Monitor with its current operational state and the
// timestamp of the last verification performed by the pipeline.
type MonitorView struct {
	Monitor       monitor.Monitor
	CurrentState  CurrentState
	LastCheckedAt *time.Time
}

type QueryMonitors struct {
	monitors  MonitorRepository
	incidents IncidentRepository
	samples   SampleRepository
}

func NewQueryMonitors(monitors MonitorRepository, incidents IncidentRepository, samples SampleRepository) *QueryMonitors {
	return &QueryMonitors{monitors: monitors, incidents: incidents, samples: samples}
}

func (uc *QueryMonitors) Execute(ctx context.Context, accountID string) ([]MonitorView, error) {
	monitors, err := uc.monitors.FindAllByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if len(monitors) == 0 {
		return nil, nil
	}

	openIncidents, err := uc.incidents.FindAllOpen(ctx)
	if err != nil {
		return nil, err
	}
	openByMonitor := make(map[string]*incident.Incident, len(openIncidents))
	for i := range openIncidents {
		openByMonitor[openIncidents[i].MonitorID] = &openIncidents[i]
	}

	ids := make([]string, len(monitors))
	for i, m := range monitors {
		ids[i] = m.ID
	}
	lastTimestamps, err := uc.samples.FindLastTimestamps(ctx, ids)
	if err != nil {
		return nil, err
	}

	views := make([]MonitorView, len(monitors))
	for i, m := range monitors {
		open := openByMonitor[m.ID]
		state := deriveCurrentState(m, open)

		var lastCheckedAt *time.Time
		if t, ok := lastTimestamps[m.ID]; ok {
			lastCheckedAt = &t
		}

		views[i] = MonitorView{
			Monitor:       m,
			CurrentState:  state,
			LastCheckedAt: lastCheckedAt,
		}
	}
	return views, nil
}

func deriveCurrentState(m monitor.Monitor, open *incident.Incident) CurrentState {
	if open == nil {
		return CurrentStateFunctioning
	}
	if m.Type == monitor.TypeCheckout {
		return CurrentStateSlow
	}
	return CurrentStateDown
}
