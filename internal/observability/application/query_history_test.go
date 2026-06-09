package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/persistence/memory"
)

// RN-016: history is scoped to the requested period — incidents and availability.
func TestQueryHistory_ReturnsIncidentsAndAvailabilityForPeriod(t *testing.T) {
	ctx := context.Background()
	incidents := memory.NewIncidentRepository()
	samples := memory.NewSampleRepository()

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)

	if err := incidents.Save(ctx, incident.Open("inc-in-period", "mon-1", from.Add(2*time.Hour))); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := incidents.Save(ctx, incident.Open("inc-out-of-period", "mon-1", from.Add(-48*time.Hour))); err != nil {
		t.Fatalf("setup: %v", err)
	}

	for i, success := range []bool{true, true, true, false} {
		s := collector.Sample{Timestamp: from.Add(time.Duration(i) * time.Hour), Success: success}
		if err := samples.Save(ctx, "mon-1", s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	uc := application.NewQueryHistory(incidents, samples)
	got, err := uc.Execute(ctx, application.QueryHistoryInput{MonitorID: "mon-1", From: from, To: to})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Incidents) != 1 || got.Incidents[0].ID != "inc-in-period" {
		t.Fatalf("expected only the Incident within the period, got %+v", got.Incidents)
	}
	if got.AvailabilityPercentage != 75 {
		t.Fatalf("expected 75%% availability (3 of 4 samples successful), got %.2f", got.AvailabilityPercentage)
	}
}
