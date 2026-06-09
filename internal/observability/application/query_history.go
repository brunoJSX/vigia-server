package application

import (
	"context"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/observability/incident"
)

// QueryHistory is the workflow QueryHistory — read-only consultation of a
// Monitor's incidents and availability over a period (RN-016).
type QueryHistory struct {
	incidents IncidentRepository
	samples   SampleRepository
}

func NewQueryHistory(incidents IncidentRepository, samples SampleRepository) *QueryHistory {
	return &QueryHistory{incidents: incidents, samples: samples}
}

type QueryHistoryInput struct {
	MonitorID string
	From, To  time.Time
}

type QueryHistoryResult struct {
	Incidents              []incident.Incident
	AvailabilityPercentage float64
}

func (uc *QueryHistory) Execute(ctx context.Context, in QueryHistoryInput) (QueryHistoryResult, error) {
	incidents, err := uc.incidents.FindByMonitorAndPeriod(ctx, in.MonitorID, in.From, in.To)
	if err != nil {
		return QueryHistoryResult{}, err
	}

	samples, err := uc.samples.FindByMonitorAndPeriod(ctx, in.MonitorID, in.From, in.To)
	if err != nil {
		return QueryHistoryResult{}, err
	}

	return QueryHistoryResult{
		Incidents:              incidents,
		AvailabilityPercentage: availability(samples),
	}, nil
}

func availability(samples []collector.Sample) float64 {
	if len(samples) == 0 {
		return 100
	}

	successful := 0
	for _, s := range samples {
		if s.Success {
			successful++
		}
	}

	return float64(successful) / float64(len(samples)) * 100
}
