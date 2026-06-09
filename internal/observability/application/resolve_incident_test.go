package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/persistence/memory"
)

// RN-038 / RN-012: resolving closes the Incident and fixes its duration.
func TestResolveIncident_ClosesOpenIncidentAndPublishesEvent(t *testing.T) {
	ctx := context.Background()
	incidents := memory.NewIncidentRepository()
	publisher := &spyPublisher{}

	openedAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	resolvedAt := openedAt.Add(45 * time.Minute)
	if err := incidents.Save(ctx, incident.Open("inc-1", "mon-1", openedAt)); err != nil {
		t.Fatalf("setup: %v", err)
	}

	uc := application.NewResolveIncident(incidents, publisher, fixedClock(resolvedAt))
	if err := uc.Execute(ctx, "mon-1", "Test Monitor", "acc-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	open, err := incidents.FindOpenByMonitorID(ctx, "mon-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if open != nil {
		t.Fatalf("expected no Open Incident after resolving, got %+v", open)
	}

	all, err := incidents.FindByMonitorAndPeriod(ctx, "mon-1", openedAt.Add(-time.Hour), resolvedAt.Add(time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 1 || all[0].Status != incident.StatusResolved {
		t.Fatalf("expected the Incident to be Resolved, got %+v", all)
	}
	if all[0].Duration() != 45*time.Minute {
		t.Fatalf("expected duration of 45m (RN-012), got %s", all[0].Duration())
	}

	if len(publisher.events) != 1 || publisher.events[0].Kind != application.EventIncidentResolved {
		t.Fatalf("expected a single IncidentResolved event, got %+v", publisher.events)
	}
}

func TestResolveIncident_NoOpenIncident_DoesNothing(t *testing.T) {
	ctx := context.Background()
	incidents := memory.NewIncidentRepository()
	publisher := &spyPublisher{}

	uc := application.NewResolveIncident(incidents, publisher, fixedClock(time.Now()))
	if err := uc.Execute(ctx, "mon-1", "Test Monitor", "acc-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(publisher.events) != 0 {
		t.Fatalf("expected no event when there is nothing to resolve, got %+v", publisher.events)
	}
}
