package monitor_test

import (
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

// RN-037: every Monitor has a status.
func TestNew_AlwaysHasAStatus(t *testing.T) {
	m := monitor.New("mon-1", "https://example.com", monitor.TypeUptime, 3, time.Minute, 0)

	if m.Status == "" {
		t.Fatal("expected Monitor to be created with a status, got empty string")
	}
	if m.Status != monitor.StatusActive {
		t.Fatalf("expected new Monitor to start Active, got %q", m.Status)
	}
}

func TestPauseResumeDisable_TransitionStatus(t *testing.T) {
	m := monitor.New("mon-1", "https://example.com", monitor.TypeUptime, 3, time.Minute, 0)

	m.Pause()
	if m.Status != monitor.StatusPaused {
		t.Fatalf("expected Paused after Pause(), got %q", m.Status)
	}

	m.Resume()
	if m.Status != monitor.StatusActive {
		t.Fatalf("expected Active after Resume(), got %q", m.Status)
	}

	m.Disable()
	if m.Status != monitor.StatusDisabled {
		t.Fatalf("expected Disabled after Disable(), got %q", m.Status)
	}
}

// RN-024: Checkout is a specialized Monitor type — modeled as a Type value,
// behaving like any other Monitor (no parallel struct).
func TestCheckoutType_BehavesAsRegularMonitor(t *testing.T) {
	m := monitor.New("mon-2", "https://example.com/checkout", monitor.TypeCheckout, 3, time.Minute, 3*time.Second)

	if m.Type != monitor.TypeCheckout {
		t.Fatalf("expected Type Checkout, got %q", m.Type)
	}
	if m.Status != monitor.StatusActive {
		t.Fatalf("expected Checkout Monitor to follow the same status rules, got %q", m.Status)
	}
}

// RN-025: Checkout Monitor carries AcceptableResponseTime; Uptime does not
// require it (zero is valid for non-Checkout types).
func TestCheckoutMonitor_CarriesAcceptableResponseTime(t *testing.T) {
	art := 3 * time.Second
	m := monitor.New("mon-3", "https://example.com/checkout", monitor.TypeCheckout, 3, time.Minute, art)

	if m.AcceptableResponseTime != art {
		t.Fatalf("expected AcceptableResponseTime %v, got %v", art, m.AcceptableResponseTime)
	}
}

func TestUptimeMonitor_ZeroAcceptableResponseTime_IsValid(t *testing.T) {
	m := monitor.New("mon-4", "https://example.com", monitor.TypeUptime, 3, time.Minute, 0)

	if m.AcceptableResponseTime != 0 {
		t.Fatalf("expected AcceptableResponseTime to be zero for Uptime, got %v", m.AcceptableResponseTime)
	}
}
