// Package analyzer holds the Analyzer mechanism — part of the
// processing_pipeline (see domain-discovery.md / spec.yaml). It is an
// internal, provisional label, not a domain concept: it interprets collected
// data against a threshold and consecutiveness rule (RN-027) and produces a
// Decision about an Incident.
package analyzer

import (
	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
)

// Analyzer interprets the most recent collected samples for a Monitor and
// decides what should happen to its Incident. hasOpenIncident tells it
// whether the Monitor currently has an Incident in state Open (RN-002).
// The Analyzer is type-aware: for Checkout it evaluates latency against
// AcceptableResponseTime (Lentidão — RN-025); for other types it evaluates
// reachability (Sample.Success).
type Analyzer interface {
	Analyze(m monitor.Monitor, samples []collector.Sample, hasOpenIncident bool) Decision
}

// ThresholdAnalyzer opens an Incident once `threshold` consecutive samples
// are problematic, and resolves it once `threshold` consecutive samples are
// healthy (RN-027). Anything short of that consecutiveness is NoAction.
type ThresholdAnalyzer struct{}

func NewThresholdAnalyzer() *ThresholdAnalyzer {
	return &ThresholdAnalyzer{}
}

func (a *ThresholdAnalyzer) Analyze(m monitor.Monitor, samples []collector.Sample, hasOpenIncident bool) Decision {
	if m.Threshold <= 0 || len(samples) < m.Threshold {
		return DecisionNoAction
	}

	recent := samples[len(samples)-m.Threshold:]

	switch {
	case allBad(m, recent) && !hasOpenIncident:
		return DecisionOpenIncident
	case allGood(m, recent) && hasOpenIncident:
		return DecisionResolveIncident
	default:
		return DecisionNoAction
	}
}

// sampleIsBad returns true when a sample represents a problem for the given
// Monitor type. For Checkout: slow (latency > AcceptableResponseTime) or
// failed. For Uptime/Dependency: failed only.
func sampleIsBad(m monitor.Monitor, s collector.Sample) bool {
	if m.Type == monitor.TypeCheckout {
		return !s.Success || s.Latency > m.AcceptableResponseTime
	}
	return !s.Success
}

func allBad(m monitor.Monitor, samples []collector.Sample) bool {
	for _, s := range samples {
		if !sampleIsBad(m, s) {
			return false
		}
	}
	return true
}

func allGood(m monitor.Monitor, samples []collector.Sample) bool {
	for _, s := range samples {
		if sampleIsBad(m, s) {
			return false
		}
	}
	return true
}
