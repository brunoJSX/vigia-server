package analyzer_test

import (
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/analyzer"
	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

func uptimeMonitor() monitor.Monitor {
	return monitor.New("m", "acc-1", "Test Monitor", "", "https://example.com", monitor.TypeUptime, 3, time.Minute, 0)
}

func checkoutMonitor(art time.Duration) monitor.Monitor {
	return monitor.New("m", "acc-1", "Test Checkout", "", "https://example.com/checkout", monitor.TypeCheckout, 3, time.Minute, art)
}

func sample(success bool) collector.Sample {
	return collector.Sample{Timestamp: time.Now(), Success: success, Latency: 100 * time.Millisecond}
}

func sampleWithLatency(success bool, latency time.Duration) collector.Sample {
	return collector.Sample{Timestamp: time.Now(), Success: success, Latency: latency}
}

func samples(results ...bool) []collector.Sample {
	out := make([]collector.Sample, 0, len(results))
	for _, ok := range results {
		out = append(out, sample(ok))
	}
	return out
}

// RN-027: opening an Incident requires `threshold` consecutive failures.
func TestAnalyze_NotEnoughConsecutiveFailures_NoAction(t *testing.T) {
	a := analyzer.NewThresholdAnalyzer()

	got := a.Analyze(uptimeMonitor(), samples(false, false), false)

	if got != analyzer.DecisionNoAction {
		t.Fatalf("expected NoAction with samples below threshold, got %q", got)
	}
}

func TestAnalyze_ThresholdConsecutiveFailures_OpensIncident(t *testing.T) {
	a := analyzer.NewThresholdAnalyzer()

	got := a.Analyze(uptimeMonitor(), samples(true, false, false, false), false)

	if got != analyzer.DecisionOpenIncident {
		t.Fatalf("expected OpenIncident at threshold consecutive failures, got %q", got)
	}
}

// RN-002: a Monitor already has an Incident Open — analyzer must not try to
// open another one (no UpdateIncident exists either; stays NoAction).
func TestAnalyze_AlreadyHasOpenIncident_DoesNotOpenAnother(t *testing.T) {
	a := analyzer.NewThresholdAnalyzer()

	got := a.Analyze(uptimeMonitor(), samples(false, false, false), true)

	if got != analyzer.DecisionNoAction {
		t.Fatalf("expected NoAction when an Incident is already Open, got %q", got)
	}
}

func TestAnalyze_ThresholdConsecutiveSuccessesWithOpenIncident_ResolvesIncident(t *testing.T) {
	a := analyzer.NewThresholdAnalyzer()

	got := a.Analyze(uptimeMonitor(), samples(false, true, true, true), true)

	if got != analyzer.DecisionResolveIncident {
		t.Fatalf("expected ResolveIncident at threshold consecutive successes, got %q", got)
	}
}

func TestAnalyze_ThresholdConsecutiveSuccessesWithoutOpenIncident_NoAction(t *testing.T) {
	a := analyzer.NewThresholdAnalyzer()

	got := a.Analyze(uptimeMonitor(), samples(true, true, true), false)

	if got != analyzer.DecisionNoAction {
		t.Fatalf("expected NoAction — nothing to resolve without an open Incident, got %q", got)
	}
}

// RN-025: Checkout Analyzer evaluates latency against AcceptableResponseTime.
// A sample that succeeded (Success: true) but exceeded the threshold counts
// as a problem (Lentidão).
func TestAnalyze_Checkout_SlowButSuccessful_CountsAsBad(t *testing.T) {
	a := analyzer.NewThresholdAnalyzer()
	m := checkoutMonitor(2 * time.Second)

	slow := []collector.Sample{
		sampleWithLatency(true, 3*time.Second),
		sampleWithLatency(true, 3*time.Second),
		sampleWithLatency(true, 3*time.Second),
	}

	got := a.Analyze(m, slow, false)

	if got != analyzer.DecisionOpenIncident {
		t.Fatalf("expected OpenIncident for consecutive slow Checkout samples, got %q", got)
	}
}

func TestAnalyze_Checkout_FastSamples_NoAction(t *testing.T) {
	a := analyzer.NewThresholdAnalyzer()
	m := checkoutMonitor(2 * time.Second)

	fast := []collector.Sample{
		sampleWithLatency(true, 500*time.Millisecond),
		sampleWithLatency(true, 800*time.Millisecond),
		sampleWithLatency(true, 1*time.Second),
	}

	got := a.Analyze(m, fast, false)

	if got != analyzer.DecisionNoAction {
		t.Fatalf("expected NoAction for fast Checkout samples, got %q", got)
	}
}

// Uptime monitors must NOT be penalised for high latency — they only care
// about reachability (Sample.Success). RN-025 applies to Checkout only.
func TestAnalyze_Uptime_HighLatencyButSuccessful_NoAction(t *testing.T) {
	a := analyzer.NewThresholdAnalyzer()
	m := uptimeMonitor()

	highLatency := []collector.Sample{
		sampleWithLatency(true, 10*time.Second),
		sampleWithLatency(true, 10*time.Second),
		sampleWithLatency(true, 10*time.Second),
	}

	got := a.Analyze(m, highLatency, false)

	if got != analyzer.DecisionNoAction {
		t.Fatalf("Uptime should not penalise high latency (RN-025 is Checkout-only), got %q", got)
	}
}
