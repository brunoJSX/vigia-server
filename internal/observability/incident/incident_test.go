package incident_test

import (
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/incident"
)

var openedAt = time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)

// RN-038: an Incident has two possible states, starting in Open.
func TestOpen_StartsInOpenState(t *testing.T) {
	inc := incident.Open("inc-1", "mon-1", openedAt)

	if inc.Status != incident.StatusOpen {
		t.Fatalf("expected new Incident to be Open, got %q", inc.Status)
	}
	if inc.ResolvedAt != nil {
		t.Fatal("expected ResolvedAt to be unset for an Open Incident")
	}
}

// RN-012 + RN-038: resolving transitions to Resolved and fixes the duration.
func TestResolve_TransitionsToResolvedAndCalculatesDuration(t *testing.T) {
	inc := incident.Open("inc-1", "mon-1", openedAt)
	resolvedAt := openedAt.Add(15 * time.Minute)

	inc.Resolve(resolvedAt)

	if inc.Status != incident.StatusResolved {
		t.Fatalf("expected Resolved after Resolve(), got %q", inc.Status)
	}
	if inc.Duration() != 15*time.Minute {
		t.Fatalf("expected duration of 15m, got %s", inc.Duration())
	}
}

// RN-038: a Resolved Incident never returns to Open — Resolve again is a no-op.
func TestResolve_ResolvedNeverReturnsToOpen(t *testing.T) {
	inc := incident.Open("inc-1", "mon-1", openedAt)
	firstResolve := openedAt.Add(15 * time.Minute)
	inc.Resolve(firstResolve)

	laterResolve := openedAt.Add(2 * time.Hour)
	inc.Resolve(laterResolve)

	if inc.Status != incident.StatusResolved {
		t.Fatalf("expected Incident to remain Resolved, got %q", inc.Status)
	}
	if !inc.ResolvedAt.Equal(firstResolve) {
		t.Fatalf("expected ResolvedAt to stay at first resolution %s, got %s", firstResolve, inc.ResolvedAt)
	}
}

// RN-012: an Incident still Open carries no duration yet.
func TestDuration_IsZeroWhileOpen(t *testing.T) {
	inc := incident.Open("inc-1", "mon-1", openedAt)

	if inc.Duration() != 0 {
		t.Fatalf("expected zero duration for an Open Incident, got %s", inc.Duration())
	}
}
