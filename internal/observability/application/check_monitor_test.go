package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/analyzer"
	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/persistence/memory"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

func newCheckMonitor(
	monitors *memory.MonitorRepository,
	incidents *memory.IncidentRepository,
	samples *memory.SampleRepository,
	c *fakeCollector,
	publisher *spyPublisher,
	now time.Time,
) *application.CheckMonitor {
	resolveIncident := application.NewResolveIncident(incidents, publisher, fixedClock(now))
	return application.NewCheckMonitor(
		monitors, incidents, samples,
		c, analyzer.NewThresholdAnalyzer(),
		resolveIncident, publisher,
		fixedClock(now), sequentialIDs("inc"),
	)
}

// RN-001 / RN-027: threshold consecutive failures open an Incident.
func TestCheckMonitor_OpensIncidentAfterConsecutiveFailures(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	incidents := memory.NewIncidentRepository()
	samples := memory.NewSampleRepository()
	publisher := &spyPublisher{}

	m := monitor.New("mon-1", "acc-1", "Test Monitor", "", "https://example.com", monitor.TypeUptime, 3, time.Minute, 5*time.Second)
	if err := monitors.Save(ctx, m); err != nil {
		t.Fatalf("setup: %v", err)
	}

	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	fc := &fakeCollector{results: []bool{false, false, false}, start: now}
	uc := newCheckMonitor(monitors, incidents, samples, fc, publisher, now)

	for i := 0; i < 3; i++ {
		if err := uc.Execute(ctx, m.ID); err != nil {
			t.Fatalf("unexpected error on check %d: %v", i, err)
		}
	}

	open, err := incidents.FindOpenByMonitorID(ctx, m.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if open == nil {
		t.Fatalf("expected an Open Incident after %d consecutive failures (RN-027)", m.Threshold)
	}

	if len(publisher.events) != 1 || publisher.events[0].Kind != application.EventIncidentOpened {
		t.Fatalf("expected a single IncidentOpened event, got %+v", publisher.events)
	}
}

// RN-002: a Monitor never has more than one Incident Open simultaneously.
func TestCheckMonitor_DoesNotOpenSecondIncidentWhileOneIsOpen(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	incidents := memory.NewIncidentRepository()
	samples := memory.NewSampleRepository()
	publisher := &spyPublisher{}

	m := monitor.New("mon-1", "acc-1", "Test Monitor", "", "https://example.com", monitor.TypeUptime, 2, time.Minute, 5*time.Second)
	if err := monitors.Save(ctx, m); err != nil {
		t.Fatalf("setup: %v", err)
	}

	openedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := incidents.Save(ctx, incident.Open("inc-existing", m.ID, openedAt)); err != nil {
		t.Fatalf("setup: %v", err)
	}

	now := openedAt.Add(time.Hour)
	fc := &fakeCollector{results: []bool{false, false, false, false}, start: now}
	uc := newCheckMonitor(monitors, incidents, samples, fc, publisher, now)

	for i := 0; i < 4; i++ {
		if err := uc.Execute(ctx, m.ID); err != nil {
			t.Fatalf("unexpected error on check %d: %v", i, err)
		}
	}

	all, err := incidents.FindByMonitorAndPeriod(ctx, m.ID, time.Time{}, time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected at most one Incident per Monitor (RN-002), got %d", len(all))
	}
}

// RN-027: threshold consecutive successes resolve the Open Incident.
func TestCheckMonitor_ResolvesIncidentAfterConsecutiveSuccesses(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	incidents := memory.NewIncidentRepository()
	samples := memory.NewSampleRepository()
	publisher := &spyPublisher{}

	m := monitor.New("mon-1", "acc-1", "Test Monitor", "", "https://example.com", monitor.TypeUptime, 2, time.Minute, 5*time.Second)
	if err := monitors.Save(ctx, m); err != nil {
		t.Fatalf("setup: %v", err)
	}

	openedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := incidents.Save(ctx, incident.Open("inc-existing", m.ID, openedAt)); err != nil {
		t.Fatalf("setup: %v", err)
	}

	now := openedAt.Add(time.Hour)
	fc := &fakeCollector{results: []bool{true, true}, start: now}
	uc := newCheckMonitor(monitors, incidents, samples, fc, publisher, now)

	for i := 0; i < 2; i++ {
		if err := uc.Execute(ctx, m.ID); err != nil {
			t.Fatalf("unexpected error on check %d: %v", i, err)
		}
	}

	open, err := incidents.FindOpenByMonitorID(ctx, m.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if open != nil {
		t.Fatalf("expected the Incident to be resolved, but it is still Open: %+v", open)
	}

	if len(publisher.events) != 1 || publisher.events[0].Kind != application.EventIncidentResolved {
		t.Fatalf("expected a single IncidentResolved event, got %+v", publisher.events)
	}
}

// Pré-condição do processo de Verificação de Monitor: Monitors fora de Active não são checados.
func TestCheckMonitor_SkipsInactiveMonitor(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	incidents := memory.NewIncidentRepository()
	samples := memory.NewSampleRepository()
	publisher := &spyPublisher{}

	m := monitor.New("mon-1", "acc-1", "Test Monitor", "", "https://example.com", monitor.TypeUptime, 3, time.Minute, 5*time.Second)
	m.Pause()
	if err := monitors.Save(ctx, m); err != nil {
		t.Fatalf("setup: %v", err)
	}

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fc := &fakeCollector{start: now}
	uc := newCheckMonitor(monitors, incidents, samples, fc, publisher, now)

	if err := uc.Execute(ctx, m.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fc.calls != 0 {
		t.Fatalf("expected the Collector not to run for a non-Active Monitor, got %d calls", fc.calls)
	}
	if len(publisher.events) != 0 {
		t.Fatalf("expected no events for a non-Active Monitor, got %+v", publisher.events)
	}
}
