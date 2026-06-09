package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/persistence/memory"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

func newTestMonitor(id string) monitor.Monitor {
	return monitor.New(id, "Test Monitor", "", "https://example.com", monitor.TypeUptime, 3, time.Minute, 5*time.Second)
}

// RN-037: status transitions move a Monitor between Active, Paused and Disabled.
func TestPauseMonitor_TransitionsToPaused(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	m := newTestMonitor("mon-1")
	if err := monitors.Save(ctx, m); err != nil {
		t.Fatalf("setup: %v", err)
	}

	uc := application.NewPauseMonitor(monitors)
	if err := uc.Execute(ctx, m.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := monitors.FindByID(ctx, m.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != monitor.StatusPaused {
		t.Fatalf("expected status Paused, got %q", got.Status)
	}
}

// RN-028: changing a Monitor's status must not disturb the at-most-one-Open-
// Incident link the system maintains (RN-002) — Pause/Resume/Disable only
// touch Monitor.Status, never Incidents.
func TestPauseMonitor_DoesNotAffectExistingOpenIncident(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	incidents := memory.NewIncidentRepository()

	m := newTestMonitor("mon-1")
	if err := monitors.Save(ctx, m); err != nil {
		t.Fatalf("setup: %v", err)
	}

	openedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := incidents.Save(ctx, incident.Open("inc-1", m.ID, openedAt)); err != nil {
		t.Fatalf("setup: %v", err)
	}

	uc := application.NewPauseMonitor(monitors)
	if err := uc.Execute(ctx, m.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	open, err := incidents.FindOpenByMonitorID(ctx, m.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if open == nil || open.ID != "inc-1" {
		t.Fatalf("expected the existing Open Incident to remain untouched (RN-028), got %+v", open)
	}
}

func TestResumeMonitor_TransitionsToActive(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	m := newTestMonitor("mon-1")
	m.Pause()
	if err := monitors.Save(ctx, m); err != nil {
		t.Fatalf("setup: %v", err)
	}

	uc := application.NewResumeMonitor(monitors)
	if err := uc.Execute(ctx, m.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := monitors.FindByID(ctx, m.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != monitor.StatusActive {
		t.Fatalf("expected status Active, got %q", got.Status)
	}
}

func TestDisableMonitor_TransitionsToDisabled(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	m := newTestMonitor("mon-1")
	if err := monitors.Save(ctx, m); err != nil {
		t.Fatalf("setup: %v", err)
	}

	uc := application.NewDisableMonitor(monitors)
	if err := uc.Execute(ctx, m.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := monitors.FindByID(ctx, m.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != monitor.StatusDisabled {
		t.Fatalf("expected status Disabled, got %q", got.Status)
	}
}
