// Package scheduler triggers CheckMonitor for Active Monitors whose
// configured Interval has elapsed. A simple time.Ticker poll is enough for
// this increment — the simplest architecture that works (CLAUDE.md).
// Temporal stays an infra option for later, should retries/durability
// become a real need; the domain does not depend on this choice.
package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/shared/clock"
)

type TickerScheduler struct {
	monitors     application.MonitorRepository
	checkMonitor *application.CheckMonitor
	pollInterval time.Duration
	clock        clock.Clock
	logger       *log.Logger

	lastChecked map[string]time.Time
}

func NewTickerScheduler(
	monitors application.MonitorRepository,
	checkMonitor *application.CheckMonitor,
	pollInterval time.Duration,
	clk clock.Clock,
	logger *log.Logger,
) *TickerScheduler {
	if logger == nil {
		logger = log.Default()
	}
	return &TickerScheduler{
		monitors:     monitors,
		checkMonitor: checkMonitor,
		pollInterval: pollInterval,
		clock:        clk,
		logger:       logger,
		lastChecked:  make(map[string]time.Time),
	}
}

// Run polls until ctx is cancelled. Call it in its own goroutine.
func (s *TickerScheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *TickerScheduler) tick(ctx context.Context) {
	active, err := s.monitors.FindActive(ctx)
	if err != nil {
		s.logger.Printf("scheduler: failed to list active monitors: %v", err)
		return
	}

	now := s.clock()
	for _, m := range active {
		last, checked := s.lastChecked[m.ID]
		if checked && now.Sub(last) < m.Interval {
			continue
		}

		if err := s.checkMonitor.Execute(ctx, m.ID); err != nil {
			s.logger.Printf("scheduler: check failed for monitor %q: %v", m.ID, err)
			continue
		}
		s.lastChecked[m.ID] = now
	}
}
