package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/persistence/memory"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

// RN-037: a Monitor always has a defined status — newly created ones start Active.
func TestCreateMonitor_PersistsActiveMonitor(t *testing.T) {
	ctx := context.Background()
	monitors := memory.NewMonitorRepository()
	uc := application.NewCreateMonitor(monitors, sequentialIDs("mon"))

	got, err := uc.Execute(ctx, application.CreateMonitorInput{
		Target:    "https://example.com/health",
		Type:      monitor.TypeUptime,
		Threshold: 3,
		Interval:  time.Minute,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !got.IsActive() {
		t.Fatalf("expected new Monitor to be Active (RN-037), got status %q", got.Status)
	}

	saved, err := monitors.FindByID(ctx, got.ID)
	if err != nil {
		t.Fatalf("monitor was not persisted: %v", err)
	}
	if saved.Target != "https://example.com/health" {
		t.Fatalf("persisted monitor has wrong target: %q", saved.Target)
	}
}
