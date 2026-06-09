package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/persistence/memory"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

// Consolidates the previous day's incidents and availability per active Monitor.
func TestBuildDailySummary_ConsolidatesPreviousDayPerActiveMonitor(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	incidents := memory.NewIncidentRepository()
	samples := memory.NewSampleRepository()
	publisher := &spyPublisher{}

	now := time.Date(2026, 6, 8, 3, 0, 0, 0, time.UTC)
	yesterday := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)

	active := monitor.New("mon-active", "Active Monitor", "", "https://example.com", monitor.TypeUptime, 3, time.Minute, 5*time.Second)
	if err := monitors.Save(ctx, active); err != nil {
		t.Fatalf("setup: %v", err)
	}

	paused := monitor.New("mon-paused", "Paused Monitor", "", "https://example.org", monitor.TypeUptime, 3, time.Minute, 5*time.Second)
	paused.Pause()
	if err := monitors.Save(ctx, paused); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := incidents.Save(ctx, incident.Open("inc-yesterday", active.ID, yesterday.Add(2*time.Hour))); err != nil {
		t.Fatalf("setup: %v", err)
	}

	for i, success := range []bool{true, false, true} {
		s := collector.Sample{Timestamp: yesterday.Add(time.Duration(i) * time.Hour), Success: success}
		if err := samples.Save(ctx, active.ID, s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	uc := application.NewBuildDailySummary(monitors, incidents, samples, publisher, fixedClock(now))
	got, err := uc.Execute(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !got.Date.Equal(yesterday) {
		t.Fatalf("expected summary date %s, got %s", yesterday, got.Date)
	}

	if len(got.Entries) != 1 || got.Entries[0].MonitorID != active.ID {
		t.Fatalf("expected a single entry for the Active Monitor only, got %+v", got.Entries)
	}

	entry := got.Entries[0]
	if len(entry.Incidents) != 1 || entry.Incidents[0].ID != "inc-yesterday" {
		t.Fatalf("expected the previous day's Incident in the entry, got %+v", entry.Incidents)
	}

	want := float64(2) / float64(3) * 100
	if entry.AvailabilityPercentage != want {
		t.Fatalf("expected availability %.4f, got %.4f", want, entry.AvailabilityPercentage)
	}

	if len(publisher.events) != 1 || publisher.events[0].Kind != application.EventDailySummaryReady {
		t.Fatalf("expected a DailySummaryReady event, got %+v", publisher.events)
	}
}
