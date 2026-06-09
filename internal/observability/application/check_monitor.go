package application

import (
	"context"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/analyzer"
	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/observability/incident"
	"github.com/vigia/vigia-v1/internal/observability/monitor"
	"github.com/vigia/vigia-v1/internal/shared/clock"
	"github.com/vigia/vigia-v1/internal/shared/id"
)

// CheckMonitor is the workflow CheckMonitor — it runs the processing pipeline
// for a single Monitor: Collector coleta → Analyzer avalia → Decision →
// Incident é consequência (RN-001, RN-002, RN-027).
type CheckMonitor struct {
	monitors        MonitorRepository
	incidents       IncidentRepository
	samples         SampleRepository
	collector       collector.Collector
	analyzer        analyzer.Analyzer
	resolveIncident *ResolveIncident
	notifications   NotificationPublisher
	clock           clock.Clock
	ids             id.Generator
}

func NewCheckMonitor(
	monitors MonitorRepository,
	incidents IncidentRepository,
	samples SampleRepository,
	c collector.Collector,
	a analyzer.Analyzer,
	resolveIncident *ResolveIncident,
	notifications NotificationPublisher,
	clk clock.Clock,
	ids id.Generator,
) *CheckMonitor {
	return &CheckMonitor{
		monitors:        monitors,
		incidents:       incidents,
		samples:         samples,
		collector:       c,
		analyzer:        a,
		resolveIncident: resolveIncident,
		notifications:   notifications,
		clock:           clk,
		ids:             ids,
	}
}

func (uc *CheckMonitor) Execute(ctx context.Context, monitorID string) error {
	m, err := uc.monitors.FindByID(ctx, monitorID)
	if err != nil {
		return err
	}

	// Pré-condição do processo de Verificação de Monitor: só Monitors Active
	// são verificados.
	if !m.IsActive() {
		return nil
	}

	// Collector timeout is a technical detail derived from Monitor config.
	// Checkout uses AcceptableResponseTime * factor so slow-but-responding
	// requests can still be measured against the threshold (RN-025).
	// Other types use a fixed internal default — no domain-level time config.
	const collectorTimeoutFactor = 3
	const defaultCollectorTimeout = 10 * time.Second

	collectorTimeout := defaultCollectorTimeout
	if m.Type == monitor.TypeCheckout && m.AcceptableResponseTime > 0 {
		collectorTimeout = m.AcceptableResponseTime * collectorTimeoutFactor
	}

	sample, err := uc.collector.Collect(ctx, m.Target, collectorTimeout)
	if err != nil {
		return err
	}

	if err := uc.samples.Save(ctx, m.ID, sample); err != nil {
		return err
	}

	recent, err := uc.samples.FindRecent(ctx, m.ID, m.Threshold)
	if err != nil {
		return err
	}

	open, err := uc.incidents.FindOpenByMonitorID(ctx, m.ID)
	if err != nil {
		return err
	}

	decision := uc.analyzer.Analyze(m, recent, open != nil)
	now := uc.clock()

	switch decision {
	case analyzer.DecisionOpenIncident:
		inc := incident.Open(uc.ids(), m.ID, now)
		if err := uc.incidents.Save(ctx, inc); err != nil {
			return err
		}
		return uc.notifications.Publish(ctx, Event{
			Kind:       EventIncidentOpened,
			MonitorID:  m.ID,
			IncidentID: inc.ID,
			OccurredAt: now,
		})
	case analyzer.DecisionResolveIncident:
		return uc.resolveIncident.Execute(ctx, m.ID)
	default:
		return nil
	}
}
